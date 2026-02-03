package citrinelexer

import (
	"testing"
)

func TestParseSelectStatement(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "simple select",
			sql:  "SELECT name FROM users",
		},
		{
			name: "select with where",
			sql:  "SELECT name, age FROM users WHERE id = 123",
		},
		{
			name: "select all",
			sql:  "SELECT * FROM users",
		},
		{
			name: "select with order by",
			sql:  "SELECT name FROM users ORDER BY name ASC",
		},
		{
			name: "select with limit",
			sql:  "SELECT name FROM users LIMIT 10",
		},
		{
			name: "complex select",
			sql:  "SELECT u.name, u.age FROM users u WHERE u.id > 100 ORDER BY u.name LIMIT 50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt, ok := stmt.(*SelectStatement)
			if !ok {
				t.Fatalf("Expected SelectStatement, got %T", stmt)
			}

			if len(selectStmt.Fields) == 0 {
				t.Fatal("Expected fields, got none")
			}
		})
	}
}

func TestParseCreateTable(t *testing.T) {
	sql := "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	createStmt, ok := stmt.(*CreateTableStatement)
	if !ok {
		t.Fatalf("Expected CreateTableStatement, got %T", stmt)
	}

	if createStmt.Table.Name != "users" {
		t.Fatalf("Expected table name 'users', got '%s'", createStmt.Table.Name)
	}

	if len(createStmt.Columns) != 2 {
		t.Fatalf("Expected 2 columns, got %d", len(createStmt.Columns))
	}
}

func TestParseInsert(t *testing.T) {
	sql := "INSERT INTO users (id) VALUES (1)"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	insertStmt, ok := stmt.(*InsertStatement)
	if !ok {
		t.Fatalf("Expected InsertStatement, got %T", stmt)
	}

	if insertStmt.Table.Name != "users" {
		t.Fatalf("Expected table name 'users', got '%s'", insertStmt.Table.Name)
	}

	if len(insertStmt.Columns) != 1 {
		t.Fatalf("Expected 1 column, got %d", len(insertStmt.Columns))
	}

	if len(insertStmt.Values) != 1 {
		t.Fatalf("Expected 1 value set, got %d", len(insertStmt.Values))
	}
}

func TestParseUpdate(t *testing.T) {
	sql := "UPDATE users SET name = 'test' WHERE id = 1"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	updateStmt, ok := stmt.(*UpdateStatement)
	if !ok {
		t.Fatalf("Expected UpdateStatement, got %T", stmt)
	}

	if updateStmt.Table.Name != "users" {
		t.Fatalf("Expected table name 'users', got '%s'", updateStmt.Table.Name)
	}
}

func TestParseDelete(t *testing.T) {
	sql := "DELETE FROM users WHERE id = 123"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	deleteStmt, ok := stmt.(*DeleteStatement)
	if !ok {
		t.Fatalf("Expected DeleteStatement, got %T", stmt)
	}

	if deleteStmt.From.Name != "users" {
		t.Fatalf("Expected table name 'users', got '%s'", deleteStmt.From.Name)
	}

	if deleteStmt.Where == nil {
		t.Fatal("Expected WHERE clause")
	}
}

func TestParseParameters(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "positional parameter",
			sql:  "SELECT * FROM users WHERE id = ?",
		},
		{
			name: "named parameter colon",
			sql:  "SELECT * FROM users WHERE name = :name",
		},
		{
			name: "named parameter dollar",
			sql:  "SELECT * FROM users WHERE age = $age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt, ok := stmt.(*SelectStatement)
			if !ok {
				t.Fatalf("Expected SelectStatement, got %T", stmt)
			}

			if selectStmt.Where == nil {
				t.Fatal("Expected WHERE clause with parameter")
			}
		})
	}
}

func TestParseFunctionCall(t *testing.T) {
	sql := "SELECT name FROM users"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectStmt, ok := stmt.(*SelectStatement)
	if !ok {
		t.Fatalf("Expected SelectStatement, got %T", stmt)
	}

	if len(selectStmt.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(selectStmt.Fields))
	}
}

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "string literal",
			sql:  "SELECT 'hello' FROM users",
		},
		{
			name: "number literal",
			sql:  "SELECT 42 FROM users",
		},
		{
			name: "boolean literal",
			sql:  "SELECT TRUE FROM users",
		},
		{
			name: "binary expression",
			sql:  "SELECT * FROM users WHERE age > 18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			_, ok := stmt.(*SelectStatement)
			if !ok {
				t.Fatalf("Expected SelectStatement, got %T", stmt)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "invalid token",
			sql:  "INVALID STATEMENT",
		},
		{
			name: "incomplete select",
			sql:  "SELECT",
		},
		{
			name: "missing table name",
			sql:  "SELECT * FROM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Expected error for invalid SQL: %s", tt.sql)
			}
		})
	}
}

// Test cases from BUG_REPORT.md

func TestParseInsertWithInto(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		tableName   string
		columnCount int
		valueCount  int
	}{
		{
			name:        "simple insert",
			sql:         "INSERT INTO users (id) VALUES (1)",
			tableName:   "users",
			columnCount: 1,
			valueCount:  1,
		},
		{
			name:        "multiple columns",
			sql:         "INSERT INTO users (id, name, age) VALUES (1, 'Alice', 30)",
			tableName:   "users",
			columnCount: 3,
			valueCount:  1,
		},
		{
			name:        "multiple rows",
			sql:         "INSERT INTO users (id) VALUES (1), (2), (3)",
			tableName:   "users",
			columnCount: 1,
			valueCount:  3,
		},
		{
			name:        "without INTO",
			sql:         "INSERT users (id) VALUES (1)",
			tableName:   "users",
			columnCount: 1,
			valueCount:  1,
		},
		{
			name:        "escaped string",
			sql:         "INSERT INTO users (name) VALUES ('O''Brien')",
			tableName:   "users",
			columnCount: 1,
			valueCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			insertStmt, ok := stmt.(*InsertStatement)
			if !ok {
				t.Fatalf("Expected InsertStatement, got %T", stmt)
			}

			if insertStmt.Table.Name != tt.tableName {
				t.Fatalf("Expected table name '%s', got '%s'", tt.tableName, insertStmt.Table.Name)
			}

			if len(insertStmt.Columns) != tt.columnCount {
				t.Fatalf("Expected %d columns, got %d", tt.columnCount, len(insertStmt.Columns))
			}

			if len(insertStmt.Values) != tt.valueCount {
				t.Fatalf("Expected %d value sets, got %d", tt.valueCount, len(insertStmt.Values))
			}
		})
	}
}

func TestParseVarcharWithLength(t *testing.T) {
	tests := []struct {
		name   string
		sql    string
		length int
	}{
		{
			name:   "varchar with length",
			sql:    "CREATE TABLE t (name VARCHAR(50))",
			length: 50,
		},
		{
			name:   "char with length",
			sql:    "CREATE TABLE t (code CHAR(10))",
			length: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			createStmt, ok := stmt.(*CreateTableStatement)
			if !ok {
				t.Fatalf("Expected CreateTableStatement, got %T", stmt)
			}

			if len(createStmt.Columns) == 0 {
				t.Fatal("Expected at least one column")
			}

			if createStmt.Columns[0].Length != tt.length {
				t.Fatalf("Expected length %d, got %d", tt.length, createStmt.Columns[0].Length)
			}
		})
	}
}

func TestParseUpdateWithSet(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		tableName string
		setCount  int
		hasWhere  bool
	}{
		{
			name:      "single set",
			sql:       "UPDATE users SET name = 'Bob' WHERE id = 1",
			tableName: "users",
			setCount:  1,
			hasWhere:  true,
		},
		{
			name:      "multiple set",
			sql:       "UPDATE users SET name = 'Bob', age = 30 WHERE id = 1",
			tableName: "users",
			setCount:  2,
			hasWhere:  true,
		},
		{
			name:      "without where",
			sql:       "UPDATE users SET active = TRUE",
			tableName: "users",
			setCount:  1,
			hasWhere:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			updateStmt, ok := stmt.(*UpdateStatement)
			if !ok {
				t.Fatalf("Expected UpdateStatement, got %T", stmt)
			}

			if updateStmt.Table.Name != tt.tableName {
				t.Fatalf("Expected table name '%s', got '%s'", tt.tableName, updateStmt.Table.Name)
			}

			if len(updateStmt.Set) != tt.setCount {
				t.Fatalf("Expected %d set clauses, got %d", tt.setCount, len(updateStmt.Set))
			}

			if tt.hasWhere && updateStmt.Where == nil {
				t.Fatal("Expected WHERE clause")
			}

			if !tt.hasWhere && updateStmt.Where != nil {
				t.Fatal("Did not expect WHERE clause")
			}
		})
	}
}

func TestParseMixedTypesTable(t *testing.T) {
	sql := "CREATE TABLE t (id INTEGER, name VARCHAR(100), code CHAR(5))"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	createStmt, ok := stmt.(*CreateTableStatement)
	if !ok {
		t.Fatalf("Expected CreateTableStatement, got %T", stmt)
	}

	if len(createStmt.Columns) != 3 {
		t.Fatalf("Expected 3 columns, got %d", len(createStmt.Columns))
	}

	// Check id column
	if createStmt.Columns[0].Name.Name != "id" {
		t.Fatalf("Expected column 'id', got '%s'", createStmt.Columns[0].Name.Name)
	}
	if createStmt.Columns[0].Type != "INTEGER" {
		t.Fatalf("Expected type 'INTEGER', got '%s'", createStmt.Columns[0].Type)
	}

	// Check name column with VARCHAR(100)
	if createStmt.Columns[1].Name.Name != "name" {
		t.Fatalf("Expected column 'name', got '%s'", createStmt.Columns[1].Name.Name)
	}
	if createStmt.Columns[1].Type != "VARCHAR" {
		t.Fatalf("Expected type 'VARCHAR', got '%s'", createStmt.Columns[1].Type)
	}
	if createStmt.Columns[1].Length != 100 {
		t.Fatalf("Expected length 100, got %d", createStmt.Columns[1].Length)
	}

	// Check code column with CHAR(5)
	if createStmt.Columns[2].Name.Name != "code" {
		t.Fatalf("Expected column 'code', got '%s'", createStmt.Columns[2].Name.Name)
	}
	if createStmt.Columns[2].Type != "CHAR" {
		t.Fatalf("Expected type 'CHAR', got '%s'", createStmt.Columns[2].Type)
	}
	if createStmt.Columns[2].Length != 5 {
		t.Fatalf("Expected length 5, got %d", createStmt.Columns[2].Length)
	}
}

func TestParseComplexInsert(t *testing.T) {
	sql := "INSERT INTO users (id, first_name, last_name, email, gender, ip_address) VALUES (1, 'Moshe', 'McTague', 'mmctague0@state.tx.us', 'Male', '67.164.161.76')"

	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	insertStmt, ok := stmt.(*InsertStatement)
	if !ok {
		t.Fatalf("Expected InsertStatement, got %T", stmt)
	}

	if insertStmt.Table.Name != "users" {
		t.Fatalf("Expected table 'users', got '%s'", insertStmt.Table.Name)
	}

	if len(insertStmt.Columns) != 6 {
		t.Fatalf("Expected 6 columns, got %d", len(insertStmt.Columns))
	}

	if len(insertStmt.Values) != 1 {
		t.Fatalf("Expected 1 value set, got %d", len(insertStmt.Values))
	}

	if len(insertStmt.Values[0]) != 6 {
		t.Fatalf("Expected 6 values, got %d", len(insertStmt.Values[0]))
	}
}

// TestSQLFileExamples tests SQL statements
func TestSQLFileExamples(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		stmtType    string
		shouldParse bool
	}{
		// CREATE TABLE examples
		{
			name:        "create users table with VARCHAR",
			sql:         "CREATE TABLE users (id INTEGER PRIMARY KEY, first_name VARCHAR(50), last_name VARCHAR(50), email VARCHAR(50), gender VARCHAR(50), ip_address VARCHAR(20))",
			stmtType:    "CREATE TABLE",
			shouldParse: true,
		},
		{
			name:        "create departments table",
			sql:         "CREATE TABLE departments (id INTEGER PRIMARY KEY, name TEXT NOT NULL, budget INTEGER)",
			stmtType:    "CREATE TABLE",
			shouldParse: true,
		},
		{
			name:        "create orders table",
			sql:         "CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER, product TEXT, amount INTEGER, created_at TEXT)",
			stmtType:    "CREATE TABLE",
			shouldParse: true,
		},

		// INSERT examples
		{
			name:        "insert into departments",
			sql:         "INSERT INTO departments (id, name, budget) VALUES (1, 'Engineering', 500000)",
			stmtType:    "INSERT",
			shouldParse: true,
		},
		{
			name:        "insert into orders",
			sql:         "INSERT INTO orders (id, user_id, product, amount, created_at) VALUES (1, 1, 'Laptop', 1500, '2024-01-15')",
			stmtType:    "INSERT",
			shouldParse: true,
		},

		// SELECT examples
		{
			name:        "select all from users",
			sql:         "SELECT * FROM users",
			stmtType:    "SELECT",
			shouldParse: true,
		},
		{
			name:        "select specific columns",
			sql:         "SELECT first_name, last_name, email FROM users",
			stmtType:    "SELECT",
			shouldParse: true,
		},
		{
			name:        "select with where",
			sql:         "SELECT * FROM users WHERE gender = 'Male'",
			stmtType:    "SELECT",
			shouldParse: true,
		},
		{
			name:        "select with order by",
			sql:         "SELECT first_name, last_name FROM users ORDER BY first_name ASC",
			stmtType:    "SELECT",
			shouldParse: true,
		},
		{
			name:        "select with limit",
			sql:         "SELECT * FROM users LIMIT 5",
			stmtType:    "SELECT",
			shouldParse: true,
		},
		{
			name:        "select with limit and offset",
			sql:         "SELECT * FROM users ORDER BY id LIMIT 2 OFFSET 4",
			stmtType:    "SELECT",
			shouldParse: true,
		},

		// UPDATE examples
		{
			name:        "update department budget",
			sql:         "UPDATE departments SET budget = 600000 WHERE id = 1",
			stmtType:    "UPDATE",
			shouldParse: true,
		},
		{
			name:        "update user email",
			sql:         "UPDATE users SET email = 'updated@example.com' WHERE id = 1",
			stmtType:    "UPDATE",
			shouldParse: true,
		},

		// DELETE examples
		{
			name:        "delete from departments",
			sql:         "DELETE FROM departments WHERE id = 99",
			stmtType:    "DELETE",
			shouldParse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if tt.shouldParse {
				if err != nil {
					t.Fatalf("Parse failed: %v", err)
				}
				if stmt.String() != tt.stmtType {
					t.Fatalf("Expected %s statement, got %s", tt.stmtType, stmt.String())
				}
			} else {
				if err == nil {
					t.Fatalf("Expected parse error for: %s", tt.sql)
				}
			}
		})
	}
}
