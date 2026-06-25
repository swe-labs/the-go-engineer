package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenAPISpec struct {
	OpenAPI string `json:"openapi"`
	Info    Info   `json:"info"`
	Paths   Paths  `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Paths map[string]PathItem

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Type       string            `json:"type"`
	Items      *Schema           `json:"items,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

func buildSpec() OpenAPISpec {
	return OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       "Products API",
			Description: "A simple API for managing products",
			Version:     "1.0.0",
		},
		Paths: Paths{
			"/products": PathItem{
				Get: &Operation{
					Summary:     "List products",
					Description: "Returns a paginated list of products",
					Parameters: []Parameter{
						{Name: "page", In: "query", Required: false, Description: "Page number", Schema: Schema{Type: "integer"}},
						{Name: "limit", In: "query", Required: false, Description: "Items per page", Schema: Schema{Type: "integer"}},
						{Name: "category", In: "query", Required: false, Description: "Filter by category", Schema: Schema{Type: "string"}},
					},
					Responses: map[string]Response{
						"200": {
							Description: "A paginated list of products",
							Content: map[string]MediaType{
								"application/json": {
									Schema: Schema{
										Type: "array",
										Items: &Schema{
											Type: "object",
											Properties: map[string]Schema{
												"id":       {Type: "integer"},
												"name":     {Type: "string"},
												"category": {Type: "string"},
												"price":    {Type: "number"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/products/{id}": PathItem{
				Get: &Operation{
					Summary:     "Get a product by ID",
					Description: "Returns a single product",
					Parameters: []Parameter{
						{Name: "id", In: "path", Required: true, Description: "Product ID", Schema: Schema{Type: "integer"}},
					},
					Responses: map[string]Response{
						"200": {Description: "The requested product"},
						"404": {Description: "Product not found"},
					},
				},
			},
		},
	}
}

func specHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	spec := buildSpec()
	json.NewEncoder(w).Encode(spec)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/openapi.json", specHandler)

	log.Println("OpenAPI spec served at /openapi.json on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func getEndpoint(spec OpenAPISpec, path, method string) *Operation {
	pathItem, ok := spec.Paths[path]
	if !ok {
		return nil
	}
	switch method {
	case "GET":
		return pathItem.Get
	case "POST":
		return pathItem.Post
	case "PUT":
		return pathItem.Put
	case "DELETE":
		return pathItem.Delete
	}
	return nil
}

var _ = fmt.Sprintf
