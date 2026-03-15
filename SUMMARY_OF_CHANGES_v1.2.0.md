# Summary of Changes - v1.2.0 (March 15, 2026)

This release introduces the **Deploy Navigator** tool and a comprehensive **E2E CLI Testing** framework to the `win-audit-2026` suite.

---

## 🚀 New Features & Enhancements

### 1. Deploy Navigator Tool (`deploy_navigator.exe`) [NEW!]
- **Interactive Navigation**: A new program that allows you to browse your filesystem using numbered menus.
- **Recursive Selection**: 
    - `0`: Go up one level (`..`).
    - `1, 2, 3...`: Enter subfolders.
    - `999`: Select the current folder as the target.
- **Manual Path Entry**: Option `888` allows you to type a custom path directly.
- **Auto-Create Directory**: Automatically prompts to create the destination path if it doesn't exist.
- **Confirmation Screen**: Shows source, destination, and the number of files to be replaced before proceeding.

### 2. E2E CLI Testing Guide (`E2E_CLI_TESTING.md`) [NEW!]
- **Architecture Documentation**: Explains the L1-L5 agent hierarchy used for automated testing.
- **Implementation Guide**: Technical details on using Go's `os/exec` and `StdinPipe/StdoutPipe` for reactive CLI testing.
- **Automated Agents**: Documentation for `tester_agent.go` and `tester_agent_nav.go`.

### 3. Improvements to Go Diff Packer (`go-diff-packer.exe`)
- **Full Path Visibility**: Now displays the absolute path (Full Path) of the output/destination folder at the end of every operation.
- **Folder Counting**: The deployment menu now shows both the number of **files** and **folders** (recursive) for each `_outdiff` directory.
- **Progressive UI**: Changed the deployment status from multiple lines to a single-line progressive update (``) to reduce console clutter.

### 4. Improvements to Unarchiver (`unarchiver.exe`)
- **Destination Transparency**: Added the absolute destination path to the final summary after successful extraction.

---

## 🛠️ Verification & Quality Assurance
- **Tester Agents**: Successfully ran `tester_agent_nav.exe` to perform automated end-to-end verification of the new navigation logic.
- **Status**: All tests passed for file integrity, recursive navigation, and automated input handling.

---

## 🔗 Related Links
- [**Main README**](./README.md)
- [**Go Diff Packer README**](./README.go-diff-packer.md)
- [**E2E CLI Testing Guide**](./go-diff-packer/E2E_CLI_TESTING.md)

---
*Developed by tps2015gh with Gemini CLI Assistance - 2026*
