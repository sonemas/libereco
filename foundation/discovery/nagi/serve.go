package nagi

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
	"github.com/sonemas/libereco/foundation/discovery"
	"google.golang.org/grpc"
)

func (n *Nagi) dialFirstValidPeer(except ...*Peer) (*PeerConnection, *Peer, error) {
	for _, peer := range n.peers {
		ok := true
		for _, e := range except {
			if peer.id == e.id {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		conn, err := peer.Dial()
		if err != nil {
			continue
		}

		return conn, peer, nil
	}

	return nil, nil, fmt.Errorf("peers exhausted")
}

// checkPotentiallyFaulty will do up to two attempts to ping the provided
// potentionally faulty nodes via different nodes. When both attempts fail,
// the peer will be marked as faulty.
func (n *Nagi) checkPotentiallyFaulty(peers ...*Peer) {
	// Loop over the potentially faulty peers.
	for _, faulty := range peers {
		conn, peer, err := n.dialFirstValidPeer(peers...)
		if err != nil {
			// Use this node in case there's no other peer available.
			// This is not ideal, but can happen on small irreliable networks.
			peer = &Peer{id: n.caddr.ID, addr: n.caddr.Addr()}
			c, err := peer.Dial()
			if err != nil {
				n.logger.Printf("all nodes failed")
				continue
			}
			conn = c
		}
		defer conn.Close()

		// Attempt 1 to ping
		{
			ctx, cancel := context.WithTimeout(context.Background(), n.RequestTimeout)
			defer cancel()
			res, err := conn.client.Ping(ctx, &networking.PingRequest{Node: &networking.Node{Id: faulty.id, Addr: faulty.addr}})
			if err != nil {
				n.logger.Printf("Ping to %s/%s failed", faulty.addr, faulty.id)
			} else if res.Success {
				// We got a successful ping, so this peer is NOT faulty.
				continue
			}
		}

		// Attempt 2 to ping
		{
			// Add the previous peer to the list of exceptions and get another peer.
			conn, peer, err := n.dialFirstValidPeer(append(peers, peer)...)
			if err != nil {
				// Use this node in case there's no other peer available.
				// This is not ideal, but can happen on small irreliable networks.
				peer = &Peer{id: n.caddr.ID, addr: n.caddr.Addr()}
				c, err := peer.Dial()
				if err != nil {
					n.logger.Printf("all nodes failed")
					continue
				}
				conn = c
			}
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), n.RequestTimeout)
			defer cancel()
			res, err := conn.client.Ping(ctx, &networking.PingRequest{Node: &networking.Node{Id: faulty.id, Addr: faulty.addr}})
			if err != nil {
				n.logger.Printf("Ping to %s/%s failed", faulty.addr, faulty.id)
			} else if res.Success {
				// We got a successful ping, so this peer is NOT faulty.
				continue
			}

			// Peer failed two ping attempts via different nodes. Mark peer as faulty.
			n.MarkPeerFaulty(faulty.id)
		}
	}
}

// sync initiates a bidirectional stream for syncing known peers between
// nodes.
func (n *Nagi) sync() {
	wg := sync.WaitGroup{}

	faulty := []*Peer{}

	for _, peer := range n.peers {
		p, err := peer.Dial()
		if err != nil {
			n.logger.Printf("Peer %s/%s is potentionally faulty", peer.addr, peer.id)
			faulty = append(faulty, peer)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), n.RequestTimeout)
		defer cancel()

		stream, err := p.client.Sync(ctx)
		if err != nil {
			n.logger.Printf("Getting sync stream failed: %v", err)
		}

		joinedPeers := n.JoinedPeers()
		faultyPeers := n.FaultyPeers()

		wg.Add(2)
		go func(stream networking.NetworkingService_SyncClient, joinedPeers []*Peer, faultyPeers []*Peer) {
			for _, peer := range joinedPeers {
				if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_JOINED}); err != nil {
					n.logger.Printf("Failed to send to stream: %s", err)
					continue
				}
			}

			for _, peer := range faultyPeers {
				if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_FAILED}); err != nil {
					n.logger.Printf("Failed to send to stream: %s", err)
					continue
				}
			}

			wg.Done()
		}(stream, joinedPeers, faultyPeers)

		go func(stream networking.NetworkingService_SyncClient) {
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					n.logger.Printf("Stream error: %v", err)
				}

				switch msg.Status {
				case networking.Node_NODE_STATUS_JOINED:
					n.AddPeer(&Peer{id: msg.Id, addr: msg.Addr})
				case networking.Node_NODE_STATUS_FAILED:
					if n.HasPeer(msg.Id) {
						n.RemovePeer(msg.Id)
					}
				default:
					n.logger.Printf("Unexpected status: %v", msg.Status)
				}
			}
			wg.Done()
		}(stream)

		wg.Wait()
		stream.CloseSend() // TODO: Should perhaps be moved to sending goroutine.

		n.EmptyNewAndfaultyPeers()

		if len(faulty) > 0 {
			n.checkPotentiallyFaulty(faulty...)
		}
	}
}

// Serve starts the server and blocks on server requests.
// It will also initialize the periodic ping requests to check whether
// peers are still alive.
func (n *Nagi) Serve() error {
	// Return an error if the node is shutting down or shutdown.
	if n.inShutdown.IsSet() {
		return discovery.ErrServiceInShutdown
	}

	if n.stopped.IsSet() {
		return discovery.ErrServiceIsStopped
	}

	l, err := net.Listen(n.caddr.Proto, n.caddr.Addr())
	if err != nil {
		return errors.Wrap(err, "creating listener")
	}

	s := grpc.NewServer()
	networking.RegisterNetworkingServiceServer(s, n)

	var errChan chan error
	go func() {
		if err := s.Serve(l); err != nil {
			errChan <- errors.Wrap(err, "starting server")
		}
	}()

	go func() {
		t := time.NewTicker(n.PingInterval)
		for {
			select {
			case <-t.C:
				n.sync()
			case <-n.stopChan:
				return
			}
		}
	}()

	stop := false
	for !stop {
		select {
		case <-n.stopChan:
			s.Stop()
			n.stopped.SetTrue()
			stop = true
			break
		case err := <-errChan:
			return err
		}
	}

	return nil
}

// Shutdown will gracefully shutdown the server.
func (n *Nagi) Shutdown(ctx context.Context) error {
	n.inShutdown.SetTrue()
	n.stopChan <- struct{}{}

	for !n.stopped.IsSet() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
