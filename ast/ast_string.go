package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

// StringConcatExpr ::= RangeExpr ( "||" RangeExpr )*
type StringConcatExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (sce *StringConcatExpr) exprSingle() {}
func (sce *StringConcatExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(sce.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(sce.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(sce.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// EQName ::= QName | URIQualifiedName
type EQName struct {
	value string
}

// Value is a getter for the value field
func (eqn *EQName) Value() string {
	return eqn.value
}

// SetValue is a setter for the value field
func (eqn *EQName) SetValue(name string) {
	if util.IsEQName(name) {
		eqn.value = name
	} else {
		// TODO occur error
	}
}

// NCName ::= Name - (Char* ':' Char*)
type NCName struct {
	value string
}

// Value is a getter for the value field
func (ncn *NCName) Value() string {
	return ncn.value
}

// SetValue is a setter for the value field
func (ncn *NCName) SetValue(name string) {
	if util.IsNCName(name) {
		ncn.value = name
	} else {
		// TODO occur error
	}
}

// BracedURILiteral ::= "Q" "{" [^{}]* "}"
type BracedURILiteral struct {
	value string
}

// Value is a getter for the value field
func (b *BracedURILiteral) Value() string {
	return b.value
}

// SetValue is a setter for the value field
func (b *BracedURILiteral) SetValue(name string) {
	if util.IsBracedURILiteral(name) {
		b.value = name
	} else {
		// TODO occur error
	}
}

// SimpleTypeName ::= TypeName
type SimpleTypeName = TypeName
