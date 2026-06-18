package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
}

func main() {
	// This is still a TCP server.
	// HTTP is just the text format we read/write on top of TCP.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	fmt.Println("HTTP server listening on http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read error:", err)
		return
	}

	rawRequest := string(buf[:n])
	fmt.Println("----- raw request -----")
	fmt.Println(rawRequest)

	req, ok := parseRequest(rawRequest)
	if !ok {
		writeResponse(conn, "400 Bad Request", "text/plain", "bad request\n")
		return
	}

	if req.Method != "GET" {
		writeResponse(conn, "405 Method Not Allowed", "text/plain", "only GET is supported\n")
		return
	}

	switch req.Path {
	case "/":
		body := "<h1>Hello from Go</h1><p>This HTTP response was written manually.</p>\n"
		writeResponse(conn, "200 OK", "text/html", body)

	case "/about":
		body := "This is a tiny HTTP server built in Go without net/http.\n"
		writeResponse(conn, "200 OK", "text/plain", body)

	case "/health":
		writeResponse(conn, "200 OK", "text/plain", "ok\n")

	default:
		writeResponse(conn, "404 Not Found", "text/plain", "404 page not found\n")
	}
}

func parseRequest(raw string) (Request, bool) {
	lines := strings.Split(raw, "\r\n")
	if len(lines) == 0 {
		return Request{}, false
	}

	// Example request line:
	// GET /about HTTP/1.1
	requestLine := lines[0]
	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return Request{}, false
	}

	req := Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
	}

	return req, true
}

func writeResponse(conn net.Conn, status string, contentType string, body string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
			"Content-Type: %s\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			"%s",
		status,
		contentType,
		len(body),
		body,
	)

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("write error:", err)
	}
}
