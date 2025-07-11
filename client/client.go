package main

import (
	"fmt"
	"net"
	"log"
	"time"
)


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connection established to localhost:8080")

	// Received data from server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("From Server (%s) : %s\n", conn.RemoteAddr(), string(buf[:n]))

	time.Sleep(2 * time.Second)

	// Send data to server
	_, err = conn.Write([]byte("Hello from client."))
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("From server (%s) : %s\n", conn.RemoteAddr(), string(buf[:n]))

		time.Sleep(1 * time.Second)
	}
}