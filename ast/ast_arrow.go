package ast

import "strings"

// ArrowExpr ::= UnaryExpr ( "=>" ArrowFunctionSpecifier ArgumentList )*
type ArrowExpr struct {
	ExprSingle
	Bindings []ArrowBinding
}

func (ae *ArrowExpr) exprSingle() {}
func (ae *ArrowExpr) String() string {
	var sb strings.Builder

	sb.WriteString(ae.ExprSingle.String())
	for _, b := range ae.Bindings {
		sb.WriteString(" ")
		sb.WriteString(b.String())
	}

	return sb.String()
}

// ArrowBinding ::= "=>" ArrowFunctionSpecifier ArgumentList **custom**
type ArrowBinding struct {
	ArrowFunctionSpecifier
	ArgumentList
}

func (ab *ArrowBinding) String() string {
	var sb strings.Builder

	sb.WriteString("=>")
	sb.WriteString(" ")
	sb.WriteString(ab.ArrowFunctionSpecifier.String())
	sb.WriteString(ab.ArgumentList.String())

	return sb.String()
}

// ArrowFunctionSpecifier ::= EQName | VarRef | ParenthesizedExpr
// TypeID ::= 								1			 | 2			| 3
type ArrowFunctionSpecifier struct {
	EQName
	VarRef
	ParenthesizedExpr
	TypeID byte
}

func (afs *ArrowFunctionSpecifier) String() string {
	var sb strings.Builder

	switch afs.TypeID {
	case 1:
		sb.WriteString(afs.EQName.Value())
	case 2:
		sb.WriteString(afs.VarRef.String())
	case 3:
		sb.WriteString(afs.ParenthesizedExpr.String())
	default:
		sb.WriteString("")
	}

	return sb.String()
}
