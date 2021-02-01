package ast

import "strings"

// MapConstructor ::= "map" "{" (MapConstructorEntry ("," MapConstructorEntry)*)? "}"
type MapConstructor struct {
	Entries []MapConstructorEntry
}

func (mc *MapConstructor) exprSingle()  {}
func (mc *MapConstructor) primaryExpr() {}
func (mc *MapConstructor) String() string {
	var sb strings.Builder

	sb.WriteString("map{")
	for i, e := range mc.Entries {
		sb.WriteString(e.String())
		if i < len(mc.Entries)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")

	return sb.String()
}

// MapConstructorEntry ::= MapKeyExpr ":" MapValueExpr
type MapConstructorEntry struct {
	MapKeyExpr
	MapValueExpr
}

func (mce *MapConstructorEntry) String() string {
	var sb strings.Builder

	sb.WriteString(mce.MapKeyExpr.String())
	sb.WriteString(":")
	sb.WriteString(" ")
	sb.WriteString(mce.MapValueExpr.String())

	return sb.String()
}

// MapKeyExpr ::= ExprSingle
type MapKeyExpr struct {
	ExprSingle
}

func (mke *MapKeyExpr) exprSingle() {}
func (mke *MapKeyExpr) String() string {
	return mke.ExprSingle.String()
}

// MapValueExpr ::= ExprSingle
type MapValueExpr struct {
	ExprSingle
}

func (mve *MapValueExpr) exprSingle() {}
func (mve *MapValueExpr) String() string {
	return mve.ExprSingle.String()
}
