package main

import (
	"bufio"
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

var ConnectionClosedErr = errors.New("Connection closed")

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientObj{address, timeout, in, out, nil, false}
}

func (tco *TelnetClientObj) Connect() error {
	conn, err := net.DialTimeout("tcp", tco.address, tco.timeout)
	tco.conn = conn
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
		return ConnectionClosedErr
	}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return scanner.Err()
	}
	_, err := w.Write([]byte(fmt.Sprintf("%s\n", scanner.Text())))
	return err
}

// Place your code here
// P.S. Author's solution takes no more than 50 lines
