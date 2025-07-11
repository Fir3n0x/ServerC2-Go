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
	logFile, err = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logInfo = log.New(logFile, "INFO : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

func main() {

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

		connMutex.Lock()
		e := l_conn.PushBack(conn)
		conn_map[conn] = e
		connMutex.Unlock()

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
	defer conn.Close()

	logInfo.Printf("Connection established with %s.\n", conn.RemoteAddr())

	// Send message to client
	_, err := conn.Write([]byte("Connected."))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Receive message from client
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	logInfo.Printf("Received: %s, from %s.\n", string(buf[:n]), conn.RemoteAddr())

	for {
		_, err = conn.Read(buf)
		if err != nil {
			logInfo.Printf("Client disconnected (%s), %v\n", conn.RemoteAddr(), err)
			break
		}	
	}

	if elem, ok := conn_map[conn]; ok {
		connMutex.Lock()
		l_conn.Remove(elem)
		delete(conn_map, conn)
		connMutex.Unlock()
	}
}


func broadcastMessage(msg string) {
	connMutex.Lock()
	defer connMutex.Unlock()

	for e := l_conn.Front(); e != nil; e = e.Next() {
		if conn, ok := e.Value.(net.Conn); ok {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				logInfo.Printf("Error when sending \"%s\" to %s: %v\n", msg, conn.RemoteAddr(), err)
			}
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