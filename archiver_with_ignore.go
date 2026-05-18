package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// IgnorePattern represents a single ignore pattern
type IgnorePattern struct {
	Pattern string
	IsDir   bool
}

// IgnorePatterns holds all patterns from the ignore config file
type IgnorePatterns struct {
	Patterns []IgnorePattern
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Archiver with Ignore - Robust Unicode-aware Compression Tool with Ignore Support")
		fmt.Println("\nUsage:")
		fmt.Println("  archiver_with_ignore.exe <folder_to_compress> [config_file]")
		fmt.Println("\nArguments:")
		fmt.Println("  folder_to_compress    Path to the folder you want to compress")
		fmt.Println("  config_file          Path to ignore config file (default: archiver_with_ignore_config1.txt)")
		fmt.Println("\nFlags:")
		fmt.Println("  -h, --help            Show this help message")
		fmt.Println("  -init                 Create a default ignore config file")
		fmt.Println("\nFeatures:")
		fmt.Println("  - Full Unicode support (Thai, Chinese, etc.)")
		fmt.Println("  - Ignore patterns support (like .gitignore)")
		fmt.Println("  - Wildcard pattern matching (*)")
		fmt.Println("  - Directory-specific exclusion (append /)")
		fmt.Println("  - Comment support (lines starting with #)")
		fmt.Println("  - Automatic timestamping to prevent overwrites")
		fmt.Println("\nIgnore Config Format:")
		fmt.Println("  - One pattern per line")
		fmt.Println("  - Empty lines and comments (#) are ignored")
		fmt.Println("  - Append / to match only directories (e.g., node_modules/)")
		fmt.Println("  - Use * for wildcard patterns (e.g., *.log, *.tmp)")
		fmt.Println("\nExample ignore config (archiver_with_ignore_config1.txt):")
		fmt.Println("  .git/")
		fmt.Println("  node_modules/")
		fmt.Println("  __pycache__/")
		fmt.Println("  *.log")
		fmt.Println("  .env")
		fmt.Println("  .vscode/")
		os.Exit(0)
	}

	// Check for init flag
	if len(os.Args) >= 2 && os.Args[1] == "-init" {
		configFile := "archiver_with_ignore_config1.txt"
		if len(os.Args) >= 3 {
			configFile = os.Args[2]
		}
		err := createDefaultConfigFile(configFile)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created default config file: %s\n", configFile)
		os.Exit(0)
	}

	targetDir := os.Args[1]
	configFile := "archiver_with_ignore_config1.txt"

	if len(os.Args) >= 3 {
		configFile = os.Args[2]
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("Error: %s is not a directory.\n", targetDir)
		os.Exit(1)
	}

	// Load ignore patterns
	patterns, err := loadIgnorePatterns(configFile)
	if err != nil {
		fmt.Printf("Warning: Could not load ignore config: %v\n", err)
		fmt.Println("Proceeding without ignore patterns...")
		patterns = &IgnorePatterns{}
	} else {
		fmt.Printf("Loaded ignore config: %s\n", configFile)
		fmt.Printf("Loaded %d ignore patterns\n", len(patterns.Patterns))
	}

	folderName := filepath.Base(absPath)
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("%s_%s.zip", folderName, timestamp)

	fmt.Printf("\nCompressing folder: %s\n", absPath)
	fmt.Printf("Output file: %s\n", zipFileName)
	fmt.Println("-------------------------------------------")

	err = zipFolderWithIgnore(absPath, zipFileName, patterns)
	if err != nil {
		fmt.Printf("\nCompression failed: %v\n", err)
		os.Exit(1)
	}

	finalPath, _ := filepath.Abs(zipFileName)
	fmt.Printf("\n-------------------------------------------\n")
	fmt.Printf("Successfully compressed to:\n%s\n", finalPath)
}

// LoadIgnorePatterns reads and parses the ignore config file
func loadIgnorePatterns(configFile string) (*IgnorePatterns, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open ignore config file: %w", err)
	}
	defer file.Close()

	patterns := &IgnorePatterns{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if pattern is for directories (ends with /)
		isDir := strings.HasSuffix(line, "/")
		if isDir {
			line = strings.TrimSuffix(line, "/")
		}

		patterns.Patterns = append(patterns.Patterns, IgnorePattern{
			Pattern: line,
			IsDir:   isDir,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ignore config file: %w", err)
	}

	return patterns, nil
}

// ShouldIgnore checks if a path should be ignored
func (ip *IgnorePatterns) ShouldIgnore(path string, isDir bool) bool {
	// Convert to forward slashes for consistent comparison
	path = filepath.ToSlash(path)

	for _, pattern := range ip.Patterns {
		if matchesPattern(path, pattern.Pattern, isDir || pattern.IsDir) {
			return true
		}
	}
	return false
}

// matchesPattern performs glob-style pattern matching (similar to .gitignore)
func matchesPattern(path, pattern string, isDir bool) bool {
	path = filepath.ToSlash(path)
	pattern = filepath.ToSlash(pattern)

	// Handle exact matches
	if path == pattern {
		return true
	}

	// Handle directory patterns - match if path starts with the directory name
	if isDir && strings.HasPrefix(path, pattern+"/") {
		return true
	}

	// Handle wildcard patterns (e.g., *.log, *.tmp)
	if strings.Contains(pattern, "*") {
		// Simple wildcard matching
		if matchWildcard(path, pattern) {
			return true
		}

		// Also check if any path component matches (for patterns like *.log)
		pathParts := strings.Split(path, "/")
		for _, part := range pathParts {
			if matchWildcard(part, pattern) {
				return true
			}
		}
	}

	// Check if path contains the pattern as a directory component
	pathParts := strings.Split(path, "/")
	for _, part := range pathParts {
		if part == pattern && isDir {
			return true
		}
	}

	return false
}

// matchWildcard performs simple wildcard matching
func matchWildcard(text, pattern string) bool {
	// Handle * wildcard
	if !strings.Contains(pattern, "*") {
		return text == pattern
	}

	parts := strings.Split(pattern, "*")
	if len(parts) == 0 {
		return false
	}

	// Check if text starts with the part before the first *
	if parts[0] != "" && !strings.HasPrefix(text, parts[0]) {
		return false
	}

	// Check if text ends with the part after the last *
	if parts[len(parts)-1] != "" && !strings.HasSuffix(text, parts[len(parts)-1]) {
		return false
	}

	// For patterns like *.ext, just check if it ends with .ext
	if strings.HasPrefix(pattern, "*") && len(parts) == 2 {
		return strings.HasSuffix(text, parts[1])
	}

	return true
}

// CreateDefaultConfigFile creates a sample ignore config file
func createDefaultConfigFile(configFile string) error {
	content := `# Archiver Ignore Configuration
# List folders and patterns to exclude from archiving
# Use one pattern per line
# Patterns support wildcards (*)
# Append / to exclude only directories

# Common directories to ignore
.git/
.gitignore
node_modules/
__pycache__/
.env
.venv/
venv/
env/

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~

# Build artifacts
bin/
obj/
dist/
build/

# OS files
.DS_Store
Thumbs.db
desktop.ini

# Package manager files
package-lock.json
yarn.lock

# Temporary files
*.tmp
*.temp
*.log
.tmp/
temp/

# Add your own patterns below:
`

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// zipFolderWithIgnore zips a folder while respecting ignore patterns
func zipFolderWithIgnore(source, target string, patterns *IgnorePatterns) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	var filesAdded int
	var filesSkipped int

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Check if this path should be ignored
		if patterns.ShouldIgnore(relPath, info.IsDir()) {
			fmt.Printf("\rSkipping: %-60s", truncateString(relPath, 60))
			filesSkipped++
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath = filepath.ToSlash(relPath)

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

		if info.IsDir() {
			if !strings.HasSuffix(header.Name, "/") {
				header.Name += "/"
			}
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			filesAdded++
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Printf("\rAdding: %-60s", truncateString(relPath, 60))
		filesAdded++
		_, err = io.Copy(writer, file)
		return err
	})

	fmt.Printf("\n\nStatistics:\n")
	fmt.Printf("  Files/Folders Added: %d\n", filesAdded)
	fmt.Printf("  Files/Folders Skipped: %d\n", filesSkipped)

	return err
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 10 {
		return s
	}
	return "..." + s[len(s)-(maxLen-3):]
}
