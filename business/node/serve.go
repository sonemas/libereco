package node

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/protobuf/networking"
	"google.golang.org/grpc"
)

func (n *Node) Serve() error {
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

	for {
		select {
		case <-n.stopChan:
			s.Stop()
			break
		case err := <-errChan:
			return err
		}
	}

	return nil
}

func (n *Node) Shutdown(ctx context.Context) error {
	// TODO: Implement context
	n.stopChan <- struct{}{}

	return nil
}
