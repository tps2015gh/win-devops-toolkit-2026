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

	if *deployMode {
		runDeployMode()
		return
	}

	if len(args) < 2 {
		if *jsonMode {
			output := jsonOutput{Success: false, ErrorMessage: "Usage: go-diff-picker [--json] <original_dir> <modified_dir>"}
			printJSON(output)
		} else {
			fmt.Println("Usage: go-diff-picker [-v] [--json] <original_dir> <modified_dir>")
			fmt.Println("       go-diff-picker --deploy")
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

		relPath, err := filepath.Rel(modDir, path)
		if err != nil {
			output.Stats.Errors++
			return nil
		}

		origPath := filepath.Join(origDir, relPath)
		shouldCopy := false
		change := changeRecord{File: relPath}

		origInfo, err := os.Stat(origPath)
		if os.IsNotExist(err) {
			shouldCopy = true
			change.Status = "NEW"
			change.Reason = "new file"
		} else if err == nil {
			if info.Size() != origInfo.Size() {
				shouldCopy = true
				change.Status = "MOD"
				change.Reason = "size changed"
				change.Size = fmt.Sprintf("%d -> %d", origInfo.Size(), info.Size())
			} else {
				modHash, modErr := getFileHash(path)
				origHash, origErr := getFileHash(origPath)
				if modErr != nil || origErr != nil {
					output.Stats.Errors++
					return nil
				}
				if modHash != origHash {
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
		}

		if shouldCopy {
			destPath := filepath.Join(outDir, relPath)
			err = copyFile(path, destPath)
			if err != nil {
				output.Stats.Errors++
			} else {
				output.Stats.Copied++
				output.Changes = append(output.Changes, change)
			}
		}
		return nil
	})

	if err != nil {
		os.Exit(1)
	}

	if !*jsonMode {
		fmt.Println("\n--- Summary ---")
		fmt.Printf("Files compared: %d\n", output.Stats.Compared)
		fmt.Printf("Files copied:   %d\n", output.Stats.Copied)
		absOutDir, _ := filepath.Abs(outDir)
		fmt.Printf("Output path:    %s\n", absOutDir)
		fmt.Println("Done.")
	}
}

func runDeployMode() {
	fmt.Println("=== Go Diff Picker - Enhanced Mode ===")
	reader := bufio.NewReader(os.Stdin)

	outDirs, err := listOutDiffDirs()
	if err != nil || len(outDirs) == 0 {
		fmt.Println("No output directories found.")
		return
	}

	var selectedDir string
	for {
		fmt.Println("\nAvailable source directories:")
		for i, dir := range outDirs {
			fmt.Printf("  %d. %s (%d files, %d folders)\n", i+1, dir, countFiles(dir), countFolders(dir))
		}
		fmt.Print("\nSelect source number (q to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "q" { return }
		num, _ := strconv.Atoi(input)
		if num >= 1 && num <= len(outDirs) {
			selectedDir = outDirs[num-1]
			break
		}
	}

	fmt.Printf("\nSelected Source: %s\n", selectedDir)
	destDir := navigateDestination(reader, ".")
	if destDir == "" { return }

	filesToDeploy := listFiles(selectedDir)
	filesToReplace := 0
	for _, file := range filesToDeploy {
		if _, err := os.Stat(filepath.Join(destDir, file)); err == nil {
			filesToReplace++
		}
	}

	fmt.Println("\n--- Final Confirmation ---")
	fmt.Printf("Source:      %s\n", selectedDir)
	fmt.Printf("Destination: %s\n", destDir)
	fmt.Printf("Files:       %d (%d to replace)\n", len(filesToDeploy), filesToReplace)
	fmt.Print("\nConfirm deployment? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" { return }

	fmt.Println("\n--- Deploying ---")
	for i, file := range filesToDeploy {
		fmt.Printf("\rProgress: [%d/%d] %s", i+1, len(filesToDeploy), strings.Repeat(" ", 20))
		copyFile(filepath.Join(selectedDir, file), filepath.Join(destDir, file))
	}
	fmt.Println("\n\nDone.")
	absPath, _ := filepath.Abs(destDir)
	fmt.Printf("Full Path: %s\n", absPath)
}

func navigateDestination(reader *bufio.Reader, startDir string) string {
	currentDir, _ := filepath.Abs(startDir)
	for {
		fmt.Printf("\nCurrent Destination: %s\n", currentDir)
		
		subDirs, _ := listSubDirs(currentDir)
		fmt.Println("Options:")
		fmt.Println("  0.  [..] Go up one level")
		for i, dir := range subDirs {
			fmt.Printf("  %d. %s\n", i+1, dir)
		}
		fmt.Println("  999. SELECT THIS FOLDER")
		fmt.Println("  888. Enter custom manual path")
		fmt.Println("  q.   Cancel")
		
		fmt.Print("\nSelect option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" { return "" }
		if input == "999" { return currentDir }
		
		if input == "888" {
			fmt.Print("Enter custom path: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("Path does not exist. Create it? (y/n): ")
				conf, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(conf)) == "y" {
					os.MkdirAll(path, 0755)
					currentDir, _ = filepath.Abs(path)
				}
			} else {
				currentDir, _ = filepath.Abs(path)
			}
			continue
		}

		num, _ := strconv.Atoi(input)
		if input == "0" {
			currentDir = filepath.Dir(currentDir)
		} else if num >= 1 && num <= len(subDirs) {
			currentDir = filepath.Join(currentDir, subDirs[num-1])
		} else {
			fmt.Println("Invalid selection.")
		}
	}
}

func listSubDirs(dir string) ([]string, error) {
	var dirs []string
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

func listOutDiffDirs() ([]string, error) {
	var dirs []string
	entries, _ := os.ReadDir(".")
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "_outdiff_") {
			dirs = append(dirs, entry.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

func listFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(dir, path)
			files = append(files, rel)
		}
		return nil
	})
	return files
}

func countFiles(dir string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() { count++ }
		return nil
	})
	return count
}

func countFolders(dir string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && path != dir { count++ }
		return nil
	})
	return count
}

func getFileHash(path string) (string, error) {
	f, _ := os.Open(path)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	os.MkdirAll(filepath.Dir(dst), 0755)
	source, _ := os.Open(src)
	defer source.Close()
	destination, _ := os.Create(dst)
	defer destination.Close()
	io.Copy(destination, source)
	return nil
}

func findNextOutDir() string {
	for i := 1; ; i++ {
		dir := fmt.Sprintf("_outdiff_%02d", i)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
			return dir
		}
	}
}

func printJSON(output jsonOutput) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(output)
}
