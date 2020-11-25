package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type TelnetClientObj struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	conn       net.Conn
	closedFlag bool
}

var ErrConnectionClosed = errors.New("connection closed")

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientObj{address, timeout, in, out, nil, false}
}

func (tco *TelnetClientObj) Connect() error {
	conn, err := net.DialTimeout("tcp", tco.address, tco.timeout)
	tco.conn = conn
	if err != nil {
		err = fmt.Errorf("error while connecting: %w", err)
	}
	return err
}

func (tco *TelnetClientObj) Send() error {
	return Transfer(tco.in, tco.conn, tco.closedFlag)
}

func (tco *TelnetClientObj) Receive() error {
	return Transfer(tco.conn, tco.out, tco.closedFlag)
}

func (tco *TelnetClientObj) Close() error {
	tco.closedFlag = true
	return tco.conn.Close()
}

func Transfer(r io.Reader, w io.Writer, closedFlag bool) error {
	if closedFlag {
		return ErrConnectionClosed
	}
	_, err := io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("error while transferring: %w", err)
	}
	return nil
}

// Place your code here
// P.S. Author's solution takes no more than 50 lines
