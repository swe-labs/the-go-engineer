package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestServer(t *testing.T) *WikiServer {
	t.Helper()
	server, err := NewWikiServer(defaultTemplatesFS, defaultStaticFS, defaultContentFS)
	if err != nil {
		t.Fatalf("failed to create test wiki server: %v", err)
	}
	return server
}

func TestWikiServer_RootRedirect(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected status 302 Found, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if loc != "/wiki/welcome" {
		t.Errorf("expected redirect to '/wiki/welcome', got %q", loc)
	}
}

func TestWikiServer_RenderPage(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/wiki/welcome", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type text/html, got %q", contentType)
	}

	body := rr.Body.String()

	// Verify semantic tags
	expectedTags := []string{
		"<header class=\"wiki-header\">",
		"<nav class=\"sidebar-nav\"",
		"<main class=\"wiki-main\"",
		"<article class=\"wiki-article\">",
		"<aside class=\"wiki-rail\"",
		"<footer class=\"wiki-footer\">",
	}

	for _, tag := range expectedTags {
		if !strings.Contains(body, tag) {
			t.Errorf("expected rendered HTML to contain %q", tag)
		}
	}

	// Verify pure CSS theme toggle checkbox exists
	if !strings.Contains(body, "id=\"theme-toggle\"") {
		t.Errorf("expected rendered HTML to contain theme toggle checkbox element")
	}

	// Verify active navigation item highlighting
	if !strings.Contains(body, "active") {
		t.Errorf("expected active class on current navigation item")
	}
}

func TestWikiServer_DocPages(t *testing.T) {
	server := setupTestServer(t)

	slugs := []string{"goroutines", "memory-model", "tcp-handshake"}
	for _, slug := range slugs {
		req := httptest.NewRequest(http.MethodGet, "/wiki/"+slug, nil)
		rr := httptest.NewRecorder()

		server.Handler().ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200 for slug %q, got %d", slug, rr.Code)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "Document Specs") {
			t.Errorf("expected metadata rail on slug %q", slug)
		}
	}
}

func TestWikiServer_NotFound(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/wiki/unknown-nonexistent-document", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "404") || !strings.Contains(body, "Document Not Found") {
		t.Errorf("expected 404 message in body, got: %s", body)
	}
}

func TestWikiServer_StaticCSS(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/static/style.css", nil)
	rr := httptest.NewRecorder()

	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for static CSS, got %d", rr.Code)
	}

	css := rr.Body.String()
	if !strings.Contains(css, "--bg-canvas") || !strings.Contains(css, "grid-template-columns") {
		t.Errorf("expected CSS variables and Grid rules in style.css")
	}
}

func TestWikiServer_SafeEscaping(t *testing.T) {
	server := setupTestServer(t)

	raw := "## Injected <script>alert('xss')</script>\nMalicious [link](javascript:void(0))"
	rendered := server.renderMarkdownToHTML(raw)

	if strings.Contains(rendered, "<script>") {
		t.Errorf("expected raw script tags to be escaped, got: %s", rendered)
	}
	if !strings.Contains(rendered, "&lt;script&gt;") {
		t.Errorf("expected script tag to be escaped to &lt;script&gt;, got: %s", rendered)
	}
}
