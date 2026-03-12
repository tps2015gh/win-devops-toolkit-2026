//go:build mcp_server

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	mysqlPath string
)

type DBConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database,omitempty"`
}

func main() {
	setupMySQLPath()

	s := server.NewMCPServer(
		"database-mcp",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)

	// Register tools
	s.AddTool(mcp.NewTool(
		"list_databases",
		mcp.WithDescription("List all databases accessible to the current user"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
	), handleListDatabases)

	s.AddTool(mcp.NewTool(
		"list_tables",
		mcp.WithDescription("List all tables in a specific database"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
	), handleListTables)

	s.AddTool(mcp.NewTool(
		"get_table_schema",
		mcp.WithDescription("Get detailed schema/columns info for a table"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
		mcp.WithString("table", mcp.Description("Table name"), mcp.Required()),
	), handleGetTableSchema)

	s.AddTool(mcp.NewTool(
		"query_data",
		mcp.WithDescription("Execute a SELECT query and return results"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("query", mcp.Description("SQL SELECT query to execute"), mcp.Required()),
	), handleQueryData)

	s.AddTool(mcp.NewTool(
		"get_table_stats",
		mcp.WithDescription("Get table statistics including row count, encoding, data length"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
	), handleGetTableStats)

	s.AddTool(mcp.NewTool(
		"get_database_info",
		mcp.WithDescription("Get comprehensive database information including version, variables, status"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
	), handleGetDatabaseInfo)

	s.AddTool(mcp.NewTool(
		"search_data",
		mcp.WithDescription("Search data across tables with complex filters, sorting, and pagination"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
		mcp.WithString("table", mcp.Description("Table name to search"), mcp.Required()),
		mcp.WithString("columns", mcp.Description("Comma-separated list of columns to return (default: *)")),
		mcp.WithString("where", mcp.Description("WHERE clause conditions (without WHERE keyword)")),
		mcp.WithString("orderBy", mcp.Description("ORDER BY clause (without ORDER BY keyword)")),
		mcp.WithNumber("limit", mcp.Description("LIMIT for results (default: 100)")),
		mcp.WithNumber("offset", mcp.Description("OFFSET for pagination (default: 0)")),
	), handleSearchData)

	s.AddTool(mcp.NewTool(
		"get_foreign_keys",
		mcp.WithDescription("Get foreign key relationships for a table"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
		mcp.WithString("table", mcp.Description("Table name"), mcp.Required()),
	), handleGetForeignKeys)

	s.AddTool(mcp.NewTool(
		"get_indexes",
		mcp.WithDescription("Get index information for a table"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
		mcp.WithString("table", mcp.Description("Table name"), mcp.Required()),
	), handleGetIndexes)

	s.AddTool(mcp.NewTool(
		"aggregate_data",
		mcp.WithDescription("Perform aggregate operations (COUNT, SUM, AVG, MIN, MAX) with GROUP BY"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name"), mcp.Required()),
		mcp.WithString("table", mcp.Description("Table name"), mcp.Required()),
		mcp.WithString("aggregations", mcp.Description("JSON object of aggregation functions, e.g., {\"total\": \"SUM(amount)\", \"count\": \"COUNT(*)\"}"), mcp.Required()),
		mcp.WithString("groupBy", mcp.Description("Column(s) to group by")),
		mcp.WithString("where", mcp.Description("WHERE clause conditions")),
	), handleAggregateData)

	s.AddTool(mcp.NewTool(
		"switch_database",
		mcp.WithDescription("Switch to a different database context for subsequent operations"),
		mcp.WithString("config", mcp.Description("Database connection config as JSON"), mcp.Required()),
		mcp.WithString("database", mcp.Description("Database name to switch to"), mcp.Required()),
	), handleSwitchDatabase)

	ctx := context.Background()
	server.ServeStdio(s)
	_ = ctx
}

func setupMySQLPath() {
	mysqlPath, _ = exec.LookPath("mysql.exe")

	if mysqlPath == "" {
		commonPaths := []string{
			"C:\\xampp\\mysql\\bin",
			"D:\\xampp\\mysql\\bin",
			"C:\\xampp_v8_1_25\\mysql\\bin",
		}
		for _, p := range commonPaths {
			m := filepath.Join(p, "mysql.exe")
			if _, err := os.Stat(m); err == nil {
				mysqlPath = m
				break
			}
		}
	}

	if mysqlPath == "" {
		fmt.Fprintln(os.Stderr, "mysql.exe not found. Please ensure MySQL/MariaDB is installed and in PATH")
		os.Exit(1)
	}
}

func parseConfig(configJSON string) (*DBConfig, error) {
	var config DBConfig
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("invalid config JSON: %w", err)
	}
	if config.User == "" {
		config.User = "root"
	}
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	return &config, nil
}

func connectDB(config *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// First try to connect without database to allow listing databases
	if config.Database == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&loc=Local",
			config.User,
			config.Password,
			config.Host,
			config.Port,
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return db, nil
}

func handleListDatabases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, ok := args["config"].(string)
	if !ok {
		return nil, fmt.Errorf("config is required")
	}

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		// Skip system databases
		if dbName != "information_schema" && dbName != "performance_schema" && dbName != "mysql" && dbName != "phpmyadmin" {
			databases = append(databases, dbName)
		}
	}

	result, _ := json.Marshal(map[string]interface{}{
		"databases": databases,
		"count":     len(databases),
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleListTables(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database": database,
		"tables":   tables,
		"count":    len(tables),
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTableSchema(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)
	table, _ := args["table"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("DESCRIBE `%s`", table))
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var schema []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		schema = append(schema, row)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database": database,
		"table":    table,
		"schema":   schema,
		"columns":  len(schema),
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleQueryData(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	query, _ := args["query"].(string)

	// Security: only allow SELECT queries
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	if !strings.HasPrefix(upperQuery, "SELECT") && !strings.HasPrefix(upperQuery, "SHOW") && !strings.HasPrefix(upperQuery, "DESCRIBE") && !strings.HasPrefix(upperQuery, "EXPLAIN") {
		return nil, fmt.Errorf("only SELECT, SHOW, DESCRIBE, and EXPLAIN queries are allowed for security")
	}

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var results []map[string]interface{}
	var rowCount int

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		results = append(results, row)
		rowCount++
	}

	result, _ := json.Marshal(map[string]interface{}{
		"columns": columns,
		"rows":    results,
		"count":   rowCount,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTableStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT 
		TABLE_NAME,
		TABLE_ROWS,
		TABLE_COLLATION,
		DATA_LENGTH,
		INDEX_LENGTH,
		DATA_FREE,
		AUTO_INCREMENT,
		CREATE_TIME,
		UPDATE_TIME
	FROM information_schema.TABLES 
	WHERE TABLE_SCHEMA = ?`

	rows, err := db.Query(query, database)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var stats []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		stats = append(stats, row)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database": database,
		"stats":    stats,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetDatabaseInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	info := make(map[string]interface{})

	// Get version
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err == nil {
		info["version"] = version
	}

	// Get current user
	var currentUser string
	err = db.QueryRow("SELECT CURRENT_USER()").Scan(&currentUser)
	if err == nil {
		info["current_user"] = currentUser
	}

	// Get database name
	var dbName string
	err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err == nil {
		info["current_database"] = dbName
	}

	// Get server status
	var status string
	err = db.QueryRow("SHOW STATUS LIKE 'Uptime'").Scan(&dbName, &status)
	if err == nil {
		info["uptime_seconds"] = status
	}

	// Get important variables
	vars := make(map[string]string)
	varRows, err := db.Query("SHOW VARIABLES WHERE Variable_name IN ('max_connections', 'character_set_server', 'collation_server', 'version_comment')")
	if err == nil {
		defer varRows.Close()
		for varRows.Next() {
			var name, value string
			if err := varRows.Scan(&name, &value); err == nil {
				vars[name] = value
			}
		}
		info["variables"] = vars
	}

	// Get connection count
	var connCount string
	err = db.QueryRow("SHOW STATUS LIKE 'Threads_connected'").Scan(&dbName, &connCount)
	if err == nil {
		info["active_connections"] = connCount
	}

	result, _ := json.Marshal(info)
	return mcp.NewToolResultText(string(result)), nil
}

func handleSearchData(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)
	table, _ := args["table"].(string)
	columns, _ := args["columns"].(string)
	where, _ := args["where"].(string)
	orderBy, _ := args["orderBy"].(string)
	limitFloat, _ := args["limit"].(float64)
	offsetFloat, _ := args["offset"].(float64)

	if limitFloat == 0 {
		limitFloat = 100
	}
	limit := int(limitFloat)
	offset := int(offsetFloat)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Build query
	selectCols := "*"
	if columns != "" {
		selectCols = columns
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", selectCols, table)

	if where != "" {
		query += " WHERE " + where
	}

	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columnsList, _ := rows.Columns()
	var results []map[string]interface{}
	var rowCount int

	for rows.Next() {
		values := make([]interface{}, len(columnsList))
		valuePtrs := make([]interface{}, len(columnsList))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columnsList {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		results = append(results, row)
		rowCount++
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database": database,
		"table":    table,
		"columns":  columnsList,
		"rows":     results,
		"count":    rowCount,
		"limit":    limit,
		"offset":   offset,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetForeignKeys(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)
	table, _ := args["table"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT 
		COLUMN_NAME,
		CONSTRAINT_NAME,
		REFERENCED_TABLE_NAME,
		REFERENCED_COLUMN_NAME
	FROM information_schema.KEY_COLUMN_USAGE
	WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND REFERENCED_COLUMN_NAME IS NOT NULL`

	rows, err := db.Query(query, database, table)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var fks []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		fks = append(fks, row)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database":     database,
		"table":        table,
		"foreign_keys": fks,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetIndexes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)
	table, _ := args["table"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SHOW INDEX FROM `%s`", table))
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var indexes []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		indexes = append(indexes, row)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database": database,
		"table":    table,
		"indexes":  indexes,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleAggregateData(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)
	table, _ := args["table"].(string)
	aggregationsJSON, _ := args["aggregations"].(string)
	groupBy, _ := args["groupBy"].(string)
	where, _ := args["where"].(string)

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}
	config.Database = database

	// Parse aggregations
	var aggMap map[string]string
	if err := json.Unmarshal([]byte(aggregationsJSON), &aggMap); err != nil {
		return nil, fmt.Errorf("invalid aggregations JSON: %w", err)
	}

	// Build aggregation part of query
	var aggParts []string
	for alias, expr := range aggMap {
		aggParts = append(aggParts, fmt.Sprintf("%s AS %s", expr, alias))
	}
	aggClause := strings.Join(aggParts, ", ")

	// Build query
	query := fmt.Sprintf("SELECT %s FROM `%s`", aggClause, table)

	if where != "" {
		query += " WHERE " + where
	}

	if groupBy != "" {
		query += " GROUP BY " + groupBy
	}

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		results = append(results, row)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"database":    database,
		"table":       table,
		"aggregations": aggMap,
		"group_by":    groupBy,
		"results":     results,
	})

	return mcp.NewToolResultText(string(result)), nil
}

func handleSwitchDatabase(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	configJSON, _ := args["config"].(string)
	database, _ := args["database"].(string)

	if database == "" {
		return nil, fmt.Errorf("database name is required")
	}

	config, err := parseConfig(configJSON)
	if err != nil {
		return nil, err
	}

	// Test connection to the new database
	config.Database = database
	db, err := connectDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to switch to database '%s': %w", database, err)
	}
	defer db.Close()

	// Verify database exists and is accessible
	var dbName string
	err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to verify database switch: %w", err)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Successfully switched to database '%s'", database),
		"current_database": dbName,
		"config": map[string]interface{}{
			"user":     config.User,
			"host":     config.Host,
			"port":     config.Port,
			"database": database,
		},
	})

	return mcp.NewToolResultText(string(result)), nil
}
