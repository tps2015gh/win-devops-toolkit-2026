🌐 **语言:** [**English**](./README.md) | [**ภาษาไทย (Thai)**](./README.th.md) | [简体中文] | [**压缩工具 (Archiver)**](./README.archiver.md) | [**数据库管理 (DB Manager)**](./README.db.md)

# win-audit-2026 (abcd-setup-server)

一套基于 Go 语言开发的强大 Windows 环境发现与审计工具包。该工具集专为系统管理员和开发人员设计，旨在快速审计 Windows Server/11 实例及 XAMPP 安装情况。

---

### 👨‍💻 作者与智能助手 (Author & Intelligent Assistant)
- **首席开发人员与所有者:** **tps2015gh** (人类)
- **智能助手:** **Gemini AI** (CLI 代理)

**关于著作权的法律声明:** 本项目由 **tps2015gh** 独立拥有并享有全部权利。Gemini AI 作为高级开发工具和智能助手，在项目所有者的直接指示和监督下提供实现支持和技术建议。所有知识产权、版权和最终架构决策均归属于人类作者。

---

## 🚀 核心功能

### 1. 系统信息采集器 (`system_info.exe`)
- **内存:** 通过 Windows `syscall` 获取高精度 RAM 信息（总量、剩余、负载）。
- **CPU:** 详细的处理器规格，包括 **物理插槽 (Sockets)**、**每个插槽的核心数 (Cores per Socket)** 和 **总核心数**。
- **存储:** 逻辑磁盘分析（总量/剩余/已用空间、文件系统）。
- **操作系统:** 版本和生成 (Build) 标识。
- **安全:** 通过 WMI 检测活跃的杀毒软件和防火墙产品。
- **网络:** 完整的适配器配置，包括 IP 地址和 DNS 搜索顺序。

### 2. XAMPP 生态扫描器 (`xampp_collector.exe`)
- **发现:** 自动检测 `C:`、`D:` 和 `E:` 盘上的 `xampp*` 目录。
- **容量统计:** 计算 XAMPP 总磁盘占用及特定的 `htdocs` 文件夹大小。
- **版本审计:** 通过直接查询二进制文件提取精确的 **PHP** 和 **MariaDB/MySQL** 版本。
- **输出:** 将结果以结构化 JSON 和易于阅读的 TXT 格式保存至 `./output/` 目录。

### 3. Windows 补丁采集器 (`patch_collector.exe`)
- **审计:** 使用 PowerShell `Get-HotFix` 命令获取所有已安装的 Windows 更新和热补丁。
- **详情:** 包括 HotFixID (KB 编号)、描述、安装日期和安装人信息。
- **输出:** 以 JSON 和 TXT 格式保存结果，用于审计合规。

### 4. Windows 防火墙审计器 (`firewall_collector.exe`)
- **发现:** 扫描并列出所有活动（启用）的 Windows 防火墙规则。
- **详情:** 提供每个规则的协议 (TCP/UDP)、本地端口、动作 (Allow/Block) 以及特定的程序路径。
- **洞察:** 帮助识别开放端口以及出厂默认与自定义应用程序规则。
- **输出:** 将结果以 JSON 和 TXT 格式保存至 `./output/` 目录。

### 5. Unicode 压缩与解压工具 (`archiver.exe` & `unarchiver.exe`)
- **解决痛点:** 修复了 Windows 自带工具不支持泰文、中文等 Unicode 字符的问题。
- **交互式解压:** `unarchiver.exe` 提供交互菜单，支持选择 ZIP 文件并自动解压到递增文件夹（`_unzip01`, `_unzip02`...）。
- **详情:** 请参阅 [**压缩工具文档**](./README.archiver.md) 获取详细使用说明。

### 6. MariaDB 数据库管理器 (`db_manager.exe`)
- **备份/还原:** 压缩 SQL 备份 (`.sql.zip`) 并支持快速还原。
- **数据库审计:** 列出所有数据库、表，并提供行数和字符集编码（Collation）统计。
- **自动发现:** 自动检测 XAMPP 或系统环境中的 MariaDB/MySQL 路径。
- **详情:** 请参阅 [**数据库管理文档**](./README.db.md) 获取详细使用说明。

---

## 🛠️ 使用方法

### 编译 (Build)
要将源代码编译为独立的 Windows 可执行文件：
```powershell
.\build.bat
```

### 运行 (Run)
1. 运行 `system_info.exe` 进行常规系统审计。
2. 运行 `xampp_collector.exe` 进行 Web 服务器发现。
3. 在 `./output/` 目录中查看详细报告。

---

## 🏗️ 开发规格 (Development Specifications)
- **开发语言:** Go (Golang)
- **编译器版本:** **Go 1.21.1** (目标平台: Windows/AMD64)
- **依赖库:** 使用 WMI 和 Windows syscalls 进行高保真发现。
- **兼容性:** 针对 Windows Server 2016、2019、2022 以及 Windows 10/11 进行优化。

---

## 📄 许可证 (License)
本项目采用 **MIT 许可证**。版权所有者为 **tps2015gh**。详情请参阅 `LICENSE` 文件。

---
*由 tps2015gh 在 Gemini AI 的协助下开发 - 2024*
