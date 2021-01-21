package ast

import (
	"github.com/zzossig/xpath/token"
)

// QuantifiedExpr ::= ("some" | "every") "$" VarName "in" ExprSingle ("," "$" VarName "in" ExprSingle)* "satisfies" ExprSingle
type QuantifiedExpr struct {
	Token         token.Token
	SimpleQClause SimpleQClause
	SatisExpr     ExprSingle
}

func (qe *QuantifiedExpr) exprSingle() {}
func (qe *QuantifiedExpr) argument()   {}

// SimpleQClause ::= "$" VarName "in" ExprSingle ("," "$" VarName "in" ExprSingle)*
type SimpleQClause struct {
	Bindings []SimpleQBinding
}

// SimpleQBinding ::= "$" VarName "in" ExprSingle
type SimpleQBinding struct {
	VarName    VarName
	ExprSingle ExprSingle
}
