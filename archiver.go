package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: archiver.exe <folder_to_compress>")
		fmt.Println("Example: archiver.exe my_project")
		os.Exit(1)
	}

	targetDir := os.Args[1]
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

	folderName := filepath.Base(absPath)
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("%s_%s.zip", folderName, timestamp)

	fmt.Printf("Compressing folder: %s\n", absPath)
	fmt.Printf("Output file: %s\n", zipFileName)

	err = zipFolder(absPath, zipFileName)
	if err != nil {
		fmt.Printf("\nCompression failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully compressed to %s\n", zipFileName)
}

func zipFolder(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

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
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Printf("\rAdding: %-60s", truncateString(relPath, 60))
		_, err = io.Copy(writer, file)
		return err
	})

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
