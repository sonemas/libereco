package zone

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/sonemas/libereco/business/protobuf/messaging"
	"github.com/sonemas/libereco/foundation/storage"
	"google.golang.org/grpc"
)

type ZoneProvider interface {
	Domain() string
	Invoke(string, map[string]interface{}) error
}

type Request struct {
	Type string
	Args map[string][]byte
}

func (r *Request) IsGetRequest() bool {
	return r.Type == "GET"
}

func (r *Request) IsSetRequest() bool {
	return r.Type == "SET"
}

type Response struct {
	State map[string][]byte
}

type Handler interface {
	Version() []string
	Handle(Request) (Response, error)
}

// Zone is an independent data zone on the network.
type Zone struct {
	logger         *log.Logger
	storage        storage.BasicStorer
	dialOptions    []grpc.DialOption
	RequestTimeout time.Duration
	Handlers       []Handler
	domain         string
}

type ZoneOption func(*Zone) error

// WithDialOptions is an option to provide dialoptions that will be used
// when making GRPC dial requests to peers.
func WithDialOptions(v ...grpc.DialOption) ZoneOption {
	return func(z *Zone) error {
		z.dialOptions = v
		return nil
	}
}

// WithRequestTimeout is an option to define requests' timeout duration.
func WithRequestTimeout(v time.Duration) ZoneOption {
	return func(z *Zone) error {
		z.RequestTimeout = v
		return nil
	}
}

// New is a factory function to initialize a new Zone.
func New(l *log.Logger, domain string, s storage.BasicStorer, opts ...ZoneOption) (*Zone, error) {
	z := Zone{
		logger:  l,
		storage: s,
		domain:  domain,
	}

	for _, opt := range opts {
		if err := opt(&z); err != nil {
			return nil, err
		}
	}

	return &z, nil
}

func (z *Zone) Domain() string {
	return z.domain
}

type InvocationError struct {
	err string
}

func (err *InvocationError) Error() string {
	return err.err
}

var (
	ErrInvalidPayloadType = errors.New("payload type should be GET or SET")
)

func (z *Zone) Invoke(ctx context.Context, req *messaging.Message) (*messaging.Response, error) {
	var handler Handler
	for _, h := range z.Handlers {
		ok := false
		for _, ver := range h.Version() {
			if strings.ToLower(ver) == strings.ToLower(req.Version) {
				ok = true
			}
		}

		if ok {
			handler = h
			break
		}
	}

	r := Request{
		Args: req.Payload,
	}

	switch req.PayloadType {
	case messaging.Message_PAYLOAD_TYPE_GET:
		r.Type = "GET"
	case messaging.Message_PAYLOAD_TYPE_SET:
		r.Type = "SET"
	default:
		return nil, ErrInvalidPayloadType
	}

	res, err := handler.Handle(r)
	if err != nil {
		return nil, &InvocationError{err.Error()}
	}

	if r.IsSetRequest() {
		changes, err := z.storage.Set(res.State, storage.Arguments{})
		if err != nil {
			return nil, err
		}
		return &messaging.Response{State: changes}, nil
	}

	return &messaging.Response{State: res.State}, nil
}
