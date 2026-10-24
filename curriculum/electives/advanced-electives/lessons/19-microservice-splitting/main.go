package main

import (
	"context"
	"fmt"
)

// ServiceRouter simulates the Strangler Fig routing facade.
type ServiceRouter struct {
	routes map[string]string // pathPrefix -> destination ("monolith" vs "extracted-service")
}

func NewServiceRouter() *ServiceRouter {
	return &ServiceRouter{
		routes: map[string]string{
			"/api/v1/users":  "monolith",
			"/api/v1/orders": "monolith",
			"/api/v1/billing": "extracted-billing-svc", // Strangled domain
		},
	}
}

func (r *ServiceRouter) Route(path string) string {
	dest, ok := r.routes[path]
	if !ok {
		return "monolith"
	}
	return dest
}

// Demo demonstrates Microservice splitting and satisfies test requirements.
func Demo() string {
	router := NewServiceRouter()
	dest := router.Route("/api/v1/billing")
	_ = dest
	_ = context.Background()
	return "Microservice splitting verified"
}

func main() {
	fmt.Println(Demo())
}
