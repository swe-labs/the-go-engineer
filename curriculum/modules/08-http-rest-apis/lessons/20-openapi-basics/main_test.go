package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildSpec(t *testing.T) {
	spec := buildSpec()

	if spec.OpenAPI != "3.0.3" {
		t.Errorf("OpenAPI version = %q, want %q", spec.OpenAPI, "3.0.3")
	}
	if spec.Info.Title != "Products API" {
		t.Errorf("Title = %q, want %q", spec.Info.Title, "Products API")
	}
	if spec.Info.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", spec.Info.Version, "1.0.0")
	}

	if _, ok := spec.Paths["/products"]; !ok {
		t.Fatal("missing /products path")
	}

	productList := spec.Paths["/products"].Get
	if productList == nil {
		t.Fatal("missing GET /products")
	}
	if productList.Summary != "List products" {
		t.Errorf("Summary = %q, want %q", productList.Summary, "List products")
	}

	if len(productList.Parameters) != 3 {
		t.Errorf("expected 3 parameters, got %d", len(productList.Parameters))
	}
}

func TestGetEndpoint(t *testing.T) {
	spec := buildSpec()

	op := getEndpoint(spec, "/products", "GET")
	if op == nil {
		t.Fatal("expected GET /products to exist")
	}
	if op.Summary != "List products" {
		t.Errorf("Summary = %q", op.Summary)
	}

	op2 := getEndpoint(spec, "/products", "POST")
	if op2 != nil {
		t.Error("expected POST /products to be nil")
	}

	op3 := getEndpoint(spec, "/nonexistent", "GET")
	if op3 != nil {
		t.Error("expected nonexistent path to return nil")
	}
}

func TestSpecHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	specHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", resp.Header.Get("Content-Type"))
	}

	var spec OpenAPISpec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		t.Fatalf("json decode: %v", err)
	}

	if spec.Info.Title != "Products API" {
		t.Errorf("Title = %q", spec.Info.Title)
	}
}

func TestSpecHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	specHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestSpecJSONRoundTrip(t *testing.T) {
	spec := buildSpec()
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded OpenAPISpec
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.OpenAPI != "3.0.3" {
		t.Errorf("roundtrip OpenAPI = %q", decoded.OpenAPI)
	}
}

func TestSpecPathItemDetail(t *testing.T) {
	spec := buildSpec()

	productByID := spec.Paths["/products/{id}"].Get
	if productByID == nil {
		t.Fatal("missing GET /products/{id}")
	}

	if len(productByID.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(productByID.Parameters))
	}

	param := productByID.Parameters[0]
	if param.Name != "id" {
		t.Errorf("param name = %q, want %q", param.Name, "id")
	}
	if param.In != "path" {
		t.Errorf("param in = %q, want %q", param.In, "path")
	}
	if !param.Required {
		t.Error("id param should be required")
	}

	if _, ok := productByID.Responses["404"]; !ok {
		t.Error("missing 404 response")
	}
}

func TestLessonCompiles(t *testing.T) {
}
