# Guide: End-to-End (E2E) Testing for Interactive CLI Programs in Go

This document explains how to build a "Tester Agent" to automate the testing of CLI applications that require user input (prompts, menus, and confirmations).

## 1. The Multi-Level Agent Architecture

Our testing workflow uses a hierarchical approach:

*   **L1 (Orchestrator):** The AI/Developer who writes the test logic.
*   **L2 (Codebase Investigator):** Tools used to map the CLI prompts and menu structures.
*   **L3 (Real Agent):** A compiled Go binary (`tester_agent.exe`) that executes the target program in a real OS environment.
*   **L4 (Input Automator):** The reactive logic inside the L3 agent that "simulates" a human by reading `stdout` and writing to `stdin`.

---

## 2. Implementation Strategy

To test an interactive program, the Tester Agent must perform three main roles:

### A. Environment Setup
Before running the test, the agent must create a "Sandbox" (e.g., `test_space/`) with known files to ensure tests are deterministic and repeatable.

### B. Process Orchestration
Use the `os/exec` package to start the target program as a sub-process.

```go
cmd := exec.Command("./target_program.exe", "--args")
stdin, _ := cmd.StdinPipe()   // To send inputs (1, y, path)
stdout, _ := cmd.StdoutPipe() // To read prompts (Select option:)
cmd.Start()
```

### C. The Reactive Loop (The "Brain")
Since CLI programs often use `fmt.Print` (without newlines) for prompts, a standard `bufio.Scanner` might hang. Instead, use a byte-buffer loop to detect prompts in real-time:

```go
go func() {
    defer stdin.Close()
    buf := make([]byte, 4096)
    for {
        n, _ := stdout.Read(buf)
        output := string(buf[:n])
        
        // React to specific prompts
        if strings.Contains(output, "Select option:") {
            io.WriteString(stdin, "999
") // Select current folder
        }
        if strings.Contains(output, "Confirm? (y/n):") {
            io.WriteString(stdin, "y
")
        }
    }
}()
```

---

## 3. Best Practices for CLI E2E

1.  **Use Absolute Paths:** Always resolve paths using `filepath.Abs()` to avoid confusion between the tester's working directory and the target program's directory.
2.  **Clean Up:** Use `defer os.RemoveAll(testDir)` to ensure the system stays clean even if a test fails.
3.  **Timeout Handling:** In a CI/CD environment, always wrap the `cmd.Wait()` in a timer to prevent infinite hangs if a prompt is missed.
4.  **Verification:** After the process exits, physically check the disk using `os.Stat()` or `os.ReadFile()` to ensure the files were moved/created correctly.

## 4. Example Workflow: `deploy_navigator`

In this project, the `tester_agent_nav.go` successfully automated the following:
1.  **Step 1:** Created `orig/file1.txt` and `mod/file1.txt`.
2.  **Step 2:** Ran `compare` mode to generate a diff folder.
3.  **Step 3:** Started `deploy` mode, detected the "Source" menu, navigated the "Destination" menu using code `888`, created a new folder, and confirmed deployment with `999` and `y`.
4.  **Step 4:** Verified `deploy/file1.txt` contained the modified content.

---
*Created by Gemini CLI Agent for win-audit-2026 project.*
