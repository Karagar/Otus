package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeoutStr := flag.String("timeout", "10s", "Timeout for connection to server. Default 10s")
	flag.Parse()
	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		println(fmt.Errorf("error while parse timeout duration: %w", err))
		return
	}
	if len(flag.Args()) != 2 {
		println("Not enought arguments")
		return
	}
	targetAdress := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	errorChan := make(chan error)
	signlsChan := make(chan os.Signal, 1)
	signal.Notify(signlsChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	client := NewTelnetClient(targetAdress, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
		close(signlsChan)
		close(errorChan)
	}()

	go func() {
		errorChan <- client.Send()
	}()
	go func() {
		errorChan <- client.Receive()
	}()

	select {
	case <-signlsChan:
		return
	case err = <-errorChan:
		if err != nil {
			log.Println(err)
			fmt.Fprintf(os.Stderr, "Connection closed\n")
			return
		}
	}
}
