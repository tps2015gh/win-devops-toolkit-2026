# WSL Distribution Manager

[**🏠 Back to Home**](./README.md) | [**English**](#english) | [**ภาษาไทย**](#thai) | [**简体中文**](#chinese) | [**Archiver Tool**](./README.archiver.md) | [**DB Manager**](./README.db.md)

---

<a name="english"></a>
## 🇺🇸 English
**WSL Manager** is a Go-based interactive utility to manage Windows Subsystem for Linux (WSL) distributions. It simplifies listing, installing, and removing distributions while monitoring your system's disk space.

### Features
- **Interactive Menu:** Easy-to-use numbered selection for all operations.
- **Auto-Discovery:** Lists all installed distributions (Local) and available ones (Online).
- **Disk Space Monitor:** Displays available and used disk space before and after operations.
- **One-Click Install/Remove:** Select a number to install or unregister a distribution.
- **Clean Output:** Automatically handles and cleans UTF-16/null-byte output from WSL commands.

### Usage
```cmd
wsl_manager.exe
```

---

<a name="thai"></a>
## 🇹🇭 ภาษาไทย
**WSL Manager** คือเครื่องมือสำหรับจัดการการติดตั้งและลบ Linux Distributions บน Windows Subsystem for Linux (WSL) โดยเน้นความง่ายในการใช้งานผ่านเมนูแบบเลือกตัวเลข พร้อมทั้งตรวจสอบพื้นที่ว่างบนดิสก์ให้อัตโนมัติ

### คุณสมบัติ
- **เมนูแบบเลือกตัวเลข:** เลือกการทำงานต่างๆ ได้ง่ายผ่านการกดเลข 1-6
- **ตรวจสอบรายชื่อ:** แสดงรายชื่อ Distribution ทั้งที่ติดตั้งแล้วในเครื่อง และที่มีให้เลือกติดตั้งใหม่ทางออนไลน์
- **ตรวจสอบพื้นที่ดิสก์:** แสดงข้อมูลพื้นที่ดิสก์ที่เหลืออยู่ก่อนและหลังการติดตั้งหรือลบข้อมูล
- **ติดตั้งและลบง่าย:** เลือกหมายเลขเพื่อติดตั้งหรือถอนการติดตั้ง (Unregister) ได้ทันที
- **จัดการภาษาอัตโนมัติ:** แก้ไขปัญหาการแสดงผลตัวอักษรผิดเพี้ยนจากคำสั่ง WSL ให้ถูกต้อง

### วิธีใช้งาน
```cmd
wsl_manager.exe
```

---

<a name="chinese"></a>
## 🇨🇳 简体中文
**WSL Manager** 是一款基于 Go 语言开发的交互式工具，用于管理 Windows Subsystem for Linux (WSL) 分发版。它通过编号菜单简化了列出、安装和删除分发版的操作，同时还能监控系统的磁盘空间。

### 功能特点
- **交互式菜单:** 所有操作均可通过简单的编号选择完成。
- **自动发现:** 列出所有本地已安装的分发版以及在线可用的分发版。
- **磁盘空间监控:** 在操作前后显示可用和已用的磁盘空间。
- **一键安装/删除:** 通过选择编号即可安装或注销分发版。
- **清洁输出:** 自动处理并清理 WSL 命令产生的 UTF-16/空字节输出，确保显示正常。

### 使用方法
```cmd
wsl_manager.exe
```

---

**Disclaimer:** I think OK To Use, but no guarantee or warranty.
