package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FirewallRule struct {
	DisplayName string `json:"DisplayName"`
	Direction   string `json:"Direction"`
	Action      string `json:"Action"`
	Enabled     string `json:"Enabled"`
	Protocol    string `json:"Protocol"`
	LocalPort   string `json:"LocalPort"`
	Program     string `json:"Program"`
}

type FirewallResults struct {
	TotalRules int            `json:"total_rules"`
	Rules      []FirewallRule `json:"rules"`
}

func getFirewallRules() ([]FirewallRule, error) {
	// Simple PowerShell command to fetch enabled rules as JSON
	psCmd := "Get-NetFirewallRule -Enabled True | Select-Object DisplayName, Direction, Action, Enabled, @{Name='Protocol';Expression={($PSItem | Get-NetFirewallPortFilter).Protocol}}, @{Name='LocalPort';Expression={($PSItem | Get-NetFirewallPortFilter).LocalPort}}, @{Name='Program';Expression={($PSItem | Get-NetFirewallApplicationFilter).Program}} | ConvertTo-Json"
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	
	output, err := cmd.Output()
	if err != nil {
		return []FirewallRule{}, nil
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []FirewallRule{}, nil
	}

	var rules []FirewallRule
	if strings.HasPrefix(trimmed, "[") {
		_ = json.Unmarshal([]byte(trimmed), &rules)
	} else if strings.HasPrefix(trimmed, "{") {
		var single FirewallRule
		if err := json.Unmarshal([]byte(trimmed), &single); err == nil {
			rules = append(rules, single)
		}
	}

	return rules, nil
}

func main() {
	fmt.Println("Analyzing Windows Firewall Rules (Enabled only)...")
	fmt.Println("This may take a moment to fetch ports and programs...")

	rules, _ := getFirewallRules()

	results := FirewallResults{
		TotalRules: len(rules),
		Rules:      rules,
	}

	outputDir := "./output"
	_ = os.MkdirAll(outputDir, 0755)

	var buf bytes.Buffer
	buf.WriteString("=== Windows Firewall Active Rules Report ===\n")
	buf.WriteString(fmt.Sprintf("Total Enabled Rules Found: %d\n\n", results.TotalRules))
	
	format := "%-45s | %-10s | %-8s | %-10s | %-12s | %s\n"
	buf.WriteString(fmt.Sprintf(format, "Display Name", "Direction", "Action", "Protocol", "Port", "Program/Path"))
	buf.WriteString(strings.Repeat("-", 125) + "\n")

	for _, r := range results.Rules {
		// Clean up fields for display
		port := r.LocalPort
		if port == "" || strings.Contains(port, "System.Object") { port = "Any" }
		proto := r.Protocol
		if proto == "" || strings.Contains(proto, "System.Object") { proto = "Any" }
		prog := r.Program
		if prog == "" || strings.Contains(prog, "System.Object") { prog = "Any" }

		name := r.DisplayName
		if len(name) > 43 { name = name[:40] + "..." }

		buf.WriteString(fmt.Sprintf(format, 
			name, r.Direction, r.Action, proto, port, prog))
	}

	// Print to screen
	fmt.Print(buf.String())

	// Save files
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(filepath.Join(outputDir, "firewall_report.json"), jsonData, 0644)
	_ = os.WriteFile(filepath.Join(outputDir, "firewall_report.txt"), buf.Bytes(), 0644)

	fmt.Printf("\nReport saved to %s/firewall_report.json and %s/firewall_report.txt\n", outputDir, outputDir)
}
