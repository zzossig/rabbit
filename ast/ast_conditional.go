package ast

// IfExpr ::= "if" "(" Expr ")" "then" ExprSingle "else" ExprSingle
type IfExpr struct {
	TestExpr Expr
	ThenExpr ExprSingle
	ElseExpr ExprSingle
}

func (ie *IfExpr) exprSingle() {}
func (ie *IfExpr) argument()   {}
