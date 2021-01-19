package ast

// XPath ::= Expr
type XPath Expr

// Expr ::= ExprSingle ("," ExprSingle)*
type Expr struct {
	Sequence []ExprSingle
}

// ExprSingle ::= ForExpr | LetExpr | QuantifiedExpr | IfExpr | OrExpr
type ExprSingle interface {
	exprSingle()
}

// Identifier is a string of alphanumeric characters.
type Identifier struct {
	Value string
}

// VarName ::= EQName
type VarName EQName

// EQName ::= QName | URIQualifiedName
type EQName struct {
	QName QName
}
