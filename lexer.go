package citrinelexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	SELECT TokenType = iota
	FROM
	WHERE
	INSERT
	UPDATE
	DELETE
	CREATE
	TABLE
	TRUNCATE
	DROP
	ALTER
	INDEX
	PRIMARY
	KEY
	FOREIGN
	REFERENCES
	NOT
	NULL
	DEFAULT
	AUTO_INCREMENT
	UNIQUE

	INT
	INTEGER
	VARCHAR
	TEXT
	CHAR
	BOOLEAN
	DATETIME
	TIMESTAMP

	ORDER
	BY
	GROUP
	HAVING
	LIMIT
	OFFSET
	INNER
	LEFT
	RIGHT
	JOIN
	ON
	AS
	DISTINCT
	COUNT
	SUM
	AVG
	MAX
	MIN

	AND
	OR
	IN
	LIKE
	BETWEEN
	IS

	IDENTIFIER
	STRING
	NUMBER
	BOOLEAN_LITERAL

	EQUAL
	GREATER
	LESS
	GREATER_EQUAL
	LESS_EQUAL
	NOT_EQUAL
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULO

	SEMICOLON
	COMMA
	LPAREN
	RPAREN
	DOT
	ASTERISK

	EOF
	ILLEGAL
)

var (
	equalStr        = "="
	doubleEqualStr  = "=="
	greaterStr      = ">"
	greaterEqualStr = ">="
	lessStr         = "<"
	lessEqualStr    = "<="
	notEqualStr1    = "!="
	notEqualStr2    = "<>"
	emptyStr        = ""
)

var singleCharTokens = map[rune]struct {
	TokenType
	Value string
}{
	';': {SEMICOLON, ";"},
	',': {COMMA, ","},
	'(': {LPAREN, "("},
	')': {RPAREN, ")"},
	'.': {DOT, "."},
	'*': {ASTERISK, "*"},
	'+': {PLUS, "+"},
	'-': {MINUS, "-"},
	'/': {DIVIDE, "/"},
	'%': {MODULO, "%"},
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, '%s', %d:%d}", t.Type.String(), t.Value, t.Line, t.Col)
}

func (tt TokenType) String() string {
	switch tt {
	case SELECT:
		return "SELECT"
	case FROM:
		return "FROM"
	case WHERE:
		return "WHERE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}

type Lexer struct {
	input    string
	position int
	readPos  int
	ch       rune
	line     int
	col      int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.position = l.readPos
	if l.readPos >= len(l.input) {
		l.ch = 0
		l.readPos++
		return
	}

	l.ch = rune(l.input[l.readPos])
	l.readPos++

	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPos])
}

func (l *Lexer) IsAtEnd() bool {
	return l.ch == 0
}

func (l *Lexer) GetCurrentPosition() (int, int) {
	return l.line, l.col
}

func (l *Lexer) MakeToken(tokenType TokenType, value string) Token {
	return Token{
		Type:  tokenType,
		Value: value,
		Line:  l.line,
		Col:   l.col,
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.MakeToken(EQUAL, doubleEqualStr)
			tok.Col = l.col - 1
		} else {
			tok = l.MakeToken(EQUAL, equalStr)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: GREATER_EQUAL, Value: greaterEqualStr, Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: GREATER, Value: greaterStr, Line: l.line, Col: l.col}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: LESS_EQUAL, Value: lessEqualStr, Line: l.line, Col: l.col - 1}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: NOT_EQUAL, Value: notEqualStr2, Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: LESS, Value: lessStr, Line: l.line, Col: l.col}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NOT_EQUAL, Value: notEqualStr1, Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case ';', ',', '(', ')', '.', '*', '+', '-', '/', '%':
		if charToken, ok := singleCharTokens[l.ch]; ok {
			tok = Token{Type: charToken.TokenType, Value: charToken.Value, Line: l.line, Col: l.col}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case '\'':
		tok.Type = STRING
		tok.Value = l.readString()
		tok.Line = l.line
		tok.Col = l.col
		return tok
	case 0:
		tok = Token{Type: EOF, Value: emptyStr, Line: l.line, Col: l.col}
	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Type = lookupIdent(tok.Value)
			tok.Line = l.line
			tok.Col = l.col
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Value = l.readNumber()
			tok.Line = l.line
			tok.Col = l.col
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		} else if unicode.IsSpace(l.ch) {
			l.readChar()
		} else {
			break
		}
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	start := l.position + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	value := l.input[start:l.position]
	if l.ch == '\'' {
		l.readChar()
	}
	return value
}

func isLetter(ch rune) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		return true
	}
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return unicode.IsDigit(ch)
}

var keywords = map[string]TokenType{
	"SELECT":   SELECT,
	"FROM":     FROM,
	"WHERE":    WHERE,
	"INSERT":   INSERT,
	"UPDATE":   UPDATE,
	"DELETE":   DELETE,
	"CREATE":   CREATE,
	"TABLE":    TABLE,
	"TRUNCATE": TRUNCATE,
	"DROP":     DROP,
	"ALTER":    ALTER,
	"INDEX":    INDEX,

	"PRIMARY":        PRIMARY,
	"KEY":            KEY,
	"FOREIGN":        FOREIGN,
	"REFERENCES":     REFERENCES,
	"NOT":            NOT,
	"NULL":           NULL,
	"DEFAULT":        DEFAULT,
	"AUTO_INCREMENT": AUTO_INCREMENT,
	"UNIQUE":         UNIQUE,

	"INT":       INT,
	"INTEGER":   INTEGER,
	"VARCHAR":   VARCHAR,
	"TEXT":      TEXT,
	"CHAR":      CHAR,
	"BOOLEAN":   BOOLEAN,
	"DATETIME":  DATETIME,
	"TIMESTAMP": TIMESTAMP,

	"ORDER":    ORDER,
	"BY":       BY,
	"GROUP":    GROUP,
	"HAVING":   HAVING,
	"LIMIT":    LIMIT,
	"OFFSET":   OFFSET,
	"INNER":    INNER,
	"LEFT":     LEFT,
	"RIGHT":    RIGHT,
	"JOIN":     JOIN,
	"ON":       ON,
	"AS":       AS,
	"DISTINCT": DISTINCT,

	"COUNT": COUNT,
	"SUM":   SUM,
	"AVG":   AVG,
	"MAX":   MAX,
	"MIN":   MIN,

	"AND":     AND,
	"OR":      OR,
	"IN":      IN,
	"LIKE":    LIKE,
	"BETWEEN": BETWEEN,
	"IS":      IS,

	"TRUE":  BOOLEAN_LITERAL,
	"FALSE": BOOLEAN_LITERAL,
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToUpper(ident)]; ok {
		return tok
	}
	return IDENTIFIER
}

func (l *Lexer) GetAllTokens() []Token {
	var tokens []Token
	for {
		token := l.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}
	return tokens
}

func (tt TokenType) IsKeyword() bool {
	return tt >= SELECT && tt <= IS
}

func (tt TokenType) IsOperator() bool {
	return tt >= EQUAL && tt <= MODULO
}
