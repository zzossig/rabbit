package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// Expr ::= ExprSingle ("," ExprSingle)*
type Expr struct {
	Exprs []ExprSingle
}

func (e *Expr) exprSingle() {}
func (e *Expr) String() string {
	var sb strings.Builder

	for i, expr := range e.Exprs {
		sb.WriteString(expr.String())
		if i < len(e.Exprs)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

// RangeExpr ::= AdditiveExpr ( "to" AdditiveExpr )?
type RangeExpr struct {
	LeftExpr  ExprSingle
	Token     token.Token // token.TO
	RightExpr ExprSingle
}

func (re *RangeExpr) exprSingle() {}
func (re *RangeExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(re.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(re.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(re.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// UnionExpr ::= IntersectExceptExpr ( ("union" | "|") IntersectExceptExpr )*
type UnionExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (ue *UnionExpr) exprSingle() {}
func (ue *UnionExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ue.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(ue.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(ue.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// IntersectExceptExpr ::= InstanceofExpr ( ("intersect" | "except") InstanceofExpr )*
type IntersectExceptExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (iee *IntersectExceptExpr) exprSingle() {}
func (iee *IntersectExceptExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(iee.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(iee.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(iee.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}
