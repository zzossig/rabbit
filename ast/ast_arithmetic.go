package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// AdditiveExpr ::= MultiplicativeExpr ( ("+" | "-") MultiplicativeExpr )*
type AdditiveExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (ae *AdditiveExpr) exprSingle() {}
func (ae *AdditiveExpr) String() string {
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

// MultiplicativeExpr ::= UnionExpr ( ("*" | "div" | "idiv" | "mod") UnionExpr )*
type MultiplicativeExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (me *MultiplicativeExpr) exprSingle() {}
func (me *MultiplicativeExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(me.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(me.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(me.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// UnaryExpr ::= ("-" | "+")* ValueExpr
type UnaryExpr struct {
	ExprSingle
	Token token.Token
}

func (ue *UnaryExpr) exprSingle() {}
func (ue *UnaryExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ue.Token.Literal)
	sb.WriteString(ue.ExprSingle.String())
	sb.WriteString(")")

	return sb.String()
}

// ValueExpr ::= SimpleMapExpr
type ValueExpr = SimpleMapExpr
