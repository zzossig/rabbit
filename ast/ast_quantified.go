package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// QuantifiedExpr ::= ("some" | "every") "$" VarName "in" ExprSingle ("," "$" VarName "in" ExprSingle)* "satisfies" ExprSingle
type QuantifiedExpr struct {
	Token      token.Token
	ExprSingle // satisfies ExprSingle
	SimpleQClause
}

func (qe *QuantifiedExpr) exprSingle() {}
func (qe *QuantifiedExpr) String() string {
	var sb strings.Builder

	sb.WriteString(qe.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(qe.SimpleQClause.String())
	sb.WriteString(" ")
	sb.WriteString("satisfies")
	sb.WriteString(" ")
	sb.WriteString(qe.ExprSingle.String())

	return sb.String()
}

// SimpleQClause ::= "$" VarName "in" ExprSingle ("," "$" VarName "in" ExprSingle)*
type SimpleQClause struct {
	Bindings []SimpleQBinding
}

func (sqc *SimpleQClause) String() string {
	var sb strings.Builder

	for i, b := range sqc.Bindings {
		sb.WriteString(b.String())
		if i < len(sqc.Bindings)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

// SimpleQBinding ::= "$" VarName "in" ExprSingle
type SimpleQBinding struct {
	VarName
	ExprSingle
}

func (sqb *SimpleQBinding) String() string {
	var sb strings.Builder

	sb.WriteString("$")
	sb.WriteString(sqb.VarName.Value())
	sb.WriteString(" ")
	sb.WriteString("in")
	sb.WriteString(" ")
	sb.WriteString(sqb.ExprSingle.String())

	return sb.String()
}
