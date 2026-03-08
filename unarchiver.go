package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	zips, _ := filepath.Glob("*.zip")

	var selectedZip string

	// Check if a zip file was passed as an argument
	if len(os.Args) >= 2 {
		selectedZip = os.Args[1]
	} else {
		// Interactive Mode
		fmt.Println("Unarchiver - Robust Unicode-aware Extraction Tool")
		fmt.Println("\nAvailable ZIP files in this folder:")
		
		if len(zips) == 0 {
			fmt.Println("  (No .zip files found in current directory)")
			fmt.Println("\nPress Enter to exit...")
			bufio.NewReader(os.Stdin).ReadString('\n')
			os.Exit(0)
		}

		for i, z := range zips {
			fmt.Printf("  [%d] %s\n", i+1, z)
		}

		fmt.Printf("\nSelect a file number to extract (1-%d): ", len(zips))
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(zips) {
			fmt.Println("Invalid selection. Exiting.")
			os.Exit(1)
		}

		selectedZip = zips[choice-1]
	}

	destFolder := "_unzip"

	fmt.Printf("\nExtracting: %s\n", selectedZip)
	fmt.Printf("Destination: %s\n", destFolder)

	err := unzip(selectedZip, destFolder)
	if err != nil {
		fmt.Printf("\nExtraction failed: %v\n", err)
		os.Exit(1)
	}

	finalPath, _ := filepath.Abs(destFolder)
	fmt.Printf("\nSuccessfully extracted to:\n%s\n", finalPath)
	
	fmt.Println("\nPress Enter to close...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		
		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) && fpath != filepath.Clean(dest) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		fmt.Printf("\rExtracting: %-60s", truncateString(f.Name, 60))

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
