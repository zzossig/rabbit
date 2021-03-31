package ast

import (
	"strings"

	"github.com/zzossig/rabbit/token"
	"github.com/zzossig/rabbit/util"
)

// ComparisonExpr ::= StringConcatExpr ( (ValueComp | GeneralComp | NodeComp) StringConcatExpr )?
type ComparisonExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (ce *ComparisonExpr) exprSingle() {}
func (ce *ComparisonExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ce.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(ce.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(ce.RightExpr.String())
	sb.WriteString(")")

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
