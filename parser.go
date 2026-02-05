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
	case BEGIN:
		return p.parseBeginStatement()
	case COMMIT:
		return p.parseCommitStatement()
	case ROLLBACK:
		return p.parseRollbackStatement()
	case DROP:
		return p.parseDropStatement()
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

	// Check for DISTINCT
	if p.currentToken.Type == DISTINCT {
		stmt.Distinct = true
		p.nextToken()
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

		for p.isJoinKeyword() {
			join, err := p.parseJoinClause()
			if err != nil {
				return nil, err
			}
			stmt.Joins = append(stmt.Joins, join)
		}
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

			// Check for column alias: expr AS alias or expr alias
			if p.currentToken.Type == AS {
				pos := token.Pos(p.currentToken.Col)
				p.nextToken() // skip AS
				// Allow keywords (COUNT, SUM, etc.) as aliases too
				if p.currentToken.Type != IDENTIFIER && !p.isValidAliasToken() {
					return nil, fmt.Errorf("expected alias after AS")
				}
				alias := p.currentToken.Value
				p.nextToken()
				expr = &AliasedExpression{
					Expr:  expr,
					Alias: alias,
					Pos_:  pos,
				}
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

func (p *Parser) parseCreateStatement() (Statement, error) {
	pos := token.Pos(p.currentToken.Col)

	if !p.expectToken(CREATE) {
		return nil, fmt.Errorf("expected CREATE")
	}

	// Check for UNIQUE INDEX
	if p.currentToken.Type == UNIQUE {
		p.nextToken()
		if p.currentToken.Type == INDEX {
			return p.parseCreateIndexStatement(pos, true)
		}
		return nil, fmt.Errorf("expected INDEX after UNIQUE")
	}

	// Check for INDEX
	if p.currentToken.Type == INDEX {
		return p.parseCreateIndexStatement(pos, false)
	}

	// Parse CREATE TABLE
	if !p.expectToken(TABLE) {
		return nil, fmt.Errorf("expected TABLE or INDEX")
	}

	stmt := &CreateTableStatement{
		Create: pos,
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

func (p *Parser) parseBeginStatement() (*BeginStatement, error) {
	stmt := &BeginStatement{
		Begin: token.Pos(p.currentToken.Col),
	}
	p.nextToken() // skip BEGIN

	// TRANSACTION is optional
	if p.currentToken.Type == TRANSACTION {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseCommitStatement() (*CommitStatement, error) {
	stmt := &CommitStatement{
		Commit: token.Pos(p.currentToken.Col),
	}
	p.nextToken() // skip COMMIT
	return stmt, nil
}

func (p *Parser) parseRollbackStatement() (*RollbackStatement, error) {
	stmt := &RollbackStatement{
		Rollback: token.Pos(p.currentToken.Col),
	}
	p.nextToken() // skip ROLLBACK
	return stmt, nil
}

func (p *Parser) parseCreateIndexStatement(pos token.Pos, unique bool) (*CreateIndexStatement, error) {
	stmt := &CreateIndexStatement{
		Create: pos,
		Unique: unique,
	}

	if !p.expectToken(INDEX) {
		return nil, fmt.Errorf("expected INDEX")
	}

	// Check for IF NOT EXISTS
	if p.currentToken.Type == IF {
		p.nextToken()
		if !p.expectToken(NOT) {
			return nil, fmt.Errorf("expected NOT after IF")
		}
		if !p.expectToken(EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after IF NOT")
		}
		stmt.IfNotExists = true
	}

	// Parse index name
	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected index name")
	}
	stmt.Name = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	// Parse ON table_name
	if !p.expectToken(ON) {
		return nil, fmt.Errorf("expected ON")
	}

	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected table name")
	}
	stmt.Table = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

	// Parse (column1 [ASC|DESC], column2 [ASC|DESC], ...)
	if !p.expectToken(LPAREN) {
		return nil, fmt.Errorf("expected (")
	}

	columns, err := p.parseIndexColumns()
	if err != nil {
		return nil, err
	}
	stmt.Columns = columns

	if !p.expectToken(RPAREN) {
		return nil, fmt.Errorf("expected )")
	}

	return stmt, nil
}

func (p *Parser) parseIndexColumns() ([]*IndexColumn, error) {
	var columns []*IndexColumn

	for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
		if p.currentToken.Type != IDENTIFIER {
			return nil, fmt.Errorf("expected column name")
		}

		col := &IndexColumn{
			Column: &Identifier{
				Name: p.currentToken.Value,
				Pos_: token.Pos(p.currentToken.Col),
			},
		}
		p.nextToken()

		// Check for ASC or DESC
		if p.currentToken.Type == ASC {
			col.Direction = "ASC"
			p.nextToken()
		} else if p.currentToken.Type == DESC {
			col.Direction = "DESC"
			p.nextToken()
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

func (p *Parser) parseDropStatement() (Statement, error) {
	pos := token.Pos(p.currentToken.Col)

	if !p.expectToken(DROP) {
		return nil, fmt.Errorf("expected DROP")
	}

	if p.currentToken.Type == INDEX {
		return p.parseDropIndexStatement(pos)
	}

	return nil, fmt.Errorf("expected INDEX after DROP")
}

func (p *Parser) parseDropIndexStatement(pos token.Pos) (*DropIndexStatement, error) {
	stmt := &DropIndexStatement{
		Drop: pos,
	}

	if !p.expectToken(INDEX) {
		return nil, fmt.Errorf("expected INDEX")
	}

	// Check for IF EXISTS
	if p.currentToken.Type == IF {
		p.nextToken()
		if !p.expectToken(EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after IF")
		}
		stmt.IfExists = true
	}

	// Parse index name
	if p.currentToken.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected index name")
	}
	stmt.Name = &Identifier{
		Name: p.currentToken.Value,
		Pos_: token.Pos(p.currentToken.Col),
	}
	p.nextToken()

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

	pos := token.Pos(p.currentToken.Col)

	// Handle NOT IN, NOT BETWEEN, NOT LIKE
	if p.currentToken.Type == NOT {
		p.nextToken()
		switch p.currentToken.Type {
		case IN:
			return p.parseInExpression(left, pos, true)
		case BETWEEN:
			return p.parseBetweenExpression(left, pos, true)
		case LIKE, GLOB:
			operator := "NOT " + p.currentToken.Value
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
		default:
			return nil, fmt.Errorf("expected IN, BETWEEN, or LIKE after NOT")
		}
	}

	// Handle IN
	if p.currentToken.Type == IN {
		return p.parseInExpression(left, pos, false)
	}

	// Handle BETWEEN
	if p.currentToken.Type == BETWEEN {
		return p.parseBetweenExpression(left, pos, false)
	}

	if p.isComparisonOperator() {
		operator := p.currentToken.Value
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

func (p *Parser) parseInExpression(left Expression, pos token.Pos, not bool) (*InExpression, error) {
	p.nextToken() // skip IN

	if !p.expectToken(LPAREN) {
		return nil, fmt.Errorf("expected ( after IN")
	}

	var values []Expression
	for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
		val, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		values = append(values, val)

		if p.currentToken.Type == COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	if !p.expectToken(RPAREN) {
		return nil, fmt.Errorf("expected ) after IN values")
	}

	return &InExpression{
		Expr:   left,
		Values: values,
		Not:    not,
		Pos_:   pos,
	}, nil
}

func (p *Parser) parseBetweenExpression(left Expression, pos token.Pos, not bool) (*BetweenExpression, error) {
	p.nextToken() // skip BETWEEN

	low, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if !p.expectToken(AND) {
		return nil, fmt.Errorf("expected AND in BETWEEN expression")
	}

	high, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return &BetweenExpression{
		Expr: left,
		Low:  low,
		High: high,
		Not:  not,
		Pos_: pos,
	}, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	var name string
	var pos token.Pos

	switch p.currentToken.Type {
	case IDENTIFIER, COUNT, SUM, AVG, MIN, MAX:
		name = p.currentToken.Value
		pos = token.Pos(p.currentToken.Col)
		p.nextToken() // Advance past the identifier

		// Handle function call: COUNT(*), SUM(col), etc.
		if p.currentToken.Type == LPAREN {
			p.nextToken()
			args := []Expression{}

			// Handle COUNT(*) special case
			if p.currentToken.Type == ASTERISK {
				args = append(args, &Identifier{Name: "*", Pos_: token.Pos(p.currentToken.Col)})
				p.nextToken()
			} else if p.currentToken.Type != RPAREN {
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

		// Handle qualified identifier: table.column or table.*
		if p.currentToken.Type == DOT {
			p.nextToken() // skip .
			if p.currentToken.Type == ASTERISK {
				// table.*
				p.nextToken()
				return &QualifiedAsterisk{
					Table: name,
					Pos_:  pos,
				}, nil
			} else if p.currentToken.Type == IDENTIFIER {
				// table.column
				colName := p.currentToken.Value
				p.nextToken()
				return &QualifiedIdentifier{
					Table:  name,
					Column: colName,
					Pos_:   pos,
				}, nil
			} else {
				return nil, fmt.Errorf("expected column name or * after .")
			}
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

	case NULL:
		pos := token.Pos(p.currentToken.Col)
		p.nextToken()
		return &NullLiteral{
			Pos_: pos,
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

	// Handle alias: "table AS alias" or "table alias"
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
	} else if p.currentToken.Type == IDENTIFIER && !p.isReservedKeyword() {
		// Alias without AS keyword: "users u"
		table.Alias = &Identifier{
			Name: p.currentToken.Value,
			Pos_: token.Pos(p.currentToken.Col),
		}
		p.nextToken()
	}

	return table, nil
}

func (p *Parser) isJoinKeyword() bool {
	switch p.currentToken.Type {
	case JOIN, INNER, LEFT, RIGHT, FULL, CROSS:
		return true
	default:
		return false
	}
}

func (p *Parser) parseJoinClause() (*JoinClause, error) {
	join := &JoinClause{}

	// Determine join type
	switch p.currentToken.Type {
	case INNER:
		join.Type = "INNER"
		p.nextToken()
		if p.currentToken.Type != JOIN {
			return nil, fmt.Errorf("expected JOIN after INNER")
		}
		p.nextToken()
	case LEFT:
		join.Type = "LEFT"
		p.nextToken()
		if p.currentToken.Type == OUTER {
			p.nextToken()
		}
		if p.currentToken.Type != JOIN {
			return nil, fmt.Errorf("expected JOIN after LEFT")
		}
		p.nextToken()
	case RIGHT:
		join.Type = "RIGHT"
		p.nextToken()
		if p.currentToken.Type == OUTER {
			p.nextToken()
		}
		if p.currentToken.Type != JOIN {
			return nil, fmt.Errorf("expected JOIN after RIGHT")
		}
		p.nextToken()
	case FULL:
		join.Type = "FULL"
		p.nextToken()
		if p.currentToken.Type == OUTER {
			p.nextToken()
		}
		if p.currentToken.Type != JOIN {
			return nil, fmt.Errorf("expected JOIN after FULL")
		}
		p.nextToken()
	case CROSS:
		join.Type = "CROSS"
		p.nextToken()
		if p.currentToken.Type != JOIN {
			return nil, fmt.Errorf("expected JOIN after CROSS")
		}
		p.nextToken()
	case JOIN:
		join.Type = "INNER" // Default to INNER JOIN
		p.nextToken()
	default:
		return nil, fmt.Errorf("expected JOIN keyword")
	}

	// Parse table reference
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	join.Table = table

	// Parse ON condition (not required for CROSS JOIN)
	if p.currentToken.Type == ON {
		p.nextToken()
		condition, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		join.Condition = condition
	}

	return join, nil
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
	case EQUAL, NOT_EQUAL, NOT_EQUAL2, GREATER, LESS, GREATER_EQUAL, LESS_EQUAL, LIKE, GLOB:
		return true
	default:
		return false
	}
}

// isValidAliasToken returns true if current token can be used as an alias
// This includes IDENTIFIER and common keywords that users might want as aliases
func (p *Parser) isValidAliasToken() bool {
	switch p.currentToken.Type {
	case IDENTIFIER,
		// Aggregate function names commonly used as aliases
		COUNT, SUM, AVG, MIN, MAX,
		// Other common words that might be used as aliases
		TEXT, INTEGER, REAL, BOOLEAN,
		// Query-related words
		ORDER, GROUP, LIMIT, OFFSET,
		// Value-related
		TRUE, FALSE, NULL,
		// Type names
		INT, CHAR, VARCHAR, BLOB, DATETIME, TIMESTAMP:
		return true
	default:
		return false
	}
}

func (p *Parser) isReservedKeyword() bool {
	switch p.currentToken.Type {
	case SELECT, FROM, WHERE, INSERT, INTO, VALUES, UPDATE, SET, DELETE,
		CREATE, TABLE, DROP, ALTER, INDEX, PRIMARY, KEY, FOREIGN, REFERENCES,
		NOT, NULL, DEFAULT, UNIQUE, CHECK, CONSTRAINT,
		ORDER, BY, GROUP, HAVING, LIMIT, OFFSET,
		INNER, LEFT, RIGHT, FULL, OUTER, CROSS, JOIN, ON, AS,
		AND, OR, IN, LIKE, BETWEEN, IS, EXISTS,
		BEGIN, COMMIT, ROLLBACK, TRANSACTION,
		TRUE, FALSE:
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
