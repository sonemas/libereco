package caddr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidAddr is an error indicating an address is not in the correct format.
var ErrInvalidAddr = errors.New("address should be in the format host:port/id[/protocol]")

// CAddr is a combined address.
type CAddr struct {
	ID    string
	Proto string
	Host  string
	Port  uint
}

// String implements the stringer interface.
func (a CAddr) String() string {
	return fmt.Sprintf("%s:%d/%s/%s", a.Host, a.Port, a.ID, a.Proto)
}

// Addr returns the address string of the CAddr.
func (a CAddr) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// FromString splits the provided address into the host/port, id and optional protocol.
// The format of addresses is: host:post/id[/protocpl]. The default protocol is TCP.
func FromString(addr string) (CAddr, error) {
	p := strings.Split(addr, "/")
	if len(p) < 2 {
		return CAddr{}, ErrInvalidAddr
	}

	h := strings.Split(p[0], ":")
	if len(h) != 2 {
		return CAddr{}, ErrInvalidAddr
	}

	port, err := strconv.Atoi(h[1])
	if err != nil {
		return CAddr{}, ErrInvalidAddr
	}

	pr := "tcp"
	if len(p) == 3 {
		pr = p[2]
	}

	return CAddr{
		Host:  h[0],
		Port:  uint(port),
		ID:    p[1],
		Proto: pr,
	}, nil
}
