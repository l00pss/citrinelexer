package main

import (
	"fmt"
	"testing"

	"github.com/l00pss/citrinelexer"
)

func ExampleLexer_basic() {
	input := "SELECT name FROM users WHERE id = 123;"
	lexer := citrinelexer.NewLexer(input)

	fmt.Println("Tokenizing:", input)
	fmt.Println()

	for {
		token := lexer.NextToken()
		fmt.Printf("%s\n", token)

		if token.Type == citrinelexer.EOF {
			break
		}
	}
}

func ExampleLexer_getAllTokens() {
	input := "SELECT * FROM users;"
	lexer := citrinelexer.NewLexer(input)

	tokens := lexer.GetAllTokens()

	fmt.Printf("Found %d tokens:\n", len(tokens))
	for i, token := range tokens {
		fmt.Printf("%d: %s\n", i+1, token)
	}
}

func TestTokenTypeHelpers(t *testing.T) {
	// Test keyword detection
	if !citrinelexer.SELECT.IsKeyword() {
		t.Error("SELECT should be identified as keyword")
	}

	if citrinelexer.IDENTIFIER.IsKeyword() {
		t.Error("IDENTIFIER should not be identified as keyword")
	}

	// Test operator detection
	if !citrinelexer.EQUAL.IsOperator() {
		t.Error("EQUAL should be identified as operator")
	}

	if citrinelexer.STRING.IsOperator() {
		t.Error("STRING should not be identified as operator")
	}
}

func TestLexerHelpers(t *testing.T) {
	lexer := citrinelexer.NewLexer("SELECT")

	// Test position tracking
	line, col := lexer.GetCurrentPosition()
	if line != 1 || col != 1 {
		t.Errorf("Expected position (1,1), got (%d,%d)", line, col)
	}

	// Test end detection
	if lexer.IsAtEnd() {
		t.Error("Lexer should not be at end initially")
	}

	// Process all tokens
	lexer.GetAllTokens()

	if !lexer.IsAtEnd() {
		t.Error("Lexer should be at end after processing all tokens")
	}
}
