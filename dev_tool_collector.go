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
	Status   string `json:"status"`
	Version  string `json:"version"`
	Path     string `json:"path"`
}

type DevToolsReport struct {
	Platform string     `json:"platform"`
	Tools    []ToolInfo `json:"tools"`
}

func main() {
	fmt.Println("Development Tools Discovery Engine")
	fmt.Println("==================================")

	outputDir := "./output"
	os.MkdirAll(outputDir, 0755)

	toolsToCheck := []struct {
		name string
		cmd  string
		args []string
	}{
		{"Go", "go", []string{"version"}},
		{"Rust", "rustc", []string{"--version"}},
		{"Java Runtime", "java", []string{"-version"}},
		{"Java Compiler", "javac", []string{"-version"}},
		{"Python", "python", []string{"--version"}},
		{"Node.js", "node", []string{"-v"}},
		{"Bun", "bun", []string{"-v"}},
		{"NPM", "npm", []string{"-v"}},
		{"Yarn", "yarn", []string{"-v"}},
		{"Git", "git", []string{"--version"}},
		{"SSH", "ssh", []string{"-V"}},
		{"GitHub CLI", "gh", []string{"--version"}},
		{"PHP", "php", []string{"-v"}},
		{"Perl", "perl", []string{"-v"}},
		{"TypeScript", "tsc", []string{"-v"}},
		{"Playwright", "npx", []string{"playwright", "--version"}},
		{"VS Code", "code", []string{"--version"}},
		{"Visual Studio", "vswhere", []string{"-latest", "-property", "displayName"}},
		{"Node Packer (pkg)", "pkg", []string{"-v"}},
		{"Node Packer (nexe)", "nexe", []string{"--version"}},
		{"Python Pandas", "python", []string{"-c", "import pandas; print(pandas.__version__)"}},
		{"Python Tkinter", "python", []string{"-c", "import tkinter; print('Installed')"}},
		{"MySQL Client", "mysql", []string{"--version"}},
		{"PostgreSQL Client", "psql", []string{"--version"}},
		{"C# Compiler", "csc", []string{"-version"}},
		{"Dotnet SDK", "dotnet", []string{"--version"}},
		{"Android Debug Bridge", "adb", []string{"version"}},
		{"Nginx", "nginx", []string{"-v"}},
		{"Docker", "docker", []string{"--version"}},
		{"Kubernetes (k8s)", "kubectl", []string{"version", "--client"}},
		{"k3s", "k3s", []string{"--version"}},
		{"Apache (httpd)", "httpd", []string{"-v"}},
		{"Lighttpd", "lighttpd", []string{"-v"}},
		{"Kong", "kong", []string{"version"}},
	}

	var results []ToolInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	total := len(toolsToCheck)
	fmt.Printf("Starting discovery for %d tools...
", total)

	for i, t := range toolsToCheck {
		wg.Add(1)
		go func(name, command string, args []string, index int) {
			defer wg.Done()
			
			fmt.Printf("[%d/%d] Checking %s...
", index+1, total, name)

			info := ToolInfo{Name: name, Status: "Not Found", Version: "N/A"}
			
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
					versionStr := strings.TrimSpace(out.String() + stderr.String())
					if versionStr != "" {
						lines := strings.Split(versionStr, "
")
						info.Version = strings.TrimSpace(lines[0])
					}
				}
			}

			mu.Lock()
			results = append(results, info)
			mu.Unlock()
		}(t.name, t.cmd, t.args, i)
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

	fmt.Println("
Discovery Complete!")
	fmt.Printf("Report saved to %s/dev_tools.json and %s/dev_tools.html
", outputDir, outputDir)
}

func checkAndroidStudio(results *[]ToolInfo) {
	paths := []string{
		filepath.Join(os.Getenv("PROGRAMFILES"), "Android", "Android Studio", "bin", "studio64.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Android Studio", "bin", "studio64.exe"),
	}
	
	info := ToolInfo{Name: "Android Studio", Status: "Not Found", Version: "N/A"}
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
	const htmlTmpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Dev Tools Report</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 40px; background: #f4f7f6; }
        table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 0 20px rgba(0,0,0,0.1); }
        th, td { padding: 12px 15px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #009879; color: white; }
        tr:hover { background-color: #f5f5f5; }
        .status-installed { color: #009879; font-weight: bold; }
        .status-missing { color: #e74c3c; }
        .platform { margin-bottom: 20px; font-style: italic; color: #666; }
    </style>
</head>
<body>
    <h1>Development Tools Audit</h1>
    <div class="platform">Platform: {{.Platform}}</div>
    <table>
        <thead>
            <tr>
                <th>Tool Name</th>
                <th>Status</th>
                <th>Version</th>
                <th>System Path</th>
            </tr>
        </thead>
        <tbody>
            {{range .Tools}}
            <tr>
                <td>{{.Name}}</td>
                <td class="{{if eq .Status "Installed"}}status-installed{{else}}status-missing{{end}}">{{.Status}}</td>
                <td>{{.Version}}</td>
                <td><small>{{.Path}}</small></td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`

	tmpl, _ := template.New("report").Parse(htmlTmpl)
	f, _ := os.Create(filename)
	defer f.Close()
	tmpl.Execute(f, report)
}
