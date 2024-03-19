package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Telnet struct {
	host    string
	port    string
	timeout time.Duration
	conn    net.Conn
}

func Start() error {
	telnet := &Telnet{}
	telnet.timeout = *flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=<timeout>] <host> <port>")
		os.Exit(1)
	}

	telnet.host = flag.Arg(0)
	telnet.port = flag.Arg(1)

	fmt.Println(telnet)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	dialer := net.Dialer{Timeout: 10 * time.Second}

	// Подключение к серверу
	var err error
	telnet.conn, err = dialer.Dial("tcp", net.JoinHostPort(telnet.host, telnet.port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer telnet.conn.Close()

	errChan := make(chan error)
	go telnet.SocketWriter(errChan)
	go telnet.SocketReader(errChan)

	select {
	case <-sigChan:
		return nil
	case err = <-errChan:
		return err
	}

}

func (tel *Telnet) SocketWriter(errChan chan error) {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		buff, err := inputReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err.Error())
			errChan <- err
			return
		}
		if _, err := tel.conn.Write(buff); err != nil {
			errChan <- err
			return
		}
	}

}

func (tel *Telnet) SocketReader(errChan chan error) {
	for {
		serverReader := bufio.NewReader(tel.conn)
		for {
			buff, err := serverReader.ReadBytes('\n')
			if err != nil {
				errChan <- err
				return
			}
			fmt.Println(string(buff))
		}
	}

}

func main() {
	err := Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}
