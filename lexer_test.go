package citrinelexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `SELECT name, age FROM users WHERE id = 123;`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{SELECT, "SELECT"},
		{IDENTIFIER, "name"},
		{COMMA, ","},
		{IDENTIFIER, "age"},
		{FROM, "FROM"},
		{IDENTIFIER, "users"},
		{WHERE, "WHERE"},
		{IDENTIFIER, "id"},
		{EQUAL, "="},
		{NUMBER, "123"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := `123 456.78 0 999.99`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{NUMBER, "123"},
		{NUMBER, "456.78"},
		{NUMBER, "0"},
		{NUMBER, "999.99"},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestStrings(t *testing.T) {
	input := `'hello' 'world' 'John Doe'`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{STRING, "hello"},
		{STRING, "world"},
		{STRING, "John Doe"},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `= == > < >= <= != <>`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{EQUAL, "="},
		{EQUAL, "=="},
		{GREATER, ">"},
		{LESS, "<"},
		{GREATER_EQUAL, ">="},
		{LESS_EQUAL, "<="},
		{NOT_EQUAL, "!="},
		{NOT_EQUAL, "<>"},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}