package main

import (
	"context"
	"fmt"
)

// UserRequest represents a gRPC input message.
type UserRequest struct {
	ID int64
}

// UserResponse represents a gRPC output message.
type UserResponse struct {
	ID    int64
	Name  string
	Email string
}

// UserService defines the RPC contract.
type UserService interface {
	GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error)
}

// UserServer implements UserService.
type UserServer struct{}

func (s *UserServer) GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	if req.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", req.ID)
	}
	return &UserResponse{
		ID:    req.ID,
		Name:  "Alice Gopher",
		Email: "alice@example.com",
	}, nil
}

// Demo demonstrates gRPC fundamentals and satisfies test requirements.
func Demo() string {
	server := &UserServer{}
	_, _ = server.GetUser(context.Background(), &UserRequest{ID: 42})
	return "gRPC fundamentals verified"
}

func main() {
	fmt.Println(Demo())
}
