package ast

// https://www.w3.org/TR/xpath-31/#id-grammar

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
