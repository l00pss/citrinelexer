# Citrine Lexer 

A simple and fast SQL lexer written in Go. Perfect for parsing SQL queries into tokens!

## What does it do?

This lexer takes SQL code and breaks it down into meaningful pieces (tokens). Think of it like taking apart a sentence to understand each word.

**Input:** `SELECT name FROM users WHERE id = 123;`

**Output:** `SELECT`, `name`, `FROM`, `users`, `WHERE`, `id`, `=`, `123`, `;`

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/l00pss/citrinelexer"
)

func main() {
    sql := "SELECT name FROM users WHERE id = 123;"
    lexer := citrinelexer.NewLexer(sql)
    
    for {
        token := lexer.NextToken()
        fmt.Printf("%s\n", token)
        
        if token.Type == citrinelexer.EOF {
            break
        }
    }
}
```

## API Reference

### Basic Methods

```go
// Create a new lexer
lexer := citrinelexer.NewLexer("SELECT * FROM users;")

// Get next token
token := lexer.NextToken()

// Get all tokens at once
tokens := lexer.GetAllTokens()

// Check current position
line, col := lexer.GetCurrentPosition()

// Check if finished
if lexer.IsAtEnd() {
    // Done parsing
}
```

## Running Tests

```bash
go test -v .
```

## Running Examples

```bash
cd example
go run demo.go
```

## Performance

**System:** Apple M1 Pro, macOS

```
goos: darwin
goarch: arm64
pkg: github.com/l00pss/citrinelexer
cpu: Apple M1 Pro

BenchmarkLexer-10                         534313              2226 ns/op           152 B/op            18 allocs/op
BenchmarkSingleCharTokens-10             5115759               237.2 ns/op           0 B/op             0 allocs/op
BenchmarkKeywordLookup-10                2361190               521.1 ns/op           0 B/op             0 allocs/op
```

**What this means:**
- **BenchmarkLexer**: Can parse ~450,000 complex SQL queries per second
- **BenchmarkSingleCharTokens**: Can process ~4.2 million punctuation tokens per second with **zero allocations**
- **BenchmarkKeywordLookup**: Can identify ~1.9 million keywords per second with **zero allocations**

Pretty fast! 

## Contributing

Found a bug? Want to add a feature? Pull requests are welcome!

## License

MIT License - feel free to use this in your projects!

---

Made with ❤️ by [l00pss](https://github.com/l00pss)