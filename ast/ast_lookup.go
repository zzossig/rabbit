package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// UnaryLookup ::= "?" KeySpecifier
// Unary lookup is used in predicates (e.g. $map[?name='Mike'] or with the simple map operator (e.g. $maps ! ?name='Mike').
type UnaryLookup struct {
	Token token.Token // token.UQUESTION
	KeySpecifier
}

func (ul *UnaryLookup) exprSingle() {}
func (ul *UnaryLookup) String() string {
	var sb strings.Builder

	sb.WriteString("?")
	sb.WriteString(ul.KeySpecifier.String())

	return sb.String()
}

// KeySpecifier ::= NCName | IntegerLiteral | ParenthesizedExpr | "*"
// TypeID ::=				1			 | 2							| 3									| 4
type KeySpecifier struct {
	NCName
	IntegerLiteral
	ParenthesizedExpr
	TypeID byte
}

func (ks *KeySpecifier) String() string {
	switch ks.TypeID {
	case 1:
		return ks.NCName.Value()
	case 2:
		return ks.IntegerLiteral.String()
	case 3:
		return ks.ParenthesizedExpr.String()
	case 4:
		return "*"
	default:
		return ""
	}
}
