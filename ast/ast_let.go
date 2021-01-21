package ast

// LetExpr ::= SimpleLetClause "return" ExprSingle
type LetExpr struct {
	SimpleLetClause SimpleLetClause
	ReturnExpr      ExprSingle
}

func (le *LetExpr) exprSingle() {}
func (le *LetExpr) argument()   {}

// SimpleLetClause ::= "let" SimpleLetBinding ("," SimpleLetBinding)*
type SimpleLetClause struct {
	Bindings []SimpleLetBinding
}

// SimpleLetBinding ::= "$" VarName ":=" ExprSingle
type SimpleLetBinding struct {
	VarName    VarName
	ExprSingle ExprSingle
}
