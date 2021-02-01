package ast

// XPath ::= Expr
type XPath struct {
	Expr
}

func (x *XPath) exprSingle() {}
func (x *XPath) String() string {
	return x.Expr.String()
}

// ExprSingle ::= ForExpr | LetExpr | QuantifiedExpr | IfExpr | OrExpr
type ExprSingle interface {
	exprSingle()
	String() string
}
