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

func TestParseWhereClause(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		operator string
	}{
		{"equals string", "SELECT * FROM t WHERE name = 'Alice'", "="},
		{"equals number", "SELECT * FROM t WHERE id = 1", "="},
		{"not equals", "SELECT * FROM t WHERE id != 1", "!="},
		{"greater than", "SELECT * FROM t WHERE age > 18", ">"},
		{"less than", "SELECT * FROM t WHERE age < 65", "<"},
		{"greater or equal", "SELECT * FROM t WHERE age >= 18", ">="},
		{"less or equal", "SELECT * FROM t WHERE age <= 65", "<="},
		{"LIKE pattern", "SELECT * FROM t WHERE name LIKE 'A%'", "LIKE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.Where == nil {
				t.Fatal("Expected WHERE clause")
			}

			be, ok := selectStmt.Where.(*BinaryExpression)
			if !ok {
				t.Fatalf("Expected BinaryExpression, got %T", selectStmt.Where)
			}

			if be.Operator != tt.operator {
				t.Fatalf("Expected operator %q, got %q", tt.operator, be.Operator)
			}
		})
	}
}

func TestDebugWhereClause(t *testing.T) {
	sql := "SELECT * FROM users WHERE gender = 'Male'"
	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	sel := stmt.(*SelectStatement)

	be, ok := sel.Where.(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected BinaryExpression, got %T", sel.Where)
	}

	// Check left side is Identifier
	left, ok := be.Left.(*Identifier)
	if !ok {
		t.Fatalf("Expected left to be Identifier, got %T", be.Left)
	}
	if left.Name != "gender" {
		t.Fatalf("Expected left.Name='gender', got %q", left.Name)
	}

	// Check operator
	if be.Operator != "=" {
		t.Fatalf("Expected operator '=', got %q", be.Operator)
	}

	// Check right side is StringLiteral
	right, ok := be.Right.(*StringLiteral)
	if !ok {
		t.Fatalf("Expected right to be StringLiteral, got %T", be.Right)
	}
	if right.Value != "Male" {
		t.Fatalf("Expected right.Value='Male', got %q", right.Value)
	}
}

func TestTableAlias(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		table string
		alias string
	}{
		{"alias without AS", "SELECT * FROM users u", "users", "u"},
		{"alias with AS", "SELECT * FROM users AS u", "users", "u"},
		{"alias for orders", "SELECT * FROM orders o", "orders", "o"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.From == nil {
				t.Fatal("Expected FROM clause")
			}

			if selectStmt.From.Name.Name != tt.table {
				t.Fatalf("Expected table %q, got %q", tt.table, selectStmt.From.Name.Name)
			}

			if selectStmt.From.Alias == nil {
				t.Fatalf("Expected alias %q, got nil", tt.alias)
			}

			if selectStmt.From.Alias.Name != tt.alias {
				t.Fatalf("Expected alias %q, got %q", tt.alias, selectStmt.From.Alias.Name)
			}
		})
	}
}

func TestCountStar(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"count star", "SELECT COUNT(*) FROM users"},
		{"count star with where", "SELECT COUNT(*) FROM users WHERE gender = 'Male'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if len(selectStmt.Fields) != 1 {
				t.Fatalf("Expected 1 field, got %d", len(selectStmt.Fields))
			}

			fn, ok := selectStmt.Fields[0].(*FunctionCall)
			if !ok {
				t.Fatalf("Expected FunctionCall, got %T", selectStmt.Fields[0])
			}

			if fn.Name != "COUNT" {
				t.Fatalf("Expected function name COUNT, got %s", fn.Name)
			}

			if len(fn.Args) != 1 {
				t.Fatalf("Expected 1 argument, got %d", len(fn.Args))
			}

			// Check that arg is * identifier
			ident, ok := fn.Args[0].(*Identifier)
			if !ok {
				t.Fatalf("Expected Identifier for *, got %T", fn.Args[0])
			}
			if ident.Name != "*" {
				t.Fatalf("Expected * argument, got %s", ident.Name)
			}
		})
	}
}

func TestQualifiedAsterisk(t *testing.T) {
	sql := "SELECT users.* FROM users"
	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectStmt := stmt.(*SelectStatement)
	if len(selectStmt.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(selectStmt.Fields))
	}

	qa, ok := selectStmt.Fields[0].(*QualifiedAsterisk)
	if !ok {
		t.Fatalf("Expected QualifiedAsterisk, got %T", selectStmt.Fields[0])
	}

	if qa.Table != "users" {
		t.Fatalf("Expected table 'users', got %q", qa.Table)
	}
}

func TestQualifiedIdentifier(t *testing.T) {
	sql := "SELECT u.first_name FROM users u"
	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectStmt := stmt.(*SelectStatement)
	if len(selectStmt.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(selectStmt.Fields))
	}

	qi, ok := selectStmt.Fields[0].(*QualifiedIdentifier)
	if !ok {
		t.Fatalf("Expected QualifiedIdentifier, got %T", selectStmt.Fields[0])
	}

	if qi.Table != "u" {
		t.Fatalf("Expected table 'u', got %q", qi.Table)
	}

	if qi.Column != "first_name" {
		t.Fatalf("Expected column 'first_name', got %q", qi.Column)
	}
}

func TestTransactionStatements(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		stmtType string
	}{
		{"begin", "BEGIN", "BEGIN"},
		{"begin transaction", "BEGIN TRANSACTION", "BEGIN"},
		{"commit", "COMMIT", "COMMIT"},
		{"rollback", "ROLLBACK", "ROLLBACK"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if stmt.String() != tt.stmtType {
				t.Fatalf("Expected %s, got %s", tt.stmtType, stmt.String())
			}
		})
	}
}

func TestColumnAlias(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		table string
	}{
		{"count with alias", "SELECT COUNT(*) as total_users FROM users", "users"},
		{"count with AS alias", "SELECT COUNT(*) AS total FROM users", "users"},
		{"sum with alias", "SELECT SUM(amount) as total_revenue FROM orders", "orders"},
		{"simple column alias", "SELECT name as n FROM users", "users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.From == nil {
				t.Fatal("FROM clause is nil - AS alias broke FROM parsing")
			}

			if selectStmt.From.Name.Name != tt.table {
				t.Fatalf("Expected table %q, got %q", tt.table, selectStmt.From.Name.Name)
			}
		})
	}
}

func TestJoinClause(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		joinType  string
		joinTable string
		joinAlias string
	}{
		{
			name:      "inner join with alias",
			sql:       "SELECT u.id, o.amount FROM users u INNER JOIN orders o ON u.id = o.user_id",
			joinType:  "INNER",
			joinTable: "orders",
			joinAlias: "o",
		},
		{
			name:      "left join",
			sql:       "SELECT u.name, d.name FROM users u LEFT JOIN departments d ON u.dept_id = d.id",
			joinType:  "LEFT",
			joinTable: "departments",
			joinAlias: "d",
		},
		{
			name:      "simple join",
			sql:       "SELECT * FROM users JOIN orders ON users.id = orders.user_id",
			joinType:  "INNER",
			joinTable: "orders",
			joinAlias: "",
		},
		{
			name:      "right join",
			sql:       "SELECT * FROM users u RIGHT JOIN orders o ON u.id = o.user_id",
			joinType:  "RIGHT",
			joinTable: "orders",
			joinAlias: "o",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.From == nil {
				t.Fatal("FROM clause is nil")
			}

			if len(selectStmt.Joins) == 0 {
				t.Fatal("Expected at least one JOIN clause")
			}

			join := selectStmt.Joins[0]
			if join.Type != tt.joinType {
				t.Fatalf("Expected join type %q, got %q", tt.joinType, join.Type)
			}

			if join.Table.Name.Name != tt.joinTable {
				t.Fatalf("Expected join table %q, got %q", tt.joinTable, join.Table.Name.Name)
			}

			if tt.joinAlias != "" {
				if join.Table.Alias == nil {
					t.Fatalf("Expected join alias %q, got nil", tt.joinAlias)
				}
				if join.Table.Alias.Name != tt.joinAlias {
					t.Fatalf("Expected join alias %q, got %q", tt.joinAlias, join.Table.Alias.Name)
				}
			}

			if join.Condition == nil {
				t.Fatal("Expected ON condition")
			}
		})
	}
}

func TestKeywordAsAlias(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		alias string
	}{
		{"count as alias", "SELECT COUNT(*) as count FROM users", "count"},
		{"sum as alias", "SELECT SUM(amount) as sum FROM orders", "sum"},
		{"avg as alias", "SELECT AVG(price) as avg FROM items", "avg"},
		{"min as alias", "SELECT MIN(age) as min FROM users", "min"},
		{"max as alias", "SELECT MAX(age) as max FROM users", "max"},
		{"text as alias", "SELECT name as text FROM users", "text"},
		{"integer as alias", "SELECT id as integer FROM users", "integer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if len(selectStmt.Fields) == 0 {
				t.Fatal("Expected at least one field")
			}

			aliased, ok := selectStmt.Fields[0].(*AliasedExpression)
			if !ok {
				t.Fatalf("Expected AliasedExpression, got %T", selectStmt.Fields[0])
			}

			if aliased.Alias != tt.alias {
				t.Fatalf("Expected alias %q, got %q", tt.alias, aliased.Alias)
			}
		})
	}
}

func TestNullLiteral(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"select null", "SELECT NULL"},
		{"insert with null", "INSERT INTO users (name, email) VALUES ('John', NULL)"},
		{"select null alias", "SELECT NULL as empty FROM users"},
		{"null in where", "SELECT * FROM users WHERE email = NULL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
		})
	}

	// Test that NullLiteral is correctly parsed
	t.Run("verify null literal type", func(t *testing.T) {
		stmt, err := Parse("SELECT NULL")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		selectStmt := stmt.(*SelectStatement)
		if len(selectStmt.Fields) == 0 {
			t.Fatal("Expected at least one field")
		}

		nullLit, ok := selectStmt.Fields[0].(*NullLiteral)
		if !ok {
			t.Fatalf("Expected NullLiteral, got %T", selectStmt.Fields[0])
		}

		if nullLit.String() != "NULL" {
			t.Fatalf("Expected 'NULL', got %q", nullLit.String())
		}
	})
}

func TestCreateIndex(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		indexName   string
		tableName   string
		unique      bool
		ifNotExists bool
		columns     []string
		directions  []string
	}{
		{
			name:       "simple create index",
			sql:        "CREATE INDEX idx_name ON users (name)",
			indexName:  "idx_name",
			tableName:  "users",
			unique:     false,
			columns:    []string{"name"},
			directions: []string{""},
		},
		{
			name:       "create unique index",
			sql:        "CREATE UNIQUE INDEX idx_email ON users (email)",
			indexName:  "idx_email",
			tableName:  "users",
			unique:     true,
			columns:    []string{"email"},
			directions: []string{""},
		},
		{
			name:        "create index if not exists",
			sql:         "CREATE INDEX IF NOT EXISTS idx_age ON users (age)",
			indexName:   "idx_age",
			tableName:   "users",
			ifNotExists: true,
			columns:     []string{"age"},
			directions:  []string{""},
		},
		{
			name:        "create unique index if not exists",
			sql:         "CREATE UNIQUE INDEX IF NOT EXISTS idx_phone ON users (phone)",
			indexName:   "idx_phone",
			tableName:   "users",
			unique:      true,
			ifNotExists: true,
			columns:     []string{"phone"},
			directions:  []string{""},
		},
		{
			name:       "multi column index",
			sql:        "CREATE INDEX idx_name_age ON users (name, age)",
			indexName:  "idx_name_age",
			tableName:  "users",
			columns:    []string{"name", "age"},
			directions: []string{"", ""},
		},
		{
			name:       "index with asc direction",
			sql:        "CREATE INDEX idx_name_asc ON users (name ASC)",
			indexName:  "idx_name_asc",
			tableName:  "users",
			columns:    []string{"name"},
			directions: []string{"ASC"},
		},
		{
			name:       "index with desc direction",
			sql:        "CREATE INDEX idx_name_desc ON users (name DESC)",
			indexName:  "idx_name_desc",
			tableName:  "users",
			columns:    []string{"name"},
			directions: []string{"DESC"},
		},
		{
			name:       "multi column index with directions",
			sql:        "CREATE INDEX idx_composite ON users (name ASC, age DESC)",
			indexName:  "idx_composite",
			tableName:  "users",
			columns:    []string{"name", "age"},
			directions: []string{"ASC", "DESC"},
		},
		{
			name:       "multi column index with mixed directions",
			sql:        "CREATE INDEX idx_mixed ON orders (customer_id, order_date DESC, amount ASC)",
			indexName:  "idx_mixed",
			tableName:  "orders",
			columns:    []string{"customer_id", "order_date", "amount"},
			directions: []string{"", "DESC", "ASC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			createIdx, ok := stmt.(*CreateIndexStatement)
			if !ok {
				t.Fatalf("Expected CreateIndexStatement, got %T", stmt)
			}

			if createIdx.Name.Name != tt.indexName {
				t.Errorf("Expected index name %q, got %q", tt.indexName, createIdx.Name.Name)
			}

			if createIdx.Table.Name != tt.tableName {
				t.Errorf("Expected table name %q, got %q", tt.tableName, createIdx.Table.Name)
			}

			if createIdx.Unique != tt.unique {
				t.Errorf("Expected unique=%v, got %v", tt.unique, createIdx.Unique)
			}

			if createIdx.IfNotExists != tt.ifNotExists {
				t.Errorf("Expected ifNotExists=%v, got %v", tt.ifNotExists, createIdx.IfNotExists)
			}

			if len(createIdx.Columns) != len(tt.columns) {
				t.Fatalf("Expected %d columns, got %d", len(tt.columns), len(createIdx.Columns))
			}

			for i, col := range createIdx.Columns {
				if col.Column.Name != tt.columns[i] {
					t.Errorf("Column %d: expected name %q, got %q", i, tt.columns[i], col.Column.Name)
				}
				if col.Direction != tt.directions[i] {
					t.Errorf("Column %d: expected direction %q, got %q", i, tt.directions[i], col.Direction)
				}
			}
		})
	}

	// Test String() method
	t.Run("string method", func(t *testing.T) {
		stmt, err := Parse("CREATE INDEX idx_test ON users (name)")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		createIdx := stmt.(*CreateIndexStatement)
		if createIdx.String() != "CREATE INDEX" {
			t.Errorf("Expected 'CREATE INDEX', got %q", createIdx.String())
		}

		stmt2, err := Parse("CREATE UNIQUE INDEX idx_test ON users (name)")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		createIdx2 := stmt2.(*CreateIndexStatement)
		if createIdx2.String() != "CREATE UNIQUE INDEX" {
			t.Errorf("Expected 'CREATE UNIQUE INDEX', got %q", createIdx2.String())
		}
	})
}

func TestDropIndex(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		indexName string
		ifExists  bool
	}{
		{
			name:      "simple drop index",
			sql:       "DROP INDEX idx_name",
			indexName: "idx_name",
			ifExists:  false,
		},
		{
			name:      "drop index if exists",
			sql:       "DROP INDEX IF EXISTS idx_email",
			indexName: "idx_email",
			ifExists:  true,
		},
		{
			name:      "drop index with underscore name",
			sql:       "DROP INDEX idx_users_name_age",
			indexName: "idx_users_name_age",
			ifExists:  false,
		},
		{
			name:      "drop index if exists complex name",
			sql:       "DROP INDEX IF EXISTS idx_orders_customer_date",
			indexName: "idx_orders_customer_date",
			ifExists:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			dropIdx, ok := stmt.(*DropIndexStatement)
			if !ok {
				t.Fatalf("Expected DropIndexStatement, got %T", stmt)
			}

			if dropIdx.Name.Name != tt.indexName {
				t.Errorf("Expected index name %q, got %q", tt.indexName, dropIdx.Name.Name)
			}

			if dropIdx.IfExists != tt.ifExists {
				t.Errorf("Expected ifExists=%v, got %v", tt.ifExists, dropIdx.IfExists)
			}
		})
	}

	// Test String() method
	t.Run("string method", func(t *testing.T) {
		stmt, err := Parse("DROP INDEX idx_test")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		dropIdx := stmt.(*DropIndexStatement)
		if dropIdx.String() != "DROP INDEX" {
			t.Errorf("Expected 'DROP INDEX', got %q", dropIdx.String())
		}
	})
}

func TestSelectDistinct(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		distinct bool
	}{
		{
			name:     "select without distinct",
			sql:      "SELECT name FROM users",
			distinct: false,
		},
		{
			name:     "select distinct single column",
			sql:      "SELECT DISTINCT name FROM users",
			distinct: true,
		},
		{
			name:     "select distinct multiple columns",
			sql:      "SELECT DISTINCT name, email FROM users",
			distinct: true,
		},
		{
			name:     "select distinct with where",
			sql:      "SELECT DISTINCT status FROM orders WHERE amount > 100",
			distinct: true,
		},
		{
			name:     "select distinct all columns",
			sql:      "SELECT DISTINCT * FROM users",
			distinct: true,
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

			if selectStmt.Distinct != tt.distinct {
				t.Errorf("Expected distinct=%v, got %v", tt.distinct, selectStmt.Distinct)
			}
		})
	}
}

func TestLikeExpression(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		operator string
	}{
		{
			name:     "simple like",
			sql:      "SELECT * FROM users WHERE name LIKE '%John%'",
			operator: "LIKE",
		},
		{
			name:     "like with prefix",
			sql:      "SELECT * FROM users WHERE email LIKE 'admin%'",
			operator: "LIKE",
		},
		{
			name:     "like with suffix",
			sql:      "SELECT * FROM users WHERE name LIKE '%son'",
			operator: "LIKE",
		},
		{
			name:     "not like",
			sql:      "SELECT * FROM users WHERE name NOT LIKE '%test%'",
			operator: "NOT LIKE",
		},
		{
			name:     "glob pattern",
			sql:      "SELECT * FROM files WHERE path GLOB '*.txt'",
			operator: "GLOB",
		},
		{
			name:     "not glob",
			sql:      "SELECT * FROM files WHERE path NOT GLOB '*.tmp'",
			operator: "NOT GLOB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.Where == nil {
				t.Fatal("Expected WHERE clause")
			}

			binExpr, ok := selectStmt.Where.(*BinaryExpression)
			if !ok {
				t.Fatalf("Expected BinaryExpression, got %T", selectStmt.Where)
			}

			if binExpr.Operator != tt.operator {
				t.Errorf("Expected operator %q, got %q", tt.operator, binExpr.Operator)
			}
		})
	}
}

func TestInExpression(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		not        bool
		valueCount int
	}{
		{
			name:       "in with numbers",
			sql:        "SELECT * FROM users WHERE id IN (1, 2, 3)",
			not:        false,
			valueCount: 3,
		},
		{
			name:       "in with single value",
			sql:        "SELECT * FROM users WHERE status IN (1)",
			not:        false,
			valueCount: 1,
		},
		{
			name:       "in with strings",
			sql:        "SELECT * FROM users WHERE status IN ('active', 'pending', 'approved')",
			not:        false,
			valueCount: 3,
		},
		{
			name:       "not in with numbers",
			sql:        "SELECT * FROM users WHERE id NOT IN (1, 2, 3, 4, 5)",
			not:        true,
			valueCount: 5,
		},
		{
			name:       "not in with strings",
			sql:        "SELECT * FROM orders WHERE status NOT IN ('cancelled', 'refunded')",
			not:        true,
			valueCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.Where == nil {
				t.Fatal("Expected WHERE clause")
			}

			inExpr, ok := selectStmt.Where.(*InExpression)
			if !ok {
				t.Fatalf("Expected InExpression, got %T", selectStmt.Where)
			}

			if inExpr.Not != tt.not {
				t.Errorf("Expected not=%v, got %v", tt.not, inExpr.Not)
			}

			if len(inExpr.Values) != tt.valueCount {
				t.Errorf("Expected %d values, got %d", tt.valueCount, len(inExpr.Values))
			}
		})
	}

	// Test String() method
	t.Run("string method in", func(t *testing.T) {
		stmt, _ := Parse("SELECT * FROM users WHERE id IN (1, 2)")
		inExpr := stmt.(*SelectStatement).Where.(*InExpression)
		if inExpr.String() != "id IN (...)" {
			t.Errorf("Expected 'id IN (...)', got %q", inExpr.String())
		}
	})

	t.Run("string method not in", func(t *testing.T) {
		stmt, _ := Parse("SELECT * FROM users WHERE id NOT IN (1, 2)")
		inExpr := stmt.(*SelectStatement).Where.(*InExpression)
		if inExpr.String() != "id NOT IN (...)" {
			t.Errorf("Expected 'id NOT IN (...)', got %q", inExpr.String())
		}
	})
}

func TestBetweenExpression(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		not  bool
	}{
		{
			name: "between with numbers",
			sql:  "SELECT * FROM users WHERE age BETWEEN 20 AND 30",
			not:  false,
		},
		{
			name: "between with large range",
			sql:  "SELECT * FROM products WHERE price BETWEEN 100 AND 1000",
			not:  false,
		},
		{
			name: "between with strings",
			sql:  "SELECT * FROM users WHERE name BETWEEN 'A' AND 'M'",
			not:  false,
		},
		{
			name: "not between with numbers",
			sql:  "SELECT * FROM users WHERE age NOT BETWEEN 18 AND 65",
			not:  true,
		},
		{
			name: "not between with dates",
			sql:  "SELECT * FROM orders WHERE created_at NOT BETWEEN '2024-01-01' AND '2024-12-31'",
			not:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if selectStmt.Where == nil {
				t.Fatal("Expected WHERE clause")
			}

			betweenExpr, ok := selectStmt.Where.(*BetweenExpression)
			if !ok {
				t.Fatalf("Expected BetweenExpression, got %T", selectStmt.Where)
			}

			if betweenExpr.Not != tt.not {
				t.Errorf("Expected not=%v, got %v", tt.not, betweenExpr.Not)
			}

			if betweenExpr.Low == nil {
				t.Error("Expected Low value")
			}

			if betweenExpr.High == nil {
				t.Error("Expected High value")
			}
		})
	}

	// Test String() method
	t.Run("string method between", func(t *testing.T) {
		stmt, _ := Parse("SELECT * FROM users WHERE age BETWEEN 20 AND 30")
		betweenExpr := stmt.(*SelectStatement).Where.(*BetweenExpression)
		expected := "age BETWEEN 20 AND 30"
		if betweenExpr.String() != expected {
			t.Errorf("Expected %q, got %q", expected, betweenExpr.String())
		}
	})

	t.Run("string method not between", func(t *testing.T) {
		stmt, _ := Parse("SELECT * FROM users WHERE age NOT BETWEEN 10 AND 20")
		betweenExpr := stmt.(*SelectStatement).Where.(*BetweenExpression)
		expected := "age NOT BETWEEN 10 AND 20"
		if betweenExpr.String() != expected {
			t.Errorf("Expected %q, got %q", expected, betweenExpr.String())
		}
	})
}

func TestRightJoin(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		joinType string
	}{
		{
			name:     "right join",
			sql:      "SELECT * FROM orders RIGHT JOIN users ON orders.user_id = users.id",
			joinType: "RIGHT",
		},
		{
			name:     "right outer join",
			sql:      "SELECT * FROM orders RIGHT OUTER JOIN users ON orders.user_id = users.id",
			joinType: "RIGHT",
		},
		{
			name:     "right join with alias",
			sql:      "SELECT * FROM orders o RIGHT JOIN users u ON o.user_id = u.id",
			joinType: "RIGHT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if len(selectStmt.Joins) == 0 {
				t.Fatal("Expected at least one join")
			}

			join := selectStmt.Joins[0]
			if join.Type != tt.joinType {
				t.Errorf("Expected join type %q, got %q", tt.joinType, join.Type)
			}
		})
	}
}

func TestFullOuterJoin(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		joinType string
	}{
		{
			name:     "full join",
			sql:      "SELECT * FROM users FULL JOIN orders ON users.id = orders.user_id",
			joinType: "FULL",
		},
		{
			name:     "full outer join",
			sql:      "SELECT * FROM users FULL OUTER JOIN orders ON users.id = orders.user_id",
			joinType: "FULL",
		},
		{
			name:     "full join with alias",
			sql:      "SELECT * FROM users u FULL JOIN orders o ON u.id = o.user_id",
			joinType: "FULL",
		},
		{
			name:     "full outer join with where",
			sql:      "SELECT * FROM users FULL OUTER JOIN orders ON users.id = orders.user_id WHERE orders.amount > 100",
			joinType: "FULL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			selectStmt := stmt.(*SelectStatement)
			if len(selectStmt.Joins) == 0 {
				t.Fatal("Expected at least one join")
			}

			join := selectStmt.Joins[0]
			if join.Type != tt.joinType {
				t.Errorf("Expected join type %q, got %q", tt.joinType, join.Type)
			}
		})
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "simple AND",
			sql:  "SELECT * FROM users WHERE age > 18 AND status = 1",
		},
		{
			name: "simple OR",
			sql:  "SELECT * FROM users WHERE status = 1 OR status = 2",
		},
		{
			name: "AND with OR",
			sql:  "SELECT * FROM users WHERE age > 18 AND status = 1 OR role = 'admin'",
		},
		{
			name: "multiple AND",
			sql:  "SELECT * FROM users WHERE age > 18 AND status = 1 AND role = 'user'",
		},
		{
			name: "multiple OR",
			sql:  "SELECT * FROM users WHERE status = 1 OR status = 2 OR status = 3",
		},
		{
			name: "BETWEEN with AND",
			sql:  "SELECT * FROM users WHERE age BETWEEN 20 AND 30 AND status = 1",
		},
		{
			name: "IN with AND",
			sql:  "SELECT * FROM users WHERE id IN (1, 2, 3) AND status = 1",
		},
		{
			name: "IN with OR",
			sql:  "SELECT * FROM users WHERE id IN (1, 2) OR id IN (5, 6)",
		},
		{
			name: "BETWEEN with OR",
			sql:  "SELECT * FROM users WHERE age BETWEEN 20 AND 30 OR age BETWEEN 50 AND 60",
		},
		{
			name: "complex BETWEEN AND combination",
			sql:  "SELECT * FROM orders WHERE amount BETWEEN 100 AND 500 AND status = 'active' AND user_id = 1",
		},
		{
			name: "LIKE with AND",
			sql:  "SELECT * FROM users WHERE name LIKE '%john%' AND email LIKE '%@gmail.com'",
		},
		{
			name: "NOT IN with AND",
			sql:  "SELECT * FROM users WHERE id NOT IN (1, 2) AND status = 1",
		},
		{
			name: "NOT BETWEEN with AND",
			sql:  "SELECT * FROM users WHERE age NOT BETWEEN 0 AND 18 AND status = 'active'",
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
				t.Fatal("Expected WHERE clause")
			}
		})
	}

	// Verify BETWEEN AND doesn't conflict with logical AND
	t.Run("verify BETWEEN AND structure", func(t *testing.T) {
		stmt, err := Parse("SELECT * FROM users WHERE age BETWEEN 20 AND 30 AND status = 1")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		selectStmt := stmt.(*SelectStatement)

		// The top level should be a BinaryExpression with AND operator
		binExpr, ok := selectStmt.Where.(*BinaryExpression)
		if !ok {
			t.Fatalf("Expected BinaryExpression at top level, got %T", selectStmt.Where)
		}

		if binExpr.Operator != "AND" {
			t.Errorf("Expected top level operator 'AND', got %q", binExpr.Operator)
		}

		// Left side should be BetweenExpression
		_, ok = binExpr.Left.(*BetweenExpression)
		if !ok {
			t.Errorf("Expected BetweenExpression on left side, got %T", binExpr.Left)
		}

		// Right side should be BinaryExpression (status = 1)
		rightBin, ok := binExpr.Right.(*BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression on right side, got %T", binExpr.Right)
		}

		if rightBin.Operator != "=" {
			t.Errorf("Expected right side operator '=', got %q", rightBin.Operator)
		}
	})

	// Verify IN AND combination
	t.Run("verify IN AND structure", func(t *testing.T) {
		stmt, err := Parse("SELECT * FROM users WHERE id IN (1, 2, 3) AND name = 'test'")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		selectStmt := stmt.(*SelectStatement)

		binExpr, ok := selectStmt.Where.(*BinaryExpression)
		if !ok {
			t.Fatalf("Expected BinaryExpression at top level, got %T", selectStmt.Where)
		}

		if binExpr.Operator != "AND" {
			t.Errorf("Expected top level operator 'AND', got %q", binExpr.Operator)
		}

		// Left side should be InExpression
		inExpr, ok := binExpr.Left.(*InExpression)
		if !ok {
			t.Errorf("Expected InExpression on left side, got %T", binExpr.Left)
		}

		if len(inExpr.Values) != 3 {
			t.Errorf("Expected 3 values in IN expression, got %d", len(inExpr.Values))
		}
	})

	// Verify OR precedence
	t.Run("verify OR precedence", func(t *testing.T) {
		stmt, err := Parse("SELECT * FROM users WHERE a = 1 AND b = 2 OR c = 3")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		selectStmt := stmt.(*SelectStatement)

		// Top level should be OR (lower precedence)
		binExpr, ok := selectStmt.Where.(*BinaryExpression)
		if !ok {
			t.Fatalf("Expected BinaryExpression at top level, got %T", selectStmt.Where)
		}

		if binExpr.Operator != "OR" {
			t.Errorf("Expected top level operator 'OR', got %q", binExpr.Operator)
		}

		// Left side of OR should be AND expression
		leftAnd, ok := binExpr.Left.(*BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression on left side of OR, got %T", binExpr.Left)
		}

		if leftAnd.Operator != "AND" {
			t.Errorf("Expected left side operator 'AND', got %q", leftAnd.Operator)
		}
	})
}
