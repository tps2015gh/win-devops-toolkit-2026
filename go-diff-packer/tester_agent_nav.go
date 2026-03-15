package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Tester Agent: Starting End-to-End Tests for Deploy Navigator ===")

	testRoot := "test_space_nav"
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

	fmt.Println("[STEP 1] Test environment prepared.")

	fmt.Println("[STEP 2] Running compare mode...")
	cmdCompare := exec.Command("./deploy_navigator.exe", origDir, modDir)
	cmdCompare.Run()

	fmt.Println("[STEP 3] Running deploy mode with navigation automation...")
	cmdDeploy := exec.Command("./deploy_navigator.exe", "--deploy")
	stdin, _ := cmdDeploy.StdinPipe()
	stdout, _ := cmdDeploy.StdoutPipe()
	cmdDeploy.Start()

	go func() {
		defer stdin.Close()
		buf := make([]byte, 4096)
		hasSentManual := false
		hasSentConfirmSelection := false

		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			chunk := string(buf[:n])
			fmt.Print(chunk)

			if strings.Contains(chunk, "Select source number") {
				// Send the last one (assuming it's the one we just made)
				io.WriteString(stdin, "12\n") 
			} else if strings.Contains(chunk, "Select option:") {
				if !hasSentManual {
					io.WriteString(stdin, "888\n")
					hasSentManual = true
				} else if !hasSentConfirmSelection {
					io.WriteString(stdin, "999\n")
					hasSentConfirmSelection = true
				}
			} else if strings.Contains(chunk, "Enter custom path:") {
				io.WriteString(stdin, deployDir+"\n")
			} else if strings.Contains(chunk, "Create it? (y/n):") {
				io.WriteString(stdin, "y\n")
			} else if strings.Contains(chunk, "Confirm deployment? (y/n):") {
				io.WriteString(stdin, "y\n")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	cmdDeploy.Wait()

	fmt.Println("\n[STEP 4] Verifying...")
	content, err := os.ReadFile(filepath.Join(deployDir, "file1.txt"))
	if err != nil || string(content) != "modified content" {
		fmt.Printf("Verification failed! Content: %s, Error: %v\n", string(content), err)
		os.Exit(1)
	}

	fmt.Println("\n=== DEPLOY NAVIGATOR TEST PASSED! ===")
}
