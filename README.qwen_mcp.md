# Using Database MCP Server with Qwen Code CLI

This guide explains how to configure and use the Database MCP Server with Qwen Code CLI for AI-driven database operations.

---

## 📋 Prerequisites

1. **Built MCP Server**: Ensure `db_mcp_server.exe` is compiled and exists
2. **MySQL/MariaDB**: Database server must be running
3. **Qwen Code CLI**: Installed and configured

---

## 🔧 Configuration

### Step 1: Create MCP Configuration

#### Option A: Project-Level Configuration (Recommended)

Create `.qwen/mcp.json` in your project root:

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

#### Option B: Global Configuration

**Windows:** Create or edit `%APPDATA%\Qwen\mcp.json`

**macOS/Linux:** Create or edit `~/.qwen/mcp.json`

```json
{
  "mcpServers": {
    "database": {
      "command": "C:\\dev\\abcd_setup_server\\db_mcp_server.exe"
    }
  }
}
```

### Step 2: Update Path for Your Environment

Replace the path with your actual location:

| Your Location | Config Path |
|---------------|-------------|
| Current project | `C:\\dev\\abcd_setup_server\\db_mcp_server.exe` |
| Custom folder | `D:\\tools\\db_mcp_server.exe` |
| In PATH | `db_mcp_server.exe` |

### Step 3: Restart Qwen Code

After adding the configuration, **restart Qwen Code CLI** for changes to take effect.

---

## ✅ Verify Connection

### Test 1: Check Available Tools

Ask Qwen:
```
What MCP tools do you have available?
```

Expected response should include:
- `list_databases`
- `list_tables`
- `get_table_schema`
- `query_data`
- `get_table_stats`
- `get_database_info`
- `search_data`
- `get_foreign_keys`
- `get_indexes`
- `aggregate_data`
- `switch_database`

### Test 2: List Databases

Ask Qwen:
```
List all available databases
```

Expected response:
```json
{
  "databases": ["mydb", "testdb", "production"],
  "count": 3
}
```

---

## 📖 Usage Examples

### 1. List Available Databases

**Prompt:**
```
Show me all databases I can access
```

**What Qwen does:**
- Calls `list_databases` tool
- Returns list of user databases (excludes system DBs)

---

### 2. Switch to a Database

**Prompt:**
```
Switch to the 'production' database
```

**What Qwen does:**
- Calls `switch_database` tool
- Returns connection config with database locked
- Uses this config for subsequent operations

---

### 3. List Tables

**Prompt:**
```
Show all tables in the current database
```

**What Qwen does:**
- Calls `list_tables` tool
- Returns table names and count

---

### 4. View Table Schema

**Prompt:**
```
What columns does the 'users' table have?
```

**What Qwen does:**
- Calls `get_table_schema` tool
- Returns field names, types, null constraints, keys, defaults

---

### 5. Query Data

**Prompt:**
```
Show me the first 10 users from the users table
```

**What Qwen does:**
- Calls `query_data` tool with `SELECT * FROM users LIMIT 10`
- Returns rows with column names

---

### 6. Complex Search with Filters

**Prompt:**
```
Find all orders where status is 'pending' and total is greater than 1000, 
sorted by created date descending, limit 50 results
```

**What Qwen does:**
- Calls `search_data` tool with:
  - `where`: "status = 'pending' AND total > 1000"
  - `orderBy`: "created_at DESC"
  - `limit`: 50

---

### 7. Get Table Statistics

**Prompt:**
```
Show me statistics for all tables in the database
```

**What Qwen does:**
- Calls `get_table_stats` tool
- Returns row counts, collation, data length, index length, auto-increment values

---

### 8. Get Database Information

**Prompt:**
```
What version of MySQL/MariaDB is running?
```

**What Qwen does:**
- Calls `get_database_info` tool
- Returns version, current user, uptime, variables, active connections

---

### 9. View Foreign Keys

**Prompt:**
```
What foreign keys does the 'orders' table have?
```

**What Qwen does:**
- Calls `get_foreign_keys` tool
- Returns constraint names, referenced tables/columns

---

### 10. View Indexes

**Prompt:**
```
Show all indexes on the 'users' table
```

**What Qwen does:**
- Calls `get_indexes` tool
- Returns index names, columns, uniqueness, index type

---

### 11. Aggregate Data

**Prompt:**
```
Calculate total revenue, order count, and average order value grouped by status
for orders created this year
```

**What Qwen does:**
- Calls `aggregate_data` tool with:
  - `aggregations`: `{"total_revenue": "SUM(total)", "order_count": "COUNT(*)", "avg_order": "AVG(total)"}`
  - `groupBy`: "status"
  - `where`: "created_at >= '2026-01-01'"

---

## 🔒 Security Notes

### Read-Only Access

The MCP server **only allows**:
- `SELECT` queries
- `SHOW` commands
- `DESCRIBE` commands
- `EXPLAIN` commands

**Blocked operations:**
- `INSERT`, `UPDATE`, `DELETE`
- `DROP`, `TRUNCATE`, `ALTER`
- Any write operations

### Credential Handling

- Credentials are passed directly to MySQL driver
- Passwords are **not logged** or stored
- Use least-privilege database users when possible

---

## 🐛 Troubleshooting

### Issue: MCP Server Not Found

**Error:** `Failed to start MCP server: db_mcp_server.exe not found`

**Solution:**
1. Verify the path in `.qwen/mcp.json` is correct
2. Use absolute path with escaped backslashes (Windows):
   ```json
   "command": "C:\\\\dev\\\\abcd_setup_server\\\\db_mcp_server.exe"
   ```
3. Ensure the executable exists and is not blocked by antivirus

---

### Issue: mysql.exe Not Found

**Error:** `mysql.exe not found. Please ensure MySQL/MariaDB is installed`

**Solution:**
1. Add MySQL to PATH:
   ```powershell
   setx PATH "%PATH%;C:\xampp\mysql\bin"
   ```
2. Or set `MYSQL_HOME`:
   ```powershell
   setx MYSQL_HOME "C:\xampp\mysql\bin"
   ```
3. Restart Qwen Code after changes

---

### Issue: Connection Refused

**Error:** `connection failed: dial tcp 127.0.0.1:3306: connectex: No connection could be made`

**Solution:**
1. Verify MySQL/MariaDB service is running:
   ```powershell
   Get-Service -Name MySQL*
   ```
2. Check host and port in your connection
3. Ensure firewall allows connections on port 3306

---

### Issue: Access Denied

**Error:** `Access denied for user 'root'@'localhost'`

**Solution:**
1. Verify username and password
2. Check user privileges:
   ```sql
   SHOW GRANTS FOR 'root'@'localhost';
   ```
3. Grant necessary permissions if needed

---

### Issue: No Tools Available

**Symptom:** Qwen doesn't list database tools

**Solution:**
1. Check `.qwen/mcp.json` syntax is valid JSON
2. Restart Qwen Code CLI completely
3. Check Qwen logs for MCP connection errors
4. Verify MCP server runs manually:
   ```powershell
   C:\dev\abcd_setup_server\db_mcp_server.exe
   ```

---

## 📝 Example Workflow

### Complete Database Exploration

```
User: List all databases
Qwen: [Calls list_databases] Found: mydb, testdb, production

User: Switch to production database
Qwen: [Calls switch_database] Successfully switched to 'production'

User: Show all tables
Qwen: [Calls list_tables] Found: users, orders, products, categories

User: What columns does the orders table have?
Qwen: [Calls get_table_schema] Returns: id, customer_id, total, status, created_at, updated_at

User: Show pending orders over $500
Qwen: [Calls search_data] Returns filtered results

User: Calculate total revenue by status
Qwen: [Calls aggregate_data] Returns grouped sums
```

---

## 📚 Related Documentation

- [English MCP Documentation](./README.db_mcp.md) - Full tool reference
- [Thai MCP Documentation](./README.db_mcp.th.md) - คู่มือภาษาไทย
- [Main README](./README.md) - Project overview

---

*Developed by tps2015gh with Gemini AI Assistance - 2026*
