package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("--- STARTING CI3 SEARCHER TEST ---")
	
	// Create mock project
	os.MkdirAll("tests/mock_ci3/application/models", 0755)
	os.WriteFile("tests/mock_ci3/application/models/User_model.php", []byte("class User_model extends CI_Model { }"), 0644)

	// Simulate "user model" search followed by "q"
	input := "user model\nq\n"
	cmd := exec.Command("go", "run", "ci3_searcher.go", "tests/mock_ci3")
	cmd.Stdin = strings.NewReader(input)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[ERROR] Command failed: %v\n", err)
		os.Exit(1)
	}

	output := out.String()
	fmt.Println(output)

	// Verify
	fmt.Println("\n--- VERIFYING RESULTS ---")
	if strings.Contains(output, "model") && strings.Contains(output, "User_model.php") {
		fmt.Println("[PASS] Search results found User_model.php")
	} else {
		fmt.Println("[FAIL] Search results did not contain User_model.php")
		os.Exit(1)
	}
	fmt.Println("--- TEST SUCCESSFUL ---")
}
