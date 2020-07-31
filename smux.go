package smux

import (
	"net"
	"sync"
)

// Listen announces on the local network address.
func Listen(network, address string) (*Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return &Listener{listener}, nil
}

// Dial connects to the address on the named network.
func Dial(network, address string) (*Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn:    conn,
		streams: sync.Map{},
		ch:      make(chan Stream, 1),
		counter: NewCounter(START_STREAM_ID_OF_CLIENT),
	}, nil
}
