package ast

import (
	"strings"

	"github.com/zzossig/rabbit/token"
)

// OrExpr ::= AndExpr ( "or" AndExpr )*
type OrExpr struct {
	LeftExpr  ExprSingle
	Token     token.Token // token.OR
	RightExpr ExprSingle
}

func (oe *OrExpr) exprSingle() {}
func (oe *OrExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(oe.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(oe.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(oe.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// AndExpr ::= ComparisonExpr ( "and" ComparisonExpr )*
type AndExpr struct {
	LeftExpr  ExprSingle
	Token     token.Token // token.AND
	RightExpr ExprSingle
}

func (ae *AndExpr) exprSingle() {}
func (ae *AndExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ae.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(ae.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(ae.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}
