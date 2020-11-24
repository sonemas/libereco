package node

import (
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
	"google.golang.org/grpc"
)

var (
	// ErrPeerNotFound is an error indicating that a peer can't be found.
	ErrPeerNotFound = errors.New("peer not found")

	// ErrPeerExists is an error indicating that a peer already exists in
	// a node's finger table.
	ErrPeerExists = errors.New("peer already exists")

	// ErrInvalidAddr is an error indicating an address is not in the correct format.
	ErrInvalidAddr = errors.New("address should be in the format host:port/id[/protocol]")

	// ErrBootstrapingFailed is an error indicating that bootstrapping failed.
	ErrBootstrapingFailed = errors.New("bootstrap process failed")
)

// SplitAddr splits the provided address into the host/port, id and optional protocol.
// The format of addresses is: host:post/id[/protocpl]. The default protocol is TCP.
func SplitAddr(addr string) (string, string, string, error) {
	p := strings.Split(addr, "/")
	if len(p) < 2 {
		return "", "", "", ErrInvalidAddr
	}

	pr := "tcp"
	if len(p) == 3 {
		pr = p[2]
	}

	return p[0], p[1], pr, nil
}

// Node is a networking node.
type Node struct {
	mu            sync.RWMutex
	id            string
	addr          string
	logger        *log.Logger
	peers         map[string]*Peer
	newPeers      map[string]*Peer
	inactivePeers map[string]*Peer
	stopChan      chan struct{}
	dialOptions   []grpc.DialOption

	// PingInterval is the interval between pings.
	PingInterval time.Duration

	// Request timeout is the maximum time for requests to last.
	RequestTimeout time.Duration

	// BootstrapNodes are the nodes used to bootstrap to the network.
	BootstrapNodes []string
}

// NodeOption is an option that can be provided to New to customize
// the node's internal values.
type NodeOption func(*Node) error

// WithBootstrapNodes is an option to provide the bootstrap nodes for the node
// to use for the bootstrap process.
func WithBootstrapNodes(v ...string) NodeOption {
	return func(n *Node) error {
		n.BootstrapNodes = v
		return nil
	}
}

// WithDialOptions is an option to provide dialoptions that will be used
// when making GRPC dial requests to peers.
func WithDialOptions(v ...grpc.DialOption) NodeOption {
	return func(n *Node) error {
		n.dialOptions = v
		return nil
	}
}

// WithPingInterval is an option to set a node's duration of ping sessions
// to peers.
func WithPingInterval(v time.Duration) NodeOption {
	return func(n *Node) error {
		n.PingInterval = v
		return nil
	}
}

// WithRequestTimeout is an option to define requests' timeout duration.
func WithRequestTimeout(v time.Duration) NodeOption {
	return func(n *Node) error {
		n.PingInterval = v
		return nil
	}
}

// New returns an initialized node.
func New(logger *log.Logger, addr string, opts ...NodeOption) (*Node, error) {
	addr, id, _, err := SplitAddr(addr)
	if err != nil {
		return nil, err
	}

	n := Node{
		logger:         logger,
		id:             id,
		addr:           addr,
		peers:          make(map[string]*Peer),
		newPeers:       make(map[string]*Peer),
		inactivePeers:  make(map[string]*Peer),
		stopChan:       make(chan struct{}, 1),
		PingInterval:   60 * time.Second,
		RequestTimeout: 20 * time.Second,
	}

	for _, opt := range opts {
		if err := opt(&n); err != nil {
			return nil, errors.Wrapf(err, "executing option %T", opt)
		}
	}

	if len(n.BootstrapNodes) > 0 {
		success := false

		for _, node := range n.BootstrapNodes {
			n.logger.Printf("Bootstrapping via %q.", node)
			addr, _, _, err := SplitAddr(node)
			if err != nil {
				n.logger.Printf("Bootstrapping via %q failed: %s.", node, err)
				continue
			}

			peer, err := DialPeer(addr, n.dialOptions...)
			if err != nil {
				n.logger.Printf("Bootstrapping via %q failed: %s.", node, err)
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), n.RequestTimeout)
			defer cancel()

			// n.logger.Printf("ID: %s, Addr: %s", id, addr)
			stream, err := peer.client.Register(ctx, &networking.RegisterRequest{Id: n.id, Addr: n.addr})
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					n.logger.Printf("Stream error: %s", err)
					continue
				}

				// Don't put own details in the finter table.
				if msg.Id == n.id {
					continue
				}

				switch msg.Status {
				case networking.Node_NODE_STATUS_FAILED:
					if n.HasPeer(msg.Id) {
						n.RemovePeer(msg.Id)
					}
				default:
					// Any other status than failed means that the peer
					// should be added.
					n.AddPeer(&Peer{id: msg.Id, addr: msg.Addr})
				}
			}
			success = true
		}

		if !success {
			n.logger.Printf("Could't successfully bootstrap node.")
			return nil, ErrBootstrapingFailed
		}
		n.logger.Printf("Successfully bootstrapped node.")
	}

	return &n, nil
}

// HasPeer returns true if a peer exists in the finger table.
func (n *Node) HasPeer(id string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, ok := n.peers[id]
	return ok
}

// GetPeer returns a pointer to Peer or and error if the peer is or isn't in
// the node's finger table.
func (n *Node) GetPeer(id string) (*Peer, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	p, ok := n.peers[id]
	if !ok {
		return nil, ErrPeerNotFound
	}

	return p, nil
}

// AddPeer adds a peer to the node's finger table and marks the peer as a new
// peer. Unlike SetPeer, AddPeer does NOT return an error if the peer already
// exists in the finger table.
func (n *Node) AddPeer(p *Peer) {
	if n.HasPeer(p.id) {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.peers[p.id] = p
	n.newPeers[p.id] = p

	// Remove from inactive list if there.
	if _, ok := n.inactivePeers[p.id]; ok {
		delete(n.inactivePeers, p.id)
	}

	n.logger.Printf("Node %s/%s has been added.", p.addr, p.id)
}

// SetPeer adds a peer to the node's finger table.
// Returns an error in case the peer already exists in the table.
func (n *Node) SetPeer(p *Peer) error {
	if n.HasPeer(p.id) {
		return ErrPeerExists
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.peers[p.id] = p

	// Remove from inactive list if there.
	if _, ok := n.inactivePeers[p.id]; ok {
		delete(n.inactivePeers, p.id)
	}

	return nil
}

// UpdatePeer updates a peer in the node's fiunger table. Use inactive to
// mark the peer as infactive.
func (n *Node) UpdatePeer(p *Peer, inactive bool) error {
	if !n.HasPeer(p.id) {
		return ErrPeerNotFound
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if inactive {
		n.inactivePeers[p.id] = p
		delete(n.peers, p.id)

		// Remove from new list if there.
		if _, ok := n.newPeers[p.id]; ok {
			delete(n.newPeers, p.id)
		}

		n.logger.Printf("Node %s/%s has been listed as inactive.", p.addr, p.id)

		return nil
	}

	n.peers[p.id] = p
	if _, ok := n.newPeers[p.id]; ok {
		n.newPeers[p.id] = p
	}

	// Remove from inactive list if there.
	if _, ok := n.inactivePeers[p.id]; ok {
		delete(n.inactivePeers, p.id)
	}

	return nil
}

// MarkPeerInactive makrs a peer as inactive via the provided id.
func (n *Node) MarkPeerInactive(id string) error {
	p, err := n.GetPeer(id)
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.inactivePeers[id] = p
	delete(n.peers, id)

	n.logger.Printf("Node %s/%s has been listed as inactive.", p.addr, p.id)

	return nil
}

// RemovePeer will remove a peer form the node's finger table, as well as
// the tables of new/inactive peers.
func (n *Node) RemovePeer(id string) error {
	p, err := n.GetPeer(id)
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.peers, id)

	// Remove from newPeers if there.
	if _, ok := n.newPeers[id]; ok {
		delete(n.newPeers, id)
	}

	// Remove from inactivePeers if there.
	if _, ok := n.inactivePeers[id]; ok {
		delete(n.inactivePeers, id)
	}

	n.logger.Printf("Node %s/%s has been listed as inactive.", p.addr, p.id)

	return nil
}

// NewPeers returns a slice of new peers and resets the interal table
// of new peers.
func (n *Node) NewPeers() []*Peer {
	n.mu.Lock()
	peers := n.newPeers
	n.mu.Unlock()

	r := []*Peer{}

	for _, peer := range peers {
		r = append(r, peer)
	}

	return r
}

// InactivePeers returns a slice of inactive peers and resets the interal table
// of inactive peers.
func (n *Node) InactivePeers() []*Peer {
	n.mu.Lock()
	peers := n.inactivePeers
	n.mu.Unlock()

	r := []*Peer{}

	for _, peer := range peers {
		r = append(r, peer)
	}

	return r
}

// EmptyNewAndInactivePeers resets the node's lists of
// new and inactive peers.
func (n *Node) EmptyNewAndInactivePeers() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.newPeers = make(map[string]*Peer)
	n.inactivePeers = make(map[string]*Peer)
}
