# CLAUDE.md - Project Intelligence & Development Guide

## Project Overview: `win-devops-toolkit-2026` (DevOps Server Admin Suite)
A specialized suite of Go-based system discovery, audit, and deployment tools designed for Windows Server and Windows 11 environments. The project provides high-fidelity system insights, XAMPP ecosystem analysis, and DevOps deployment capabilities.

## Core Philosophy
1. **Precision & Transparency:** Tools must provide exact technical details (versions, paths, sizes) with zero "guesswork."
2. **Zero-Footprint Deployment:** Prioritize standalone Windows executables (`.exe`) that require no external runtime or complex installation.
3. **Machine & Human Readable:** Every discovery operation must produce both a structured JSON (for automation) and a formatted TXT (for human audit) in the `./output/` directory.
4. **Concurrency by Default:** Use Goroutines to parallelize WMI and filesystem operations, ensuring fast execution even on hardware with high latency or large disk arrays.
5. **DevOps-First Design:** Tools should integrate with deployment pipelines and infrastructure-as-code workflows.

## System Architecture & Flow
- **`main.go` (System Info Collector):**
  - Uses `WMI` (`root\CIMV2`, `root\SecurityCenter2`, `root\StandardCimv2`) for OS, CPU, Disk, and Security status.
  - Uses `syscall` (`GlobalMemoryStatusEx`) for high-precision RAM metrics.
  - **Output:** `./output/system_info.{json,txt}`.
- **`xampp_collector.go` (Ecosystem Scanner):**
  - Scans root drives for `xampp*` patterns.
  - Recursively calculates directory sizes (Total vs. `htdocs`).
  - Probes binaries (`php.exe`, `mysql.exe`) for internal version strings (e.g., MariaDB versions).
  - **Output:** `./output/xampp_report.{json,txt}`.
- **`patch_collector.go` (Patch Auditor):**
  - Executes PowerShell `Get-HotFix` for precise update tracking.
  - Formats date strings for standardized JSON parsing.
  - **Output:** `./output/patch_report.{json,txt}`.
- **`firewall_collector.go` (Firewall Auditor):**
  - Uses `Get-NetFirewallRule` to map active rules to ports and programs.
  - Handles complex object filtering for clean JSON output.
  - **Output:** `./output/firewall_report.{json,txt}`.
- **`go-diff-packer.go` (Deployment Tool):**
  - Compares two directories and identifies changed/new files.
  - Uses size-first comparison, then SHA-256 hash verification.
  - Supports JSON output for CI/CD integration.
  - **Output:** Auto-incrementing diff folders and structured reports.
- **`db_manager.go` (Database Management):**
  - Automated MariaDB/MySQL discovery and auditing.
  - Backup/restore with compression.
  - Table and collation analysis.
  - **Output:** `./output/db_*.{json,txt}`.
- **`build.bat`:** Centralized compilation script targeting `windows/amd64`.

## Development Direction
- **Phase 1 (Current):** System and XAMPP discovery, database management, directory diffing and deployment.
- **Phase 2 (Future):** Discovery of IIS sites, MSSQL instances, Windows Task Scheduler entries, and container management.
- **Phase 3 (Future):** Automated "Setup" capabilities—using the discovered data to generate configuration scripts or Infrastructure-as-Code (IaC) templates. Integration with Terraform, Ansible, and Docker.

## Rules for AI Assistants
- **Strict Typing:** Always use Go structs for WMI queries to ensure type safety.
- **Path Handling:** Always use `path/filepath` for Windows compatibility.
- **Security:** Never commit the `./output/` directory or `*.exe` files. Ensure `.gitignore` is always respected.
- **Error Handling:** WMI queries can fail depending on user permissions (especially `SecurityCenter2` on Servers). Always provide fallbacks or "Unknown" status instead of crashing.
- **DevOps Integration:** Consider how tools can integrate with CI/CD pipelines and automation frameworks.
- **Performance:** Optimize for large deployments (1000+ servers) and high-latency networks.

---
*This document serves as the foundational context for any AI assisting in the development of this project.*
