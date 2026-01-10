package citrinelexer

import (
	"go/ast"
	"go/token"
)

type Node interface {
	ast.Node
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// SELECT statement
type SelectStatement struct {
	Select  token.Pos
	Fields  []Expression
	From    *TableRef
	Where   Expression
	GroupBy []Expression
	Having  Expression
	OrderBy []OrderByItem
	Limit   *LimitClause
}

func (s *SelectStatement) Pos() token.Pos { return s.Select }
func (s *SelectStatement) End() token.Pos { return token.NoPos }
func (s *SelectStatement) String() string { return "SELECT" }
func (s *SelectStatement) statementNode() {}

// CREATE TABLE statement
type CreateTableStatement struct {
	Create      token.Pos
	Table       *Identifier
	Columns     []*ColumnDef
	Constraints []Constraint
}

func (c *CreateTableStatement) Pos() token.Pos { return c.Create }
func (c *CreateTableStatement) End() token.Pos { return token.NoPos }
func (c *CreateTableStatement) String() string { return "CREATE TABLE" }
func (c *CreateTableStatement) statementNode() {}

// INSERT statement
type InsertStatement struct {
	Insert  token.Pos
	Table   *Identifier
	Columns []*Identifier
	Values  [][]Expression
}

func (i *InsertStatement) Pos() token.Pos { return i.Insert }
func (i *InsertStatement) End() token.Pos { return token.NoPos }
func (i *InsertStatement) String() string { return "INSERT" }
func (i *InsertStatement) statementNode() {}

// UPDATE statement
type UpdateStatement struct {
	Update token.Pos
	Table  *Identifier
	Set    []*Assignment
	Where  Expression
}

func (u *UpdateStatement) Pos() token.Pos { return u.Update }
func (u *UpdateStatement) End() token.Pos { return token.NoPos }
func (u *UpdateStatement) String() string { return "UPDATE" }
func (u *UpdateStatement) statementNode() {}

// DELETE statement
type DeleteStatement struct {
	Delete token.Pos
	From   *Identifier
	Where  Expression
}

func (d *DeleteStatement) Pos() token.Pos { return d.Delete }
func (d *DeleteStatement) End() token.Pos { return token.NoPos }
func (d *DeleteStatement) String() string { return "DELETE" }
func (d *DeleteStatement) statementNode() {}

// Expressions
type Identifier struct {
	Name string
	Pos_ token.Pos
}

func (i *Identifier) Pos() token.Pos  { return i.Pos_ }
func (i *Identifier) End() token.Pos  { return token.NoPos }
func (i *Identifier) String() string  { return i.Name }
func (i *Identifier) expressionNode() {}

type StringLiteral struct {
	Value string
	Pos_  token.Pos
}

func (s *StringLiteral) Pos() token.Pos  { return s.Pos_ }
func (s *StringLiteral) End() token.Pos  { return token.NoPos }
func (s *StringLiteral) String() string  { return "'" + s.Value + "'" }
func (s *StringLiteral) expressionNode() {}

type NumberLiteral struct {
	Value string
	Pos_  token.Pos
}

func (n *NumberLiteral) Pos() token.Pos  { return n.Pos_ }
func (n *NumberLiteral) End() token.Pos  { return token.NoPos }
func (n *NumberLiteral) String() string  { return n.Value }
func (n *NumberLiteral) expressionNode() {}

type BooleanLiteral struct {
	Value bool
	Pos_  token.Pos
}

func (b *BooleanLiteral) Pos() token.Pos { return b.Pos_ }
func (b *BooleanLiteral) End() token.Pos { return token.NoPos }
func (b *BooleanLiteral) String() string {
	if b.Value {
		return "TRUE"
	}
	return "FALSE"
}
func (b *BooleanLiteral) expressionNode() {}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Pos_     token.Pos
}

func (b *BinaryExpression) Pos() token.Pos { return b.Pos_ }
func (b *BinaryExpression) End() token.Pos { return token.NoPos }
func (b *BinaryExpression) String() string {
	return b.Left.String() + " " + b.Operator + " " + b.Right.String()
}
func (b *BinaryExpression) expressionNode() {}

type FunctionCall struct {
	Name string
	Args []Expression
	Pos_ token.Pos
}

func (f *FunctionCall) Pos() token.Pos  { return f.Pos_ }
func (f *FunctionCall) End() token.Pos  { return token.NoPos }
func (f *FunctionCall) String() string  { return f.Name + "()" }
func (f *FunctionCall) expressionNode() {}

// Supporting types
type TableRef struct {
	Name  *Identifier
	Alias *Identifier
}

type ColumnDef struct {
	Name        *Identifier
	Type        string
	Constraints []Constraint
}

type Constraint interface {
	Node
	constraintNode()
}

type PrimaryKeyConstraint struct {
	Pos_ token.Pos
}

func (p *PrimaryKeyConstraint) Pos() token.Pos  { return p.Pos_ }
func (p *PrimaryKeyConstraint) End() token.Pos  { return token.NoPos }
func (p *PrimaryKeyConstraint) String() string  { return "PRIMARY KEY" }
func (p *PrimaryKeyConstraint) constraintNode() {}

type NotNullConstraint struct {
	Pos_ token.Pos
}

func (n *NotNullConstraint) Pos() token.Pos  { return n.Pos_ }
func (n *NotNullConstraint) End() token.Pos  { return token.NoPos }
func (n *NotNullConstraint) String() string  { return "NOT NULL" }
func (n *NotNullConstraint) constraintNode() {}

type Assignment struct {
	Column *Identifier
	Value  Expression
}

type OrderByItem struct {
	Expression Expression
	Direction  string // "ASC" or "DESC"
}

type LimitClause struct {
	Count  Expression
	Offset Expression
}

// Parameter placeholder
type Parameter struct {
	Name string // for named parameters (:name, $name)
	Pos_ token.Pos
}

func (p *Parameter) Pos() token.Pos { return p.Pos_ }
func (p *Parameter) End() token.Pos { return token.NoPos }
func (p *Parameter) String() string {
	if p.Name == "" {
		return "?"
	}
	return p.Name
}
func (p *Parameter) expressionNode() {}
