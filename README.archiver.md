# Archiver - Robust Unicode-aware Compression Tool

[**🏠 Back to Home**](./README.md) | [**English**](#english) | [**ภาษาไทย**](#thai) | [**简体中文**](#chinese) | [**DB Manager**](./README.db.md)

---

<a name="english"></a>
## 🇺🇸 English
**Archiver** is a high-performance Go-based compression utility specifically designed to overcome common issues on Windows Server and Windows 11. It solves the problem of "characters that cannot be used in compressed folders" (like Thai or Chinese Unicode) that often crash the built-in Windows Zip utility.

### Features
- **Unicode Support:** Native support for UTF-8 filenames (Thai, Chinese, and others).
- **Hard Situation Ready:** Optimized for deep directory structures and large repositories (like `.git` history).
- **Auto-Timestamp:** Generates unique filenames (e.g., `folder_20260308_120000.zip`).
- **Zero Dependencies:** Compiles into a single standalone `.exe` for Windows Server.

### Usage
```cmd
archiver.exe <folder_name>
```

### Unarchive (Restore)
To extract a ZIP file into an auto-incrementing folder (e.g., `_unzip01`, `_unzip02`...):
```cmd
unarchiver.exe <zip_file_name>
```

---

<a name="thai"></a>
## 🇹🇭 ภาษาไทย
**Archiver** คือโปรแกรมบีบอัดไฟล์ประสิทธิภาพสูงที่พัฒนาด้วยภาษา Go ออกแบบมาเพื่อแก้ปัญหาที่พบบ่อยบน Windows Server และ Windows 11 โดยเฉพาะ โดยเน้นแก้ปัญหา "ไม่สามารถบีบอัดได้เนื่องจากมีตัวอักษรที่ไม่รองรับ" (เช่น ภาษาไทย) ซึ่งมักเกิดขึ้นกับโปรแกรม Zip มาตรฐานของ Windows

### คุณสมบัติ
- **รองรับ Unicode:** บีบอัดไฟล์ที่มีชื่อเป็นภาษาไทยหรือภาษาอื่นๆ ได้ 100%
- **ทนทานสูง:** รองรับการบีบอัดในโฟลเดอร์ที่มีโครงสร้างซับซ้อน (เช่น ประวัติใน `.git`)
- **ใส่เวลาให้อัตโนมัติ:** สร้างชื่อไฟล์ Zip ตามเวลาปัจจุบัน (เช่น `folder_20260308_120000.zip`)
- **ไม่ต้องติดตั้ง:** รวมเป็นไฟล์ `.exe` ไฟล์เดียว ใช้งานได้ทันทีบน Windows Server

### วิธีใช้งาน
```cmd
archiver.exe <ชื่อโฟลเดอร์>
```

### วิธีแตกไฟล์ (Restore)
เพื่อแตกไฟล์ Zip ไปยังโฟลเดอร์แบบรันเลขให้อัตโนมัติ (เช่น `_unzip01`, `_unzip02`...):
```cmd
unarchiver.exe <ชื่อไฟล์_zip>
```

---

<a name="chinese"></a>
## 🇨🇳 简体中文
**Archiver** 是一款基于 Go 语言开发的高性能压缩工具，专门用于解决 Windows Server 和 Windows 11 上的常见问题。它解决了 Windows 自带压缩工具经常遇到的“包含无法在压缩文件夹中使用的字符”（如中文或泰文 Unicode）而导致失败的问题。

### 功能特点
- **完美支持 Unicode:** 原生支持 UTF-8 文件名（中文、泰文等）。
- **适应极端环境:** 优化了对深度目录结构和大容量库（如 `.git` 历史记录）的处理。
- **自动时间戳:** 自动生成带时间的唯一文件名（如 `folder_20260308_120000.zip`）。
- **零依赖:** 编译为单个独立的 `.exe` 文件，适用于 Windows Server 各个版本。

### 使用方法
```cmd
archiver.exe <文件夹名称>
```

### 还原文件 (Restore)
将 ZIP 文件解压到自动递增的文件夹中（例如 `_unzip01`、`_unzip02` 等）：
```cmd
unarchiver.exe <zip_文件名>
```

---

**Disclaimer:** I think OK To Use, but no guarantee or warranty.
