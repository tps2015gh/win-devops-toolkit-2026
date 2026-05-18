🌐 **Languages:** [English] | [**ภาษาไทย (Thai)**](./README.th.md) | [**简体中文 (Chinese)**](./README.zh.md) | [**Archiver Tool**](./README.archiver.md) | [**DB Manager**](./README.db.md) | [**DB MCP Server**](./README.db_mcp.md) | [**ไทย**](./README.th.md) | [**Qwen MCP Guide**](./README.qwen_mcp.md) | [**WSL Manager**](./README.wsl.md) | [**Go Diff Packer**](./README.go-diff-packer.md) | [**Usage**](./README.go-diff-packer.usage.md) | [**E2E Testing**](./go-diff-packer/E2E_CLI_TESTING.md) | [**Changes v1.2.0**](./RELEASE_v1.2.0.txt)

# win-devops-toolkit-2026

A robust Go-based system discovery, audit, and deployment suite for Windows environments. This comprehensive DevOps toolkit is designed for system administrators and developers to quickly discover Windows Server/11 instances, audit system configurations, manage databases, and orchestrate deployment operations with automated tools.

---

### 👨‍💻 Project Team & Contributors
- **Director & Supervisor:** **tps2015gh** (Human)
- **Programming & Testing:** tps2015gh, Qwen Code
- **DevOps & Architecture:** Claude Haiku 4.5 (AI Assistant)
- **Intelligent Assistants:** Qwen Code (CLI Agent), Claude Haiku 4.5 (DevOps & Architecture)

**Legal Note on Authorship:** This project is solely owned and directed by **tps2015gh**. Qwen Code and Claude Haiku 4.5 provide intelligent assistance, implementation support, and architectural guidance under direct instruction and oversight. All intellectual property and final decisions reside with the human author.

---

## 🚀 Key Features

### 1. System Information Collector (`system_info.exe`)
- **Memory:** High-precision RAM (Total, Free, Load) via Windows `syscall`.
- **CPU:** Detailed processor specs, including **Physical Sockets**, **Cores per Socket**, and **Total Core** counts.
- **Storage:** Logical disk analysis (Total/Free/Used space, File System).
- **OS:** Version and Build identification.
- **Security:** Detection of active Antivirus and Firewall products via WMI.
- **Network:** Full adapter configuration, including IP addresses and DNS search order.

### 2. XAMPP Ecosystem Scanner (`xampp_collector.exe`)
- **Discovery:** Automatic detection of `xampp*` directories on `C:`, `D:`, and `E:`.
- **Sizing:** Calculates total XAMPP disk usage and specific `htdocs` folder size.
- **Version Auditing:** Extracts precise **PHP** and **MariaDB/MySQL** versions by directly querying the binaries.
- **Output:** Saves results as structured JSON and human-readable TXT in `./output/`.

### 3. Windows Patch Collector (`patch_collector.exe`)
- **Audit:** Retrieves all installed Windows Updates and HotFixes using the PowerShell `Get-HotFix` cmdlet.
- **Details:** Includes HotFixID (KB number), Description, Install Date, and InstalledBy info.
- **Output:** Saves results in JSON and TXT formats for audit compliance.

### 4. Windows Firewall Auditor (`firewall_collector.exe`)
- **Discovery:** Scans and lists all active (Enabled) Windows Firewall rules.
- **Details:** Provides protocol (TCP/UDP), local ports, action (Allow/Block), and the specific program path for each rule.
- **Insight:** Helps identify open ports and factory-default vs. custom application rules.
- **Output:** Saves results in JSON and TXT formats in `./output/`.

### 5. Unicode Archiver & Unarchiver (`archiver.exe` & `unarchiver.exe`)
- **Problem Solver:** Fixes the "characters that cannot be used in compressed folders" issue on Windows (supporting Thai, Chinese, etc.).
- **Interactive Unarchive:** `unarchiver.exe` provides an interactive menu to select and extract ZIP files to auto-incrementing folders (`_unzip01`, `_unzip02`...).
- **Details:** See the [**Archiver Documentation**](./README.archiver.md) for full usage instructions.

### 6. MariaDB Database Manager (`db_manager.exe`)
- **Backup/Restore:** Compressed SQL backups (`.sql.zip`) and easy restoration.
- **Database Audit:** Lists all databases, tables, and provides stats on row counts and character encoding (Collation).
- **Auto-Discovery:** Automatically detects MariaDB/MySQL paths in XAMPP or system environment.
- **Details:** See the [**DB Manager Documentation**](./README.db.md) for full usage instructions.

### 7. Dev Tool Discovery Engine (`dev_tool_collector.exe`)
- **Comprehensive Audit:** Scans for Go, Rust, Java, Python, Node.js, PHP, Perl, VS Code, Android SDK, and more.
- **Library Check:** Detects Python libraries like Pandas and Tkinter.
- **Rich Reporting:** Generates a structured JSON and a beautiful, interactive HTML report in `./output/`.
- **Progress Tracking:** Shows real-time status updates during discovery.

### 8. Go Diff Packer (`go-diff-packer.exe`)
- **Directory Comparison**: Compares two directories and identifies new/modified files.
- **Smart Detection**: Uses size comparison first, then SHA-256 hash for accuracy.
- **Differential Backup**: Copies only changed files to auto-incrementing output folders (`_outdiff_01`, `_outdiff_02`...).
- **Deployment Ready**: Perfect for deploying only changed files to production servers.
- **Summary Statistics**: Shows files compared, copied, skipped, and errors.
- **Details:** See the [**Go Diff Packer Documentation**](./README.go-diff-packer.md), [**Usage Guide**](./README.go-diff-packer.usage.md), and [**E2E Testing Guide**](./go-diff-packer/E2E_CLI_TESTING.md).

### 9. Database MCP Server (`db_mcp_server.exe`)
- **AI-Ready:** Database access for Claude, ChatGPT, and other AI assistants.
- **Schema Inspection:** View database structure and relationships.
- **Complex Queries:** Execute sophisticated SQL queries safely.
- **Read-Only Mode:** Secure by design with configurable permissions.
- **Details:** See the [**DB MCP Server Documentation**](./README.db_mcp.md).

### 10. CodeIgniter 3 Project Searcher (`ci3_searcher.exe`)
- **Intelligent Search:** Indexes CodeIgniter 3 project components (Controllers, Models, Views, Config, JS, CSS, Database interactions) and allows for similarity-based search.
- **Vector Space Search:** Utilizes vector embeddings with attention mechanism for fast and relevant search results.
- **Component Awareness:** Prioritizes matches based on component type (e.g., searching for "user controller" will prioritize user-related controllers).
- **Interactive CLI:** Provides an interactive command-line interface for real-time search queries.
- **Result Export:** Automatically exports extensive search results to a text file.
- **Details:** See the [**CI3 Searcher Documentation**](./README.ci3_searcher.md) for full usage and technical insights.

### 11. WSL Manager (`wsl_manager.exe`)
- **WSL Management:** Discover and manage Windows Subsystem for Linux installations.
- **Distribution Audit:** List installed Linux distributions and their configurations.
- **Details:** See the [**WSL Manager Documentation**](./README.wsl.md).

---

## 🛠️ Usage

### Build
To compile the source code into standalone Windows executables:
```powershell
.\build.bat
```

### Run
1. Run `system_info.exe` for general system audit.
2. Run `xampp_collector.exe` for web server discovery.
3. Run `go-diff-packer.exe` for deployment operations.
4. Check the `./output/` directory for detailed reports.

---

## 🏗️ Development Specifications
- **Language:** Go (Golang)
- **Compiler Version:** **Go 1.21.1** (Target: Windows/AMD64)
- **Dependencies:** Uses WMI and Windows syscalls for high-fidelity discovery.
- **Compatibility:** Optimized for Windows Server 2016, 2019, 2022, and Windows 10/11.
- **Tools Included:** System Info, XAMPP Collector, Patch Collector, Firewall Collector, Archiver, DB Manager, Dev Tool Collector, Go Diff Packer, DB MCP Server, CI3 Searcher, WSL Manager

---

## 📄 License
This project is licensed under the **MIT License**. The copyright holder and owner is **tps2015gh**. See the `LICENSE` file for full details.

---

*Developed by tps2015gh with Qwen Code and Claude Haiku 4.5 Assistance - 2026*
