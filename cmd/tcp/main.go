package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 9000")

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

	fmt.Println("client connected:", conn.RemoteAddr())

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read error:", err)
		return
	}

	msg := string(buf[:n])

	fmt.Println("client said:", msg)

	response := "server received: " + msg

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Println("write error:", err)
		return
	}
}
