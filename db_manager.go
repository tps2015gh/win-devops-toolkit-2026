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
	fmt.Println("MariaDB Database Manager (Backup/Restore/Audit)")
	fmt.Println("==============================================")

	// 1. Setup Binaries
	setupBinaries()

	// 2. Get Credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter credentials (format: user/pass): ")
	creds, _ := reader.ReadString('\n')
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
		fmt.Println("\nChoose action:")
		fmt.Println("1. List Databases & Backup")
		fmt.Println("2. Restore from .sql.zip")
		fmt.Println("3. List Tables in a Database")
		fmt.Println("4. List Tables with Stats (Rows & Encoding)")
		fmt.Println("5. Exit")
		fmt.Print("Select (1-5): ")

		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "1":
			handleBackup(user, pass)
		case "2":
			handleRestore(user, pass)
		case "3":
			handleListTables(user, pass)
		case "4":
			handleTableStats(user, pass)
		case "5":
			return
		default:
			fmt.Println("Invalid selection.")
		}
	}
}

func setupBinaries() {
	mysqlPath, _ = exec.LookPath("mysql.exe")
	mysqldumpPath, _ = exec.LookPath("mysqldump.exe")

	if mysqlPath == "" || mysqldumpPath == "" {
		commonPaths := []string{
			"C:\\xampp\\mysql\\bin",
			"D:\\xampp\\mysql\\bin",
			"C:\\xampp_v8_1_25\\mysql\\bin",
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
		p, _ := reader.ReadString('\n')
		mysqlPath = strings.TrimSpace(p)
	}
	if mysqldumpPath == "" {
		fmt.Print("mysqldump.exe not found. Enter full path: ")
		p, _ := reader.ReadString('\n')
		mysqldumpPath = strings.TrimSpace(p)
	}
	
	fmt.Printf("Using mysql:     %s\n", mysqlPath)
	fmt.Printf("Using mysqldump: %s\n", mysqldumpPath)
}

func selectDatabase(user, pass string) string {
	dbs := listDatabases(user, pass)
	if len(dbs) == 0 {
		fmt.Println("No databases found or connection failed.")
		return ""
	}

	fmt.Println("\nAvailable Databases:")
	for i, db := range dbs {
		fmt.Printf("[%d] %s\n", i+1, db)
	}

	fmt.Printf("\nSelect database (1-%d): ", len(dbs))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	idx := strings.TrimSpace(input)

	i, err := strconv.Atoi(idx)
	if err != nil || i < 1 || i > len(dbs) {
		fmt.Println("Invalid selection.")
		return ""
	}
	return dbs[i-1]
}

func handleBackup(user, pass string) {
	dbName := selectDatabase(user, pass)
	if dbName != "" {
		backupDB(dbName, user, pass)
	}
}

func handleListTables(user, pass string) {
	dbName := selectDatabase(user, pass)
	if dbName == "" {
		return
	}

	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, dbName, "-e", "SHOW TABLES;")

	cmd := exec.Command(mysqlPath, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error listing tables: %v\n", err)
		return
	}

	fmt.Printf("\nTables in %s:\n", dbName)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Tables_in") {
			continue
		}
		fmt.Printf("  - %s\n", line)
	}
}

func handleTableStats(user, pass string) {
	dbName := selectDatabase(user, pass)
	if dbName == "" {
		return
	}

	query := "SELECT TABLE_NAME, TABLE_ROWS, TABLE_COLLATION FROM information_schema.TABLES WHERE TABLE_SCHEMA = '" + dbName + "';"
	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, "-e", query)

	cmd := exec.Command(mysqlPath, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error fetching table stats: %v\n", err)
		return
	}

	fmt.Printf("\nStats for %s:\n", dbName)
	fmt.Printf("%-30s | %-10s | %-20s\n", "Table Name", "Rows", "Collation/Encoding")
	fmt.Println(strings.Repeat("-", 65))

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" { // Skip header and empty lines
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			fmt.Printf("%-30s | %-10s | %-20s\n", parts[0], parts[1], parts[2])
		} else if len(parts) == 2 {
			fmt.Printf("%-30s | %-10s | %-20s\n", parts[0], parts[1], "N/A")
		}
	}
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
		fmt.Printf("Error listing databases: %v\n", err)
		return nil
	}

	lines := strings.Split(string(out), "\n")
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

	fmt.Printf("Backing up %s to %s...\n", dbName, sqlFile)

	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, dbName)

	cmd := exec.Command(mysqldumpPath, args...)
	file, err := os.Create(sqlFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	cmd.Stdout = file
	err = cmd.Run()
	file.Close()

	if err != nil {
		fmt.Printf("Error running mysqldump: %v\n", err)
		os.Remove(sqlFile)
		return
	}

	fmt.Printf("Compressing to %s...\n", zipFile)
	if err := zipSingleFile(sqlFile, zipFile); err != nil {
		fmt.Printf("Error zipping: %v\n", err)
		return
	}

	os.Remove(sqlFile)
	absPath, _ := filepath.Abs(zipFile)
	fmt.Printf("Backup successful!\nPath: %s\n", absPath)
}

func handleRestore(user, pass string) {
	zips, _ := filepath.Glob("*.sql.zip")
	if len(zips) == 0 {
		fmt.Println("No .sql.zip files found in current directory.")
		return
	}

	fmt.Println("\nAvailable Backups:")
	for i, z := range zips {
		fmt.Printf("[%d] %s\n", i+1, z)
	}

	fmt.Printf("\nSelect backup to restore (1-%d): ", len(zips))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	idx := strings.TrimSpace(input)

	i, err := strconv.Atoi(idx)
	if err != nil || i < 1 || i > len(zips) {
		fmt.Println("Invalid selection.")
		return
	}

	selectedZip := zips[i-1]
	
	fmt.Print("Enter target database name (will be created if not exists): ")
	dbName, _ := reader.ReadString('\n')
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		fmt.Println("Database name cannot be empty.")
		return
	}

	createDB(dbName, user, pass)

	sqlFile, err := unzipSingleFile(selectedZip)
	if err != nil {
		fmt.Printf("Error unzipping: %v\n", err)
		return
	}
	defer os.Remove(sqlFile)

	fmt.Printf("Restoring %s to %s...\n", sqlFile, dbName)
	
	escapedPath := strings.ReplaceAll(sqlFile, "\\", "/")
	sourceCmd := fmt.Sprintf("source %s", escapedPath)

	args := []string{"-u", user}
	if pass != "" {
		args = append(args, "-p"+pass)
	}
	args = append(args, dbName, "-e", sourceCmd)

	cmd := exec.Command(mysqlPath, args...)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Restore failed: %v\n", err)
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
