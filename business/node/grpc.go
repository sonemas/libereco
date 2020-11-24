package node

import (
	"context"
	"io"

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
		if _, exists := n.newPeers[peer.id]; exists {
			status = networking.Node_NODE_STATUS_JOINED
		}
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: status}); err != nil {
			continue
		}
	}

	// Send known inactive peers.
	// It's important to send these as well on registrtion, in case a peer reconnects after a network failure.
	for _, peer := range n.InactivePeers() {
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_FAILED}); err != nil {
			continue
		}
	}

	return nil
}

// Ping implements the NetworkingServiceServer interface.
func (n *Node) Ping(_ *networking.EmptyRequest, stream networking.NetworkingService_PingServer) error {
	for _, peer := range n.NewPeers() {
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_JOINED}); err != nil {
			continue
		}
	}

	for _, peer := range n.InactivePeers() {
		if err := stream.Send(&networking.Node{Id: peer.id, Addr: peer.addr, Status: networking.Node_NODE_STATUS_FAILED}); err != nil {
			continue
		}
	}

	return nil
}

// PingReq implements the NetworkingServiceServer interface.
func (n *Node) PingReq(ctx context.Context, req *networking.PingRequest) (*networking.PingReqResponse, error) {
	peer, err := DialPeer(req.Node.Addr, n.dialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "dialing peer")
	}

	stream, err := peer.client.Ping(ctx, &networking.EmptyRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "ping request to third node")
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
	}

	return &networking.PingReqResponse{Success: true}, nil
}
