package main

import (
	"fmt"
	"net"
	"log"
	"container/list"
	"os"
	"bufio"
	"strings"
	"sync"
)


var(
	// Stored connections
	l_conn = list.New()
	conn_map = make(map[net.Conn]*list.Element)
	connMutex = sync.Mutex{}
	logInfo *log.Logger
)

func main() {

	// Handle log file
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logInfo = log.New(logFile, "INFO : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	log.SetOutput(logFile)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started on port 8080")

	go showOptions()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}


func displayConnections() {
	connMutex.Lock()
	defer connMutex.Unlock()

	fmt.Println("Active connections : ")
	for e := l_conn.Front(); e != nil; e = e.Next() {
		if conn, ok := e.Value.(net.Conn); ok {
			fmt.Printf("- %s\n", conn.RemoteAddr())
		}
	}
}


func handleConnection(conn net.Conn) {

	connMutex.Lock()
	e := l_conn.PushBack(conn)
	conn_map[conn] = e
	connMutex.Unlock()


	logInfo.Printf("Connection established with %s\n", conn.RemoteAddr())

	// Send message to client
	_, err := conn.Write([]byte("Connected."))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Asynchrone reading from client
	go func(c net.Conn) {
		reader := bufio.NewReader(c)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				logInfo.Printf("Client disconnected (%s): %v\n", c.RemoteAddr(), err)
				break
			}
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "BEGIN_FILE:") {
				filename := strings.TrimPrefix(line, "BEGIN_FILE:")
				file, err := os.Create("uploads/" + filename)
				if err != nil {
					logInfo.Printf("[!] Could not create file %s: %v", filename, err)
					continue
				}
				for {
					dataLine, err := reader.ReadString('\n')
					if err != nil {
						logInfo.Printf("[!] Error while reading file: %v", err)
						break
					}
					if strings.TrimSpace(dataLine) == "END_FILE" {
						break
					}
					file.WriteString(dataLine)
				}
				file.Close()
				logInfo.Printf("[+] File %s received from %s", filename, c.RemoteAddr())
				continue
			}

			logInfo.Printf("[Response from %s] %s", c.RemoteAddr(), line)
		}

		connMutex.Lock()
		if elem, ok := conn_map[c]; ok {
			l_conn.Remove(elem)
			delete(conn_map, c)
		}
		connMutex.Unlock()
		c.Close()
	}(conn)
}


func broadcastMessage(msg string) {
	connMutex.Lock()
	defer connMutex.Unlock()

	for e := l_conn.Front(); e != nil; e = e.Next() {
		if conn, ok := e.Value.(net.Conn); ok {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				logInfo.Printf("Error when sending \"%s\" to %s: %v\n", msg, conn.RemoteAddr(), err)
			}
			logInfo.Printf("Send command \"%s\" to %s\n", msg, conn.RemoteAddr())
		}
	}
}


func showOptions() {
	for {
		fmt.Println("\nServer options : ")
		fmt.Println("[+] 1 - Broadcast a message")
		fmt.Println("[+] 2 - Display active connections")
		fmt.Println("[+] 3 - Quit server")
		fmt.Print("> ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Message to send : ")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			broadcastMessage(msg)
		case 2:
			displayConnections()
		case 3:
			fmt.Println("Shutting down server...")
			os.Exit(0)
		default:
			fmt.Println("[!] Invalid choice, please try again.")
		}
	}
}


