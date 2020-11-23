package node

import (
	"errors"
	"log"
	"sync"

	"google.golang.org/grpc"
)

var (
	// ErrPeerNotFound is an error indicating that a peer can't be found.
	ErrPeerNotFound = errors.New("peer not found")

	// ErrPeerExists is an error indicating that a peer already exists in
	// a node's finger table.
	ErrPeerExists = errors.New("peer already exists")
)

// Node is a networking node.
type Node struct {
	mu   sync.RWMutex
	addr string

	logger        *log.Logger
	peers         map[string]*Peer
	newPeers      map[string]*Peer
	inactivePeers map[string]*Peer
	initialized   bool
	stopChan      chan struct{}
	dialOptions   []grpc.DialOption
}

// New returns an initialized node.
func New(logger *log.Logger, addr string, bootstrapNodes ...string) (*Node, error) {
	// TODO: Check if addr is valid/
	n := Node{
		logger:        logger,
		addr:          addr,
		peers:         make(map[string]*Peer),
		newPeers:      make(map[string]*Peer),
		inactivePeers: make(map[string]*Peer),
		stopChan:      make(chan struct{}, 1),
	}

	return &n, nil
}

// DialOptions sets the GRPC dial options used when making GRPC calls to peers.
func (n *Node) DialOptions(opts ...grpc.DialOption) {
	n.dialOptions = append(n.dialOptions, opts...)
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

	return nil
}

// RemovePeer will remove a peer form the node's finger table, as well as
// the tables of new/inactive peers.
func (n *Node) RemovePeer(id string) error {
	if !n.HasPeer(id) {
		return ErrPeerNotFound
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

	return nil
}

// NewPeers returns a slice of new peers and resets the interal table
// of new peers.
func (n *Node) NewPeers() []*Peer {
	peers := n.newPeers
	r := []*Peer{}

	for _, peer := range peers {
		r = append(r, peer)
	}

	n.newPeers = make(map[string]*Peer)
	return r
}

// InactivePeers returns a slice of inactive peers and resets the interal table
// of inactive peers.
func (n *Node) InactivePeers() []*Peer {
	peers := n.inactivePeers
	r := []*Peer{}

	for _, peer := range peers {
		r = append(r, peer)
	}

	n.newPeers = make(map[string]*Peer)
	return r
}