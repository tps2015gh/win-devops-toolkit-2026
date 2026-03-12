# คู่มือการติดตั้ง MCP Server

🌐 **Languages:** [English](./README.mcp_setup.md) | [**ภาษาไทย (Thai)**]

คู่มือการติดตั้ง Database MCP Server สำหรับ **Qwen Code** และ **Gemini CLI**

---

## 📋 ข้อกำหนดเบื้องต้น

- Go 1.23+
- MySQL/MariaDB ติดตั้งแล้ว (XAMPP หรือแบบแยก)
- Git ตั้งค่าแล้ว

---

## 🚀 การติดตั้งอย่างรวดเร็ว

### ขั้นตอนที่ 1: สร้าง MCP Server

```powershell
cd C:\dev\abcd_setup_server
go get github.com/go-sql-driver/mysql
go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
```

### ขั้นตอนที่ 2: ตั้งค่า MCP Client

#### สำหรับ Qwen Code

สร้างหรือแก้ไขไฟล์ตั้งค่า global ที่ `%USERPROFILE%\.qwen\mcp.json`:

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**คำสั่ง:**
```powershell
# สร้างโฟลเดอร์ถ้ายังไม่มี
if not exist "%USERPROFILE%\.qwen" mkdir "%USERPROFILE%\.qwen"

# สร้างไฟล์ตั้งค่า (รันใน PowerShell)
@'
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
'@ | Out-File -FilePath "$env:USERPROFILE\.qwen\mcp.json" -Encoding UTF8
```

#### สำหรับ Gemini CLI

สร้างหรือแก้ไขไฟล์ตั้งค่า global ที่ `%APPDATA%\gemini-cli\mcp.json`:

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**คำสั่ง:**
```powershell
# สร้างโฟลเดอร์ถ้ายังไม่มี
if not exist "$env:APPDATA\gemini-cli" mkdir "$env:APPDATA\gemini-cli"

# สร้างไฟล์ตั้งค่า
@'
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
'@ | Out-File -FilePath "$env:APPDATA\gemini-cli\mcp.json" -Encoding UTF8
```

---

## ✅ ตรวจสอบการติดตั้ง

### ใน Qwen Code

1. **ปิดและเปิด Qwen Code ใหม่**
2. รัน `/mcp list` - ควรเห็น:
   ```
   database - C:\dev\abcd_setup_server\db_mcp_server.exe
   ```
3. ทดสอบการทำงาน:
   ```
   /mcp database list_databases {"config": {"user": "root", "password": ""}}
   ```

### ใน Gemini CLI

1. **ปิดและเปิด Gemini CLI ใหม่**
2. รัน `/mcp list` - ควรเห็น database server
3. ทดสอบการทำงานคล้ายกัน

---

## 🔧 เครื่องมือ MCP ที่มี

| เครื่องมือ | คำอธิบาย |
|-----------|----------|
| `list_databases` | แสดงฐานข้อมูลทั้งหมดที่เข้าถึงได้ |
| `list_tables` | แสดงตารางในฐานข้อมูล |
| `get_table_schema` | แสดงรายละเอียดคอลัมน์ |
| `query_data` | รันคำสั่ง SELECT |
| `get_table_stats` | แสดงสถิติของตาราง |
| `get_database_info` | แสดงเวอร์ชันและตัวแปร |
| `search_data` | ค้นหาข้อมูลแบบซับซ้อน |
| `get_foreign_keys` | แสดงความสัมพันธ์ foreign key |
| `get_indexes` | แสดงข้อมูล index |
| `aggregate_data` | ทำ COUNT, SUM, AVG, MIN, MAX |
| `switch_database` | สลับฐานข้อมูล |

---

## 📝 ตัวอย่างการใช้งาน

### แสดงฐานข้อมูล
```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

### แสดงตาราง
```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb"
}
```

### Query ข้อมูล
```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "query": "SELECT * FROM users LIMIT 10"
}
```

### สลับฐานข้อมูล (Workflow แนะนำ)
```json
{
  "config": {
    "user": "root",
    "password": ""
  },
  "database": "production_db"
}
```

---

## 🐛 การแก้ปัญหา

### "No MCP servers configured"

- ตรวจสอบว่ามีไฟล์ตั้งค่าในตำแหน่งที่ถูกต้อง
- ปิดและเปิดแอปใหม่หลังสร้าง/แก้ไขไฟล์ตั้งค่า
- ตรวจสอบ path ของ `db_mcp_server.exe`

### "mysql.exe not found"

- ตรวจสอบว่า MySQL/MariaDB ติดตั้งแล้ว
- เพิ่มใน PATH หรือตั้งค่า `MYSQL_HOME`:
  ```powershell
  setx MYSQL_HOME "C:\xampp\mysql\bin"
  ```

### "Access denied"

- ตรวจสอบ username และ password
- ตรวจสอบสิทธิ์ของ user

### "SQL unknown driver"

- สร้างใหม่พร้อม MySQL driver:
  ```powershell
  go get github.com/go-sql-driver/mysql
  go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
  ```
- ปิดและเปิด Qwen Code / Gemini CLI ใหม่

---

## 📁 ตำแหน่งไฟล์ตั้งค่า

| แอปพลิเคชัน | ตำแหน่งไฟล์ตั้งค่า |
|-------------|-------------------|
| Qwen Code | `%USERPROFILE%\.qwen\mcp.json` |
| Gemini CLI | `%APPDATA%\gemini-cli\mcp.json` |

---

## 🔒 หมายเหตุด้านความปลอดภัย

- อนุญาตเฉพาะ SELECT, SHOW, DESCRIBE, EXPLAIN
- ไม่允许 INSERT, UPDATE, DELETE, DROP
- รหัสผ่านส่งตรงไปยัง MySQL driver (ไม่บันทึก)
- ตรวจสอบข้อมูลใน MCP client

---

## 📖 เอกสารที่เกี่ยวข้อง

- [Main README](./README.md) - ภาพรวมโปรเจกต์
- [DB Manager](./README.db.md) - เครื่องมือ backup/restore
- [Database MCP Server](./README.db_mcp.md) - เอกสารเครื่องมือโดยละเอียด
- [Archiver](./README.archiver.md) - ZIP archiver รองรับ Unicode

---

*Developed by tps2015gh - 2026*
