# Database MCP Server

🌐 **Languages:** [English] | [**ภาษาไทย (Thai)**](./README.db_mcp.th.md)

A **Model Context Protocol (MCP)** server for MariaDB/MySQL database operations. This tool provides AI assistants with direct database access capabilities for querying, schema inspection, and complex data analysis.

---

## 🚀 Features

### Core Capabilities

| Tool | Description |
|------|-------------|
| `list_databases` | List all databases accessible to the current user |
| `list_tables` | List all tables in a specific database |
| `get_table_schema` | Get detailed column information for a table |
| `query_data` | Execute SELECT queries and return results |
| `get_table_stats` | Get table statistics (rows, encoding, data length) |
| `get_database_info` | Get database version, variables, and status |
| `search_data` | Complex queries with WHERE, ORDER BY, LIMIT, OFFSET |
| `get_foreign_keys` | Get foreign key relationships for a table |
| `get_indexes` | Get index information for a table |
| `aggregate_data` | Perform COUNT, SUM, AVG, MIN, MAX with GROUP BY |
| `switch_database` | **Switch to a different database** for subsequent operations |

---

## 📦 Installation

### Build from Source

```powershell
go build -tags mcp_server -o db_mcp_server.exe db_mcp_server.go
```

### Requirements

- Go 1.23+
- MySQL/MariaDB client binaries (`mysql.exe`) in PATH or XAMPP installation
- MCP-compatible client (Claude Desktop, VS Code, etc.)

---

## 🔧 Configuration

### MCP Client Setup

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

## 📖 Usage Examples

### Connection Config Format

All tools require a `config` parameter in JSON format:

```json
{
  "user": "root",
  "password": "your_password",
  "host": "127.0.0.1",
  "port": 3306,
  "database": "mydb"
}
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `user` | ✅ Yes | - | Database username |
| `password` | ❌ No | `""` | Database password |
| `host` | ❌ No | `127.0.0.1` | Database host |
| `port` | ❌ No | `3306` | Database port |
| `database` | ❌ No | - | Default database (optional for `list_databases`) |

---

### Tool Examples

#### 0. Complete Workflow: List and Switch Database

**Step 1: List all available databases**

```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

**Response:**
```json
{
  "databases": ["mydb", "testdb", "production", "shop_db"],
  "count": 4
}
```

**Step 2: User selects a database** (e.g., `production`)

**Step 3: Switch to the selected database**

```json
{
  "config": {
    "user": "root",
    "password": ""
  },
  "database": "production"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Successfully switched to database 'production'",
  "current_database": "production",
  "config": {
    "user": "root",
    "host": "127.0.0.1",
    "port": 3306,
    "database": "production"
  }
}
```

> 💡 **Workflow Tip:** The agent should first call `list_databases` to show available databases to the user, then use `switch_database` to connect to the user's choice. The returned config can be cached for subsequent operations.

---

#### 1. List Databases (Individual Call)

```json
{
  "config": {
    "user": "root",
    "password": ""
  }
}
```

**Response:**
```json
{
  "databases": ["mydb", "testdb", "production"],
  "count": 3
}
```

---

#### 2. List Tables

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

**Response:**
```json
{
  "database": "mydb",
  "tables": ["users", "orders", "products"],
  "count": 3
}
```

---

#### 3. Get Table Schema

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

**Response:**
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

#### 4. Query Data (SELECT)

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

**Response:**
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

> ⚠️ **Security:** Only SELECT, SHOW, DESCRIBE, and EXPLAIN queries are allowed.

---

#### 5. Get Table Statistics

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

**Response:**
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

#### 6. Get Database Info

```json
{
  "config": {
    "user": "root",
    "password": "",
    "database": "mydb"
  }
}
```

**Response:**
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

#### 7. Search Data (Complex Query)

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

**Response:**
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

#### 8. Get Foreign Keys

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

**Response:**
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

#### 9. Get Indexes

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

**Response:**
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

#### 10. Aggregate Data

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

**Response:**
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

#### 11. Switch Database

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

**Response:**
```json
{
  "success": true,
  "message": "Successfully switched to database 'production_db'",
  "current_database": "production_db",
  "config": {
    "user": "root",
    "host": "127.0.0.1",
    "port": 3306,
    "database": "production_db"
  }
}
```

> 💡 **Tip:** Use this tool to switch between databases without re-configuring the entire connection. The returned config can be reused for subsequent operations.

---

## 🔒 Security Notes

- **Read-Only:** The server only allows SELECT, SHOW, DESCRIBE, and EXPLAIN queries
- **No Write Operations:** INSERT, UPDATE, DELETE, DROP are blocked
- **Credential Handling:** Passwords are passed directly to MySQL driver (not logged)
- **SQL Injection:** User-provided WHERE clauses are passed as-is; validate inputs in your MCP client

---

## 🐛 Troubleshooting

### mysql.exe not found

Ensure MySQL/MariaDB is installed and in PATH, or set `MYSQL_HOME` environment variable:

```powershell
# For XAMPP
setx MYSQL_HOME "C:\xampp\mysql\bin"
```

### Connection refused

- Verify MySQL/MariaDB service is running
- Check host and port in config
- Ensure firewall allows connections

### Access denied

- Verify username and password
- Check user privileges for the target database

---

## 📝 Related Documentation

- [Main README](./README.md) - Project overview
- [DB Manager](./README.db.md) - Interactive CLI backup/restore tool
- [Archiver](./README.archiver.md) - Unicode ZIP archiver

---

*Developed by tps2015gh with Gemini AI Assistance - 2026*
