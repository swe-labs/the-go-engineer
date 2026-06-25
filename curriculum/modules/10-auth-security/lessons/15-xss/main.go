package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var safeTpl = template.Must(template.New("safe").Parse(`<html><body><h1>Hello, {{.}}!</h1></body></html>`))

func renderSafe(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	safeTpl.Execute(w, name)
}

func renderUnsafe(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><body><h1>Hello, %s!</h1></body></html>", name)
}

type SafeRenderer struct{}

func (s *SafeRenderer) Render(name string) string {
	var b strings.Builder
	err := safeTpl.Execute(&b, name)
	if err != nil {
		return fmt.Sprintf("render error: %v", err)
	}
	return b.String()
}

type UnsafeRenderer struct{}

func (u *UnsafeRenderer) Render(name string) string {
	return fmt.Sprintf("<html><body><h1>Hello, %s!</h1></body></html>", name)
}

func detectXSS(input string) bool {
	dangerous := []string{"<script", "onerror=", "onload=", "javascript:", "onclick="}
	lower := strings.ToLower(input)
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Unsafe render (string interpolation) ===")
	unsafe := &UnsafeRenderer{}
	fmt.Println(unsafe.Render("<script>alert('XSS')</script>"))
	fmt.Println(unsafe.Render("<img src=x onerror=alert(1)>"))

	fmt.Println("\n=== Safe render (html/template) ===")
	safe := &SafeRenderer{}
	fmt.Println(safe.Render("<script>alert('XSS')</script>"))
	fmt.Println(safe.Render("<img src=x onerror=alert(1)>"))

	fmt.Println("\n=== XSS detection ===")
	tests := []string{
		"<script>alert(1)</script>",
		"<img src=x onerror=alert(1)>",
		"Hello, World!",
		"javascript:alert(1)",
	}
	for _, t := range tests {
		fmt.Printf("  %-45s dangerous=%v\n", t, detectXSS(t))
	}

	fmt.Println("\nNote: html/template auto-escapes HTML special characters.")
	fmt.Println("CSP headers provide defense-in-depth against XSS.")
}
