package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type DiskStatus struct {
	All  uint64
	Free uint64
	Used uint64
}

func getDiskSpace(path string) (DiskStatus, error) {
	var freeBytes, totalBytes, availBytes uint64
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpace := kernel32.NewProc("GetDiskFreeSpaceExW")
	
	ret, _, err := getDiskFreeSpace.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availBytes)),
	)
	
	if ret == 0 {
		return DiskStatus{}, err
	}
	
	return DiskStatus{
		All:  totalBytes,
		Free: availBytes,
		Used: totalBytes - availBytes,
	}, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nWSL Distribution Manager")
		fmt.Println("========================")
		fmt.Println("1. List Installed Distros (with Disk Usage)")
		fmt.Println("2. List Available Distros Online")
		fmt.Println("3. Install New Distro")
		fmt.Println("4. Remove (Unregister) Distro")
		fmt.Println("5. System Disk Summary")
		fmt.Println("6. Exit")
		fmt.Print("Select (1-6): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			listInstalled()
		case "2":
			listOnline()
		case "3":
			installDistro(reader)
		case "4":
			removeDistro(reader)
		case "5":
			showDiskSummary()
		case "6":
			return
		default:
			fmt.Println("Invalid selection.")
		}
	}
}

func listInstalled() {
	fmt.Println("\nInstalled WSL Distributions:")
	cmd := exec.Command("wsl", "--list", "--verbose")
	out, _ := cmd.CombinedOutput()
	fmt.Println(cleanOutput(out))
}

func listOnline() {
	fmt.Println("\nAvailable Distributions Online:")
	cmd := exec.Command("wsl", "--list", "--online")
	out, _ := cmd.CombinedOutput()
	fmt.Println(cleanOutput(out))
}

func installDistro(reader *bufio.Reader) {
	fmt.Println("\nFetching available distributions...")
	cmd := exec.Command("wsl", "--list", "--online")
	out, _ := cmd.CombinedOutput()
	lines := strings.Split(cleanOutput(out), "\n")
	
	var onlineDistros []string
	fmt.Println("\nSelect a distribution to install:")
	count := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "NAME") || strings.Contains(line, "FRIENDLY") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			fmt.Printf("[%d] %s\n", count, fields[0])
			onlineDistros = append(onlineDistros, fields[0])
			count++
		}
	}

	if len(onlineDistros) == 0 {
		fmt.Println("No distributions found.")
		return
	}

	showDiskSummary()
	fmt.Printf("\nEnter number to install (1-%d) or 'c' to cancel: ", len(onlineDistros))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "c" { return }
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(onlineDistros) {
		fmt.Println("Invalid selection.")
		return
	}

	selected := onlineDistros[idx-1]
	fmt.Printf("Installing %s... (This may take a few minutes)\n", selected)
	installCmd := exec.Command("wsl", "--install", "-d", selected)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Run()
}

func removeDistro(reader *bufio.Reader) {
	cmd := exec.Command("wsl", "--list", "--all")
	out, _ := cmd.CombinedOutput()
	lines := strings.Split(cleanOutput(out), "\n")
	
	var installedDistros []string
	fmt.Println("\nSelect a distribution to REMOVE (WARNING: All data will be deleted!):")
	count := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "NAME") || strings.Contains(line, "STATE") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 { continue }
		name := strings.TrimPrefix(fields[0], "*")
		fmt.Printf("[%d] %s\n", count, name)
		installedDistros = append(installedDistros, name)
		count++
	}

	if len(installedDistros) == 0 {
		fmt.Println("No distributions found.")
		return
	}

	fmt.Printf("\nEnter number to UNREGISTER (1-%d) or 'c' to cancel: ", len(installedDistros))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "c" { return }
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(installedDistros) {
		fmt.Println("Invalid selection.")
		return
	}

	selected := installedDistros[idx-1]
	fmt.Printf("Removing %s...\n", selected)
	removeCmd := exec.Command("wsl", "--unregister", selected)
	removeCmd.Run()
	fmt.Println("Removal complete.")
}

func showDiskSummary() {
	disk, err := getDiskSpace("C:\\")
	if err != nil {
		fmt.Printf("Error getting disk space: %v\n", err)
		return
	}
	fmt.Printf("\nSystem Disk Summary (C:):\n")
	fmt.Printf("  Total: %.2f GB\n", float64(disk.All)/1024/1024/1024)
	fmt.Printf("  Free:  %.2f GB\n", float64(disk.Free)/1024/1024/1024)
	fmt.Printf("  Used:  %.2f GB\n", float64(disk.Used)/1024/1024/1024)
}

func cleanOutput(b []byte) string {
	s := string(b)
	if strings.Contains(s, "\x00") {
		return strings.ReplaceAll(s, "\x00", "")
	}
	return s
}
