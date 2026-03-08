package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

type DiskStatus struct {
	All  uint64
	Free uint64
	Used uint64
}

const (
	MinSpaceRequired = 2 * 1024 * 1024 * 1024 // 2GB
	RecSpaceRequired = 8 * 1024 * 1024 * 1024 // 8GB
)

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
		fmt.Println("3. Install New Distro (with Disk Check)")
		fmt.Println("4. Remove (Unregister) Distro")
		fmt.Println("5. Login to Distro")
		fmt.Println("6. Install Dev Tools (Git, Nano, Python, Gemini-CLI, OpenCode AI)")
		fmt.Println("7. System Disk Summary")
		fmt.Println("8. Exit")
		fmt.Print("Select (1-8): ")

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
			loginDistro(reader)
		case "6":
			installToolsToDistro(reader)
		case "7":
			showDiskSummary()
		case "8":
			return
		default:
			fmt.Println("Invalid selection.")
		}
	}
}

func listInstalled() {
	fmt.Println("\nInstalled WSL Distributions (Status):")
	cmd := exec.Command("wsl", "--list", "--verbose")
	out, _ := cmd.CombinedOutput()
	fmt.Println(cleanOutput(out))

	fmt.Println("Disk Usage (ext4.vhdx):")
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Lxss`, registry.READ)
	if err != nil { return }
	defer k.Close()

	subkeys, _ := k.ReadSubKeyNames(-1)
	var totalWSL uint64
	for _, guid := range subkeys {
		if !strings.HasPrefix(guid, "{") { continue }
		sk, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Lxss\`+guid, registry.READ)
		if err != nil { continue }
		name, _, _ := sk.GetStringValue("DistributionName")
		basePath, _, _ := sk.GetStringValue("BasePath")
		sk.Close()

		vhdxPath := filepath.Join(basePath, "ext4.vhdx")
		if info, err := os.Stat(vhdxPath); err == nil {
			fmt.Printf("  - %-20s: %s\n", name, formatSize(info.Size()))
			totalWSL += uint64(info.Size())
		}
	}
	fmt.Printf("\nTotal WSL Disk Usage: %s\n", formatSize(int64(totalWSL)))
}

func listOnline() {
	fmt.Println("\nAvailable Distributions Online:")
	fmt.Println("Typical initial size: 1.5GB - 2.0GB")
	cmd := exec.Command("wsl", "--list", "--online")
	out, _ := cmd.CombinedOutput()
	fmt.Println(cleanOutput(out))
}

func installDistro(reader *bufio.Reader) {
	disk, _ := getDiskSpace("C:\\")
	
	fmt.Println("\nFetching available distributions...")
	cmd := exec.Command("wsl", "--list", "--online")
	out, _ := cmd.CombinedOutput()
	output := cleanOutput(out)
	lines := strings.Split(output, "\n")
	
	var onlineDistros []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "The following") || strings.HasPrefix(line, "NAME") || strings.HasPrefix(line, "Install") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			onlineDistros = append(onlineDistros, fields[0])
		}
	}

	if len(onlineDistros) == 0 {
		fmt.Println("No distributions found.")
		return
	}

	fmt.Println("\nSelect a distribution to install:")
	for i, d := range onlineDistros {
		fmt.Printf("[%d] %s\n", i+1, d)
	}

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
	
	fmt.Printf("\nSystem Check for %s:\n", selected)
	fmt.Printf("  Free Space: %s\n", formatSize(int64(disk.Free)))
	fmt.Printf("  Estimated Install: ~2.00 GB\n")
	fmt.Printf("  Recommended Buffer: 8.00 GB\n")

	if disk.Free < MinSpaceRequired {
		fmt.Println("\n[!] ERROR: Not enough disk space to guarantee installation.")
		return
	} else if disk.Free < RecSpaceRequired {
		fmt.Println("\n[!] WARNING: Low disk space. You may run out of space after updates.")
		fmt.Print("Proceed anyway? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" { return }
	} else {
		fmt.Println("\n[+] Disk space check passed.")
	}

	fmt.Printf("Installing %s... (This may take a few minutes)\n", selected)
	installCmd := exec.Command("wsl", "--install", "-d", selected)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Run()
}

func removeDistro(reader *bufio.Reader) {
	distros := getInstalledDistros()
	if len(distros) == 0 {
		fmt.Println("No removable distributions found.")
		return
	}

	fmt.Println("\nSelect a distribution to REMOVE (WARNING: ALL DATA DELETED!):")
	for i, d := range distros {
		fmt.Printf("[%d] %s\n", i+1, d)
	}

	fmt.Printf("\nEnter number to UNREGISTER (1-%d) or 'c' to cancel: ", len(distros))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "c" { return }
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(distros) {
		fmt.Println("Invalid selection.")
		return
	}

	selected := distros[idx-1]
	fmt.Printf("ARE YOU SURE you want to delete %s? (y/n): ", selected)
	confirm, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" { return }

	fmt.Printf("Removing %s...\n", selected)
	removeCmd := exec.Command("wsl", "--unregister", selected)
	removeCmd.Run()
	fmt.Println("Removal complete.")
}

func loginDistro(reader *bufio.Reader) {
	distros := getInstalledDistros()
	if len(distros) == 0 {
		fmt.Println("No distributions found to login.")
		return
	}

	fmt.Println("\nSelect a distribution to login:")
	for i, d := range distros {
		fmt.Printf("[%d] %s\n", i+1, d)
	}

	fmt.Printf("\nEnter number to LOGIN (1-%d) or 'c' to cancel: ", len(distros))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "c" { return }
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(distros) {
		fmt.Println("Invalid selection.")
		return
	}

	selected := distros[idx-1]
	fmt.Printf("Logging into %s... (Type 'exit' to return to manager)\n", selected)
	
	cmd := exec.Command("wsl", "-d", selected)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func installToolsToDistro(reader *bufio.Reader) {
	distros := getInstalledDistros()
	if len(distros) == 0 {
		fmt.Println("No distributions found.")
		return
	}

	fmt.Println("\nSelect a distribution to install Dev Tools (Git, Nano, Python, Gemini-CLI, OpenCode AI):")
	for i, d := range distros {
		fmt.Printf("[%d] %s\n", i+1, d)
	}

	fmt.Printf("\nEnter number (1-%d) or 'c' to cancel: ", len(distros))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "c" { return }
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(distros) {
		fmt.Println("Invalid selection.")
		return
	}

	selected := distros[idx-1]
	fmt.Printf("\nInstalling tools to %s...\n", selected)

	script := `
	# 0. Check for sudo/root
	SUDO=""
	if [ "$(id -u)" -ne 0 ]; then
		if command -v sudo >/dev/null; then
			SUDO="sudo"
		else
			echo "Warning: Not root and 'sudo' not found. Some commands may fail."
		fi
	fi

	# 1. Install Basic Tools (Bash, Git, Nano, Python, Curl, Nodejs/NPM, sudo, which)
	echo "Phase 1: Installing system packages..."
	if command -v apt-get >/dev/null; then
		echo "Detected Debian/Ubuntu (apt)..."
		$SUDO apt-get update && $SUDO apt-get install -y bash git nano python3 python3-pip curl nodejs npm sudo which
	elif command -v pacman >/dev/null; then
		echo "Detected Arch Linux (pacman)..."
		$SUDO pacman -Sy --noconfirm bash git nano python python-pip curl nodejs npm sudo which
	elif command -v dnf >/dev/null; then
		echo "Detected Fedora/RHEL (dnf)..."
		$SUDO dnf install -y bash git nano python3 python3-pip curl nodejs npm sudo which
	elif command -v zypper >/dev/null; then
		echo "Detected openSUSE (zypper)..."
		$SUDO zypper install -y bash git nano python3 curl nodejs npm sudo which
	else
		echo "Could not detect package manager. Please install bash, git, nano, python, curl, and nodejs manually."
	fi

	# 2. Install Gemini-CLI (via NPM)
	echo "Phase 2: Installing Gemini CLI..."
	if command -v npm >/dev/null; then
		$SUDO npm install -g @google/gemini-cli
	else
		echo "NPM not found. Skipping Gemini CLI installation."
	fi

	# 3. Install OpenCode (AI Agent Tool)
	echo "Phase 3: Installing OpenCode AI Agent..."
	if command -v curl >/dev/null; then
		curl -fsSL https://opencode.ai/install | bash
	else
		echo "Curl not found. Skipping OpenCode installation."
	fi
	`
	
	cmd := exec.Command("wsl", "-u", "root", "-d", selected, "sh", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n[!] Installation failed or was interrupted: %v\n", err)
	} else {
		fmt.Printf("\n[+] Setup complete for %s!\n", selected)
	}
}

func getInstalledDistros() []string {
	cmd := exec.Command("wsl", "--list", "--quiet")
	out, _ := cmd.CombinedOutput()
	lines := strings.Split(cleanOutput(out), "\n")
	
	var distros []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" || name == "Windows" || strings.Contains(name, "Distributions") || strings.Contains(name, "Default") {
			continue
		}
		distros = append(distros, name)
	}
	return distros
}

func showDiskSummary() {
	disk, err := getDiskSpace("C:\\")
	if err != nil {
		fmt.Printf("Error getting disk space: %v\n", err)
		return
	}
	fmt.Printf("\nSystem Disk Summary (C:):\n")
	fmt.Printf("  Total: %s\n", formatSize(int64(disk.All)))
	fmt.Printf("  Free:  %s\n", formatSize(int64(disk.Free)))
	fmt.Printf("  Used:  %s\n", formatSize(int64(disk.Used)))
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit { return fmt.Sprintf("%d B", bytes) }
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func cleanOutput(b []byte) string {
	s := string(b)
	if strings.Contains(s, "\x00") {
		return strings.ReplaceAll(s, "\x00", "")
	}
	return s
}
