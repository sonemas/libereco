package node

import (
	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
	"google.golang.org/grpc"
)

// PeerConnection is an open connection to a peer.
// conn is the connection itself.
// client is the service client for protocol functionality.
type PeerConnection struct {
	conn   *grpc.ClientConn
	client networking.NetworkingServiceClient
}

// DialPeer returns an active peer connection or an error if the peer can't be reached.
func DialPeer(addr string, opts ...grpc.DialOption) (*PeerConnection, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "dialing peer")
	}

	return &PeerConnection{
		conn:   conn,
		client: networking.NewNetworkingServiceClient(conn),
	}, nil
}

// Close closes an open peer connection and sets all internal values to nil.
func (pc *PeerConnection) Close() {
	pc.conn.Close()
	pc.conn = nil
	pc.client = nil
}

// Peer holds information about peer nodes.
type Peer struct {
	id   string
	addr string
}

// Dial connects to the peer and returns a connection for further interaction
// with the peer. Returns an error in case the connection can't be established.
func (p *Peer) Dial() (*PeerConnection, error) {
	return DialPeer(p.addr)
}
