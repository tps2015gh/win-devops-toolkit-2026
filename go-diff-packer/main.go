package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type stats struct {
	Compared int `json:"compared"`
	Copied   int `json:"copied"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
}

type changeRecord struct {
	File   string `json:"file"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
	Size   string `json:"size,omitempty"`
}

type jsonOutput struct {
	OriginalDir  string         `json:"original_dir"`
	ModifiedDir  string         `json:"modified_dir"`
	OutputDir    string         `json:"output_dir"`
	Stats        stats          `json:"stats"`
	Changes      []changeRecord `json:"changes"`
	Success      bool           `json:"success"`
	ErrorMessage string         `json:"error_message,omitempty"`
}

func main() {
	verbose := flag.Bool("v", false, "verbose output: show unchanged files")
	jsonMode := flag.Bool("json", false, "output results as JSON (suppresses console output)")
	deployMode := flag.Bool("deploy", false, "deploy mode: copy files from _outdiff folder to destination")
	flag.Parse()

	args := flag.Args()

	// Deploy mode
	if *deployMode {
		runDeployMode()
		return
	}

	// Compare mode
	if len(args) < 2 {
		if *jsonMode {
			output := jsonOutput{Success: false, ErrorMessage: "Usage: go-diff-packer [--json] <original_dir> <modified_dir>"}
			printJSON(output)
		} else {
			fmt.Println("Usage: go-diff-packer [-v] [--json] <original_dir> <modified_dir>")
			fmt.Println("       go-diff-packer --deploy")
		}
		return
	}

	origDir := args[0]
	modDir := args[1]

	output := jsonOutput{
		OriginalDir: origDir,
		ModifiedDir: modDir,
		Changes:     []changeRecord{},
	}

	// Validate directories
	if _, err := os.Stat(origDir); os.IsNotExist(err) {
		output.ErrorMessage = fmt.Sprintf("Original directory does not exist: %s", origDir)
		if *jsonMode {
			printJSON(output)
		} else {
			fmt.Printf("Error: %s\n", output.ErrorMessage)
		}
		os.Exit(1)
	}
	if _, err := os.Stat(modDir); os.IsNotExist(err) {
		output.ErrorMessage = fmt.Sprintf("Modified directory does not exist: %s", modDir)
		if *jsonMode {
			printJSON(output)
		} else {
			fmt.Printf("Error: %s\n", output.ErrorMessage)
		}
		os.Exit(1)
	}

	outDir := findNextOutDir()
	output.OutputDir = outDir

	if !*jsonMode {
		fmt.Printf("Comparing %s and %s\n", origDir, modDir)
		fmt.Printf("Output directory: %s\n\n", outDir)
	}

	err := filepath.Walk(modDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		output.Stats.Compared++

		// Get relative path
		relPath, err := filepath.Rel(modDir, path)
		if err != nil {
			if !*jsonMode {
				fmt.Printf("[ERROR] %s: %v\n", relPath, err)
			}
			output.Stats.Errors++
			return nil
		}

		origPath := filepath.Join(origDir, relPath)
		shouldCopy := false
		change := changeRecord{File: relPath}

		// Check if file exists in original
		origInfo, err := os.Stat(origPath)
		if os.IsNotExist(err) {
			if !*jsonMode {
				fmt.Printf("[NEW] %s\n", relPath)
			}
			shouldCopy = true
			change.Status = "NEW"
			change.Reason = "new file"
		} else if err == nil {
			// File exists, compare size first (faster check)
			if info.Size() != origInfo.Size() {
				if !*jsonMode {
					fmt.Printf("[MOD] %s (size: %d -> %d)\n", relPath, origInfo.Size(), info.Size())
				}
				shouldCopy = true
				change.Status = "MOD"
				change.Reason = "size changed"
				change.Size = fmt.Sprintf("%d -> %d", origInfo.Size(), info.Size())
			} else {
				// Sizes match, compare hash
				modHash, modErr := getFileHash(path)
				origHash, origErr := getFileHash(origPath)

				if modErr != nil {
					if !*jsonMode {
						fmt.Printf("[ERROR] %s: hash calculation failed (modified): %v\n", relPath, modErr)
					}
					output.Stats.Errors++
					return nil
				}
				if origErr != nil {
					if !*jsonMode {
						fmt.Printf("[ERROR] %s: hash calculation failed (original): %v\n", relPath, origErr)
					}
					output.Stats.Errors++
					return nil
				}

				if modHash != origHash {
					if !*jsonMode {
						fmt.Printf("[MOD] %s (content changed)\n", relPath)
					}
					shouldCopy = true
					change.Status = "MOD"
					change.Reason = "content changed"
				} else {
					output.Stats.Skipped++
					if *verbose {
						fmt.Printf("[OK] %s (unchanged)\n", relPath)
					}
				}
			}
		} else {
			if !*jsonMode {
				fmt.Printf("[ERROR] %s: %v\n", relPath, err)
			}
			output.Stats.Errors++
			return nil
		}

		if shouldCopy {
			destPath := filepath.Join(outDir, relPath)
			err = copyFile(path, destPath)
			if err != nil {
				if !*jsonMode {
					fmt.Printf("[ERROR] %s: copy failed: %v\n", relPath, err)
				}
				output.Stats.Errors++
			} else {
				output.Stats.Copied++
				output.Changes = append(output.Changes, change)
			}
		} else {
			// Also track skipped files in JSON if verbose
			if change.Status == "" {
				change.Status = "OK"
				change.Reason = "unchanged"
			}
			if *verbose || *jsonMode {
				output.Changes = append(output.Changes, change)
			}
		}

		return nil
	})

	if err != nil {
		output.ErrorMessage = fmt.Sprintf("Error walking directory: %v", err)
		if *jsonMode {
			printJSON(output)
		} else {
			fmt.Printf("Error walking directory: %v\n", err)
		}
		os.Exit(1)
	}

	output.Success = true

	if *jsonMode {
		printJSON(output)
	} else {
		// Print summary
		fmt.Println("\n--- Summary ---")
		fmt.Printf("Files compared: %d\n", output.Stats.Compared)
		fmt.Printf("Files copied:   %d\n", output.Stats.Copied)
		fmt.Printf("Files skipped:  %d\n", output.Stats.Skipped)
		if output.Stats.Errors > 0 {
			fmt.Printf("Errors:         %d\n", output.Stats.Errors)
		}
		fmt.Println("Done.")
	}
}

func runDeployMode() {
	fmt.Println("=== Go Diff Packer - Deploy Mode ===\n")

	// List all _outdiff folders
	outDirs, err := listOutDiffDirs()
	if err != nil {
		fmt.Printf("Error listing output directories: %v\n", err)
		os.Exit(1)
	}

	if len(outDirs) == 0 {
		fmt.Println("No _outdiff folders found in current directory.")
		os.Exit(0)
	}

	// Display list
	fmt.Println("Available output directories:")
	fmt.Println("-----------------------------")
	for i, dir := range outDirs {
		fileCount := countFiles(dir)
		fmt.Printf("  %d. %s (%d files)\n", i+1, dir, fileCount)
	}
	fmt.Println()

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select output folder number (or 'q' to quit): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "q" {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(outDirs) {
		fmt.Println("Invalid selection.")
		os.Exit(1)
	}

	selectedDir := outDirs[num-1]
	fmt.Printf("\nSelected: %s\n", selectedDir)

	// Get destination folder
	fmt.Print("\nEnter destination folder path: ")
	destDir, _ := reader.ReadString('\n')
	destDir = strings.TrimSpace(destDir)

	if destDir == "" {
		fmt.Println("No destination folder specified.")
		os.Exit(1)
	}

	// Validate destination
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		fmt.Printf("Destination folder does not exist: %s\n", destDir)
		fmt.Print("Create it? (y/n): ")
		create, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(create)) != "y" {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Printf("Error creating destination folder: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created: %s\n", destDir)
	}

	// Count files to deploy
	filesToDeploy := listFiles(selectedDir)
	if len(filesToDeploy) == 0 {
		fmt.Println("\nNo files to deploy in selected folder.")
		os.Exit(0)
	}

	// Check for existing files and ask once for replace confirmation
	filesToReplace := []string{}
	filesToCopy := []string{}

	for _, file := range filesToDeploy {
		destPath := filepath.Join(destDir, file)
		if _, err := os.Stat(destPath); err == nil {
			filesToReplace = append(filesToReplace, file)
		} else {
			filesToCopy = append(filesToCopy, file)
		}
	}

	if len(filesToReplace) > 0 {
		fmt.Printf("\n⚠️  %d file(s) will be replaced:\n", len(filesToReplace))
		for _, file := range filesToReplace {
			fmt.Printf("   - %s\n", file)
		}
		fmt.Print("\nConfirm replace? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
	}

	// Deploy files
	fmt.Println("\n--- Deploying ---")
	copied := 0
	replaced := 0
	errors := 0

	for _, file := range filesToDeploy {
		srcPath := filepath.Join(selectedDir, file)
		destPath := filepath.Join(destDir, file)

		// Ensure destination directory exists
		err := os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			fmt.Printf("[ERROR] %s: could not create directory: %v\n", file, err)
			errors++
			continue
		}

		// Check if replacing
		isReplace := false
		if _, err := os.Stat(destPath); err == nil {
			isReplace = true
		}

		// Copy file
		err = copyFile(srcPath, destPath)
		if err != nil {
			fmt.Printf("[ERROR] %s: copy failed: %v\n", file, err)
			errors++
		} else {
			if isReplace {
				fmt.Printf("[REPLACE] %s\n", file)
				replaced++
			} else {
				fmt.Printf("[COPY] %s\n", file)
				copied++
			}
		}
	}

	// Summary
	fmt.Println("\n--- Deployment Summary ---")
	fmt.Printf("Files copied:   %d\n", copied)
	fmt.Printf("Files replaced: %d\n", replaced)
	if errors > 0 {
		fmt.Printf("Errors:         %d\n", errors)
	}
	fmt.Printf("Destination:    %s\n", destDir)
	fmt.Println("Done.")
}

func listOutDiffDirs() ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "_outdiff_") {
			dirs = append(dirs, entry.Name())
		}
	}

	// Sort by name (natural sort for _outdiff_01, _outdiff_02, etc.)
	sort.Strings(dirs)

	return dirs, nil
}

func listFiles(dir string) []string {
	var files []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(dir, path)
			files = append(files, relPath)
		}
		return nil
	})

	return files
}

func countFiles(dir string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func printJSON(output jsonOutput) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(output)
}

func findNextOutDir() string {
	i := 1
	for {
		dir := fmt.Sprintf("_outdiff_%02d", i)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				os.Exit(1)
			}
			return dir
		}
		i++
	}
}

func getFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	// Ensure destination directory exists
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
