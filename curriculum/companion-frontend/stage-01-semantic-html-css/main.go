package main

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Embed static assets, templates, and initial markdown documentation
//
//go:embed templates/*.html
var defaultTemplatesFS embed.FS

//go:embed static/*
var defaultStaticFS embed.FS

//go:embed content/*.md
var defaultContentFS embed.FS

// NavItem represents a document in the wiki sidebar navigation.
type NavItem struct {
	Slug   string
	Title  string
	Active bool
}

// WikiPageData is the view model injected into layout.html.
type WikiPageData struct {
	Title           string
	Slug            string
	Body            template.HTML
	WordCount       int
	ReadTimeMinutes int
	NavItems        []NavItem
}

// WikiServer encapsulates the HTTP handlers and embedded assets.
type WikiServer struct {
	templates *template.Template
	contentFS fs.FS
	staticFS  fs.FS
	navItems  []NavItem
}

// NewWikiServer constructs a configured WikiServer.
func NewWikiServer(templatesFS, staticFS, contentFS fs.FS) (*WikiServer, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	ws := &WikiServer{
		templates: tmpl,
		contentFS: contentFS,
		staticFS:  staticFS,
	}

	if err := ws.loadNavItems(); err != nil {
		return nil, fmt.Errorf("failed to index navigation items: %w", err)
	}

	return ws, nil
}

// loadNavItems discovers all markdown files in contentFS to build the sidebar navigation.
func (ws *WikiServer) loadNavItems() error {
	entries, err := fs.ReadDir(ws.contentFS, "content")
	if err != nil {
		return err
	}

	var items []NavItem
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		slug := strings.TrimSuffix(entry.Name(), ".md")
		title := ws.extractTitle(slug)
		items = append(items, NavItem{
			Slug:  slug,
			Title: title,
		})
	}
	ws.navItems = items
	return nil
}

// extractTitle extracts the first H1 header from the markdown file or falls back to title-cased slug.
func (ws *WikiServer) extractTitle(slug string) string {
	path := fmt.Sprintf("content/%s.md", slug)
	data, err := fs.ReadFile(ws.contentFS, path)
	if err != nil {
		return strings.Title(strings.ReplaceAll(slug, "-", " "))
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return strings.Title(strings.ReplaceAll(slug, "-", " "))
}

// Handler returns the root http.Handler with all registered routes.
func (ws *WikiServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// Static assets handler
	staticSub, err := fs.Sub(ws.staticFS, "static")
	if err == nil {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))
	}

	// Wiki routing
	mux.HandleFunc("/wiki/", ws.handleWikiPage)
	mux.HandleFunc("/wiki", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/wiki/welcome", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/wiki/welcome", http.StatusFound)
			return
		}
		ws.handleNotFound(w, r.URL.Path)
	})

	return mux
}

// handleWikiPage renders a markdown document inside the semantic layout.
func (ws *WikiServer) handleWikiPage(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/wiki/")
	slug = strings.Trim(slug, "/")

	if slug == "" {
		http.Redirect(w, r, "/wiki/welcome", http.StatusFound)
		return
	}

	docPath := fmt.Sprintf("content/%s.md", slug)
	rawBytes, err := fs.ReadFile(ws.contentFS, docPath)
	if err != nil {
		ws.handleNotFound(w, slug)
		return
	}

	rawText := string(rawBytes)
	title := ws.extractTitle(slug)
	renderedHTML := ws.renderMarkdownToHTML(rawText)

	words := len(strings.Fields(rawText))
	readTime := words / 180
	if readTime < 1 {
		readTime = 1
	}

	// Build active navigation item slice
	navItems := make([]NavItem, len(ws.navItems))
	for i, item := range ws.navItems {
		navItems[i] = NavItem{
			Slug:   item.Slug,
			Title:  item.Title,
			Active: item.Slug == slug,
		}
	}

	data := WikiPageData{
		Title:           title,
		Slug:            slug,
		Body:            template.HTML(renderedHTML),
		WordCount:       words,
		ReadTimeMinutes: readTime,
		NavItems:        navItems,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ws.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// handleNotFound renders the 404 template.
func (ws *WikiServer) handleNotFound(w http.ResponseWriter, slug string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_ = ws.templates.ExecuteTemplate(w, "not_found.html", map[string]string{
		"Slug": slug,
	})
}

// renderMarkdownToHTML converts standard markdown constructs to semantic HTML without external dependencies.
func (ws *WikiServer) renderMarkdownToHTML(raw string) string {
	lines := strings.Split(raw, "\n")
	var buf bytes.Buffer
	inCodeBlock := false
	inList := false

	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	codeRegex := regexp.MustCompile("`([^`]+)`")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Code block toggle
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				buf.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				if inList {
					buf.WriteString("</ul>\n")
					inList = false
				}
				buf.WriteString("<pre><code>")
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			buf.WriteString(html.EscapeString(line) + "\n")
			continue
		}

		// Empty line
		if trimmed == "" {
			if inList {
				buf.WriteString("</ul>\n")
				inList = false
			}
			continue
		}

		// Headings
		if strings.HasPrefix(trimmed, "### ") {
			if inList {
				buf.WriteString("</ul>\n")
				inList = false
			}
			buf.WriteString(fmt.Sprintf("<h3>%s</h3>\n", html.EscapeString(strings.TrimPrefix(trimmed, "### "))))
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			if inList {
				buf.WriteString("</ul>\n")
				inList = false
			}
			buf.WriteString(fmt.Sprintf("<h2>%s</h2>\n", html.EscapeString(strings.TrimPrefix(trimmed, "## "))))
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			// Skip top-level # because it's rendered as article-title in header
			continue
		}

		// List items
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList {
				buf.WriteString("<ul>\n")
				inList = true
			}
			itemText := strings.TrimPrefix(trimmed, "- ")
			itemText = strings.TrimPrefix(itemText, "* ")
			formatted := formatInlineMarkdown(itemText, linkRegex, boldRegex, codeRegex)
			buf.WriteString(fmt.Sprintf("<li>%s</li>\n", formatted))
			continue
		}

		if inList {
			buf.WriteString("</ul>\n")
			inList = false
		}

		// Paragraph
		formatted := formatInlineMarkdown(trimmed, linkRegex, boldRegex, codeRegex)
		buf.WriteString(fmt.Sprintf("<p>%s</p>\n", formatted))
	}

	if inCodeBlock {
		buf.WriteString("</code></pre>\n")
	}
	if inList {
		buf.WriteString("</ul>\n")
	}

	return buf.String()
}

// formatInlineMarkdown formats bold, inline code, and links.
func formatInlineMarkdown(text string, linkRegex, boldRegex, codeRegex *regexp.Regexp) string {
	escaped := html.EscapeString(text)

	// Replace bold: **text** -> <strong>text</strong>
	escaped = boldRegex.ReplaceAllString(escaped, "<strong>$1</strong>")

	// Replace inline code: `code` -> <code>code</code>
	escaped = codeRegex.ReplaceAllString(escaped, "<code>$1</code>")

	// Replace links: [label](url) -> <a href="url">label</a>
	escaped = linkRegex.ReplaceAllString(escaped, `<a href="$2">$1</a>`)

	return escaped
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server, err := NewWikiServer(defaultTemplatesFS, defaultStaticFS, defaultContentFS)
	if err != nil {
		log.Fatalf("Fatal: failed to initialize WikiServer: %v", err)
	}

	addr := ":" + port
	log.Printf("[Stage 1 Starter] Engineering Wiki server listening at http://localhost%s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      server.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
