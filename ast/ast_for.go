package ast

import "strings"

// ForExpr ::= SimpleForClause "return" ExprSingle
type ForExpr struct {
	SimpleForClause
	ExprSingle
}

func (fe *ForExpr) exprSingle() {}
func (fe *ForExpr) String() string {
	var sb strings.Builder

	sb.WriteString(fe.SimpleForClause.String())
	sb.WriteString(" ")
	sb.WriteString("return")
	sb.WriteString(" ")
	sb.WriteString(fe.ExprSingle.String())

	return sb.String()
}

// SimpleForClause ::= "for" SimpleForBinding ("," SimpleForBinding)*
type SimpleForClause struct {
	Bindings []SimpleForBinding
}

func (sfc *SimpleForClause) String() string {
	var sb strings.Builder

	sb.WriteString("for")
	sb.WriteString(" ")
	for i, b := range sfc.Bindings {
		sb.WriteString(b.String())
		if i < len(sfc.Bindings)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

// SimpleForBinding ::= "$" VarName "in" ExprSingle
type SimpleForBinding struct {
	VarName
	ExprSingle
}

func (sfb *SimpleForBinding) String() string {
	var sb strings.Builder

	sb.WriteString("$")
	sb.WriteString(sfb.VarName.Value())
	sb.WriteString(" ")
	sb.WriteString("in")
	sb.WriteString(" ")
	sb.WriteString(sfb.ExprSingle.String())

	return sb.String()
}
