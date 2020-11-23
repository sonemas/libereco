package node

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
)

// Register implements the NetworkingServiceServer interface.
func (n *Node) Register(req *networking.RegisterRequest, stream networking.NetworkingService_RegisterServer) error {
	n.AddPeer(&Peer{id: req.Id, addr: req.Addr})

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
