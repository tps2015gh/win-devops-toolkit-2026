package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type ToolInfo struct {
	Name     string `json:"name"`
	Group    string `json:"group"`
	Status   string `json:"status"`
	Version  string `json:"version"`
	Path     string `json:"path"`
}

type DevToolsReport struct {
	Platform string     `json:"platform"`
	Tools    []ToolInfo `json:"tools"`
}

func cleanOutput(b []byte) string {
	s := string(b)
	if strings.Contains(s, "\x00") {
		return strings.ReplaceAll(s, "\x00", "")
	}
	return s
}

func main() {
	fmt.Println("Development Tools Discovery Engine")
	fmt.Println("==================================")

	outputDir := "./output"
	os.MkdirAll(outputDir, 0755)

	toolsToCheck := []struct {
		name  string
		group string
		cmd   string
		args  []string
	}{
		// Languages & Runtimes
		{"Go", "Languages", "go", []string{"version"}},
		{"Rust", "Languages", "rustc", []string{"--version"}},
		{"Java Runtime", "Languages", "java", []string{"-version"}},
		{"Java Compiler", "Languages", "javac", []string{"-version"}},
		{"Python", "Languages", "python", []string{"--version"}},
		{"Node.js", "Languages", "node", []string{"-v"}},
		{"Bun", "Languages", "bun", []string{"-v"}},
		{"PHP", "Languages", "php", []string{"-v"}},
		{"Perl", "Languages", "perl", []string{"-v"}},
		{"TypeScript", "Languages", "tsc", []string{"-v"}},
		{"C# Compiler", "Languages", "csc", []string{"-version"}},
		{"Dotnet SDK", "Languages", "dotnet", []string{"--version"}},

		// Package Managers
		{"NPM", "Package Managers", "npm", []string{"-v"}},
		{"Yarn", "Package Managers", "yarn", []string{"-v"}},

		// Cloud & Containers
		{"Docker", "Cloud & Containers", "docker", []string{"--version"}},
		{"Kubernetes (k8s)", "Cloud & Containers", "kubectl", []string{"version", "--client"}},
		{"k3s", "Cloud & Containers", "k3s", []string{"--version"}},
		{"Nginx", "Servers", "nginx", []string{"-v"}},
		{"Apache (httpd)", "Servers", "httpd", []string{"-v"}},
		{"Lighttpd", "Servers", "lighttpd", []string{"-v"}},
		{"Kong", "Servers", "kong", []string{"version"}},

		// AI & Machine Learning
		{"Ollama", "AI & ML", "ollama", []string{"--version"}},
		{"LM Studio", "AI & ML", "lms", []string{"--version"}},
		{"vLLM", "AI & ML", "vllm", []string{"--version"}},
		{"Python Pandas", "AI & ML", "python", []string{"-c", "import pandas; print(pandas.__version__)"}},
		{"Python Tkinter", "AI & ML", "AI & ML", []string{"-c", "import tkinter; print('Installed')"}},

		// Development Tools
		{"Git", "Dev Tools", "git", []string{"--version"}},
		{"SSH", "Dev Tools", "ssh", []string{"-V"}},
		{"GitHub CLI", "Dev Tools", "gh", []string{"--version"}},
		{"VS Code", "Dev Tools", "code", []string{"--version"}},
		{"Visual Studio", "Dev Tools", "vswhere", []string{"-latest", "-property", "displayName"}},
		{"Android Debug Bridge", "Dev Tools", "adb", []string{"version"}},
		{"WSL", "OS & Virtualization", "wsl", []string{"--list", "--verbose"}},
		{"Ubuntu (via WSL)", "OS & Virtualization", "wsl", []string{"-d", "Ubuntu", "grep", "PRETTY_NAME", "/etc/os-release"}},

		// Databases
		{"MySQL Client", "Databases", "mysql", []string{"--version"}},
		{"PostgreSQL Client", "Databases", "psql", []string{"--version"}},

		// Node Packers
		{"Node Packer (pkg)", "Node Packers", "pkg", []string{"-v"}},
		{"Node Packer (nexe)", "Node Packers", "nexe", []string{"--version"}},
	}

	var results []ToolInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	total := len(toolsToCheck)
	fmt.Printf("Starting discovery for %d tools...\n", total)

	for i, t := range toolsToCheck {
		wg.Add(1)
		go func(name, group, command string, args []string, index int) {
			defer wg.Done()
			
			fmt.Printf("[%d/%d] Checking %s...\n", index+1, total, name)

			info := ToolInfo{Name: name, Group: group, Status: "Not Found", Version: "N/A"}
			
			path, err := exec.LookPath(command)
			if err == nil {
				info.Path = path
				cmd := exec.Command(command, args...)
				var out bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				
				if err := cmd.Run(); err == nil {
					info.Status = "Installed"
					versionStr := cleanOutput(out.Bytes()) + cleanOutput(stderr.Bytes())
					versionStr = strings.TrimSpace(versionStr)
					if versionStr != "" {
						lines := strings.Split(versionStr, "\n")
						info.Version = strings.TrimSpace(lines[0])
					}
				}
			}

			mu.Lock()
			results = append(results, info)
			mu.Unlock()
		}(t.name, t.group, t.cmd, t.args, i)
	}

	wg.Wait()

	checkAndroidStudio(&results)

	report := DevToolsReport{
		Platform: runtime.GOOS + " " + runtime.GOARCH,
		Tools:    results,
	}

	jsonData, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(filepath.Join(outputDir, "dev_tools.json"), jsonData, 0644)

	saveHTMLReport(report, filepath.Join(outputDir, "dev_tools.html"))

	fmt.Println("\nDiscovery Complete!")
	fmt.Printf("Report saved to %s/dev_tools.json and %s/dev_tools.html\n", outputDir, outputDir)
}

func checkAndroidStudio(results *[]ToolInfo) {
	paths := []string{
		filepath.Join(os.Getenv("PROGRAMFILES"), "Android", "Android Studio", "bin", "studio64.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Android Studio", "bin", "studio64.exe"),
	}
	
	info := ToolInfo{Name: "Android Studio", Group: "Dev Tools", Status: "Not Found", Version: "N/A"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			info.Status = "Installed"
			info.Path = p
			break
		}
	}
	*results = append(*results, info)
}

func saveHTMLReport(report DevToolsReport, filename string) {
	const htmlTmpl = `<!DOCTYPE html>
<html>
<head>
    <title>Dev Tools Report</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 40px; background: #f4f7f6; }
        table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 0 20px rgba(0,0,0,0.1); margin-bottom: 30px; }
        th, td { padding: 12px 15px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #009879; color: white; }
        tr:hover { background-color: #f5f5f5; }
        .status-installed { color: #009879; font-weight: bold; }
        .status-missing { color: #e74c3c; }
        .platform { margin-bottom: 20px; font-style: italic; color: #666; }
        .group-header { background: #34495e; color: white; padding: 10px 15px; font-weight: bold; font-size: 1.2em; border-radius: 4px 4px 0 0; }
    </style>
</head>
<body>
    <h1>Development Tools Audit</h1>
    <div class="platform">Platform: {{.Platform}}</div>
    
    {{ $tools := .Tools }}
    {{ range $group := (getGroups .Tools) }}
    <div class="group-header">{{ $group }}</div>
    <table>
        <thead>
            <tr>
                <th style="width: 25%">Tool Name</th>
                <th style="width: 15%">Status</th>
                <th style="width: 30%">Version</th>
                <th style="width: 30%">System Path</th>
            </tr>
        </thead>
        <tbody>
            {{ range $tools }}
                {{ if eq .Group $group }}
                <tr>
                    <td>{{.Name}}</td>
                    <td class="{{if eq .Status "Installed"}}status-installed{{else}}status-missing{{end}}">{{.Status}}</td>
                    <td>{{.Version}}</td>
                    <td><small>{{.Path}}</small></td>
                </tr>
                {{ end }}
            {{ end }}
        </tbody>
    </table>
    {{ end }}
</body>
</html>`

	funcMap := template.FuncMap{
		"getGroups": func(tools []ToolInfo) []string {
			groupMap := make(map[string]bool)
			var groups []string
			// Deterministic order for groups
			order := []string{"Languages", "Package Managers", "AI & ML", "Cloud & Containers", "Servers", "Dev Tools", "OS & Virtualization", "Databases", "Node Packers"}
			
			for _, t := range tools {
				groupMap[t.Group] = true
			}
			
			for _, o := range order {
				if groupMap[o] {
					groups = append(groups, o)
					delete(groupMap, o)
				}
			}
			// Add any remaining groups
			for g := range groupMap {
				groups = append(groups, g)
			}
			return groups
		},
	}

	tmpl, _ := template.New("report").Funcs(funcMap).Parse(htmlTmpl)
	f, _ := os.Create(filename)
	defer f.Close()
	tmpl.Execute(f, report)
}
