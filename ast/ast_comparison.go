package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

// ComparisonExpr ::= StringConcatExpr ( (ValueComp | GeneralComp | NodeComp) StringConcatExpr )?
// TypeID ::= 														1					| 2						| 3
type ComparisonExpr struct {
	LeftExpr  ExprSingle
	Token     token.Token
	RightExpr ExprSingle
	TypeID    byte
}

func (ce *ComparisonExpr) exprSingle() {}
func (ce *ComparisonExpr) String() string {
	var sb strings.Builder

	sb.WriteString(ce.LeftExpr.String())
	sb.WriteString(ce.Token.Literal)
	sb.WriteString(ce.RightExpr.String())

	return sb.String()
}

// GeneralComp ::= "=" | "!=" | "<" | "<=" | ">" | ">="
type GeneralComp struct {
	value token.Token
}

// Value is a getter for the value field
func (gc *GeneralComp) Value() token.Token {
	return gc.value
}

// SetValue is a setter for the value field
func (gc *GeneralComp) SetValue(t token.Token) {
	if util.IsGeneralComp(t.Literal) {
		gc.value = t
	} else {
		// TODO error
	}

}

// ValueComp ::= "eq" | "ne" | "lt" | "le" | "gt" | "ge"
type ValueComp struct {
	value token.Token
}

// Value is a getter for the value field
func (vc *ValueComp) Value() token.Token {
	return vc.value
}

// SetValue is a setter for the value field
func (vc *ValueComp) SetValue(t token.Token) {
	if util.IsValueComp(t.Literal) {
		vc.value = t
	} else {
		// TODO error
	}
}

// NodeComp ::= "is" | "<<" | ">>"
type NodeComp struct {
	value token.Token
}

// Value is a getter for the value field
func (nc *NodeComp) Value() token.Token {
	return nc.value
}

// SetValue is a setter for the value field
func (nc *NodeComp) SetValue(t token.Token) {
	if util.IsNodeComp(t.Literal) {
		nc.value = t
	} else {
		// TODO error
	}
}
