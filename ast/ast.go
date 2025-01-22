package ast

import (
	"bytes"
	"github.com/jatin-malik/make-thy-interpreter/token"
)

type Node interface {
	TokenLiteral() string // for debugging
	String() string       // for debugging
}

// Statement represents a single statement in the program.
type Statement interface {
	Node                 // A statement is a node in the AST.
	statementBehaviour() //TODO: This is just to guide us during dev with compile time type checks. Remove once the parser is complete.
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Right Expression
}

func (letStmt LetStatement) TokenLiteral() string {
	return letStmt.Token.Literal
}

func (letStmt LetStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString(letStmt.TokenLiteral() + " ")
	buf.WriteString(letStmt.Name.Value + " ")
	buf.WriteString("= ")
	buf.WriteString(letStmt.Right.String())
	buf.WriteString(";")
	return buf.String()
}

func (letStmt LetStatement) statementBehaviour() {
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (returnStmt ReturnStatement) TokenLiteral() string {
	return returnStmt.Token.Literal
}

func (returnStmt ReturnStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString(returnStmt.TokenLiteral() + " ")
	buf.WriteString(returnStmt.Value.String())
	buf.WriteString(";")
	return buf.String()
}

func (returnStmt ReturnStatement) statementBehaviour() {
}

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (expStmt *ExpressionStatement) TokenLiteral() string {
	return expStmt.Token.Literal
}

func (expStmt *ExpressionStatement) String() string {
	if expStmt.Expr == nil {
		return ""
	}
	return expStmt.Expr.String()
}

func (expStmt *ExpressionStatement) statementBehaviour() {}

// Expression represents a generic expression in the program.
type Expression interface {
	Node                  // An expression is a node in the AST.
	expressionBehaviour() //TODO: This is just to guide us during dev with compile time type checks. Remove once the parser is complete.
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expressionBehaviour() {}
func (il IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl BooleanLiteral) expressionBehaviour() {}
func (bl BooleanLiteral) String() string {
	return bl.Token.Literal
}

func (bl BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (pe PrefixExpression) expressionBehaviour() {}

func (pe PrefixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("( ")
	buf.WriteString(pe.TokenLiteral())
	buf.WriteString(pe.Right.String())
	buf.WriteString(" )")
	return buf.String()
}

func (pe PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionBehaviour() {}
func (ie InfixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("( ")
	buf.WriteString(ie.Left.String())
	buf.WriteString(" ")
	buf.WriteString(ie.Operator)
	buf.WriteString(" ")
	buf.WriteString(ie.Right.String())
	buf.WriteString(" )")
	return buf.String()
}

func (ie InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

type Identifier struct {
	Token token.Token // for debugging
	Value string      // Literal value of the identifier
}

func (ident Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

func (ident Identifier) String() string {
	return ident.Value
}

func (ident Identifier) expressionBehaviour() {
}

// Program is the root node of the AST.
type Program struct {
	Statements []Statement
}

func (prg Program) TokenLiteral() string {
	if len(prg.Statements) > 0 {
		return prg.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (prg Program) String() string {
	var buf bytes.Buffer

	for _, stmt := range prg.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}
