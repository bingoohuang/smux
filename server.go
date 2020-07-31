package smux

import (
	"bufio"
	"bytes"
	"io"
)

type Server struct {
	Network string
	Address string
	Handler Handler
}

func (s Server) ListenAndServe() error {
	l, err := Listen(s.Network, s.Address)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		s.dealConn(conn)
	}
}

func (s Server) dealConn(conn Conn) {
	go conn.Listen()
	go s.dealStream(conn)
}

func (s Server) dealStream(conn Conn) {
	defer conn.Close()

	for {
		stream, err := conn.Accept()
		if err != nil {
			break
		}

		go s.readWrite(stream)
	}
}

func (s Server) readWrite(stream Stream) {
	go stream.Poll()

	var buf bytes.Buffer
	out := make([]byte, 512)
	for {
		n, err := stream.Read(out)
		if err == io.EOF {
			break
		}

		buf.Write(out[:n])
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	s.Handler.Serve(w, bytes.NewReader(buf.Bytes()))
	_ = w.Flush()
	_, _ = stream.WriteOnce(b.Bytes())
}

type Handler interface {
	Serve(io.Writer, io.Reader)
}

type HandlerFunc func(io.Writer, io.Reader)

func (f HandlerFunc) Serve(w io.Writer, r io.Reader) {
	f(w, r)
}
