package node

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/messaging"
	"github.com/sonemas/libereco/foundation/discovery"
	"google.golang.org/grpc"
)

// Serve starts the server and blocks on server requests.
// It will also initialize the periodic ping requests to check whether
// peers are still alive.
func (n *Node) Serve() error {
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
	if err := n.Discovery.Init(s); err != nil {
		return errors.Wrap(err, "initializing discovery service")
	}

	messaging.RegisterMessagingServiceServer(s, n)

	var errChan chan error
	go func() {
		if err := s.Serve(l); err != nil {
			errChan <- errors.Wrap(err, "starting server")
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
func (n *Node) Shutdown(ctx context.Context) error {
	n.inShutdown.SetTrue()
	if err := n.Discovery.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutting down discovery service")
	}
	n.stopChan <- struct{}{}

	for !n.stopped.IsSet() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
