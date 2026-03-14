# Go Diff Packer - Usage Guide

## Quick Start

### Basic Usage

```powershell
# Navigate to the go-diff-packer directory
cd go-diff-packer

# Run with Go (development mode)
go run main.go <original_directory> <modified_directory>

# Or build and run the executable
go build -o go-diff-packer.exe
.\go-diff-packer.exe <original_directory> <modified_directory>
```

### Examples

```powershell
# Compare two directories
go run main.go C:\xampp\htdocs\app_v1 C:\xampp\htdocs\app_v2

# Verbose mode (shows unchanged files)
go run main.go -v C:\xampp\htdocs\app_v1 C:\xampp\htdocs\app_v2

# JSON mode (for AI/automation)
go run main.go -json C:\xampp\htdocs\app_v1 C:\xampp\htdocs\app_v2 > diff.json

# Using the executable
.\go-diff-packer.exe C:\projects\original C:\projects\modified

# JSON output to file for AI processing
.\go-diff-packer.exe -json C:\v1 C:\v2 > changes.json

# Deploy mode (interactive)
.\go-diff-packer.exe --deploy
```

### Deploy Mode

Deploy mode allows you to interactively select an `_outdiff_*` folder and copy files to a destination:

```powershell
# Start deploy mode
.\go-diff-packer.exe --deploy

# Interactive prompts:
# 1. Select output folder number from list
# 2. Enter destination folder path
# 3. Confirm replacement if files exist (asked once)
# 4. Files are copied with progress display
```

#### Deploy Mode Example

```
=== Go Diff Packer - Deploy Mode ===

Available output directories:
-----------------------------
  1. _outdiff_01 (2 files)
  2. _outdiff_02 (5 files)

Select output folder number (or 'q' to quit): 1

Selected: _outdiff_01

Enter destination folder path: C:\xampp\htdocs\live

⚠️  2 file(s) will be replaced:
   - config.php
   - index.php

Confirm replace? (y/n): y

--- Deploying ---
[REPLACE] config.php
[COPY] new_feature.php

--- Deployment Summary ---
Files copied:   1
Files replaced: 2
Destination:    C:\xampp\htdocs\live
Done.
```

### JSON Output Example

```json
{
  "original_dir": "C:\\v1",
  "modified_dir": "C:\\v2",
  "output_dir": "_outdiff_01",
  "stats": {
    "compared": 15,
    "copied": 3,
    "skipped": 12,
    "errors": 0
  },
  "changes": [
    {
      "file": "config.php",
      "status": "MOD",
      "reason": "size changed",
      "size": "1024 -> 2048"
    },
    {
      "file": "new_feature.php",
      "status": "NEW",
      "reason": "new file"
    },
    {
      "file": "unchanged.txt",
      "status": "OK",
      "reason": "unchanged"
    }
  ],
  "success": true
}
```

### JSON Output Fields

| Field | Type | Description |
|-------|------|-------------|
| `original_dir` | string | Path to original directory |
| `modified_dir` | string | Path to modified directory |
| `output_dir` | string | Auto-created output folder |
| `stats.compared` | integer | Total files compared |
| `stats.copied` | integer | Files copied to output |
| `stats.skipped` | integer | Unchanged files |
| `stats.errors` | integer | Error count |
| `changes` | array | List of all file changes |
| `changes[].file` | string | Relative file path |
| `changes[].status` | string | `NEW`, `MOD`, or `OK` |
| `changes[].reason` | string | Why file was flagged |
| `changes[].size` | string | Size change (if applicable) |
| `success` | boolean | Operation completed successfully |
| `error_message` | string | Error details (if failed) |

---

## Output Format

The tool creates auto-incrementing output folders (`_outdiff_01`, `_outdiff_02`, ...) and copies only changed/new files.

### Status Markers

| Marker | Description |
|--------|-------------|
| `[NEW]` | File exists in modified but not in original |
| `[MOD]` | File has been modified (size or content changed) |
| `[OK]` | File is unchanged (only shown in verbose mode) |
| `[ERROR]` | An error occurred during processing |

### Example Output

```
Comparing C:\app\v1 and C:\app\v2
Output directory: _outdiff_01

[MOD] config.php (size: 1024 -> 2048)
[NEW] new_feature.php
[MOD] index.php (content changed)

--- Summary ---
Files compared: 15
Files copied:   3
Files skipped:  12
Done.
```

---

## Command Line Options

| Option | Description |
|--------|-------------|
| `-v` | Verbose mode - shows unchanged files with `[OK]` marker |
| `-json` | JSON output mode - outputs structured JSON instead of console text (for AI/automation) |
| `--deploy` | Deploy mode - interactive deployment from _outdiff folder to destination |
| `<original_dir>` | Path to the original/reference directory (required for compare mode) |
| `<modified_dir>` | Path to the modified directory to compare (required for compare mode) |

---

## Use Cases

### 1. Deploy Only Changed Files

Compare production vs development, deploy only changes:

```powershell
# Compare live site with development version
.\go-diff-packer.exe C:\xampp\htdocs\live C:\xampp\htdocs\dev

# Copy only changed files from _outdiff_01 to production
xcopy /E /I /Y _outdiff_01\* C:\xampp\htdocs\live\
```

### 2. Backup Before Updates

```powershell
# Compare current installation with new version
.\go-diff-packer.exe C:\xampp\htdocs\current C:\xampp\htdocs\update_package

# Review changes in _outdiff_01 before applying
```

### 3. Code Review Assistance

```powershell
# See what changed between branches or versions
go run main.go -v C:\project\branch_v1 C:\project\branch_v2
```

### 4. XAMPP Migration

```powershell
# Compare old XAMPP htdocs with new setup
.\go-diff-packer.exe D:\xampp_old\htdocs C:\xampp_new\htdocs

# Review and copy only necessary files
```

---

## How It Works

1. **Size Comparison First**: Quickly compares file sizes (fast check)
2. **SHA-256 Hash**: If sizes match, compares content hashes (accurate check)
3. **Smart Copy**: Only copies files that are new or modified
4. **Auto Output**: Creates `_outdiff_01`, `_outdiff_02`, etc. automatically

---

## Performance Tips

- **Large Directories**: The tool skips hash calculation if file sizes differ (optimization)
- **Verbose Mode**: Use `-v` only when needed (slower due to extra output)
- **Clean Output**: Delete old `_outdiff_*` folders to keep workspace clean

---

## Troubleshooting

### "original directory does not exist"
Ensure the first argument is the original/reference directory path.

### "modified directory does not exist"
Ensure the second argument is the modified directory path.

### "hash calculation failed"
Check file permissions - the file may be locked or inaccessible.

### "copy failed"
Ensure destination directory is writable and not locked by another process.

---

## File Structure

```
go-diff-packer/
├── main.go              # Source code
├── go.mod               # Go module definition
├── go-diff-packer.exe   # Compiled executable (after build)
└── _outdiff_01/         # Output folder (auto-created)
    └── ...              # Changed/new files copied here
```

---

## Integration with Build Process

Add to your `build.bat`:

```batch
@echo off
cd go-diff-packer
go build -o go-diff-packer.exe
echo Go Diff Packer built successfully!
```

---

*Part of abcd_setup_server project - Developed by tps2015gh with Qwen Code Assistance*
