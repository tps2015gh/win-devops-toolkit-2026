package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileHash(t *testing.T) {
	// Create a temporary file
	content := []byte("test content")
	tmpFile := "test_hash.txt"
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	hash1, err := getFileHash(tmpFile)
	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	hash2, err := getFileHash(tmpFile)
	if err != nil {
		t.Fatalf("Error getting hash second time: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hashes do not match: %s vs %s", hash1, hash2)
	}

	// Change content
	err = os.WriteFile(tmpFile, []byte("different content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	hash3, err := getFileHash(tmpFile)
	if err != nil {
		t.Fatalf("Error getting hash after change: %v", err)
	}

	if hash1 == hash3 {
		t.Errorf("Hashes should be different after content change")
	}
}

func TestCopyFile(t *testing.T) {
	src := "src_test.txt"
	dst := "dst_test.txt"
	content := []byte("hello world")

	err := os.WriteFile(src, content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(src)
	defer os.Remove(dst)

	err = copyFile(src, dst)
	if err != nil {
		t.Fatalf("Error copying file: %v", err)
	}

	dstContent, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("Error reading destination file: %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("Content mismatch: expected %s, got %s", string(content), string(dstContent))
	}
}

func TestFindNextOutDir(t *testing.T) {
	// This function uses hardcoded "_outdiff_" prefix. Let's test it in a temp dir if possible,
	// but it uses os.Stat(".") so it's tied to current dir.
	// We'll just verify it returns something starting with _outdiff_
	outDir := findNextOutDir()
	if outDir == "" {
		t.Error("Expected output directory name, got empty string")
	}
	if !filepath.HasPrefix(filepath.Base(outDir), "_outdiff_") {
		t.Errorf("Expected prefix _outdiff_, got %s", outDir)
	}
	// Note: this function creates the directory, so we should clean it up if it was created during test
	// But it's better to just check if it exists
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		t.Errorf("Output directory %s was not created", outDir)
	}
}

func TestCountingFunctions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_counts")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create structure:
	// tmpDir/file1.txt
	// tmpDir/subdir1/file2.txt
	// tmpDir/subdir1/subsubdir1/file3.txt
	// tmpDir/empty_dir

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("1"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "subdir1", "file2.txt"), []byte("2"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "subdir1", "subsubdir1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "subdir1", "subsubdir1", "file3.txt"), []byte("3"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "empty_dir"), 0755)

	fCount := countFiles(tmpDir)
	if fCount != 3 {
		t.Errorf("Expected 3 files, got %d", fCount)
	}

	dCount := countFolders(tmpDir)
	if dCount != 3 { // subdir1, subdir1/subsubdir1, empty_dir
		t.Errorf("Expected 3 folders, got %d", dCount)
	}

	files := listFiles(tmpDir)
	if len(files) != 3 {
		t.Errorf("Expected 3 files in list, got %d", len(files))
	}

	subDirs, _ := listSubDirs(tmpDir)
	if len(subDirs) != 2 { // subdir1, empty_dir (at first level)
		t.Errorf("Expected 2 first-level subdirs, got %d", len(subDirs))
	}
}
