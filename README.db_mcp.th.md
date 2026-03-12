# Database MCP Server

🌐 **ภาษา:** [**English**](./README.db_mcp.md) | [ภาษาไทย]

เซิร์ฟเวอร์ **Model Context Protocol (MCP)** สำหรับการทำงานกับฐานข้อมูล MariaDB/MySQL เครื่องมือนี้ช่วยให้ AI Assistant สามารถเข้าถึงฐานข้อมูลโดยตรงเพื่อการค้นหา ตรวจสอบโครงสร้าง และวิเคราะห์ข้อมูลอย่างซับซ้อน

---

## 🚀 คุณสมบัติ

### ความสามารถหลัก

| เครื่องมือ | คำอธิบาย |
|------|-------------|
| `list_databases` | แสดงรายการฐานข้อมูลทั้งหมดที่ผู้ใช้สามารถเข้าถึงได้ |
| `list_tables` | แสดงตารางทั้งหมดในฐานข้อมูลที่ระบุ |
| `get_table_schema` | ดูข้อมูลคอลัมน์โดยละเอียดของตาราง |
| `query_data` | รันคำสั่ง SELECT และส่งคืนผลลัพธ์ |
| `get_table_stats` | ดูสถิติของตาราง (จำนวนแถว, การเข้ารหัส, ขนาดข้อมูล) |
| `get_database_info` | ดูข้อมูลเซิร์ฟเวอร์ฐานข้อมูล เวอร์ชัน และสถานะ |
| `search_data` | ค้นหาข้อมูลแบบซับซ้อนด้วย WHERE, ORDER BY, LIMIT, OFFSET |
| `get_foreign_keys` | ดูความสัมพันธ์ Foreign Key ของตาราง |
| `get_indexes` | ดูข้อมูล Index ของตาราง |
| `aggregate_data` | คำนวณ COUNT, SUM, AVG, MIN, MAX พร้อม GROUP BY |
| `switch_database` | **เปลี่ยนฐานข้อมูล** สำหรับการทำงานต่อไป |

---

## 📦 การติดตั้ง

### คอมไพล์จากซอร์สโค้ด

```powershell
go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
```

### ความต้องการของระบบ

- Go 1.23+
- MySQL/MariaDB client binaries (`mysql.exe`) ใน PATH หรือติดตั้ง XAMPP
- MCP-compatible client (Claude Desktop, VS Code, ฯลฯ)

---

## 🔧 การตั้งค่า

### ตั้งค่า MCP Client

**Claude Desktop** (`%APPDATA%\Claude\claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**VS Code** (`.vscode/mcp.json`):

```json
{
  "servers": {
    "database": {
      "type": "stdio",
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**Cursor** (`.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

---

## 📖 ตัวอย่างการใช้งาน

### รูปแบบ Connection Config

เครื่องมือทั้งหมดต้องการพารามิเตอร์ `config` ในรูปแบบ JSON:

```json
{
  "user": "root",
  "password": "your_password",
  "host": "127.0.0.1",
  "port": 3306,
  "database": "mydb"
}
```

| ฟิลด์ | จำเป็น | ค่าเริ่มต้น | คำอธิบาย |
|-------|----------|---------|-------------|
| `user` | ✅ ใช่ | - | ชื่อผู้ใช้ฐานข้อมูล |
| `password` | ❌ ไม่ | `""` | รหัสผ่านฐานข้อมูล |
| `host` | ❌ ไม่ | `127.0.0.1` | โฮสต์ฐานข้อมูล |
| `port` | ❌ ไม่ | `3306` | พอร์ตฐานข้อมูล |
| `database` | ❌ ไม่ | - | ฐานข้อมูลเริ่มต้น (ไม่จำเป็นสำหรับ `list_databases`) |

---

### ตัวอย่างการใช้งานเครื่องมือ

#### 0. ขั้นตอนการทำงาน: แสดงรายการและเปลี่ยนฐานข้อมูล

**ขั้นตอนที่ 1: แสดงรายการฐานข้อมูลทั้งหมด**

```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

**ผลลัพธ์:**
```json
{
  "databases": ["mydb", "testdb", "production", "shop_db"],
  "count": 4
}
```

**ขั้นตอนที่ 2: ผู้ใช้เลือกฐานข้อมูล** (เช่น `production`)

**ขั้นตอนที่ 3: เปลี่ยนไปใช้ฐานข้อมูลที่เลือก**

```json
{
  "config": {
    "user": "root",
    "password": ""
  },
  "database": "production"
}
```

**ผลลัพธ์:**
```json
{
  "success": true,
  "message": "เปลี่ยนไปใช้ฐานข้อมูล 'production' สำเร็จ",
  "current_database": "production",
  "config": {
    "user": "root",
    "host": "127.0.0.1",
    "port": 3306,
    "database": "production"
  }
}
```

> 💡 **เคล็ดลับ:** Agent ควรเรียก `list_databases` ก่อนเพื่อแสดงฐานข้อมูลให้ผู้ใช้เลือก จากนั้นใช้ `switch_database` เพื่อเชื่อมต่อกับฐานข้อมูลที่ผู้ใช้เลือก Config ที่ส่งคืนมาสามารถนำไปใช้ต่อได้

---

#### 1. แสดงรายการฐานข้อมูล

```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

**ผลลัพธ์:**
```json
{
  "databases": ["mydb", "testdb", "production"],
  "count": 3
}
```

---

#### 2. แสดงรายการตาราง

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

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "tables": ["users", "orders", "products"],
  "count": 3
}
```

---

#### 3. ดูโครงสร้างตาราง (Schema)

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb",
  "table": "users"
}
```

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "table": "users",
  "schema": [
    {
      "Field": "id",
      "Type": "int(11)",
      "Null": "NO",
      "Key": "PRI",
      "Default": null,
      "Extra": "auto_increment"
    },
    {
      "Field": "username",
      "Type": "varchar(50)",
      "Null": "NO",
      "Key": "",
      "Default": null,
      "Extra": ""
    }
  ],
  "columns": 2
}
```

---

#### 4. ค้นหาข้อมูล (SELECT)

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "query": "SELECT * FROM users WHERE id > 100 LIMIT 10"
}
```

**ผลลัพธ์:**
```json
{
  "columns": ["id", "username", "email", "created_at"],
  "rows": [
    {"id": 101, "username": "john", "email": "john@example.com", "created_at": "2026-03-01T10:00:00Z"},
    {"id": 102, "username": "jane", "email": "jane@example.com", "created_at": "2026-03-02T11:30:00Z"}
  ],
  "count": 2
}
```

> ⚠️ **ความปลอดภัย:** อนุญาตเฉพาะคำสั่ง SELECT, SHOW, DESCRIBE, และ EXPLAIN เท่านั้น

---

#### 5. ดูสถิติตาราง

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

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "stats": [
    {
      "TABLE_NAME": "users",
      "TABLE_ROWS": 1500,
      "TABLE_COLLATION": "utf8mb4_general_ci",
      "DATA_LENGTH": 32768,
      "INDEX_LENGTH": 16384,
      "DATA_FREE": 0,
      "AUTO_INCREMENT": 1501,
      "CREATE_TIME": "2026-01-15T08:00:00Z",
      "UPDATE_TIME": "2026-03-12T14:30:00Z"
    }
  ]
}
```

---

#### 6. ดูข้อมูลฐานข้อมูล

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  }
}
```

**ผลลัพธ์:**
```json
{
  "version": "10.4.32-MariaDB",
  "current_user": "root@localhost",
  "current_database": "mydb",
  "uptime_seconds": "86400",
  "variables": {
    "max_connections": "151",
    "character_set_server": "utf8mb4",
    "collation_server": "utf8mb4_general_ci",
    "version_comment": "MariaDB Server"
  },
  "active_connections": "5"
}
```

---

#### 7. ค้นหาข้อมูลแบบซับซ้อน

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb",
  "table": "orders",
  "columns": "id, customer_id, total, status, created_at",
  "where": "status = 'pending' AND total > 1000",
  "orderBy": "created_at DESC",
  "limit": 50,
  "offset": 0
}
```

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "table": "orders",
  "columns": ["id", "customer_id", "total", "status", "created_at"],
  "rows": [...],
  "count": 25,
  "limit": 50,
  "offset": 0
}
```

---

#### 8. ดู Foreign Keys

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb",
  "table": "orders"
}
```

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "table": "orders",
  "foreign_keys": [
    {
      "COLUMN_NAME": "customer_id",
      "CONSTRAINT_NAME": "fk_orders_customers",
      "REFERENCED_TABLE_NAME": "customers",
      "REFERENCED_COLUMN_NAME": "id"
    }
  ]
}
```

---

#### 9. ดู Indexes

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb",
  "table": "users"
}
```

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "table": "users",
  "indexes": [
    {
      "Table": "users",
      "Non_unique": "0",
      "Key_name": "PRIMARY",
      "Column_name": "id",
      "Index_type": "BTREE"
    },
    {
      "Table": "users",
      "Non_unique": "1",
      "Key_name": "idx_email",
      "Column_name": "email",
      "Index_type": "BTREE"
    }
  ]
}
```

---

#### 10. คำนวณ Aggregate Data

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  },
  "database": "mydb",
  "table": "orders",
  "aggregations": {
    "total_revenue": "SUM(total)",
    "order_count": "COUNT(*)",
    "avg_order": "AVG(total)",
    "max_order": "MAX(total)"
  },
  "groupBy": "status",
  "where": "created_at >= '2026-01-01'"
}
```

**ผลลัพธ์:**
```json
{
  "database": "mydb",
  "table": "orders",
  "aggregations": {
    "total_revenue": "SUM(total)",
    "order_count": "COUNT(*)",
    "avg_order": "AVG(total)",
    "max_order": "MAX(total)"
  },
  "group_by": "status",
  "results": [
    {"status": "completed", "total_revenue": "50000.00", "order_count": "120", "avg_order": "416.67", "max_order": "2500.00"},
    {"status": "pending", "total_revenue": "15000.00", "order_count": "35", "avg_order": "428.57", "max_order": "1800.00"}
  ]
}
```

---

#### 11. เปลี่ยนฐานข้อมูล

```json
{
  "config": {
    "user": "root",
    "password": "",
    "host": "127.0.0.1",
    "port": 3306
  },
  "database": "production_db"
}
```

**ผลลัพธ์:**
```json
{
  "success": true,
  "message": "เปลี่ยนไปใช้ฐานข้อมูล 'production_db' สำเร็จ",
  "current_database": "production_db",
  "config": {
    "user": "root",
    "host": "127.0.0.1",
    "port": 3306,
    "database": "production_db"
  }
}
```

> 💡 **เคล็ดลับ:** ใช้เครื่องมือนี้เพื่อสลับระหว่างฐานข้อมูลโดยไม่ต้องตั้งค่าการเชื่อมต่อใหม่ Config ที่ส่งคืนมาสามารถนำไปใช้ต่อได้

---

## 🔒 หมายเหตุความปลอดภัย

- **อ่านอย่างเดียว:** เซิร์ฟเวอร์อนุญาตเฉพาะคำสั่ง SELECT, SHOW, DESCRIBE, และ EXPLAIN
- **ไม่มีการเขียน:** INSERT, UPDATE, DELETE, DROP ถูกปิดกั้น
- **การจัดการรหัสผ่าน:** รหัสผ่านถูกส่งตรงไปยัง MySQL driver (ไม่มีการบันทึก)
- **SQL Injection:** WHERE clause ที่ผู้ใช้ป้อนจะถูกส่งไปตามเดิม; ควรตรวจสอบใน MCP client

---

## 🐛 การแก้ปัญหา

### ไม่พบ mysql.exe

ตรวจสอบว่าติดตั้ง MySQL/MariaDB และอยู่ใน PATH หรือตั้งค่า environment variable `MYSQL_HOME`:

```powershell
# สำหรับ XAMPP
setx MYSQL_HOME "C:\xampp\mysql\bin"
```

### การเชื่อมต่อถูกปฏิเสธ

- ตรวจสอบว่าเซอร์วิส MySQL/MariaDB กำลังทำงาน
- ตรวจสอบ host และ port ใน config
- ตรวจสอบว่า firewall อนุญาตการเชื่อมต่อ

### Access denied

- ตรวจสอบชื่อผู้ใช้และรหัสผ่าน
- ตรวจสอบสิทธิ์ของผู้ใช้สำหรับฐานข้อมูลเป้าหมาย

---

## 📝 เอกสารที่เกี่ยวข้อง

- [README หลัก](./README.md) - ภาพรวมโครงการ
- [ตัวจัดการฐานข้อมูล](./README.db.md) - เครื่องมือ CLI สำรอง/คืนค่าข้อมูล
- [เครื่องมือบีบอัด](./README.archiver.md) - เครื่องมือบีบอัด ZIP รองรับ Unicode

---

*พัฒนาโดย tps2015gh พร้อมความช่วยเหลือจาก Gemini AI - 2026*
