package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type UserV1 struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserV2 struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var usersV1 = []UserV1{
	{ID: 1, Name: "Alice"},
	{ID: 2, Name: "Bob"},
}

var usersV2 = []UserV2{
	{ID: 1, FirstName: "Alice", LastName: "Smith", Email: "alice@example.com"},
	{ID: 2, FirstName: "Bob", LastName: "Jones", Email: "bob@example.com"},
}

func usersV1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersV1)
}

func usersV2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersV2)
}

type VersionHandler struct {
	handlers map[string]http.HandlerFunc
}

func (vh *VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	var version string

	if strings.Contains(accept, "application/vnd.api.v2") {
		version = "v2"
	} else if strings.Contains(accept, "application/vnd.api.v1") {
		version = "v1"
	} else {
		version = "v1"
	}

	if h, ok := vh.handlers[version]; ok {
		h(w, r)
		return
	}

	http.Error(w, `{"error":"unsupported version"}`, http.StatusNotImplemented)
}

func versionMiddleware(version string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", version)
		handler(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/users", versionMiddleware("v1", usersV1Handler))
	mux.Handle("/api/v2/users", versionMiddleware("v2", usersV2Handler))

	vh := &VersionHandler{
		handlers: map[string]http.HandlerFunc{
			"v1": usersV1Handler,
			"v2": usersV2Handler,
		},
	}
	mux.Handle("/api/users", vh)

	log.Println("Versioned API on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func extractVersion(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, p := range parts {
		if strings.HasPrefix(p, "v") && len(p) > 1 {
			return p
		}
	}
	return "v1"
}

func convertV1ToV2(u UserV1) UserV2 {
	parts := strings.SplitN(u.Name, " ", 2)
	firstName := parts[0]
	lastName := ""
	if len(parts) > 1 {
		lastName = parts[1]
	}
	return UserV2{
		ID:        u.ID,
		FirstName: firstName,
		LastName:  lastName,
	}
}

var _ = fmt.Sprintf
