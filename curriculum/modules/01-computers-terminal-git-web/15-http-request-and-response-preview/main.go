package main

import "fmt"

func httpResponse(statusCode int, body string) string {
	statusText := ""
	if statusCode == 200 {
		statusText = "OK"
	} else if statusCode == 404 {
		statusText = "Not Found"
	} else if statusCode == 500 {
		statusText = "Internal Server Error"
	}
	return fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, statusText, len(body), body)
}

func httpRequest(method, path, body string) string {
	if path == "/" || path == "" {
		return httpResponse(200, "Hello from Go server!\nMethod: "+method+"\nPath: "+path)
	}
	if path == "/api/hello" {
		return httpResponse(200, `{"message": "Hello, API!", "method": "`+method+`"}`)
	}
	return httpResponse(404, "Not Found: "+path)
}

func main() {
	fmt.Println("=== HTTP Request/Response Preview ===")
	fmt.Println()

	paths := []string{"/", "/api/hello", "/unknown"}
	methods := []string{"GET", "POST"}

	for _, path := range paths {
		fmt.Printf("--- %s %s ---\n", methods[0], path)
		resp := httpRequest(methods[0], path, "")
		fmt.Println(resp)
	}

	fmt.Println("HTTP (Hypertext Transfer Protocol) is the foundation of")
	fmt.Println("data communication on the web. A client sends a request")
	fmt.Println("and a server returns a response with a status code, headers, and body.")
}
