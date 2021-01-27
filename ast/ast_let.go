package ast

import "strings"

// LetExpr ::= SimpleLetClause "return" ExprSingle
type LetExpr struct {
	SimpleLetClause
	ExprSingle
}

func (le *LetExpr) exprSingle() {}
func (le *LetExpr) String() string {
	var sb strings.Builder

	sb.WriteString(le.SimpleLetClause.String())
	sb.WriteString(" ")
	sb.WriteString("return")
	sb.WriteString(" ")
	sb.WriteString(le.ExprSingle.String())

	return sb.String()
}

// SimpleLetClause ::= "let" SimpleLetBinding ("," SimpleLetBinding)*
type SimpleLetClause struct {
	Bindings []SimpleLetBinding
}

func (slc *SimpleLetClause) String() string {
	var sb strings.Builder

	sb.WriteString("let")
	sb.WriteString(" ")
	for i, b := range slc.Bindings {
		sb.WriteString(b.String())
		if i < len(slc.Bindings)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

// SimpleLetBinding ::= "$" VarName ":=" ExprSingle
type SimpleLetBinding struct {
	VarName
	ExprSingle
}

func (slb *SimpleLetBinding) String() string {
	var sb strings.Builder

	sb.WriteString("$")
	sb.WriteString(slb.VarName.Value())
	sb.WriteString(" ")
	sb.WriteString(":=")
	sb.WriteString(" ")
	sb.WriteString(slb.ExprSingle.String())

	return sb.String()
}
