package parser

import "github.com/jatin-malik/make-thy-interpreter/token"

// Operator precedences believed by the parser
const (
	LowestPrecedence = iota
	EqualsPrecedence
	LtPrecedence
	SumPrecedence
	AsteriskPrecedence
	CallPrecedence
	IndexPrecedence
)

var precedenceByToken = map[token.TokenType]int{
	token.PLUS:     SumPrecedence,
	token.MINUS:    SumPrecedence,
	token.ASTERISK: AsteriskPrecedence,
	token.SLASH:    AsteriskPrecedence,
	token.EQ:       EqualsPrecedence,
	token.NEQ:      EqualsPrecedence,
	token.LT:       LtPrecedence,
	token.GT:       LtPrecedence,
	token.LPAREN:   CallPrecedence,
	token.LBRACKET: IndexPrecedence,
}

func getTokenPrecedence(tokenType token.TokenType) int {
	precedence, ok := precedenceByToken[tokenType]
	if !ok {
		return LowestPrecedence
	}
	return precedence
}
