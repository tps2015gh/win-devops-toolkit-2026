🌐 **Languages:** [English] | [**ภาษาไทย (Thai)**](./README.th.md) | [**简体中文 (Chinese)**](./README.zh.md)

# win-audit-2026 (abcd-setup-server)

A robust Go-based discovery and audit suite for Windows environments. This toolset is designed for system administrators and developers to quickly audit Windows Server/11 instances and XAMPP installations.

---

### 👨‍💻 Author & Intelligent Assistant
- **Lead Developer & Owner:** **tps2015gh** (Human)
- **Intelligent Assistant:** **Gemini AI** (CLI Agent)

**Legal Note on Authorship:** This project is authored and owned exclusively by **tps2015gh**. Gemini AI acted as a sophisticated development tool and intelligent assistant, providing implementation support and technical suggestions under the direct instruction and oversight of the project owner. All intellectual property, copyright, and final architectural decisions reside with the human author.

---

## 🚀 Key Features

### 1. System Information Collector (`system_info.exe`)
- **Memory:** High-precision RAM (Total, Free, Load) via Windows `syscall`.
- **CPU:** Detailed processor specs and core counts.
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
3. Check the `./output/` directory for detailed reports.

---

## 📄 License
This project is licensed under the **MIT License**. The copyright holder and owner is **tps2015gh**. See the `LICENSE` file for full details.

---
*Developed by tps2015gh with Gemini AI Assistance - 2024*
