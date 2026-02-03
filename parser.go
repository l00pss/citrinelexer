package citrinelexer

import (
	"fmt"
	"go/token"
	"strconv"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func Parse(sql string) (Statement, error) {
	lexer := NewLexer(sql)
	parser := NewParser(lexer)
	return parser.ParseStatement()
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseStatement() (Statement, error) {
	switch p.currentToken.Type {
	case SELECT:
		return p.parseSelectStatement()
	case CREATE:
		return p.parseCreateStatement()
	case INSERT:
		return p.parseInsertStatement()
	case UPDATE:
		return p.parseUpdateStatement()
	case DELETE:
		return p.parseDeleteStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Type)
	}
}

func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	stmt := &SelectStatement{
		Select: token.Pos(p.currentToken.Col),
	}

	if !p.expectToken(SELECT) {
		return nil, fmt.Errorf("expected SELECT")
	}

	fields, err := p.parseSelectFields()
	if err != nil {
		return nil, err
	}
	stmt.Fields = fields

	if p.currentToken.Type == FROM {
		p.nextToken()
		from, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	if p.currentToken.Type == WHERE {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	if p.currentToken.Type == ORDER {
		p.nextToken()
		if !p.expectToken(BY) {
			return nil, fmt.Errorf("expected BY after ORDER")
		}
		orderBy, err := p.parseOrderBy()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	if p.currentToken.Type == LIMIT {
		p.nextToken()
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	return stmt, nil
}

func (p *Parser) parseSelectFields() ([]Expression, error) {
	var fields []Expression

	if p.currentToken.Type == ASTERISK {
		fields = append(fields, &Identifier{
			Name: "*",
			Pos_: token.Pos(p.currentToken.Col),
		})
		p.nextToken()
	} else {
		for {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			fields = append(fields, expr)

			if p.currentToken.Type != COMMA {
				break
			}
			p.nextToken()
		}
	}

	return fields, nil
}

func (p *Parser) parseCreateStatement() (*CreateTableStatement, error) {
	stmt := &CreateTableStatement{
		Create: token.Pos(p.currentToken.Col),
	}

	if !p.expectToken(CREATE) {
		return nil, fmt.Errorf("expected CREATE")
	}

	if !p.expectToken(TABLE) {
		return nil, fmt.Errorf("expected TABLE")
	}

	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.Table = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	if !p.expectToken(LPAREN) {
		return nil, fmt.Errorf("expected (")
	}

	columns, err := p.parseColumnDefs()
	if err != nil {
		return nil, err
	}
	stmt.Columns = columns

	if !p.expectToken(RPAREN) {
		return nil, fmt.Errorf("expected )")
	}

	return stmt, nil
}

func (p *Parser) parseColumnDefs() ([]*ColumnDef, error) {
	var columns []*ColumnDef

	for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
		col, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)

		if p.currentToken.Type == COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	return columns, nil
}

func (p *Parser) parseColumnDef() (*ColumnDef, error) {
	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected column name")
	}

	col := &ColumnDef{
		Name: &Identifier{
			Name: p.currentToken.Value,
			Pos_: token.Pos(p.currentToken.Col),
		},
	}
	p.nextToken()

	if p.isDataType() {
		col.Type = p.currentToken.Value
		p.nextToken()

		// Parse length parameter for VARCHAR(n), CHAR(n), etc.
		if p.currentToken.Type == LPAREN {
			p.nextToken() // consume (
			if p.currentToken.Type == NUMBER {
				col.Length, _ = strconv.Atoi(p.currentToken.Value)
				p.nextToken()

				// Parse precision and scale for DECIMAL(p,s)
				if p.currentToken.Type == COMMA {
					p.nextToken()
					col.Precision = col.Length
					col.Length = 0
					if p.currentToken.Type == NUMBER {
						col.Scale, _ = strconv.Atoi(p.currentToken.Value)
						p.nextToken()
					}
				}
			}
			if !p.expectToken(RPAREN) {
				return nil, fmt.Errorf("expected )")
			}
		}
	}

	for p.isConstraintKeyword() {
		constraint, err := p.parseConstraint()
		if err != nil {
			return nil, err
		}
		col.Constraints = append(col.Constraints, constraint)
	}

	return col, nil
}

func (p *Parser) parseInsertStatement() (*InsertStatement, error) {
	stmt := &InsertStatement{
		Insert: token.Pos(p.currentToken.Col),
	}

	if !p.expectToken(INSERT) {
		return nil, fmt.Errorf("expected INSERT")
	}

	// INTO is optional
	if p.currentToken.Type == INTO {
		p.nextToken()
	}

	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.Table = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	// Parse column list (optional)
	if p.currentToken.Type == LPAREN {
		p.nextToken()
		for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
			if p.currentToken.Type != IDENTIFIER {
				return nil, fmt.Errorf("expected column name")
			}
			stmt.Columns = append(stmt.Columns, &Identifier{
				Name: p.currentToken.Value,
				Pos_: token.Pos(p.currentToken.Col),
			})
			p.nextToken()
			if p.currentToken.Type == COMMA {
				p.nextToken()
			}
		}
		if !p.expectToken(RPAREN) {
			return nil, fmt.Errorf("expected )")
		}
	}

	// Parse VALUES clause
	if p.currentToken.Type == VALUES {
		p.nextToken()

		for {
			if p.currentToken.Type != LPAREN {
				return nil, fmt.Errorf("expected (")
			}
			p.nextToken()

			var values []Expression
			for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				values = append(values, expr)
				if p.currentToken.Type == COMMA {
					p.nextToken()
				}
			}
			if !p.expectToken(RPAREN) {
				return nil, fmt.Errorf("expected )")
			}
			stmt.Values = append(stmt.Values, values)

			// Multiple value sets: INSERT INTO t VALUES (1), (2), (3)
			if p.currentToken.Type != COMMA {
				break
			}
			p.nextToken()
		}
	}

	return stmt, nil
}

func (p *Parser) parseUpdateStatement() (*UpdateStatement, error) {
	stmt := &UpdateStatement{
		Update: token.Pos(p.currentToken.Col),
	}

	if !p.expectToken(UPDATE) {
		return nil, fmt.Errorf("expected UPDATE")
	}

	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.Table = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	// Parse SET clause
	if !p.expectToken(SET) {
		return nil, fmt.Errorf("expected SET")
	}

	for {
		if p.currentToken.Type != IDENTIFIER {
			return nil, fmt.Errorf("expected column name")
		}

		col := &Identifier{
			Name: p.currentToken.Value,
			Pos_: token.Pos(p.currentToken.Col),
		}
		p.nextToken()

		if !p.expectToken(EQUAL) {
			return nil, fmt.Errorf("expected =")
		}

		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		stmt.Set = append(stmt.Set, &Assignment{
			Column: col,
			Value:  value,
		})

		if p.currentToken.Type != COMMA {
			break
		}
		p.nextToken()
	}

	// Parse WHERE clause (optional)
	if p.currentToken.Type == WHERE {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

func (p *Parser) parseDeleteStatement() (*DeleteStatement, error) {
	stmt := &DeleteStatement{
		Delete: token.Pos(p.currentToken.Col),
	}

	if !p.expectToken(DELETE) {
		return nil, fmt.Errorf("expected DELETE")
	}

	if !p.expectToken(FROM) {
		return nil, fmt.Errorf("expected FROM")
	}

	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.From = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	if p.currentToken.Type == WHERE {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseComparison()
}

func (p *Parser) parseComparison() (Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.isComparisonOperator() {
		operator := p.currentToken.Value
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()

		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
			Pos_:     pos,
		}, nil
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	var name string
	var pos token.Pos

	switch p.currentToken.Type {
	case IDENTIFIER, COUNT, SUM, AVG, MIN, MAX:
		name = p.currentToken.Value
		pos = token.Pos(p.currentToken.Col)
		p.nextToken() // Advance past the identifier

		if p.currentToken.Type == LPAREN {
			p.nextToken()
			args := []Expression{}

			if p.currentToken.Type != RPAREN {
				for {
					arg, err := p.parseExpression()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)

					if p.currentToken.Type != COMMA {
						break
					}
					p.nextToken()
				}
			}

			if !p.expectToken(RPAREN) {
				return nil, fmt.Errorf("expected )")
			}

			return &FunctionCall{
				Name: name,
				Args: args,
				Pos_: pos,
			}, nil
		}

		return &Identifier{
			Name: name,
			Pos_: pos,
		}, nil

	case STRING:
		value := p.currentToken.Value
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &StringLiteral{
			Value: value,
			Pos_:  pos,
		}, nil

	case NUMBER:
		value := p.currentToken.Value
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &NumberLiteral{
			Value: value,
			Pos_:  pos,
		}, nil

	case TRUE, FALSE:
		value := p.currentToken.Type == TRUE
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &BooleanLiteral{
			Value: value,
			Pos_:  pos,
		}, nil

	case PARAMETER:
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &Parameter{
			Name: "",
			Pos_: pos,
		}, nil

	case NAMED_PARAMETER:
		name := p.currentToken.Value
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &Parameter{
			Name: name,
			Pos_: pos,
		}, nil

	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Type)
	}
}

func (p *Parser) parseTableRef() (*TableRef, error) {
	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}

	table := &TableRef{
		Name: &Identifier{
			Name: p.currentToken.Value,
			Pos_: token.Pos(p.currentToken.Col),
		},
	}
	p.nextToken()

	if p.currentToken.Type == AS {
		p.nextToken()
		if p.currentToken.Type != IDENTIFIER {
			return nil, fmt.Errorf("expected alias after AS")
		}
		table.Alias = &Identifier{
			Name: p.currentToken.Value,
			Pos_: token.Pos(p.currentToken.Col),
		}
		p.nextToken()
	}

	return table, nil
}

func (p *Parser) parseOrderBy() ([]OrderByItem, error) {
	var items []OrderByItem

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		direction := "ASC"
		if p.currentToken.Type == IDENTIFIER {
			if p.currentToken.Value == "DESC" || p.currentToken.Value == "ASC" {
				direction = p.currentToken.Value
				p.nextToken()
			}
		}

		items = append(items, OrderByItem{
			Expression: expr,
			Direction:  direction,
		})

		if p.currentToken.Type != COMMA {
			break
		}
		p.nextToken()
	}

	return items, nil
}

func (p *Parser) parseLimitClause() (*LimitClause, error) {
	count, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	clause := &LimitClause{
		Count: count,
	}

	if p.currentToken.Type == OFFSET {
		p.nextToken()
		offset, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		clause.Offset = offset
	}

	return clause, nil
}

func (p *Parser) parseConstraint() (Constraint, error) {
	switch p.currentToken.Type {
	case PRIMARY:
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		if !p.expectToken(KEY) {
			return nil, fmt.Errorf("expected KEY after PRIMARY")
		}
		return &PrimaryKeyConstraint{Pos_: pos}, nil
	case NOT:
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		if !p.expectToken(NULL) {
			return nil, fmt.Errorf("expected NULL after NOT")
		}
		return &NotNullConstraint{Pos_: pos}, nil
	default:
		return nil, fmt.Errorf("unknown constraint: %s", p.currentToken.Type)
	}
}

func (p *Parser) expectToken(expected TokenType) bool {
	if p.currentToken.Type == expected {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) isDataType() bool {
	switch p.currentToken.Type {
	case INTEGER, INT, TEXT, VARCHAR, CHAR, REAL, BLOB, BOOLEAN, DATETIME, TIMESTAMP:
		return true
	default:
		return false
	}
}

func (p *Parser) isConstraintKeyword() bool {
	switch p.currentToken.Type {
	case PRIMARY, NOT, UNIQUE, DEFAULT:
		return true
	default:
		return false
	}
}

func (p *Parser) isComparisonOperator() bool {
	switch p.currentToken.Type {
	case EQUAL, NOT_EQUAL, NOT_EQUAL2, GREATER, LESS, GREATER_EQUAL, LESS_EQUAL, LIKE:
		return true
	default:
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
