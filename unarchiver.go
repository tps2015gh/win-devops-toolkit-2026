package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	zips, _ := filepath.Glob("*.zip")

	if len(os.Args) < 2 {
		fmt.Println("Unarchiver - Robust Unicode-aware Extraction Tool")
		fmt.Println("
Available ZIP files in this folder:")
		if len(zips) == 0 {
			fmt.Println("  (No .zip files found)")
		} else {
			for _, z := range zips {
				fmt.Printf("  - %s
", z)
			}
		}

		fmt.Println("
Usage:")
		fmt.Println("  unarchiver.exe <zip_file_name>")
		fmt.Println("
Example:")
		if len(zips) > 0 {
			fmt.Printf("  unarchiver.exe %s
", zips[0])
		} else {
			fmt.Println("  unarchiver.exe my_backup.zip")
		}
		os.Exit(0)
	}

	zipFile := os.Args[1]
	destFolder := "_unzip"

	fmt.Printf("Extracting: %s
", zipFile)
	fmt.Printf("Destination: %s
", destFolder)

	err := unzip(zipFile, destFolder)
	if err != nil {
		fmt.Printf("
Extraction failed: %v
", err)
		os.Exit(1)
	}

	finalPath, _ := filepath.Abs(destFolder)
	fmt.Printf("
Successfully extracted to:
%s
", finalPath)
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Securely join paths and ensure no directory traversal
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		fmt.Printf("Extracting: %-60s", truncateString(f.Name, 60))

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
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
