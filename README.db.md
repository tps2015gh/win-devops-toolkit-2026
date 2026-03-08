# MariaDB Database Manager (Backup/Restore)

[**🏠 Back to Home**](./README.md) | [**English**](#english) | [**ภาษาไทย**](#thai) | [**简体中文**](#chinese) | [**Archiver Tool**](./README.archiver.md) | [**WSL Manager**](./README.wsl.md)

---

<a name="english"></a>
## 🇺🇸 English
**Database Manager** is a Go-based utility to list, backup, and restore MariaDB/MySQL databases. It automatically compresses backups into `.sql.zip` format.

### Features
- **Auto-Discovery:** Searches for `mysql.exe` and `mysqldump.exe` in PATH and XAMPP folders.
- **Connection Test:** Automatically tests your `user/pass` and allows retrying if the connection fails.
- **Database Audit:** List tables and view detailed stats (row counts & encoding/collation).
- **Interactive:** Select databases from a list for backup or auditing.
- **Compressed:** Backups are saved as `.sql.zip`.
- **Easy Restore:** Select a `.sql.zip` file and restore it to any database name.

### Usage
```cmd
db_manager.exe
```

---

<a name="thai"></a>
## 🇹🇭 ภาษาไทย
**Database Manager** คือเครื่องมือสำหรับแสดงรายชื่อ, สำรองข้อมูล (Backup) และคืนค่า (Restore) ฐานข้อมูล MariaDB/MySQL โดยจะบีบอัดไฟล์สำรองเป็นนามสกุล `.sql.zip` ให้อัตโนมัติ

### คุณสมบัติ
- **ค้นหาอัตโนมัติ:** ค้นหาไฟล์ `mysql.exe` และ `mysqldump.exe` ใน PATH และโฟลเดอร์ XAMPP โดยอัตโนมัติ
- **ทดสอบการเชื่อมต่อ:** ทดสอบ `user/pass` ทันทีที่กรอก และเลือกกรอกใหม่ได้หากเชื่อมต่อล้มเหลว
- **ตรวจสอบข้อมูล (Audit):** แสดงรายชื่อตาราง พร้อมสถิติจำนวนแถวและการเข้ารหัส (Encoding/Collation)
- **โต้ตอบได้:** แสดงรายชื่อฐานข้อมูลให้เลือกก่อนสำรองข้อมูลหรือตรวจสอบข้อมูล
- **ประหยัดพื้นที่:** ไฟล์สำรองข้อมูลจะถูกบีบอัดเป็น `.sql.zip`
- **คืนค่าง่าย:** เลือกไฟล์ `.sql.zip` และระบุชื่อฐานข้อมูลที่ต้องการคืนค่าได้ทันที

### วิธีใช้งาน
```cmd
db_manager.exe
```

---

<a name="chinese"></a>
## 🇨🇳 简体中文
**Database Manager** 是一款基于 Go 语言开发的数据库管理工具，用于列出、备份和还原 MariaDB/MySQL 数据库。它会自动将备份压缩为 `.sql.zip` 格式。

### 功能特点
- **自动发现:** 在 PATH 和 XAMPP 文件夹中自动搜索 `mysql.exe` 和 `mysqldump.exe`。
- **连接测试:** 自动测试您的 `user/pass` 凭据，并在连接失败时允许重试。
- **数据库审计:** 列出所有表并查看详细统计信息（行数和字符集编码/Collation）。
- **交互式操作:** 从列表中选择要备份或审计的数据库。
- **压缩备份:** 备份文件保存为 `.sql.zip`。
- **快速还原:** 选择 `.sql.zip` 文件并将其还原到任何指定的数据库名称。

### 使用方法
```cmd
db_manager.exe
```

---

**Disclaimer:** I think OK To Use, but no guarantee or warranty.
