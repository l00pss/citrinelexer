# Citrine Lexer

<p align="center">
  <img src="logo.png" alt="Citrine Lexer Logo" width="400">
</p>

<div align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/l00pss/citrinelexer"><img src="https://goreportcard.com/badge/github.com/l00pss/citrinelexer" alt="Go Report Card"></a>
  <a href="https://github.com/l00pss/citrinelexer/stargazers"><img src="https://img.shields.io/github/stars/l00pss/citrinelexer?style=flat&logo=github" alt="GitHub Stars"></a>
</div>

<div align="center">
  <a href="https://www.buymeacoffee.com/l00pss" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
</div>

<br>

A fast and flexible SQL lexer and parser library. Perfect for building your own database engine or SQL parser!

## What does it do?

Breaks SQL code into tokens and builds Abstract Syntax Trees (AST). Complete lexing and parsing solution.

**Input:** `SELECT name FROM users WHERE id = 123;`
**Lexing:** `SELECT`, `name`, `FROM`, `users`, `WHERE`, `id`, `=`, `123`, `;`
**Parsing:** AST with SelectStatement node containing fields, table reference, and WHERE clause

## Features

-  **Fast**: Processes ~450,000 SQL queries per second
-  **Flexible**: Supports different SQL dialects
-  **Efficient**: Zero allocation optimizations
-  **Comprehensive**: Recognizes 100+ SQL keywords
-  **Battle-tested**: Extensive test suite
-  **AST Support**: Full parsing with `go/ast` interface compatibility
-  **Dual Mode**: Use lexer alone or with parser

## Architecture

```
SQL String → Lexer → Tokens → Parser → AST
```

You can use:
- **Lexer only**: For tokenization
- **Full pipeline**: For complete parsing with AST

## Quick Start

### Lexer Only
```go
package main

import (
    "fmt"
    "github.com/l00pss/citrinelexer"
)

func main() {
    lexer := citrinelexer.NewLexer("SELECT name FROM users")
    
    for {
        token := lexer.NextToken()
        fmt.Printf("%-15s %s\n", token.Type, token.Value)
        
        if token.Type == citrinelexer.EOF {
            break
        }
    }
}
```

### Parser + AST
```go
package main

import (
    "fmt"
    "github.com/l00pss/citrinelexer"
)

func main() {
    // Parse SQL into AST
    stmt, err := citrinelexer.Parse("SELECT name, age FROM users WHERE id > 100")
    if err != nil {
        panic(err)
    }

    // Work with AST
    selectStmt := stmt.(*citrinelexer.SelectStatement)
    fmt.Printf("Table: %s\n", selectStmt.From.Name.Name)
    fmt.Printf("Fields: %d\n", len(selectStmt.Fields))
    fmt.Printf("Has WHERE: %t\n", selectStmt.Where != nil)
}
```

## Supported SQL Features

### Statements
```sql
SELECT name, age FROM users WHERE active = 1;
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);
INSERT INTO users VALUES (1, 'John');
UPDATE users SET name = 'Jane' WHERE id = 1;
DELETE FROM users WHERE id = 1;
```

### Database Commands
```sql
PRAGMA table_info(users);
VACUUM;
EXPLAIN QUERY PLAN SELECT * FROM users;
ATTACH DATABASE 'backup.db' AS backup;
```

### Expressions & Functions
```sql
SELECT COUNT(*), AVG(age) FROM users;
SELECT name || ' ' || surname AS full_name FROM users;
SELECT * FROM users WHERE age BETWEEN 18 AND 65;
```

### Parameters
```sql
SELECT * FROM users WHERE id = ?;           -- Positional
SELECT * FROM users WHERE name = :name;     -- Named (:name)
SELECT * FROM users WHERE age = $age;       -- Named ($age)
```

### Advanced Features
- Comments (`-- line` and `/* block */`)
- Quoted identifiers (`"column name"`, `[table name]`, `` `field` ``)
- Hexadecimal numbers (`0xFF`)
- Scientific notation (`1.23e-4`)
- String concatenation (`||`)

## API Reference

### Lexer API
```go
// Create lexer
lexer := citrinelexer.NewLexer("SELECT * FROM users")

// Token by token
token := lexer.NextToken()

// All tokens at once
tokens := lexer.GetAllTokens()

// Position info
line, col := lexer.GetCurrentPosition()

// Status check
if lexer.IsAtEnd() {
    // Done
}
```

### Parser API
```go
// Parse complete statement
stmt, err := citrinelexer.Parse("SELECT * FROM users")

// Use custom lexer
lexer := citrinelexer.NewLexer(sql)
parser := citrinelexer.NewParser(lexer)
stmt, err := parser.ParseStatement()
```

### AST Nodes

The library provides full AST nodes implementing `go/ast.Node` interface:

- **Statements**: `SelectStatement`, `CreateTableStatement`, `InsertStatement`, `UpdateStatement`, `DeleteStatement`
- **Expressions**: `Identifier`, `StringLiteral`, `NumberLiteral`, `BinaryExpression`, `FunctionCall`
- **Parameters**: `Parameter` (for `?` and named parameters)

## Testing

```bash
# Run all tests
go test -v

# Run specific tests
go test -v -run TestParse

# Benchmarks
go test -bench=.

# Coverage
go test -cover
```

## Performance

Benchmarks on M1 Pro MacBook:

```
BenchmarkLexer-10                    534313    2226 ns/op    152 B/op    18 allocs/op
BenchmarkSingleCharTokens-10        5115759     237 ns/op      0 B/op     0 allocs/op
BenchmarkKeywordLookup-10           2361190     521 ns/op      0 B/op     0 allocs/op
```

**Lexer Performance:**
- ~450K complex SQL queries per second
- ~4.2M punctuation tokens per second (zero allocation!)
- ~1.9M keyword recognition per second (zero allocation!)

**Parser adds minimal overhead** while providing full AST functionality.

## Use Cases

- **Database Engines**: SQL query parsing
- **Code Analysis**: SQL static analysis tools
- **Query Builders**: Dynamic SQL generation
- **Migration Tools**: Schema parsing
- **IDEs**: SQL syntax highlighting and validation

## License

MIT - Use freely in your projects!