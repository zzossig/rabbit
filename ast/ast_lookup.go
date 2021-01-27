package ast

import "strings"

// Lookup ::= "?" KeySpecifier
type Lookup struct {
	KeySpecifier
}

func (l *Lookup) pal() {}
func (l *Lookup) String() string {
	var sb strings.Builder

	sb.WriteString("?")
	sb.WriteString(l.KeySpecifier.String())

	return sb.String()
}

// UnaryLookup ::= "?" KeySpecifier
type UnaryLookup struct {
	KeySpecifier
}

func (ul *UnaryLookup) String() string {
	var sb strings.Builder

	sb.WriteString("?")
	sb.WriteString(ul.KeySpecifier.String())

	return sb.String()
}

// KeySpecifier ::= NCName | IntegerLiteral | ParenthesizedExpr | "*"
// TypeID ::=				1			 | 2							| 3									| 4
type KeySpecifier struct {
	NCName
	IntegerLiteral
	ParenthesizedExpr
	TypeID byte
}

func (ks *KeySpecifier) String() string {
	switch ks.TypeID {
	case 1:
		return ks.NCName.Value()
	case 2:
		return ks.IntegerLiteral.String()
	case 3:
		return ks.ParenthesizedExpr.String()
	case 4:
		return "*"
	default:
		return ""
	}
}
