# GEMINI.md - Project Intelligence & Development Guide

## Project Overview: `win-audit-2026` (dol_setup_server)
A specialized suite of Go-based audit and discovery tools designed for Windows Server and Windows 11 environments. The project provides high-fidelity system insights and XAMPP ecosystem analysis.

## Core Philosophy
1. **Precision & Transparency:** Tools must provide exact technical details (versions, paths, sizes) with zero "guesswork."
2. **Zero-Footprint Deployment:** Prioritize standalone Windows executables (`.exe`) that require no external runtime or complex installation.
3. **Machine & Human Readable:** Every discovery operation must produce both a structured JSON (for automation) and a formatted TXT (for human audit) in the `./output/` directory.
4. **Concurrency by Default:** Use Goroutines to parallelize WMI and filesystem operations, ensuring fast execution even on hardware with high latency or large disk arrays.

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
- **`build.bat`:** Centralized compilation script targeting `windows/amd64`.

## Development Direction
- **Phase 1 (Current):** System and XAMPP discovery.
- **Phase 2 (Future):** Discovery of IIS sites, MSSQL instances, and Windows Task Scheduler entries.
- **Phase 3 (Future):** Automated "Setup" capabilities—using the discovered data to generate configuration scripts or Infrastructure-as-Code (IaC) templates.

## Rules for AI Assistants
- **Strict Typing:** Always use Go structs for WMI queries to ensure type safety.
- **Path Handling:** Always use `path/filepath` for Windows compatibility.
- **Security:** Never commit the `./output/` directory or `*.exe` files. Ensure `.gitignore` is always respected.
- **Error Handling:** WMI queries can fail depending on user permissions (especially `SecurityCenter2` on Servers). Always provide fallbacks or "Unknown" status instead of crashing.

---
*This document serves as the foundational context for any AI assisting in the development of this project.*
