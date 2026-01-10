package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/l00pss/citrinelexer"
)

// ExampleParserUsage demonstrates how to use the parser
func ExampleParserUsage() {
	sql := "SELECT name, age FROM users WHERE id = 1;"

	fmt.Println("Parsing SQL:", sql)

	// Parse the SQL
	stmt, err := citrinelexer.Parse(sql)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return
	}

	// Check the statement type
	selectStmt, ok := stmt.(*citrinelexer.SelectStatement)
	if !ok {
		log.Println("Expected SELECT statement")
		return
	}

	// Print analysis
	if selectStmt.From != nil && selectStmt.From.Name != nil {
		fmt.Printf("Table: %s\n", selectStmt.From.Name.Name)
	} else {
		fmt.Println("Table: <unknown>")
	}
	fmt.Printf("Selected fields: %d\n", len(selectStmt.Fields))
	fmt.Printf("Has WHERE clause: %t\n", selectStmt.Where != nil)
	for i, field := range selectStmt.Fields {
		switch expr := field.(type) {
		case *citrinelexer.Identifier:
			fmt.Printf("Field %d: %s\n", i+1, expr.Name)
		case *citrinelexer.FunctionCall:
			fmt.Printf("Field %d: %s() function\n", i+1, expr.Name)
		default:
			fmt.Printf("Field %d: complex expression\n", i+1)
		}
	}
}

// ExampleLexerWithParser shows how to use lexer and parser separately
func ExampleLexerWithParser() {
	sql := "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);"

	// Use lexer first
	lexer := citrinelexer.NewLexer(sql)
	fmt.Println("Tokens:")

	var tokens []citrinelexer.Token
	for {
		token := lexer.NextToken()
		if token.Type == citrinelexer.EOF {
			break
		}
		tokens = append(tokens, token)
		fmt.Printf("  %-15s %s\n", token.Type, token.Value)
	}

	// Now parse the same SQL
	fmt.Println("\nParsing:")
	stmt, err := citrinelexer.Parse(sql)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return
	}

	createStmt := stmt.(*citrinelexer.CreateTableStatement)
	fmt.Printf("Creating table: %s\n", createStmt.Table.Name)
	fmt.Printf("Columns: %d\n", len(createStmt.Columns))

	for i, col := range createStmt.Columns {
		fmt.Printf("  Column %d: %s %s\n", i+1, col.Name.Name, col.Type)
	}
}

// ExampleParameterHandling shows parameter parsing
func ExampleParameterHandling() {
	sql := "SELECT * FROM users WHERE id = ? AND name = :username AND age > $min_age"

	stmt, err := citrinelexer.Parse(sql)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return
	}

	selectStmt := stmt.(*citrinelexer.SelectStatement)

	fmt.Println("Parameters found in WHERE clause:")
	findParameters(selectStmt.Where)
}

// findParameters recursively finds all parameters in an expression
func findParameters(expr citrinelexer.Expression) {
	switch e := expr.(type) {
	case *citrinelexer.Parameter:
		if e.Name != "" {
			fmt.Printf("  Named parameter: %s\n", e.Name)
		} else {
			fmt.Println("  Positional parameter: ?")
		}
	case *citrinelexer.BinaryExpression:
		findParameters(e.Left)
		findParameters(e.Right)
	}
}

func main() {
	fmt.Println("=== Citrine Lexer Parser Examples ===")
	fmt.Println()

	fmt.Println("1. Parser Usage:")
	ExampleParserUsage()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	fmt.Println("2. Lexer with Parser:")
	ExampleLexerWithParser()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	fmt.Println("3. Parameter Handling:")
	ExampleParameterHandling()
}
