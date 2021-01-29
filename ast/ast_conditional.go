package ast

import "strings"

// IfExpr ::= "if" "(" Expr ")" "then" ExprSingle "else" ExprSingle
type IfExpr struct {
	TestExpr ExprSingle
	ThenExpr ExprSingle
	ElseExpr ExprSingle
}

func (ie *IfExpr) exprSingle() {}
func (ie *IfExpr) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString("(")
	sb.WriteString(ie.TestExpr.String())
	sb.WriteString(")")
	sb.WriteString(" then ")
	sb.WriteString(ie.ThenExpr.String())
	sb.WriteString(" else ")
	sb.WriteString(ie.ElseExpr.String())

	return sb.String()
}
