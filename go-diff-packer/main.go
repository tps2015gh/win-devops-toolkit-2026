package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		if *jsonMode {
			output := jsonOutput{Success: false, ErrorMessage: "Usage: go-diff-packer [--json] <original_dir> <modified_dir>"}
			printJSON(output)
		} else {
			fmt.Println("Usage: go run main.go [-v] [--json] <original_dir> <modified_dir>")
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
