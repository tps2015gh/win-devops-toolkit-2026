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

type HotFix struct {
	Description string `json:"Description"`
	HotFixID    string `json:"HotFixID"`
	InstalledBy string `json:"InstalledBy"`
	InstalledOn string `json:"InstalledOn"`
}

type PatchResults struct {
	TotalPatches int      `json:"total_patches"`
	Patches      []HotFix `json:"patches"`
}

func getHotFixes() ([]HotFix, error) {
	// Let's get the full fields from PowerShell and try to parse them.
	// Sometimes PowerShell's ConvertTo-Json with Select-Object can be tricky with casing.
	psCmd := "Get-HotFix | Select-Object Description, HotFixID, InstalledBy, @{Name='InstalledOn';Expression={$_.InstalledOn.ToString('yyyy-MM-dd')}} | ConvertTo-Json"
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	
	output, err := cmd.Output()
	if err != nil {
		return []HotFix{}, nil
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []HotFix{}, nil
	}

	var patches []HotFix
	if strings.HasPrefix(trimmed, "[") {
		_ = json.Unmarshal([]byte(trimmed), &patches)
	} else if strings.HasPrefix(trimmed, "{") {
		var single HotFix
		if err := json.Unmarshal([]byte(trimmed), &single); err == nil {
			patches = append(patches, single)
		}
	}

	return patches, nil
}

func main() {
	fmt.Println("Retrieving Windows HotFixes (Patches)...")
	
	patches, _ := getHotFixes()

	results := PatchResults{
		TotalPatches: len(patches),
		Patches:      patches,
	}

	outputDir := "./output"
	_ = os.MkdirAll(outputDir, 0755)

	var buf bytes.Buffer
	buf.WriteString("=== Windows Patch / HotFix Report ===\n")
	buf.WriteString(fmt.Sprintf("Total Patches Installed: %d\n\n", results.TotalPatches))
	
	if results.TotalPatches > 0 {
		buf.WriteString(fmt.Sprintf("%-15s | %-20s | %-12s | %s\n", "HotFixID", "Description", "Date", "InstalledBy"))
		buf.WriteString(strings.Repeat("-", 85) + "\n")
		for _, p := range results.Patches {
			// Clean up fields
			id := p.HotFixID
			if id == "" { id = "N/A" }
			desc := p.Description
			if desc == "" { desc = "N/A" }
			date := p.InstalledOn
			if date == "" { date = "N/A" }
			by := p.InstalledBy
			if by == "" { by = "N/A" }

			buf.WriteString(fmt.Sprintf("%-15s | %-20s | %-12s | %s\n", 
				id, desc, date, by))
		}
	} else {
		buf.WriteString("No patches found.\n")
	}

	fmt.Print(buf.String())

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(filepath.Join(outputDir, "patch_report.json"), jsonData, 0644)
	_ = os.WriteFile(filepath.Join(outputDir, "patch_report.txt"), buf.Bytes(), 0644)

	fmt.Printf("\nReport saved to %s/patch_report.json and %s/patch_report.txt\n", outputDir, outputDir)
}
