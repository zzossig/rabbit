package ast

import (
	"fmt"
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
// TypeID ::= 1			| 2
type EQName struct {
	QName
	URIQualifiedName
	TypeID byte
}

// Value is a getter for the value field
func (eqn *EQName) Value() string {
	switch eqn.TypeID {
	case 1:
		return eqn.QName.Value()
	case 2:
		return eqn.URIQualifiedName.Value()
	default:
		return ""
	}
}

// SetValue is a setter for the value field
func (eqn *EQName) SetValue(name string) {
	if util.IsEQName(name) {
		if strings.HasPrefix(name, "Q{") {
			eqn.TypeID = 2
			eqn.URIQualifiedName.SetValue(name)
		} else {
			eqn.TypeID = 1
			eqn.QName.SetValue(name)
		}
	} else {
		// TODO occur error
	}
}

// QName ::= [http://www.w3.org/TR/REC-xml-names/#NT-QName]
type QName struct {
	prefix string
	local  string
}

// Prefix is a getter for the prefix field
func (qn *QName) Prefix() string {
	return qn.prefix
}

// Prefix is a getter for the local field
func (qn *QName) Local() string {
	return qn.local
}

// Value is a getter for the QName
func (qn *QName) Value() string {
	if qn.prefix != "" {
		return fmt.Sprintf("%s:%s", qn.prefix, qn.local)
	}
	return qn.local
}

// SetValue is a setter for the QName
func (qn *QName) SetValue(v string) {
	if !util.IsQName(v) {
		// TODO occur error
	}
	if strings.Contains(v, ":") {
		s := strings.Split(v, ":")
		qn.prefix = s[0]
		qn.local = s[1]
	} else {
		qn.local = v
	}
}

// SetPrefix is a setter for the prefix field
func (qn *QName) SetPrefix(v string) {
	if util.IsNCName(v) {
		qn.prefix = v
	} else {
		// TODO error
	}
}

// SetLocal is a setter for the local field
func (qn *QName) SetLocal(v string) {
	if util.IsNCName(v) {
		qn.local = v
	} else {
		// TODO error
	}
}

// URIQualifiedName ::= BracedURILiteral NCName
type URIQualifiedName struct {
	BracedURILiteral
	NCName
}

// Value is a getter for the fields
func (u *URIQualifiedName) Value() string {
	return fmt.Sprintf("%s%s", u.BracedURILiteral.Value(), u.NCName.Value())
}

// SetValue is a setter for the fields
func (u *URIQualifiedName) SetValue(name string) {
	if util.IsURIQualifiedName(name) {
		names := strings.SplitAfter(name, "}")
		u.BracedURILiteral.SetValue(names[0])
		u.NCName.SetValue(names[1])
	} else {
		// error
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

// URI returns uri in Q{uri}
func (b *BracedURILiteral) URI() string {
	if b.value != "" {
		return b.value[2 : len(b.value)-1]
	}
	return b.value
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
