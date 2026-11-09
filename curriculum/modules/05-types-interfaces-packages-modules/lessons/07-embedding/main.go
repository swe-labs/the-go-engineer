package main

import (
	"bytes"
	"fmt"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

type Logger struct {
	prefix string
}

func (l Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type Server struct {
	Logger
	Host string
}

type Pinger interface {
	Ping() error
}

type HealthChecker interface {
	Pinger
	Status() map[string]string
}

type SimplePinger struct{}

func (s SimplePinger) Ping() error {
	return nil
}

type Monitor struct {
	SimplePinger
}

func (m Monitor) Status() map[string]string {
	return map[string]string{"status": "ok"}
}

func main() {
	s := Server{
		Logger: Logger{prefix: "api"},
		Host:   "localhost:8080",
	}
	s.Log("server started")

	var buf bytes.Buffer
	var rw ReadWriter = &buf
	rw.Write([]byte("hello"))
	var readBuf [16]byte
	n, _ := rw.Read(readBuf[:])
	fmt.Println("read:", string(readBuf[:n]))

	var hc HealthChecker = Monitor{}
	if err := hc.Ping(); err != nil {
		fmt.Println("ping failed:", err)
	}
	fmt.Println("status:", hc.Status())
}
