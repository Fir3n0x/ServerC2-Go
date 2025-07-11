package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"bufio"
	"strings"
	"sync"
	"github.com/google/uuid"
)

type Client struct {
	ID string
	Conn net.Conn
}


var(
	// Stored connections
	logInfo *log.Logger
	clients = make(map[string]Client)
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

	clientsMu.Lock()
	defer clientsMu.Unlock()

	fmt.Println("Active clients:")
	for id, client := range clients {
		fmt.Printf("- ID: %s | Address: %s\n", id, client.Conn.RemoteAddr())
	}
}


func handleConnection(conn net.Conn) {

	id := uuid.New().String()

	client := Client{
		ID: id,
		Conn: conn,
	}

	clientsMu.Lock()
	clients[id] = client
	clientsMu.Unlock()


	logInfo.Printf("Client %s connected from %s", id, conn.RemoteAddr())

	// Send message to client
	_, err := conn.Write([]byte("Your ID: " + id + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Asynchrone reading from client
	go func(c Client) {
		reader := bufio.NewReader(c.Conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				logInfo.Printf("Client %s disconnected (%s): %v\n", c.ID, c.Conn.RemoteAddr(), err)
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
				continue
			}
		}

		clientsMu.Lock()
		delete(clients, c.ID)
		clientsMu.Unlock()
		c.Conn.Close()
	}(client)
}


func broadcastMessage(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	var toDelete []string

	for id, client := range clients {
		_, err := client.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			logInfo.Printf("Error when sending \"%s\" to %s client (%s): %v\n", msg, id, client.Conn.RemoteAddr(), err)
			toDelete = append(toDelete, id)
		} else {
			logInfo.Printf("[+] Sent message \"%s\" to %s client (%s)\n", msg, id, client.Conn.RemoteAddr())
		}
	}

	// Cleanup disconnected clients
	for _, id := range toDelete {
		clients[id].Conn.Close()
		delete(clients, id)
	}
}


func showOptions() {
	for {
		fmt.Println("\nServer options : ")
		fmt.Println("[+] 1 - Broadcast a message")
		fmt.Println("[+] 2 - Display active connections")
		fmt.Println("[+] 3 - Send a message to a specific client")
		fmt.Println("[+] 4 - Quit server")
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
			displayConnections()
			fmt.Print("Enter client ID: ")
			reader := bufio.NewReader(os.Stdin)
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			clientsMu.Lock()
			client, exists := clients[id]
			clientsMu.Unlock()

			if !exists {
				fmt.Println("[!] Client not found")
				break
			}

			fmt.Print("Enter message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)

			_, err := client.Conn.Write([]byte(msg + "\n"))
			if err != nil {
				logInfo.Printf("[!] Failed to send message to %s: %v", id, err)
			} else {
				logInfo.Printf("[+] Sent message \"%s\" to %s (%s)", msg, id, client.Conn.RemoteAddr())
			}
		case 4:
			fmt.Println("Shutting down server...")
			os.Exit(0)
		default:
			fmt.Println("[!] Invalid choice, please try again.")
		}
	}
}


