package ast

// ForExpr ::= SimpleForClause "return" ExprSingle
type ForExpr struct {
	SimpleForClause SimpleForClause
	ReturnExpr      ExprSingle
}

func (fe *ForExpr) exprSingle() {}
func (fe *ForExpr) argument()   {}

// SimpleForClause ::= "for" SimpleForBinding ("," SimpleForBinding)*
type SimpleForClause struct {
	Bindings []SimpleForBinding
}

// SimpleForBinding ::= "$" VarName "in" ExprSingle
type SimpleForBinding struct {
	VarName    VarName
	ExprSingle ExprSingle
}
