package ast

import (
	"strings"
)

// ArrayConstructor ::= SquareArrayConstructor | CurlyArrayConstructor
// TypeID ::= 					1											 | 2
type ArrayConstructor struct {
	SquareArrayConstructor
	CurlyArrayConstructor
	TypeID byte
}

func (ac *ArrayConstructor) exprSingle()  {}
func (ac *ArrayConstructor) primaryExpr() {}
func (ac *ArrayConstructor) String() string {
	switch ac.TypeID {
	case 1:
		return ac.SquareArrayConstructor.String()
	case 2:
		return ac.CurlyArrayConstructor.String()
	default:
		return ""
	}
}

// SquareArrayConstructor ::= "[" (ExprSingle ("," ExprSingle)*)? "]"
type SquareArrayConstructor struct {
	Exprs []ExprSingle
}

func (sac *SquareArrayConstructor) exprSingle()  {}
func (sac *SquareArrayConstructor) primaryExpr() {}
func (sac *SquareArrayConstructor) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, expr := range sac.Exprs {
		sb.WriteString(expr.String())
		if i < len(sac.Exprs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")

	return sb.String()
}

// CurlyArrayConstructor ::= "array" EnclosedExpr
type CurlyArrayConstructor struct {
	EnclosedExpr
}

func (cac *CurlyArrayConstructor) exprSingle()  {}
func (cac *CurlyArrayConstructor) primaryExpr() {}
func (cac *CurlyArrayConstructor) String() string {
	var sb strings.Builder

	sb.WriteString("array")
	sb.WriteString(cac.EnclosedExpr.String())

	return sb.String()
}
