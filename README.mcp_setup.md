# MCP Server Setup Guide

🌐 **Languages:** [**English**] | [ภาษาไทย (Thai)](./README.mcp_setup.th.md)

Complete setup guide for Database MCP Server with **Qwen Code** and **Gemini CLI**.

---

## 📋 Prerequisites

- Go 1.23+
- MySQL/MariaDB installed (XAMPP or standalone)
- Git configured

---

## 🚀 Quick Setup

### Step 1: Build the MCP Server

```powershell
cd C:\dev\abcd_setup_server
go get github.com/go-sql-driver/mysql
go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
```

### Step 2: Configure MCP Client

#### For Qwen Code

Create or edit the global MCP config file at `%USERPROFILE%\.qwen\mcp.json`:

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**Commands:**
```powershell
# Create directory if not exists
if not exist "%USERPROFILE%\.qwen" mkdir "%USERPROFILE%\.qwen"

# Create config file (run in PowerShell)
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

#### For Gemini CLI

Create or edit the global MCP config file at `%APPDATA%\gemini-cli\mcp.json`:

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

**Commands:**
```powershell
# Create directory if not exists
if not exist "$env:APPDATA\gemini-cli" mkdir "$env:APPDATA\gemini-cli"

# Create config file
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

## ✅ Verify Setup

### In Qwen Code

1. **Restart Qwen Code** (close and reopen)
2. Run `/mcp list` - you should see:
   ```
   database - C:\dev\abcd_setup_server\db_mcp_server.exe
   ```
3. Test with a database operation:
   ```
   /mcp database list_databases {"config": {"user": "root", "password": ""}}
   ```

### In Gemini CLI

1. **Restart Gemini CLI** (close and reopen)
2. Run `/mcp list` - you should see the database server
3. Test database operations similarly

---

## 🔧 Available MCP Tools

| Tool | Description |
|------|-------------|
| `list_databases` | List all accessible databases |
| `list_tables` | List tables in a database |
| `get_table_schema` | Get column details for a table |
| `query_data` | Execute SELECT queries |
| `get_table_stats` | Get table statistics |
| `get_database_info` | Get database version and variables |
| `search_data` | Complex search with filters and pagination |
| `get_foreign_keys` | Get foreign key relationships |
| `get_indexes` | Get index information |
| `aggregate_data` | Perform COUNT, SUM, AVG, MIN, MAX |
| `switch_database` | Switch to a different database |

---

## 📝 Usage Examples

### List Databases
```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

### List Tables
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

### Query Data
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

### Switch Database (Recommended Workflow)
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

## 🐛 Troubleshooting

### "No MCP servers configured"

- Ensure config file exists at correct location
- Restart the application after creating/editing config
- Verify path to `db_mcp_server.exe` is correct

### "mysql.exe not found"

- Ensure MySQL/MariaDB is installed
- Add to PATH or set `MYSQL_HOME`:
  ```powershell
  setx MYSQL_HOME "C:\xampp\mysql\bin"
  ```

### "Access denied"

- Verify MySQL username and password
- Check user privileges for target database

### "SQL unknown driver"

- Rebuild with MySQL driver:
  ```powershell
  go get github.com/go-sql-driver/mysql
  go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
  ```
- Restart Qwen Code / Gemini CLI

---

## 📁 Config File Locations

| Application | Config Path |
|-------------|-------------|
| Qwen Code | `%USERPROFILE%\.qwen\mcp.json` |
| Gemini CLI | `%APPDATA%\gemini-cli\mcp.json` |

---

## 🔒 Security Notes

- Only SELECT, SHOW, DESCRIBE, EXPLAIN queries allowed
- No INSERT, UPDATE, DELETE, DROP operations
- Passwords passed directly to MySQL driver (not logged)
- Validate inputs in your MCP client

---

## 📖 Related Documentation

- [Main README](./README.md) - Project overview
- [DB Manager](./README.db.md) - Interactive CLI backup/restore tool
- [Database MCP Server](./README.db_mcp.md) - Detailed tool documentation
- [Archiver](./README.archiver.md) - Unicode ZIP archiver

---

*Developed by tps2015gh - 2026*
