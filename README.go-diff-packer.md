# Go Diff Packer

> **Directory Comparison and Differential Backup Tool**

A lightweight Go utility that compares two directories and copies only new or modified files to an auto-incrementing output folder. Perfect for deployments, backups, and code reviews.

---

## 🌐 Languages
[**English**](./README.go-diff-packer.md) | [**ภาษาไทย**](./README.go-diff-packer.th.md)

---

## 🚀 Features

| Feature | Description |
|---------|-------------|
| **Smart Comparison** | Compares file sizes first, then SHA-256 hashes for accuracy |
| **Auto Output Folders** | Creates `_outdiff_01`, `_outdiff_02`, etc. automatically |
| **Detailed Reporting** | Shows `[NEW]`, `[MOD]`, `[OK]`, and `[ERROR]` status markers |
| **Summary Statistics** | Displays files compared, copied, skipped, and errors |
| **Verbose Mode** | Optional `-v` flag to show unchanged files |
| **JSON Output** | `-json` flag for AI/automation integration |
| **Deploy Mode** | `--deploy` flag for interactive deployment to destination |
| **Replace Confirmation** | Ask once before replacing existing files |
| **Error Handling** | Graceful handling of permission and I/O errors |
| **Cross-Platform** | Built with Go 1.21.1, runs on Windows/Linux/macOS |

---

## 📦 Installation

### Option 1: Use Pre-built Executable (Windows)

```powershell
# Navigate to go-diff-packer directory
cd go-diff-packer

# Run the executable
.\go-diff-packer.exe <original_dir> <modified_dir>
```

### Option 2: Build from Source

```powershell
# Ensure Go 1.21.1+ is installed
go version

# Navigate to source directory
cd go-diff-packer

# Build the executable
go build -o go-diff-packer.exe

# Run
.\go-diff-packer.exe C:\original C:\modified
```

---

## 💡 Usage

### Basic Command

```powershell
# Standard comparison
go run main.go <original_directory> <modified_directory>

# With verbose output
go run main.go -v <original_directory> <modified_directory>
```

### Examples

```powershell
# Compare XAMPP htdocs directories
.\go-diff-packer.exe D:\xampp_old\htdocs C:\xampp_new\htdocs

# Compare project versions with verbose output
go run main.go -v C:\projects\app_v1 C:\projects\app_v2

# JSON output for AI/automation
.\go-diff-packer.exe -json C:\v1 C:\v2 > changes.json

# Deploy changed files
.\go-diff-packer.exe C:\live C:\staging
xcopy /E /I /Y _outdiff_01\* C:\live\
```

---

## 📊 Output Explanation

### Status Markers

| Marker | Meaning | Action |
|--------|---------|--------|
| `[NEW]` | New file (doesn't exist in original) | Copied to output |
| `[MOD]` | Modified file (size or content changed) | Copied to output |
| `[OK]` | Unchanged file (verbose mode only) | Skipped |
| `[ERROR]` | Processing error | Logged, continues |

### Sample Output

```
Comparing C:\app\v1 and C:\app\v2
Output directory: _outdiff_01

[MOD] config.php (size: 1024 -> 2048)
[NEW] features/user_auth.php
[MOD] index.php (content changed)

--- Summary ---
Files compared: 25
Files copied:   3
Files skipped:  22
Done.
```

---

## 🔧 Technical Specifications

| Specification | Value |
|---------------|-------|
| **Language** | Go (Golang) |
| **Minimum Version** | Go 1.21.1 |
| **Hash Algorithm** | SHA-256 |
| **Output Pattern** | `_outdiff_XX` (auto-increment) |
| **Platform** | Windows (tested), Linux, macOS |
| **Dependencies** | None (standard library only) |

---

## 🏗️ Architecture

### Comparison Algorithm

```
1. Walk through modified directory
2. For each file:
   a. Check if exists in original
      - NO → Mark as [NEW], copy to output
      - YES → Continue to step b
   b. Compare file sizes
      - Different → Mark as [MOD], copy to output
      - Same → Continue to step c
   c. Calculate and compare SHA-256 hashes
      - Different → Mark as [MOD], copy to output
      - Same → Mark as [OK], skip
3. Generate summary statistics
```

### Performance Optimizations

- **Size-First Comparison**: Fast rejection of different files (no hash calculation needed)
- **Early Exit**: Skips hash computation when size differs
- **Streaming Hash**: Uses `io.Copy` for memory-efficient hash calculation

---

## 👨‍💻 Author & Development Team

| Role | Name |
|------|------|
| **Lead Developer / Owner** | **tps2015gh** (Human) |
| **Programming & Design** | Qwen Code (under supervision of tps2015gh) |
| **Testing** | tps2015gh, Qwen Code |
| **Intelligent Assistant** | Qwen Code (CLI Agent) |

### Development History

- **Original Concept**: Directory comparison for XAMPP migration
- **First Implementation**: Basic file comparison with hash verification
- **Enhanced Version**: Added statistics, verbose mode, error handling
- **Integration**: Added to abcd_setup_server project suite

**Design & Implementation Statement:** This program was designed and coded by **Qwen Code** under the supervision, direction, and instruction of **tps2015gh**. All feature requirements, architectural decisions, code reviews, and final approvals were directed by the human project owner. Qwen Code served as an intelligent programming assistant, implementing the vision and specifications provided by tps2015gh.

**Legal Note on Authorship:** This project is authored and owned exclusively by **tps2015gh**. Qwen Code acted as a sophisticated development tool and intelligent assistant, providing implementation support, code fixes, and technical suggestions under the direct instruction and oversight of the project owner. All intellectual property, copyright, and final architectural decisions reside with the human author.

---

## 📝 Change Log

### v1.1.0 (Current)
- ✅ Fixed: Skip hash calculation when file sizes differ (performance)
- ✅ Added: Error handling for hash calculation failures
- ✅ Added: Summary statistics (compared/copied/skipped/errors)
- ✅ Added: Verbose mode (`-v` flag) to show unchanged files
- ✅ Added: Directory validation with clear error messages
- ✅ Added: Detailed size information in output (old → new)

### v1.0.0 (Initial)
- Basic directory comparison
- SHA-256 hash verification
- Auto-incrementing output folders
- File copy functionality

---

## 🔗 Related Documentation

- [**Usage Guide**](./README.go-diff-packer.usage.md) - Detailed usage examples and tips
- [**Main README**](./README.md) - Back to main project documentation

---

## 📄 License

This project is licensed under the **MIT License**. The copyright holder and owner is **tps2015gh**.

---

## 🆘 Support

For issues or questions:
1. Check the [**Usage Guide**](./README.go-diff-packer.usage.md)
2. Review error messages in console output
3. Ensure proper file permissions on directories

---

*Part of abcd_setup_server project - Developed by tps2015gh with Qwen Code Assistance - 2024*
