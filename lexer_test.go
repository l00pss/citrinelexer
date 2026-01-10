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
		{NOT_EQUAL2, "<>"},
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

func TestDatabaseKeywords(t *testing.T) {
	input := `PRAGMA table_info(users);
VACUUM;
EXPLAIN QUERY PLAN SELECT * FROM users;
ATTACH DATABASE 'test.db' AS test;`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{PRAGMA, "PRAGMA"},
		{IDENTIFIER, "table_info"},
		{LPAREN, "("},
		{IDENTIFIER, "users"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{VACUUM, "VACUUM"},
		{SEMICOLON, ";"},
		{EXPLAIN, "EXPLAIN"},
		{QUERY, "QUERY"},
		{PLAN, "PLAN"},
		{SELECT, "SELECT"},
		{ASTERISK, "*"},
		{FROM, "FROM"},
		{IDENTIFIER, "users"},
		{SEMICOLON, ";"},
		{ATTACH, "ATTACH"},
		{DATABASE, "DATABASE"},
		{STRING, "test.db"},
		{AS, "AS"},
		{IDENTIFIER, "test"},
		{SEMICOLON, ";"},
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

func TestQuotedIdentifiers(t *testing.T) {
	input := `"column name" [table name] ` + "`backtick`"

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{IDENTIFIER, "column name"},
		{IDENTIFIER, "table name"},
		{IDENTIFIER, "backtick"},
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

func TestParameters(t *testing.T) {
	input := `? :name $param`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{PARAMETER, "?"},
		{NAMED_PARAMETER, ":name"},
		{NAMED_PARAMETER, "$param"},
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

func TestComments(t *testing.T) {
	input := `SELECT * -- line comment
FROM users /* block comment */ WHERE id = 1`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{SELECT, "SELECT"},
		{ASTERISK, "*"},
		{FROM, "FROM"},
		{IDENTIFIER, "users"},
		{WHERE, "WHERE"},
		{IDENTIFIER, "id"},
		{EQUAL, "="},
		{NUMBER, "1"},
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

func TestHexNumbers(t *testing.T) {
	input := `0x1A 0xFF 0x00`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{NUMBER, "0x1A"},
		{NUMBER, "0xFF"},
		{NUMBER, "0x00"},
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

func TestConcatenationOperator(t *testing.T) {
	input := `name || ' ' || surname`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{IDENTIFIER, "name"},
		{CONCAT, "||"},
		{STRING, " "},
		{CONCAT, "||"},
		{IDENTIFIER, "surname"},
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

func TestWindowFunctions(t *testing.T) {
	input := `ROW_NUMBER() OVER (PARTITION BY category ORDER BY price)`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{IDENTIFIER, "ROW_NUMBER"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{OVER, "OVER"},
		{LPAREN, "("},
		{PARTITION, "PARTITION"},
		{BY, "BY"},
		{IDENTIFIER, "category"},
		{ORDER, "ORDER"},
		{BY, "BY"},
		{IDENTIFIER, "price"},
		{RPAREN, ")"},
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