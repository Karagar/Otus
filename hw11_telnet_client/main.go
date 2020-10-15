package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	timeoutStr := flag.String("timeout", "10s", "Timeout for connection to server. Default 10s")
	flag.Parse()
	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		println(fmt.Errorf("Error while parse timeout duration: %v", err))
		return
	}
	if len(flag.Args()) != 2 {
		println("Not enought arguments")
		return
	}
	targetAdress := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	go func() {
		handleServer(targetAdress)
	}()
	time.Sleep(2 * time.Second)
	client := NewTelnetClient(targetAdress, timeout, os.Stdin, os.Stdout)
	client.Connect()
	go func() {
		for {
			client.Send()
			if err != nil {
				break
			}
		}
		return
	}()
	go func() {
		for {
			client.Receive()
			if err != nil {
				break
			}
		}
		return
	}()
	time.Sleep(300 * time.Second)
	client.Close()
}

func handleServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	defer l.Close()
	log.Println("LISTENED")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Cannot accept: %v", err)
		}
		log.Println("ACCEPTED")
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(fmt.Sprintf("Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr())))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("RECEIVED: %s", text)
		if text == "quit" || text == "exit" {
			break
		}

		conn.Write([]byte(fmt.Sprintf("I have received '%s'\n", text)))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
	}

	log.Printf("Closing connection with %s", conn.RemoteAddr())

}
