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
	"time"
)

type FirewallRule struct {
	DisplayName string      `json:"DisplayName"`
	Direction   interface{} `json:"Direction"`
	Action      interface{} `json:"Action"`
	Enabled     interface{} `json:"Enabled"`
	Protocol    interface{} `json:"Protocol"`
	LocalPort   interface{} `json:"LocalPort"`
	Program     interface{} `json:"Program"`
}

type FirewallResults struct {
	TotalRules int            `json:"total_rules"`
	Rules      []FirewallRule `json:"rules"`
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	s := fmt.Sprintf("%v", v)
	if s == "" || strings.Contains(s, "System.Object") || s == "[]" {
		return ""
	}
	return s
}

func getFirewallRules() ([]FirewallRule, error) {
	psCmd := "Get-NetFirewallRule -Enabled True | Select-Object DisplayName, Direction, Action, Enabled, @{Name='Protocol';Expression={($PSItem | Get-NetFirewallPortFilter).Protocol}}, @{Name='LocalPort';Expression={($PSItem | Get-NetFirewallPortFilter).LocalPort}}, @{Name='Program';Expression={($PSItem | Get-NetFirewallApplicationFilter).Program}} | ConvertTo-Json"
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	
	done := make(chan bool)
	go func() {
		fmt.Print("Filtering for Special Ports")
		for {
			select {
			case <-done:
				fmt.Println(" [Done]")
				return
			default:
				fmt.Print(".")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	output, err := cmd.Output()
	done <- true
	
	if err != nil {
		return []FirewallRule{}, nil
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []FirewallRule{}, nil
	}

	var allRules []FirewallRule
	if strings.HasPrefix(trimmed, "[") {
		_ = json.Unmarshal([]byte(trimmed), &allRules)
	} else if strings.HasPrefix(trimmed, "{") {
		var single FirewallRule
		if err := json.Unmarshal([]byte(trimmed), &single); err == nil {
			allRules = append(allRules, single)
		}
	}

	// Filter: ONLY keep rules that have a numeric Port (e.g., "80", "443", "3306", "5000-5010")
	// We use a regex to check for at least one digit in the Port field.
	digitRegex := regexp.MustCompile(`\d`)
	var filtered []FirewallRule
	for _, r := range allRules {
		port := toString(r.LocalPort)
		if digitRegex.MatchString(port) {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

func main() {
	fmt.Println("Analyzing Specific Port Rules in Windows Firewall...")
	fmt.Println("Filtering for rules with explicit port numbers.")
	
	rules, _ := getFirewallRules()

	results := FirewallResults{
		TotalRules: len(rules),
		Rules:      rules,
	}

	outputDir := "./output"
	_ = os.MkdirAll(outputDir, 0755)

	var buf bytes.Buffer
	buf.WriteString("=== Windows Firewall Port-Specific Rules Report ===\n")
	buf.WriteString(fmt.Sprintf("Total Rules with Specific Ports: %d\n\n", results.TotalRules))
	
	format := "%-45s | %-10s | %-8s | %-10s | %-12s | %s\n"
	buf.WriteString(fmt.Sprintf(format, "Display Name", "Direction", "Action", "Protocol", "Port", "Program/Path"))
	buf.WriteString(strings.Repeat("-", 125) + "\n")

	for _, r := range results.Rules {
		dir := toString(r.Direction)
		act := toString(r.Action)
		proto := toString(r.Protocol)
		port := toString(r.LocalPort)
		prog := toString(r.Program)

		if dir == "1" { dir = "Inbound" }
		if dir == "2" { dir = "Outbound" }
		if dir == "" { dir = "Any" }
		
		if act == "2" { act = "Allow" }
		if act == "4" { act = "Block" }
		if act == "" { act = "Any" }

		if proto == "" { proto = "Any" }
		if port == "" { port = "Any" }
		if prog == "" { prog = "Any" }

		name := r.DisplayName
		if len(name) > 43 { name = name[:40] + "..." }

		buf.WriteString(fmt.Sprintf(format, name, dir, act, proto, port, prog))
	}

	fmt.Print(buf.String())

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(filepath.Join(outputDir, "firewall_report.json"), jsonData, 0644)
	_ = os.WriteFile(filepath.Join(outputDir, "firewall_report.txt"), buf.Bytes(), 0644)

	fmt.Printf("\nReport saved to %s/firewall_report.json and %s/firewall_report.txt\n", outputDir, outputDir)
}
