# win_info (dol-setup-server)

A robust Go-based discovery and audit suite for Windows environments. This toolset is designed for system administrators and developers to quickly audit Windows Server/11 instances and XAMPP installations.

---

### 🤖 Human-AI Partnership
This project was co-engineered through a direct collaboration between **tps2015gh** (Human Programmer) and **Gemini AI** (CLI Interactive Agent). This partnership combined human architectural vision and specialized requirements with AI's precision in implementation, testing, and documentation.

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
This project is licensed under the **MIT License**—a permissive license that allows for commercial use, modification, and distribution with minimal restrictions. See the `LICENSE` file for full details.

---
*Created with Gemini CLI - 2024*
