package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

// ComparisonExpr ::= StringConcatExpr ( (ValueComp | GeneralComp | NodeComp) StringConcatExpr )?
// typeID ::= 														1					| 2						| 3
// The token field is a private field because typeID should be determined by the token value. So, typeID field set when token field is set.
type ComparisonExpr struct {
	LeftExpr  ExprSingle
	token     token.Token
	RightExpr ExprSingle
	typeID    byte
}

func (ce *ComparisonExpr) exprSingle() {}
func (ce *ComparisonExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ce.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(ce.token.Literal)
	sb.WriteString(" ")
	sb.WriteString(ce.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// TypeID is a getter for the typeID field
func (ce *ComparisonExpr) TypeID() byte {
	return ce.typeID
}

// Token is a getter for the token field
func (ce *ComparisonExpr) Token() token.Token {
	return ce.token
}

// SetToken is a setter for the token field
func (ce *ComparisonExpr) SetToken(t token.Token) {
	ce.token = t

	if util.IsValueComp(ce.token.Literal) {
		ce.typeID = 1
	} else if util.IsGeneralComp(ce.token.Literal) {
		ce.typeID = 2
	} else if util.IsNodeComp(ce.token.Literal) {
		ce.typeID = 3
	}
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
