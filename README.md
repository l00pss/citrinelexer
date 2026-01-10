# Citrine Lexer

A fast and flexible SQL lexer library. Perfect for building your own database engine or SQL parser!

## What does it do?

Breaks SQL code into meaningful pieces (tokens). Like separating a sentence into words.

**Input:** `SELECT name FROM users WHERE id = 123;`
**Output:** `SELECT`, `name`, `FROM`, `users`, `WHERE`, `id`, `=`, `123`, `;`

## Features

- 🚀 **Fast**: Processes ~450,000 SQL queries per second
- 🔧 **Flexible**: Supports different SQL dialects
- 💾 **Efficient**: Zero allocation optimizations
- 📝 **Comprehensive**: Recognizes 100+ SQL keywords
- 🧪 **Battle-tested**: Extensive test suite

## Supported SQL Features

### Database Commands
```sql
PRAGMA table_info(users);
VACUUM;
EXPLAIN QUERY PLAN SELECT * FROM users;
ATTACH DATABASE 'backup.db' AS backup;
```

### Identifier Types
```sql
"column name"    -- Double quoted
[table name]     -- Bracket style
`field_name`     -- Backtick style
```

### Parameters
```sql
SELECT * FROM users WHERE id = ?;           -- Positional
SELECT * FROM users WHERE name = :name;     -- Named (:name)
SELECT * FROM users WHERE age = $age;       -- Named ($age)
```

### Comments
```sql
SELECT * FROM users -- Line comment
/* Multi-line
   comment */ WHERE active = 1;
```

### Numbers
```sql
123        -- Integer
45.67      -- Decimal  
1.23e-4    -- Scientific notation
0xFF       -- Hexadecimal
```

### Operators
```sql
name || ' ' || surname  -- String concatenation
age <> 25              -- Alternative not equal
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/l00pss/citrinelexer"
)

func main() {
    sql := `SELECT u.name, u.age 
            FROM users u 
            WHERE u.id = ? AND u.active = 1;`
    
    lexer := citrinelexer.NewLexer(sql)
    
    for {
        token := lexer.NextToken()
        fmt.Printf("%-15s %s\n", token.Type, token.Value)
        
        if token.Type == citrinelexer.EOF {
            break
        }
    }
}
```

## API

### Basic Usage

```go
// Create lexer
lexer := citrinelexer.NewLexer("SELECT * FROM users;")

// Get tokens one by one
token := lexer.NextToken()

// Get all tokens at once
tokens := lexer.GetAllTokens()

// Position info
line, col := lexer.GetCurrentPosition()

// Check if done
if lexer.IsAtEnd() {
    fmt.Println("Parsing completed!")
}
```

### Token Types

```go
// Is it a keyword?
if token.Type.IsKeyword() {
    fmt.Println("This is a SQL keyword")
}

// Is it an operator?
if token.Type.IsOperator() {
    fmt.Println("This is an operator")
}
```

## Testing & Benchmarks

```bash
# Run tests
go test -v

# Run benchmarks
go test -bench=.

# Check coverage
go test -cover
```

## Performance

Benchmarks on M1 Pro MacBook:

```
BenchmarkLexer-10                    534313    2226 ns/op    152 B/op    18 allocs/op
BenchmarkSingleCharTokens-10        5115759     237 ns/op      0 B/op     0 allocs/op
BenchmarkKeywordLookup-10           2361190     521 ns/op      0 B/op     0 allocs/op
```

**What this means:**
- ~450K complex SQL queries per second
- ~4.2M punctuation tokens per second (zero allocation!)  
- ~1.9M keyword recognition per second (zero allocation!)

## Project Structure

```
citrinelexer/
├── lexer.go           # Main lexer implementation
├── lexer_test.go      # Comprehensive tests
├── benchmark_test.go  # Performance tests
└── example/           # Usage examples
```

## License

MIT - Use freely in your projects!