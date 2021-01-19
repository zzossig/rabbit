package ast

// ForExpr ::= SimpleForClause "return" ExprSingle
type ForExpr struct {
	SimpleForClause SimpleForClause
	ExprSingle      ExprSingle
}

func (fe *ForExpr) exprSingle() {}

// SimpleForClause ::= "for" SimpleForBinding ("," SimpleForBinding)*
type SimpleForClause struct {
	Bindings []SimpleForBinding
}

// SimpleForBinding ::= "$" VarName "in" ExprSingle
type SimpleForBinding struct {
	Name       VarName
	ExprSingle ExprSingle
}
