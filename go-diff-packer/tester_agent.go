package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== Tester Agent: Starting End-to-End Tests for go-diff-packer ===")

	// 1. Setup Test Environment
	testRoot := "test_space"
	os.RemoveAll(testRoot)
	os.MkdirAll(testRoot, 0755)
	defer os.RemoveAll(testRoot)

	origDir, _ := filepath.Abs(filepath.Join(testRoot, "orig"))
	modDir, _ := filepath.Abs(filepath.Join(testRoot, "mod"))
	deployDir, _ := filepath.Abs(filepath.Join(testRoot, "deploy"))

	os.MkdirAll(origDir, 0755)
	os.MkdirAll(modDir, 0755)

	os.WriteFile(filepath.Join(origDir, "file1.txt"), []byte("original content"), 0644)
	os.WriteFile(filepath.Join(modDir, "file1.txt"), []byte("modified content"), 0644)
	os.WriteFile(filepath.Join(modDir, "newfile.txt"), []byte("new file content"), 0644)

	fmt.Println("[STEP 1] Test environment prepared.")

	// 2. Run Compare Mode
	fmt.Println("[STEP 2] Running compare mode...")
	cmdCompare := exec.Command("./go-diff-packer.exe", origDir, modDir)
	output, err := cmdCompare.CombinedOutput()
	if err != nil {
		fmt.Printf("Compare mode failed: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	// Find output directory from output text
	var outDiffDir string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Output directory: ") {
			outDiffDir = strings.TrimSpace(strings.TrimPrefix(trimmed, "Output directory: "))
			break
		}
	}

	if outDiffDir == "" {
		fmt.Println("Could not find output directory in compare mode output.")
		fmt.Println(string(output))
		os.Exit(1)
	}
	// We won't cleanup yet as we need it for deploy

	fmt.Printf("Detected output directory: %s\n", outDiffDir)

	// Verify outDiffDir contents
	if _, err := os.Stat(filepath.Join(outDiffDir, "file1.txt")); err != nil {
		fmt.Println("Error: file1.txt not found in diff folder")
		os.Exit(1)
	}
	fmt.Println("Compare mode results verified.")

	// 3. Run Deploy Mode
	fmt.Println("[STEP 3] Running deploy mode with automated input...")
	
	cmdDeploy := exec.Command("./go-diff-packer.exe", "--deploy")
	stdin, err := cmdDeploy.StdinPipe()
	if err != nil {
		fmt.Printf("Failed to get stdin pipe: %v\n", err)
		os.Exit(1)
	}

	stdout, err := cmdDeploy.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to get stdout pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmdDeploy.Start(); err != nil {
		fmt.Printf("Failed to start deploy mode: %v\n", err)
		os.Exit(1)
	}

	// Capture output and react to prompts
	go func() {
		defer stdin.Close()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("  [DEPLOY] %s\n", line)
			
			if strings.Contains(line, "Select option") || strings.Contains(line, "Select a folder") {
				fmt.Println("  [AGENT] Prompt detected -> Sending '1'")
				io.WriteString(stdin, "1\n")
			} else if strings.Contains(line, "Enter destination folder path") {
				fmt.Printf("  [AGENT] Prompt detected -> Sending: %s\n", deployDir)
				io.WriteString(stdin, deployDir+"\n")
			} else if strings.Contains(line, "Create it? (y/n)") {
				fmt.Println("  [AGENT] Prompt detected -> Sending 'y'")
				io.WriteString(stdin, "y\n")
			} else if strings.Contains(line, "Confirm deployment? (y/n)") {
				fmt.Println("  [AGENT] Prompt detected -> Sending 'y'")
				io.WriteString(stdin, "y\n")
			}
		}
	}()

	if err := cmdDeploy.Wait(); err != nil {
		fmt.Printf("Deploy mode process finished with error: %v\n", err)
	}

	// 4. Verify Final Deployment
	fmt.Println("[STEP 4] Verifying deployed files...")
	
	verifyFile(filepath.Join(deployDir, "file1.txt"), "modified content")
	verifyFile(filepath.Join(deployDir, "newfile.txt"), "new file content")

	fmt.Println("\n=== ALL TESTS PASSED SUCCESSFULLY! ===")
	
	// Cleanup all generated folders
	os.RemoveAll(testRoot)
	for i := 1; i <= 20; i++ {
		os.RemoveAll(fmt.Sprintf("_outdiff_%02d", i))
	}
}

func verifyFile(path, expectedContent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Verification failed: Could not read %s: %v\n", path, err)
		os.Exit(1)
	}
	if string(content) != expectedContent {
		fmt.Printf("Verification failed: Content mismatch in %s\n  Expected: %s\n  Got:      %s\n", path, expectedContent, string(content))
		os.Exit(1)
	}
	fmt.Printf("Verified: %s contains correct content.\n", path)
}
