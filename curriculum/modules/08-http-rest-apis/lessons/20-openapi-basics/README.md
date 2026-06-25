# OpenAPI basics

## Learning objective

Write an OpenAPI 3.0 specification for a REST API, define paths, methods, parameters, and response schemas, serve the spec from a Go HTTP handler, and understand code generation from the spec.

## Why this matters

An OpenAPI spec is the single source of truth for your API. It documents every endpoint, request parameter, and response shape in a machine-readable format. Tools consume the spec to generate client SDKs, server stubs, interactive docs (Swagger UI), and tests. Without a spec, your API is undocumented; clients reverse-engineer endpoints from code or trial-and-error.

## Mental model

Think of an OpenAPI spec as a blueprint for your API:

- The **info** block is the title page: what is this API? Who made it? What version?
- The **paths** are the floor plan: every room (endpoint) with its doors (methods).
- Each **operation** is a room's layout: what goes in (parameters, request body), what comes out (responses), and what errors are possible.
- The **schemas** are the furniture catalog: reusable type definitions.

Just as a building blueprint lets electricians, plumbers, and inspectors all work from the same plan, an OpenAPI spec lets frontend, backend, QA, and devops teams coordinate.

## Core idea

OpenAPI 3.0 spec structure:

```yaml
openapi: 3.0.3
info:
  title: My API
  version: 1.0.0
paths:
  /products:
    get:
      summary: List products
      parameters:
        - name: page
          in: query
          schema:
            type: integer
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Product'
components:
  schemas:
    Product:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
```

Key elements:

- **Paths**: URLs of your API endpoints.
- **Operations**: HTTP method + path combination. Each has `summary`, `parameters`, `requestBody`, `responses`.
- **Parameters**: Can be `path` (in URL), `query` (after `?`), `header`, or `cookie`.
- **Schemas**: JSON Schema definitions for request/response bodies.
- **Components**: Reusable schemas, parameters, responses, etc.

## Under the hood

OpenAPI 3.0 uses JSON Schema (draft-07) for type definitions. The spec can be written in YAML or JSON. YAML is more human-readable; JSON is what most tools produce.

In Go, you can:
1. Write the spec as a static file (`openapi.json` or `openapi.yaml`) and serve it.
2. Build the spec programmatically using structs and `json.Marshal`.
3. Generate the spec from Go code using annotations (e.g., `swaggo/swag`).

The example below shows approach #2 -- building the spec as Go structs and serving it via `/openapi.json`. This guarantees the spec is always in sync with the code (since the code IS the spec).

## How Go uses it

Several Go tools work with OpenAPI:

- **ogen**: Code generation from OpenAPI spec to Go server and client.
- **oapi-codegen**: Generate Go server stubs, client, and types from OpenAPI spec.
- **swaggo/swag**: Annotate Go code with comments, auto-generate OpenAPI spec.
- **kin-openapi**: Parse and validate OpenAPI specs in Go.

The standard approach is: write your spec (YAML/JSON), generate Go code with `oapi-codegen`, implement the generated interfaces. This ensures the implementation never drifts from the spec.

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type OpenAPISpec struct {
	OpenAPI string            `json:"openapi"`
	Info    Info              `json:"info"`
	Paths   map[string]PathItem `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type PathItem struct {
	Get  *Operation `json:"get,omitempty"`
	Post *Operation `json:"post,omitempty"`
}

type Operation struct {
	Summary     string                `json:"summary"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	Responses   map[string]Response   `json:"responses"`
}

type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   Schema `json:"schema"`
}

type Schema struct {
	Type  string  `json:"type"`
	Items *Schema `json:"items,omitempty"`
}

type Response struct {
	Description string `json:"description"`
}

func buildSpec() OpenAPISpec {
	return OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:   "Products API",
			Version: "1.0.0",
		},
		Paths: map[string]PathItem{
			"/products": {
				Get: &Operation{
					Summary: "List products",
					Responses: map[string]Response{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}
}

func specHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildSpec())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/openapi.json", specHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

Client requests `GET /openapi.json`:

1. `specHandler` is called. It builds a new `OpenAPISpec` struct via `buildSpec()`.
2. `json.NewEncoder` marshals the struct to JSON: `{"openapi":"3.0.3","info":{"title":"Products API","version":"1.0.0"},"paths":{"/products":{"get":{"summary":"List products","responses":{"200":{"description":"OK"}}}}}}`.
3. Response is sent with `Content-Type: application/json`.
4. A tool like Swagger UI or `oapi-codegen` reads this spec to generate docs or client code.

For a more complete setup, add `components/schemas` to define reusable types:

```go
spec.Components = Components{
    Schemas: map[string]Schema{
        "Product": {
            Type: "object",
            Properties: map[string]Schema{
                "id":    {Type: "integer"},
                "name":  {Type: "string"},
                "price": {Type: "number"},
            },
        },
    },
}
```

## Common mistakes

- Writing the spec by hand and letting it drift from the implementation. Either generate code from spec or generate spec from code.
- Omitting error responses. Every endpoint should document its error status codes (4xx, 5xx).
- Using inconsistent naming: `camelCase` vs `snake_case`. Pick one and use it everywhere.
- Not using `$ref` for reusable schemas. Duplicating schema definitions makes the spec hard to maintain.
- Forgetting to include the spec in the deployment. Serve it at a well-known path (`/openapi.json`).
- Over-specifying: documenting internal implementation details that clients don't need.

## Debugging walkthrough

Swagger UI shows an error: "Failed to load API definition."

**Symptom**: The spec at `/openapi.json` is invalid JSON or invalid OpenAPI.

**Investigation**: Fetch the spec: `curl http://localhost:8080/openapi.json`. Check for JSON syntax errors. Validate against the OpenAPI 3.0 schema using a validator like `swagger-cli validate`.

**Root cause**: Missing required fields like `info.title` or `openapi` version. Or an invalid response code key (must be a string like `"200"`, not integer `200`).

**Fix**: Ensure all required fields are populated. In the Go struct, use `json:"openapi"` with a default value. Validate the spec with `kin-openapi` in a test:

```go
func TestSpecValid(t *testing.T) {
    spec := buildSpec()
    data, _ := json.Marshal(spec)
    _, err := openapi3.Load(data)
    if err != nil {
        t.Fatal(err)
    }
}
```

## Production notes

- Serve the spec at `/openapi.json` so tools can discover it automatically.
- Version the spec URL with your API version: `/v1/openapi.json`.
- Consider generating the spec from code using `swaggo/swag` annotations to prevent drift.
- Pin to a specific OpenAPI version (3.0.3 is stable). Avoid mixing 2.0 and 3.0 features.
- Use `oapi-codegen` to generate server stubs. Implement the generated interfaces for a spec-first workflow.

## Performance implications

- Serving the spec is a single JSON marshal per request -- negligible cost (microseconds).
- The spec file is typically 1-50 KB. Bandwidth is not a concern.
- Code generation from spec happens at build time, not runtime. No performance impact.
- The main cost is maintenance: keeping the spec in sync with implementation.

## Practice task

1. Write an OpenAPI 3.0 spec (as Go structs) for a `User` API with:
   - `GET /users` with `page`, `limit`, `role` query parameters.
   - `POST /users` with a JSON body containing `name`, `email`, `role`.
   - `GET /users/{id}` returning a single user.
2. Serve the spec at `/openapi.json`.
3. Write tests for: spec handler returns 200 with correct Content-Type, spec contains all expected paths and methods, parameters have correct types and locations.
4. Add a `components/schemas/User` definition and reference it with `$ref`.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/20-openapi-basics -v
```

Tests verify the spec structure, handler response, round-trip JSON marshalling, parameter details, and method validation.

## Review questions

1. What are the top-level sections of an OpenAPI 3.0 specification?
2. What is the difference between a `path` parameter and a `query` parameter in OpenAPI?
3. How does `$ref` help maintain a clean OpenAPI spec? When would you use `components/schemas`?
4. What tools can generate Go code from an OpenAPI spec? What are the benefits?
5. Why is it important to serve the OpenAPI spec from the running API server?

## NEXT UP

HTTP Client -- using `http.Client`, making requests with `http.NewRequest`, configuring timeouts and transport, and managing connection pooling.
