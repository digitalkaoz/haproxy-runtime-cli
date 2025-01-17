package socket

import (
	"bufio"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	mu sync.Mutex // Declare a mutex
)

func Listen(file string) net.Conn {
	c, err := net.Dial("unix", file)
	if err != nil {
		panic(err)
	}

	return c
}

func ExecCmd[T any](conn func() net.Conn, command string, cb func(*string) T) tea.Cmd {
	return func() tea.Msg {
		res, err := writeToSocket(conn, command)
		if err != nil {
			return err
		}
		return cb(res)
	}
}

func Exec(conn func() net.Conn, command string) (*string, error) {
	response, err := writeToSocket(conn, command)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("no response")
	}

	return response, nil
}

func writeToSocket(c func() net.Conn, command string) (*string, error) {
	mu.Lock()
	sock := c()
	_, err := sock.Write([]byte(command + "\n"))
	if err != nil {
		mu.Unlock()
		return nil, fmt.Errorf("failed to write to socket: %w", err)
	}

	//time.Sleep(time.Second)
	res, err := readFromSocket(sock)
	sock.Close()
	mu.Unlock()
	return res, err
}

func readFromSocket(r io.Reader) (*string, error) {
	reader := bufio.NewReader(r)
	response := ""
	for {
		r, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read from socket: %w", err)
		}
		response += r
	}
	trimmedResponse := strings.TrimSpace(response)

	return &trimmedResponse, nil
}
