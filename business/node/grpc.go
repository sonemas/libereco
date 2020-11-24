package node

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
)

// Register implements the NetworkingServiceServer interface.
func (n *Node) Register(req *networking.RegisterRequest, stream networking.NetworkingService_RegisterServer) error {
	n.logger.Printf("Received registration request: %+v", req)
	n.AddPeer(&Peer{id: req.Id, addr: req.Addr})

	// Send my details.
	stream.Send(&networking.Node{Id: n.id, Addr: n.addr})

	// Send known peers.
	for _, peer := range n.peers {
		status := networking.Node_NODE_STATUS_UNSET
		if _, exists := n.joinedPeers[peer.id]; exists {
			status = networking.Node_NODE_STATUS_JOINED
		}
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: status}); err != nil {
			continue
		}
	}

	// Send known inactive peers.
	// It's important to send these as well on registrtion, in case a peer reconnects after a network failure.
	for _, peer := range n.FaultyPeers() {
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_FAILED}); err != nil {
			continue
		}
	}

	return nil
}

// Sync implements the NetworkingServiceServer interface.
func (n *Node) Sync(stream networking.NetworkingService_SyncServer) error {
	joinedPeers := n.JoinedPeers()
	faultyPeers := n.FaultyPeers()

	// Receive
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(stream networking.NetworkingService_SyncServer) {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				n.logger.Printf("Failed to receive from stream: %s", err)
				continue
			}

			switch msg.Status {
			case networking.Node_NODE_STATUS_JOINED:
				n.AddPeer(&Peer{id: msg.Id, addr: msg.Addr})
			case networking.Node_NODE_STATUS_FAILED:
				if n.HasPeer(msg.Id) {
					n.RemovePeer(msg.Id)
				}
			default:
				n.logger.Printf("Received unexpected status %q: %#v", msg.Status, msg)
			}
		}

		wg.Done()
	}(stream)

	// Send
	go func(stream networking.NetworkingService_SyncServer) {
		for _, peer := range joinedPeers {
			if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_JOINED}); err != nil {
				continue
			}
		}

		for _, peer := range faultyPeers {
			if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_FAILED}); err != nil {
				continue
			}
		}

		wg.Done()
	}(stream)

	return nil
}

// Ping implements the NetworkingServiceServer interface.
func (n *Node) Ping(ctx context.Context, req *networking.PingRequest) (*networking.PingReqResponse, error) {
	// If no Node is provided, the ping request is for us.
	if req.Node == nil {
		return &networking.PingReqResponse{Success: true}, nil
	}

	// Requested to ping another node.
	peer, err := DialPeer(req.Node.Addr, n.dialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "dialing peer")
	}

	res, err := peer.client.Ping(ctx, &networking.PingRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "ping request to third node")
	}

	return res, nil
}
