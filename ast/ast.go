package ast

// XPath ::= Expr
type XPath Expr

// Item is either an atomic value, a node, or a function
type Item interface {
	item()
}

// ExprSingle ::= ForExpr | LetExpr | QuantifiedExpr | IfExpr | OrExpr
type ExprSingle interface {
	Argument
	exprSingle()
}
