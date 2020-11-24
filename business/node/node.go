package node

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/messaging"
	"github.com/sonemas/libereco/foundation/atomicvar"
	"github.com/sonemas/libereco/foundation/caddr"
	"github.com/sonemas/libereco/foundation/discovery"
	"google.golang.org/grpc"
)

type Node struct {
	logger    *log.Logger
	caddr     caddr.CAddr
	Discovery discovery.DiscoveryProvider

	// Request timeout is the maximum time for requests to last.
	RequestTimeout time.Duration

	stopChan    chan struct{}
	inShutdown  atomicvar.AtomicBool
	stopped     atomicvar.AtomicBool
	dialOptions []grpc.DialOption
}

// NodeOption is an option that can be provided to New to customize
// internal values.
type NodeOption func(*Node) error

// WithDialOptions is an option to provide dialoptions that will be used
// when making GRPC dial requests to peers.
func WithDialOptions(v ...grpc.DialOption) NodeOption {
	return func(n *Node) error {
		n.dialOptions = v
		return nil
	}
}

// WithRequestTimeout is an option to define requests' timeout duration.
func WithRequestTimeout(v time.Duration) NodeOption {
	return func(n *Node) error {
		n.RequestTimeout = v
		return nil
	}
}

func New(logger *log.Logger, addr string, opts ...NodeOption) (*Node, error) {
	caddr, err := caddr.FromString(addr)
	if err != nil {
		return nil, err
	}

	n := Node{
		logger: logger,
		caddr:  caddr,
	}

	for _, opt := range opts {
		if err := opt(&n); err != nil {
			return nil, errors.Wrapf(err, "executing option %T", opt)
		}
	}

	return &n, nil
}

// Invoke implements the MessagingServiceServer interface.
func (n *Node) Invoke(ctx context.Context, msg *messaging.Message) (*messaging.Response, error) {
	return nil, fmt.Errorf("not implemented")
}
