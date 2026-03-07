package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type XAMPPInstance struct {
	XAMPPPath    string `json:"xampp_path"`
	XAMPPSize    string `json:"xampp_size"`
	HtdocsSize   string `json:"htdocs_size"`
	PHPPath      string `json:"php_path,omitempty"`
	PHPVersion   string `json:"php_version,omitempty"`
	MySQLPath    string `json:"mysql_path,omitempty"`
	MySQLVersion string `json:"mysql_version,omitempty"`
}

type XAMPPResults struct {
	Instances []XAMPPInstance `json:"instances"`
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
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
	v := strings.TrimSpace(string(out))
	re := regexp.MustCompile(`Ver [^ ]+ Distrib ([^ ,]+)`)
	match := re.FindStringSubmatch(v)
	if len(match) > 1 {
		return match[1]
	}
	re2 := regexp.MustCompile(`Ver ([^ ,]+)`)
	match2 := re2.FindStringSubmatch(v)
	if len(match2) > 1 {
		return match2[1]
	}
	return v
}

func findXAMPP() []XAMPPInstance {
	var instances []XAMPPInstance
	drives := []string{"C:\\", "D:\\", "E:\\"}

	for _, drive := range drives {
		entries, err := os.ReadDir(drive)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := strings.ToLower(entry.Name())
			if entry.IsDir() && strings.HasPrefix(name, "xampp") {
				path := filepath.Join(drive, entry.Name())
				inst := XAMPPInstance{XAMPPPath: path}

				// Get Directory Sizes
				if size, err := getDirSize(path); err == nil {
					inst.XAMPPSize = formatSize(size)
				}
				htdocsPath := filepath.Join(path, "htdocs")
				if size, err := getDirSize(htdocsPath); err == nil {
					inst.HtdocsSize = formatSize(size)
				} else {
					inst.HtdocsSize = "N/A"
				}

				// PHP Check
				php := filepath.Join(path, "php", "php.exe")
				if _, err := os.Stat(php); err == nil {
					inst.PHPPath = php
					inst.PHPVersion = getPHPVersion(php)
				}

				// MySQL Check
				mysql := filepath.Join(path, "mysql", "bin", "mysql.exe")
				if _, err := os.Stat(mysql); err == nil {
					inst.MySQLPath = mysql
					inst.MySQLVersion = getMySQLVersion(mysql)
				}

				if inst.PHPPath != "" || inst.MySQLPath != "" {
					instances = append(instances, inst)
				}
			}
		}
	}
	return instances
}

func main() {
	fmt.Println("Searching for XAMPP installations (xampp*)...")
	fmt.Println("This may take a moment to calculate directory sizes...")
	results := XAMPPResults{
		Instances: findXAMPP(),
	}

	if len(results.Instances) == 0 {
		fmt.Println("No XAMPP installations found.")
		return
	}

	outputDir := "./output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("=== XAMPP Installation Report ===\n")
	for i, inst := range results.Instances {
		buf.WriteString(fmt.Sprintf("\n[%d] Path: %s\n", i+1, inst.XAMPPPath))
		buf.WriteString(fmt.Sprintf("    XAMPP Total Size: %s\n", inst.XAMPPSize))
		buf.WriteString(fmt.Sprintf("    htdocs Folder Size: %s\n", inst.HtdocsSize))
		if inst.PHPPath != "" {
			buf.WriteString(fmt.Sprintf("    PHP:     %s\n", inst.PHPPath))
			buf.WriteString(fmt.Sprintf("    Version: %s\n", inst.PHPVersion))
		} else {
			buf.WriteString("    PHP:     Not found\n")
		}
		if inst.MySQLPath != "" {
			buf.WriteString(fmt.Sprintf("    MySQL:   %s\n", inst.MySQLPath))
			buf.WriteString(fmt.Sprintf("    Version: %s\n", inst.MySQLVersion))
		} else {
			buf.WriteString("    MySQL:   Not found\n")
		}
	}

	fmt.Print(buf.String())

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(filepath.Join(outputDir, "xampp_report.json"), jsonData, 0644)
	_ = os.WriteFile(filepath.Join(outputDir, "xampp_report.txt"), buf.Bytes(), 0644)

	fmt.Printf("\nReport saved to %s/xampp_report.json and %s/xampp_report.txt\n", outputDir, outputDir)
}
