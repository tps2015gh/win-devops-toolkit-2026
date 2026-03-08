package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	mysqlPath     string
	mysqldumpPath string
)

func main() {
	fmt.Println("MariaDB Database Manager (Backup/Restore)")
	fmt.Println("=========================================")

	// 1. Setup Binaries
	setupBinaries()

	// 2. Get Credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter credentials (format: user/pass): ")
	creds, _ := reader.ReadString('
')
	creds = strings.TrimSpace(creds)

	parts := strings.SplitN(creds, "/", 2)
	user := "root"
	pass := ""
	if len(parts) == 2 {
		user = parts[0]
		pass = parts[1]
	} else if len(parts) == 1 && parts[0] != "" {
		user = parts[0]
	}

	for {
		fmt.Println("
Choose action:")
		fmt.Println("1. List Databases & Backup")
		fmt.Println("2. Restore from .sql.zip")
		fmt.Println("3. Exit")
		fmt.Print("Select (1-3): ")

		action, _ := reader.ReadString('
')
		action = strings.TrimSpace(action)

		switch action {
		case "1":
			handleBackup(user, pass)
		case "2":
			handleRestore(user, pass)
		case "3":
			return
		default:
			fmt.Println("Invalid selection.")
		}
	}
}

func setupBinaries() {
	// Try standard PATH first
	mysqlPath, _ = exec.LookPath("mysql.exe")
	mysqldumpPath, _ = exec.LookPath("mysqldump.exe")

	// If not in PATH, try common XAMPP paths
	if mysqlPath == "" || mysqldumpPath == "" {
		commonPaths := []string{
			"C:\xampp\mysql\bin",
			"D:\xampp\mysql\bin",
			"C:\xampp_v8_1_25\mysql\bin",
		}
		for _, p := range commonPaths {
			m := filepath.Join(p, "mysql.exe")
			md := filepath.Join(p, "mysqldump.exe")
			if _, err := os.Stat(m); err == nil && mysqlPath == "" {
				mysqlPath = m
			}
			if _, err := os.Stat(md); err == nil && mysqldumpPath == "" {
				mysqldumpPath = md
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	if mysqlPath == "" {
		fmt.Print("mysql.exe not found. Enter full path: ")
		p, _ := reader.ReadString('
')
		mysqlPath = strings.TrimSpace(p)
	}
	if mysqldumpPath == "" {
		fmt.Print("mysqldump.exe not found. Enter full path: ")
		p, _ := reader.ReadString('
')
		mysqldumpPath = strings.TrimSpace(p)
	}
	
	fmt.Printf("Using mysql:     %s
", mysqlPath)
	fmt.Printf("Using mysqldump: %s
", mysqldumpPath)
}

func handleBackup(user, pass string) {
	dbs := listDatabases(user, pass)
	if len(dbs) == 0 {
		fmt.Println("No databases found or connection failed.")
		return
	}

	fmt.Println("
Available Databases:")
	for i, db := range dbs {
		fmt.Printf("[%d] %s
", i+1, db)
	}

	fmt.Printf("
Enter number to backup (1-%d): ", len(dbs))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('
')
	idx := strings.TrimSpace(input)

	i, err := strconv.Atoi(idx)
	if err != nil || i < 1 || i > len(dbs) {
		fmt.Println("Invalid selection.")
		return
	}

	dbName := dbs[i-1]
	backupDB(dbName, user, pass)
}

func listDatabases(user, pass string) []string {
	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, "-e", "SHOW DATABASES;")

	cmd := exec.Command(mysqlPath, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error listing databases: %v
", err)
		return nil
	}

	lines := strings.Split(string(out), "
")
	var dbs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Database" || line == "information_schema" || line == "performance_schema" || line == "mysql" || line == "phpmyadmin" {
			continue
		}
		dbs = append(dbs, line)
	}
	return dbs
}

func backupDB(dbName, user, pass string) {
	timestamp := time.Now().Format("20060102_150405")
	sqlFile := fmt.Sprintf("%s_%s.sql", dbName, timestamp)
	zipFile := sqlFile + ".zip"

	fmt.Printf("Backing up %s to %s...
", dbName, sqlFile)

	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, dbName)

	cmd := exec.Command(mysqldumpPath, args...)
	file, err := os.Create(sqlFile)
	if err != nil {
		fmt.Printf("Error creating file: %v
", err)
		return
	}
	cmd.Stdout = file
	err = cmd.Run()
	file.Close()

	if err != nil {
		fmt.Printf("Error running mysqldump: %v
", err)
		os.Remove(sqlFile)
		return
	}

	// Zip it
	fmt.Printf("Compressing to %s...
", zipFile)
	if err := zipSingleFile(sqlFile, zipFile); err != nil {
		fmt.Printf("Error zipping: %v
", err)
		return
	}

	os.Remove(sqlFile)
	absPath, _ := filepath.Abs(zipFile)
	fmt.Printf("Backup successful!
Path: %s
", absPath)
}

func handleRestore(user, pass string) {
	zips, _ := filepath.Glob("*.sql.zip")
	if len(zips) == 0 {
		fmt.Println("No .sql.zip files found in current directory.")
		return
	}

	fmt.Println("
Available Backups:")
	for i, z := range zips {
		fmt.Printf("[%d] %s
", i+1, z)
	}

	fmt.Printf("
Select backup to restore (1-%d): ", len(zips))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('
')
	idx := strings.TrimSpace(input)

	i, err := strconv.Atoi(idx)
	if err != nil || i < 1 || i > len(zips) {
		fmt.Println("Invalid selection.")
		return
	}

	selectedZip := zips[i-1]
	
	// Ask for target DB name
	fmt.Print("Enter target database name (will be created if not exists): ")
	dbName, _ := reader.ReadString('
')
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		fmt.Println("Database name cannot be empty.")
		return
	}

	// Create DB if not exists
	createDB(dbName, user, pass)

	// Unzip and restore
	sqlFile, err := unzipSingleFile(selectedZip)
	if err != nil {
		fmt.Printf("Error unzipping: %v
", err)
		return
	}
	defer os.Remove(sqlFile)

	fmt.Printf("Restoring %s to %s...
", sqlFile, dbName)
	
	// Prepare source command
	escapedPath := strings.ReplaceAll(sqlFile, "", "/")
	sourceCmd := fmt.Sprintf("source %s", escapedPath)

	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, dbName, "-e", sourceCmd)

	cmd := exec.Command(mysqlPath, args...)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Restore failed: %v
", err)
	} else {
		fmt.Println("Restore successful!")
	}
}

func createDB(dbName, user, pass string) {
	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName))
	exec.Command(mysqlPath, args...).Run()
}

func zipSingleFile(src, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	info, _ := file.Stat()
	header, _ := zip.FileInfoHeader(info)
	header.Name = filepath.Base(src)
	header.Method = zip.Deflate

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func unzipSingleFile(src string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return "", fmt.Errorf("zip is empty")
	}

	f := r.File[0]
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	out, err := os.Create(f.Name)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return f.Name, err
}
