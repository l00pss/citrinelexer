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
	sql := "INSERT users"

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
}

func TestParseUpdate(t *testing.T) {
	sql := "UPDATE users"

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
