package discovery

import (
	"context"
	"errors"

	"google.golang.org/grpc"
)

// Common errors that can be returned by discovery providers.
var (
	// ErrPeerNotFound is an error indicating that a peer can't be found.
	ErrPeerNotFound = errors.New("peer not found")

	// ErrPeerExists is an error indicating that a peer already exists in
	// a node's finger table.
	ErrPeerExists = errors.New("peer already exists")

	// ErrBootstrapingFailed is an error indicating that bootstrapping failed.
	ErrBootstrapingFailed = errors.New("bootstrap process failed")

	// ErrServiceIsStopped is an error indicating that the node has been stopped.
	ErrServiceIsStopped = errors.New("service is stopped")

	// ErrServiceInShutdown is an error indicating that the node is shutting down.
	ErrServiceInShutdown = errors.New("service is shutting down")
)

// DiscoveryProvider is an interface for discovery providers.
type DiscoveryProvider interface {

	// Init initializes the discovery service and registers it
	// with the provided GRPC server.
	Init(*grpc.Server) error

	// Serve launches the provider's listener.
	Serve() error

	// Shutdown performs a graceful shutdown of the provider.
	Shutdown(context.Context) error
}
