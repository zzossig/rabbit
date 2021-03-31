package ast

import (
	"strings"

	"github.com/zzossig/rabbit/token"
)

// SimpleMapExpr ::= PathExpr ("!" PathExpr)*
type SimpleMapExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (sme *SimpleMapExpr) exprSingle() {}
func (sme *SimpleMapExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(sme.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(sme.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(sme.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}
