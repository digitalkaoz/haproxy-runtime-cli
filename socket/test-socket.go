package socket

import (
	"io"
	"net"
	"time"
)

type DummySocket struct {
	Output []byte
}

func (m *DummySocket) Read(b []byte) (int, error) {
	res := append(m.Output[:], '\n')
	n := copy(b, res)

	return n, io.EOF
}

func (m *DummySocket) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *DummySocket) Close() error {
	return nil
}

func (m *DummySocket) LocalAddr() net.Addr {
	return nil
}

func (m *DummySocket) RemoteAddr() net.Addr {
	return nil
}

func (m *DummySocket) SetDeadline(t time.Time) error {
	return nil
}

func (m *DummySocket) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *DummySocket) SetWriteDeadline(t time.Time) error {
	return nil
}
