package main

import (
	"fmt"
	"net/http"
	"strings"
)

func buildHTTPRequest(method, path string, headers map[string]string, body string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path))
	for k, v := range headers {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	if body != "" {
		b.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	}
	b.WriteString("\r\n")
	if body != "" {
		b.WriteString(body)
	}
	return b.String()
}

func classifyStatus(code int) string {
	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return "Informational"
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return "Success"
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return "Redirection"
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return "Client Error"
	case code >= http.StatusInternalServerError && code < 600:
		return "Server Error"
	default:
		return "Unknown"
	}
}

func main() {
	headers := map[string]string{
		"Host":         "api.example.com",
		"Content-Type": "application/json",
	}
	req1 := buildHTTPRequest("GET", "/api/users", headers, "")
	req2 := buildHTTPRequest("POST", "/api/users", headers, `{"name":"Alice"}`)
	req3 := buildHTTPRequest("DELETE", "/api/users/42", headers, "")

	fmt.Println("=== Raw HTTP Requests ===")
	fmt.Println("--- GET ---")
	fmt.Println(req1)
	fmt.Println("--- POST ---")
	fmt.Println(req2)
	fmt.Println("--- DELETE ---")
	fmt.Println(req3)

	codes := []int{200, 201, 301, 400, 403, 404, 500, 502, 503}
	fmt.Println("\n=== Status Code Classification ===")
	fmt.Println("Code\tClass")
	fmt.Println("----\t-----")
	for _, c := range codes {
		fmt.Printf("%d\t%s\n", c, classifyStatus(c))
	}
}
