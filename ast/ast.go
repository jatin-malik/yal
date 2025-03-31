package ast

import (
	"bytes"
	"fmt"
	"github.com/jatin-malik/yal/token"
	"strings"
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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl StringLiteral) expressionBehaviour() {}

func (sl StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

func (sl StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
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

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
	Name       string
}

func (fl FunctionLiteral) expressionBehaviour() {}

func (fl FunctionLiteral) String() string {
	var buf bytes.Buffer
	buf.WriteString(fl.TokenLiteral() + " ")
	buf.WriteString("(")
	for i, parameter := range fl.Parameters {
		if i > 0 {
			buf.WriteString(", ")

		}
		buf.WriteString(parameter.Value)
	}
	buf.WriteString(") ")
	buf.WriteString(fl.Body.String())
	return buf.String()
}

func (fl FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml MacroLiteral) expressionBehaviour() {}

func (ml MacroLiteral) String() string {
	var buf bytes.Buffer
	buf.WriteString(ml.TokenLiteral() + " ")
	buf.WriteString("(")
	for i, parameter := range ml.Parameters {
		if i > 0 {
			buf.WriteString(", ")

		}
		buf.WriteString(parameter.Value)
	}
	buf.WriteString(") ")
	buf.WriteString(ml.Body.String())
	return buf.String()
}

func (ml MacroLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

type IfElseConditional struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (iec IfElseConditional) expressionBehaviour() {}

func (iec IfElseConditional) String() string {
	var buf bytes.Buffer
	buf.WriteString(iec.TokenLiteral() + " ")
	buf.WriteString(iec.Condition.String())
	if iec.Consequence != nil {
		buf.WriteString(iec.Consequence.String())

	}
	if iec.Alternative != nil {
		buf.WriteString(" else ")
		buf.WriteString(iec.Alternative.String())
	}
	return buf.String()
}

func (iec IfElseConditional) TokenLiteral() string {
	return iec.Token.Literal
}

type LoopStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (l LoopStatement) statementBehaviour() {}

func (l LoopStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("for ")
	buf.WriteString("(")
	buf.WriteString(l.Condition.String())
	buf.WriteString(")")
	buf.WriteString("{")
	buf.WriteString(l.Body.String())
	buf.WriteString("}")
	return buf.String()
}

func (l LoopStatement) TokenLiteral() string {
	return l.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) expressionBehaviour() {}

func (pe PrefixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("( ")
	buf.WriteString(pe.Operator)
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

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (all ArrayLiteral) expressionBehaviour() {}
func (all ArrayLiteral) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, element := range all.Elements {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(element.String())

	}
	buf.WriteString("]")
	return buf.String()
}

func (all ArrayLiteral) TokenLiteral() string {
	return all.Token.Literal
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl HashLiteral) expressionBehaviour() {}

func (hl HashLiteral) String() string {
	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (hl HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce CallExpression) expressionBehaviour() {}

func (ce CallExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString(ce.Function.String())
	buf.WriteString("(")
	if len(ce.Arguments) > 0 {
		for i, a := range ce.Arguments {
			if i > 0 {
				buf.WriteString(", ")

			}
			buf.WriteString(a.String())
		}
	}
	buf.WriteString(")")
	return buf.String()
}

func (ce CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie IndexExpression) expressionBehaviour() {}

func (ie IndexExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString(ie.Left.String())
	buf.WriteString("[")
	buf.WriteString(ie.Index.String())
	buf.WriteString("]")
	return buf.String()
}

func (ie IndexExpression) TokenLiteral() string {
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

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs BlockStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("{ ")
	for _, stmt := range bs.Statements {
		buf.WriteString(stmt.String() + " ")
	}
	buf.WriteString("}")
	return buf.String()
}

// TODO: Is this a statement or an expression ?
func (bs BlockStatement) statementBehaviour() {}
