package node

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
	"google.golang.org/grpc"
)

// Serve starts the server and blocks on server requests.
// It will also initialize the periodic ping requests to check whether
// peers are still alive.
func (n *Node) Serve() error {
	// Return an error if the node is shutting down or shutdown.
	if n.inShutdown.isSet() {
		return ErrNodeInShutdown
	}

	if n.stopped.isSet() {
		return ErrNodeIsStopped
	}

	l, err := net.Listen("tcp", n.addr)
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
				// TODO: Abstraction
				wg := sync.WaitGroup{}

				for _, peer := range n.peers {
					p, err := peer.Dial()
					if err != nil {
						// Mark peer as inactive
						n.UpdatePeer(peer, true)
						continue
					}

					ctx, cancel := context.WithTimeout(context.Background(), n.RequestTimeout)
					defer cancel()

					stream, err := p.client.Sync(ctx)
					if err != nil {
						n.logger.Printf("Ping request failed: %v", err)
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
				}
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
			n.stopped.setTrue()
			stop = true
			break
		case err := <-errChan:
			return err
		}
	}

	return nil
}

// Shutdown will gracefully shutdown the server.
func (n *Node) Shutdown(ctx context.Context) error {
	n.inShutdown.setTrue()
	n.stopChan <- struct{}{}

	for !n.stopped.isSet() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
