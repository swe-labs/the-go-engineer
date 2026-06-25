package main

import (
	"bufio"
	"fmt"
	"net"
)

func startEchoServer(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				msg, err := bufio.NewReader(c).ReadString('\n')
				if err != nil {
					return
				}
				c.Write([]byte(msg))
			}(conn)
		}
	}()
	return listener, nil
}

func sendMessage(addr, msg string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", msg)
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return reply[:len(reply)-1], nil
}

func main() {
	listener, err := startEchoServer("127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	fmt.Printf("Echo server on %s\n", addr)

	reply, _ := sendMessage(addr, "hello")
	fmt.Printf("Reply 1: %s\n", reply)

	reply, _ = sendMessage(addr, "world")
	fmt.Printf("Reply 2: %s\n", reply)
}
