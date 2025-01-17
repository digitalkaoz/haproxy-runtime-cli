package socket

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestReadFromSocket(t *testing.T) {
	conn := &DummySocket{
		Output: []byte("Hello, this is a response from the socket."),
	}

	data, err := readFromSocket(conn)

	assert.Nil(t, err)
	assert.Equal(t, "Hello, this is a response from the socket.", *data)
}

func TestWriteToSocket(t *testing.T) {
	conn := &DummySocket{
		Output: []byte("Hello, this is a response from the socket."),
	}

	data, err := writeToSocket(func() net.Conn {
		return conn
	}, "help")

	assert.Nil(t, err)
	assert.Equal(t, "Hello, this is a response from the socket.", *data)
}
