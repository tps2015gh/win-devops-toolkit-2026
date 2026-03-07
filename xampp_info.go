package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type XAMPPInfo struct {
	Path         string
	PHPPath      string
	PHPVersion   string
	MySQLPath    string
	MySQLVersion string
}

func getPHPVersion(phpExe string) string {
	cmd := exec.Command(phpExe, "-v")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return "Unknown"
}

func getMySQLVersion(mysqlExe string) string {
	cmd := exec.Command(mysqlExe, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func findXAMPP() []XAMPPInfo {
	var results []XAMPPInfo
	drives := []string{"C:\\", "D:\\", "E:\\"}

	for _, drive := range drives {
		entries, err := os.ReadDir(drive)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := strings.ToLower(entry.Name())
			if entry.IsDir() && (strings.HasPrefix(name, "xampp")) {
				xamppPath := filepath.Join(drive, entry.Name())
				info := XAMPPInfo{Path: xamppPath}

				// Check PHP
				phpExe := filepath.Join(xamppPath, "php", "php.exe")
				if _, err := os.Stat(phpExe); err == nil {
					info.PHPPath = phpExe
					info.PHPVersion = getPHPVersion(phpExe)
				}

				// Check MySQL
				mysqlExe := filepath.Join(xamppPath, "mysql", "bin", "mysql.exe")
				if _, err := os.Stat(mysqlExe); err == nil {
					info.MySQLPath = mysqlExe
					info.MySQLVersion = getMySQLVersion(mysqlExe)
				}

				if info.PHPPath != "" || info.MySQLPath != "" {
					results = append(results, info)
				}
			}
		}
	}
	return results
}

func main() {
	fmt.Println("Searching for XAMPP installations (xampp*)...")
	xampps := findXAMPP()

	if len(xampps) == 0 {
		fmt.Println("No XAMPP installations found.")
		return
	}

	re := regexp.MustCompile(`Ver [^ ]+`)

	for i, x := range xampps {
		fmt.Printf("\n[%d] XAMPP Path: %s\n", i+1, x.Path)
		
		if x.PHPPath != "" {
			fmt.Printf("    - PHP:    %s\n", x.PHPPath)
			fmt.Printf("      Version: %s\n", x.PHPVersion)
		} else {
			fmt.Println("    - PHP: Not found")
		}

		if x.MySQLPath != "" {
			fmt.Printf("    - MySQL:  %s\n", x.MySQLPath)
			match := re.FindString(x.MySQLVersion)
			if match != "" {
				fmt.Printf("      Version: %s\n", match)
			} else {
				fmt.Printf("      Version: %s\n", x.MySQLVersion)
			}
		} else {
			fmt.Println("    - MySQL: Not found")
		}
	}
}
