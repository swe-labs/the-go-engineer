package main

import (
	"testing"
)

func TestEchoServer(t *testing.T) {
	listener, err := startEchoServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("startEchoServer failed: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "basic echo", input: "hello", want: "hello", wantErr: false},
		{name: "numbers", input: "12345", want: "12345", wantErr: false},
		{name: "with spaces", input: "hello world", want: "hello world", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sendMessage(addr, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendMessage() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sendMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
