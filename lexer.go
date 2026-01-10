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

	// SQLite specific keywords
	DATABASE
	SCHEMA
	CONSTRAINT
	CASCADE
	RESTRICT
	SET_NULL
	SET_DEFAULT
	CHECK
	COLLATE
	AUTOINCREMENT
	CONFLICT
	REPLACE
	IGNORE
	FAIL
	ABORT
	ROLLBACK
	WITHOUT
	ROWID

	// Pragma and maintenance
	PRAGMA
	VACUUM
	REINDEX
	ANALYZE
	ATTACH
	DETACH
	EXPLAIN
	QUERY
	PLAN

	// Data types
	INT
	INTEGER
	VARCHAR
	TEXT
	CHAR
	BOOLEAN
	REAL
	BLOB
	DATETIME
	TIMESTAMP

	// Query clauses
	ORDER
	BY
	GROUP
	HAVING
	LIMIT
	OFFSET
	INNER
	LEFT
	RIGHT
	FULL
	OUTER
	CROSS
	JOIN
	ON
	AS
	DISTINCT
	UNION
	INTERSECT
	EXCEPT

	// Window functions
	OVER
	PARTITION
	WINDOW
	ROWS
	RANGE
	UNBOUNDED
	PRECEDING
	FOLLOWING
	CURRENT
	ROW

	// Case expressions
	CASE
	WHEN
	THEN
	ELSE
	END

	// Aggregate functions
	COUNT
	SUM
	AVG
	MAX
	MIN

	// Logical operators
	AND
	OR
	IN
	LIKE
	GLOB
	MATCH
	REGEXP
	BETWEEN
	IS
	ISNULL
	NOTNULL
	EXISTS

	// Transaction
	BEGIN
	COMMIT
	TRANSACTION

	// Boolean literals
	TRUE
	FALSE

	// Literals and identifiers
	IDENTIFIER
	STRING
	NUMBER
	BOOLEAN_LITERAL
	PARAMETER
	NAMED_PARAMETER

	// Comparison operators
	EQUAL
	GREATER
	LESS
	GREATER_EQUAL
	LESS_EQUAL
	NOT_EQUAL
	NOT_EQUAL2 // <>

	// Arithmetic operators
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULO
	CONCAT // ||

	// Punctuation
	SEMICOLON
	COMMA
	LPAREN
	RPAREN
	DOT
	ASTERISK
	LBRACKET
	RBRACKET
	COLON
	PIPE
	BANG

	// Comments
	LINE_COMMENT
	BLOCK_COMMENT

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
	'[': {LBRACKET, "["},
	']': {RBRACKET, "]"},
	':': {COLON, ":"},
	'|': {PIPE, "|"},
	'!': {BANG, "!"},
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
	case INSERT:
		return "INSERT"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	case CREATE:
		return "CREATE"
	case DROP:
		return "DROP"
	case ALTER:
		return "ALTER"
	case TABLE:
		return "TABLE"
	case DATABASE:
		return "DATABASE"
	case SCHEMA:
		return "SCHEMA"
	case INDEX:
		return "INDEX"
	case UNIQUE:
		return "UNIQUE"
	case PRIMARY:
		return "PRIMARY"
	case KEY:
		return "KEY"
	case FOREIGN:
		return "FOREIGN"
	case REFERENCES:
		return "REFERENCES"
	case CONSTRAINT:
		return "CONSTRAINT"
	case CHECK:
		return "CHECK"
	case NOT:
		return "NOT"
	case NULL:
		return "NULL"
	case PRAGMA:
		return "PRAGMA"
	case VACUUM:
		return "VACUUM"
	case EXPLAIN:
		return "EXPLAIN"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case IN:
		return "IN"
	case LIKE:
		return "LIKE"
	case BETWEEN:
		return "BETWEEN"
	case IS:
		return "IS"
	case EXISTS:
		return "EXISTS"
	case CASE:
		return "CASE"
	case WHEN:
		return "WHEN"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case END:
		return "END"
	case ORDER:
		return "ORDER"
	case BY:
		return "BY"
	case GROUP:
		return "GROUP"
	case HAVING:
		return "HAVING"
	case LIMIT:
		return "LIMIT"
	case OFFSET:
		return "OFFSET"
	case DISTINCT:
		return "DISTINCT"
	case UNION:
		return "UNION"
	case INTERSECT:
		return "INTERSECT"
	case EXCEPT:
		return "EXCEPT"
	case JOIN:
		return "JOIN"
	case INNER:
		return "INNER"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case ON:
		return "ON"
	case AS:
		return "AS"
	case COUNT:
		return "COUNT"
	case SUM:
		return "SUM"
	case AVG:
		return "AVG"
	case MAX:
		return "MAX"
	case MIN:
		return "MIN"
	case INTEGER:
		return "INTEGER"
	case TEXT:
		return "TEXT"
	case REAL:
		return "REAL"
	case BLOB:
		return "BLOB"
	case VARCHAR:
		return "VARCHAR"
	case CHAR:
		return "CHAR"
	case BOOLEAN:
		return "BOOLEAN"
	case DATETIME:
		return "DATETIME"
	case TIMESTAMP:
		return "TIMESTAMP"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case PARAMETER:
		return "PARAMETER"
	case NAMED_PARAMETER:
		return "NAMED_PARAMETER"
	case EQUAL:
		return "EQUAL"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case NOT_EQUAL2:
		return "NOT_EQUAL2"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case CONCAT:
		return "CONCAT"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case DOT:
		return "DOT"
	case ASTERISK:
		return "ASTERISK"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COLON:
		return "COLON"
	case PIPE:
		return "PIPE"
	case BANG:
		return "BANG"
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
			tok = Token{Type: NOT_EQUAL2, Value: notEqualStr2, Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: LESS, Value: lessStr, Line: l.line, Col: l.col}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NOT_EQUAL, Value: notEqualStr1, Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: BANG, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: CONCAT, Value: "||", Line: l.line, Col: l.col - 1}
		} else {
			tok = Token{Type: PIPE, Value: "|", Line: l.line, Col: l.col}
		}
	case '-':
		if l.peekChar() == '-' {
			l.skipLineComment()
			return l.NextToken()
		}
		if charToken, ok := singleCharTokens[l.ch]; ok {
			tok = Token{Type: charToken.TokenType, Value: charToken.Value, Line: l.line, Col: l.col}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case '/':
		if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		}
		if charToken, ok := singleCharTokens[l.ch]; ok {
			tok = Token{Type: charToken.TokenType, Value: charToken.Value, Line: l.line, Col: l.col}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case ';', ',', '(', ')', '*', '+', '%', ']':
		if charToken, ok := singleCharTokens[l.ch]; ok {
			tok = Token{Type: charToken.TokenType, Value: charToken.Value, Line: l.line, Col: l.col}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case '.':
		if isDigit(l.peekChar()) {
			tok.Type = NUMBER
			tok.Value = l.readNumber()
			tok.Line = l.line
			tok.Col = l.col
			return tok
		}
		tok = Token{Type: DOT, Value: ".", Line: l.line, Col: l.col}
	case '\'':
		tok.Type = STRING
		tok.Value = l.readString('\'')
		tok.Line = l.line
		tok.Col = l.col
		return tok
	case '"':
		tok.Type = IDENTIFIER // SQLite uses double quotes for identifiers
		tok.Value = l.readString('"')
		tok.Line = l.line
		tok.Col = l.col
		return tok
	case '`':
		tok.Type = IDENTIFIER // MySQL style backtick identifiers
		tok.Value = l.readString('`')
		tok.Line = l.line
		tok.Col = l.col
		return tok
	case '[':
		tok.Type = IDENTIFIER // SQLite bracket identifiers
		tok.Value = l.readBracketIdentifier()
		tok.Line = l.line
		tok.Col = l.col
		return tok
	case '?':
		tok = Token{Type: PARAMETER, Value: "?", Line: l.line, Col: l.col}
	case ':':
		if isLetter(l.peekChar()) {
			tok.Type = NAMED_PARAMETER
			tok.Value = l.readNamedParameter()
			tok.Line = l.line
			tok.Col = l.col
			return tok
		}
		if charToken, ok := singleCharTokens[l.ch]; ok {
			tok = Token{Type: charToken.TokenType, Value: charToken.Value, Line: l.line, Col: l.col}
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
		}
	case '$':
		if isLetter(l.peekChar()) || isDigit(l.peekChar()) {
			tok.Type = NAMED_PARAMETER
			tok.Value = l.readNamedParameter()
			tok.Line = l.line
			tok.Col = l.col
			return tok
		}
		tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: l.line, Col: l.col}
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

	// Handle hex numbers (0x...)
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // skip 0
		l.readChar() // skip x
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return l.input[position:l.position]
	}

	// Regular numbers
	for isDigit(l.ch) {
		l.readChar()
	}

	// Handle decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Handle scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString(delimiter rune) string {
	var result strings.Builder
	l.readChar()

	for {
		if l.ch == 0 {
			break
		}

		if l.ch == delimiter {
			if l.peekChar() == delimiter {
				result.WriteRune(l.ch)
				l.readChar()
				l.readChar()
				continue
			}
			l.readChar()
			break
		}

		if l.ch == '\\' && l.peekChar() == delimiter {
			l.readChar()
			result.WriteRune(l.ch)
			l.readChar()
			continue
		}

		result.WriteRune(l.ch)
		l.readChar()
	}

	return result.String()
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

func isHexDigit(ch rune) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// skipLineComment skips -- style comments
func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// skipBlockComment skips /* */ style comments
func (l *Lexer) skipBlockComment() {
	l.readChar() // skip /
	l.readChar() // skip *

	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip *
			l.readChar() // skip /
			break
		}
		l.readChar()
	}
}

// readNamedParameter reads :name or $name style parameters
func (l *Lexer) readNamedParameter() string {
	position := l.position
	l.readChar() // skip : or $

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readBracketIdentifier reads [identifier] style identifiers
func (l *Lexer) readBracketIdentifier() string {
	position := l.position
	l.readChar() // skip [

	for l.ch != ']' && l.ch != 0 {
		l.readChar()
	}

	value := l.input[position+1 : l.position] // exclude brackets
	if l.ch == ']' {
		l.readChar()
	}
	return value
}

var keywords = map[string]TokenType{
	// Basic SQL statements
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

	// Constraints and keys
	"PRIMARY":        PRIMARY,
	"KEY":            KEY,
	"FOREIGN":        FOREIGN,
	"REFERENCES":     REFERENCES,
	"NOT":            NOT,
	"NULL":           NULL,
	"DEFAULT":        DEFAULT,
	"AUTO_INCREMENT": AUTO_INCREMENT,
	"AUTOINCREMENT":  AUTOINCREMENT,
	"UNIQUE":         UNIQUE,
	"CHECK":          CHECK,
	"CONSTRAINT":     CONSTRAINT,
	"COLLATE":        COLLATE,

	// SQLite specific
	"DATABASE": DATABASE,
	"SCHEMA":   SCHEMA,
	"CASCADE":  CASCADE,
	"RESTRICT": RESTRICT,
	"CONFLICT": CONFLICT,
	"REPLACE":  REPLACE,
	"IGNORE":   IGNORE,
	"FAIL":     FAIL,
	"ABORT":    ABORT,
	"ROLLBACK": ROLLBACK,
	"WITHOUT":  WITHOUT,
	"ROWID":    ROWID,

	// Pragma and maintenance
	"PRAGMA":  PRAGMA,
	"VACUUM":  VACUUM,
	"REINDEX": REINDEX,
	"ANALYZE": ANALYZE,
	"ATTACH":  ATTACH,
	"DETACH":  DETACH,
	"EXPLAIN": EXPLAIN,
	"QUERY":   QUERY,
	"PLAN":    PLAN,

	// Data types
	"INT":       INT,
	"INTEGER":   INTEGER,
	"VARCHAR":   VARCHAR,
	"TEXT":      TEXT,
	"CHAR":      CHAR,
	"BOOLEAN":   BOOLEAN,
	"REAL":      REAL,
	"BLOB":      BLOB,
	"DATETIME":  DATETIME,
	"TIMESTAMP": TIMESTAMP,

	// Query clauses
	"ORDER":     ORDER,
	"BY":        BY,
	"GROUP":     GROUP,
	"HAVING":    HAVING,
	"LIMIT":     LIMIT,
	"OFFSET":    OFFSET,
	"INNER":     INNER,
	"LEFT":      LEFT,
	"RIGHT":     RIGHT,
	"FULL":      FULL,
	"OUTER":     OUTER,
	"CROSS":     CROSS,
	"JOIN":      JOIN,
	"ON":        ON,
	"AS":        AS,
	"DISTINCT":  DISTINCT,
	"UNION":     UNION,
	"INTERSECT": INTERSECT,
	"EXCEPT":    EXCEPT,

	// Window functions
	"OVER":      OVER,
	"PARTITION": PARTITION,
	"WINDOW":    WINDOW,
	"ROWS":      ROWS,
	"RANGE":     RANGE,
	"UNBOUNDED": UNBOUNDED,
	"PRECEDING": PRECEDING,
	"FOLLOWING": FOLLOWING,
	"CURRENT":   CURRENT,
	"ROW":       ROW,

	// Case expressions
	"CASE": CASE,
	"WHEN": WHEN,
	"THEN": THEN,
	"ELSE": ELSE,
	"END":  END,

	// Functions
	"COUNT": COUNT,
	"SUM":   SUM,
	"AVG":   AVG,
	"MAX":   MAX,
	"MIN":   MIN,

	// Logical operators
	"AND":     AND,
	"OR":      OR,
	"IN":      IN,
	"LIKE":    LIKE,
	"GLOB":    GLOB,
	"MATCH":   MATCH,
	"REGEXP":  REGEXP,
	"BETWEEN": BETWEEN,
	"IS":      IS,
	"ISNULL":  ISNULL,
	"NOTNULL": NOTNULL,
	"EXISTS":  EXISTS,

	// Transaction
	"BEGIN":       BEGIN,
	"COMMIT":      COMMIT,
	"TRANSACTION": TRANSACTION,

	// Boolean literals
	"TRUE":  TRUE,
	"FALSE": FALSE,
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
	return tt >= SELECT && tt <= FALSE
}

func (tt TokenType) IsOperator() bool {
	return tt >= EQUAL && tt <= CONCAT
}
