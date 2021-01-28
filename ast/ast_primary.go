package ast

import (
	"fmt"
	"strings"

	"github.com/zzossig/xpath/token"
)

// PrimaryExpr ::= Literal | VarRef | ParenthesizedExpr | ContextItemExpr | FunctionCall | FunctionItemExpr | MapConstructor | ArrayConstructor | UnaryLookup
type PrimaryExpr interface {
	ExprSingle
	primaryExpr()
}

// Literal ::= NumericLiteral | StringLiteral
// TypeID ::=  1							| 2
type Literal struct {
	PrimaryExpr
	TypeID byte
}

func (l *Literal) epxrSingle()  {}
func (l *Literal) primaryExpr() {}
func (l *Literal) String() string {
	return l.PrimaryExpr.String()
}

// NumericLiteral ::= IntegerLiteral | DecimalLiteral | DoubleLiteral
// TypeID ::=					1							 | 2							| 3
type NumericLiteral struct {
	PrimaryExpr
	TypeID byte
}

func (nl *NumericLiteral) epxrSingle()  {}
func (nl *NumericLiteral) primaryExpr() {}
func (nl *NumericLiteral) String() string {
	return nl.PrimaryExpr.String()
}

// FunctionItemExpr ::= NamedFunctionRef | InlineFunctionExpr
// TypeID ::= 					1								 | 2
type FunctionItemExpr struct {
	PrimaryExpr
	TypeID byte
}

func (fie *FunctionItemExpr) epxrSingle()  {}
func (fie *FunctionItemExpr) primaryExpr() {}
func (fie *FunctionItemExpr) String() string {
	return fie.PrimaryExpr.String()
}

// IntegerLiteral ::= Digits
// Digits ::= [0-9]+
type IntegerLiteral struct {
	Value int
}

func (il *IntegerLiteral) exprSingle()  {}
func (il *IntegerLiteral) primaryExpr() {}
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

// DecimalLiteral ::= ("." Digits) | (Digits "." [0-9]*)
type DecimalLiteral struct {
	Value float64
}

func (dl *DecimalLiteral) exprSingle()  {}
func (dl *DecimalLiteral) primaryExpr() {}
func (dl *DecimalLiteral) String() string {
	return fmt.Sprintf("%f", dl.Value)
}

// DoubleLiteral ::= (("." Digits) | (Digits ("." [0-9]*)?)) [eE] [+-]? Digits
type DoubleLiteral struct {
	Value float64
}

func (dl *DoubleLiteral) exprSingle()  {}
func (dl *DoubleLiteral) primaryExpr() {}
func (dl *DoubleLiteral) String() string {
	return fmt.Sprintf("%e", dl.Value)
}

// StringLiteral ::= ('"' (EscapeQuot | [^"])* '"') | ("'" (EscapeApos | [^'])*
// EscapeQuot ::= '""'
// EscapeApos ::= "''"
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) exprSingle()  {}
func (sl *StringLiteral) primaryExpr() {}
func (sl *StringLiteral) String() string {
	var sb strings.Builder

	sb.WriteString("'")
	sb.WriteString(sl.Value)
	sb.WriteString("'")

	return sb.String()
}

// VarRef ::= "$" VarName
type VarRef struct {
	VarName
}

func (vr *VarRef) exprSingle()  {}
func (vr *VarRef) primaryExpr() {}
func (vr *VarRef) String() string {
	var sb strings.Builder
	sb.WriteString("$")
	sb.WriteString(vr.VarName.Value())
	return sb.String()
}

// VarName ::= EQName
type VarName = EQName

// ParenthesizedExpr ::= "(" Expr? ")"
type ParenthesizedExpr struct {
	Expr
}

func (pe *ParenthesizedExpr) exprSingle()  {}
func (pe *ParenthesizedExpr) primaryExpr() {}
func (pe *ParenthesizedExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.Expr.String())
	sb.WriteString(")")

	return sb.String()
}

// EnclosedExpr ::= "{" Expr? "}"
type EnclosedExpr struct {
	Expr
}

func (ee *EnclosedExpr) exprSingle()  {}
func (ee *EnclosedExpr) primaryExpr() {}
func (ee *EnclosedExpr) String() string {
	var sb strings.Builder

	sb.WriteString("{")
	sb.WriteString(ee.Expr.String())
	sb.WriteString("}")

	return sb.String()
}

// ContextItemExpr ::= "."
type ContextItemExpr struct {
	Token token.Token // token.DOT
}

func (cie *ContextItemExpr) exprSingle()  {}
func (cie *ContextItemExpr) primaryExpr() {}
func (cie *ContextItemExpr) String() string {
	return cie.Token.Literal
}

// FunctionCall ::= EQName ArgumentList
type FunctionCall struct {
	EQName
	ArgumentList
}

func (fc *FunctionCall) exprSingle()  {}
func (fc *FunctionCall) primaryExpr() {}
func (fc *FunctionCall) String() string {
	var sb strings.Builder

	sb.WriteString(fc.EQName.Value())
	sb.WriteString(fc.ArgumentList.String())

	return sb.String()
}

// Argument ::= ExprSingle | ArgumentPlaceholder
type Argument struct {
	ExprSingle
	ArgumentPlaceholder
}

func (a *Argument) String() string {
	if a.ArgumentPlaceholder.String() != "" {
		return a.ArgumentPlaceholder.String()
	}
	return a.ExprSingle.String()
}

// ArgumentList ::= "(" (Argument ("," Argument)*)? ")"
type ArgumentList struct {
	Args []Argument
}

func (al *ArgumentList) pal() {}
func (al *ArgumentList) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	for i, arg := range al.Args {
		sb.WriteString(arg.String())
		if i < len(al.Args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return sb.String()
}

// ArgumentPlaceholder ::= "?"
type ArgumentPlaceholder struct {
	Token token.Token // token.QUESTION
}

func (ap *ArgumentPlaceholder) String() string {
	return ap.Token.Literal
}

// NamedFunctionRef ::= EQName "#" IntegerLiteral
type NamedFunctionRef struct {
	EQName
	IntegerLiteral
}

func (nfr *NamedFunctionRef) exprSingle() {}
func (nfr *NamedFunctionRef) String() string {
	var sb strings.Builder

	sb.WriteString(nfr.EQName.Value())
	sb.WriteString("#")
	sb.WriteString(nfr.IntegerLiteral.String())

	return sb.String()
}

// InlineFunctionExpr ::= "function" "(" ParamList? ")" ("as" SequenceType)? FunctionBody
type InlineFunctionExpr struct {
	ParamList
	SequenceType
	FunctionBody
}

func (ifr *InlineFunctionExpr) exprSingle() {}
func (ifr *InlineFunctionExpr) String() string {
	var sb strings.Builder

	sb.WriteString("function(")
	sb.WriteString(ifr.ParamList.String())
	sb.WriteString(")")
	if ifr.SequenceType.String() != "" {
		sb.WriteString(" ")
		sb.WriteString("as")
		sb.WriteString(" ")
		sb.WriteString(ifr.SequenceType.String())
	}
	sb.WriteString(" ")
	sb.WriteString(ifr.FunctionBody.String())

	return sb.String()
}

// FunctionBody ::= EnclosedExpr
type FunctionBody = EnclosedExpr

// TypeDeclaration ::= "as" SequenceType
type TypeDeclaration struct {
	SequenceType
}

func (td *TypeDeclaration) String() string {
	var sb strings.Builder

	sb.WriteString("as")
	sb.WriteString(" ")
	sb.WriteString(td.SequenceType.String())

	return sb.String()
}

// Param ::= "$" EQName TypeDeclaration?
type Param struct {
	EQName
	TypeDeclaration
}

func (p *Param) String() string {
	var sb strings.Builder

	if p.EQName.Value() != "" {
		sb.WriteString("$")
		sb.WriteString(p.EQName.Value())
		sb.WriteString(" ")
		sb.WriteString(p.TypeDeclaration.String())
	}

	return sb.String()
}

// ParamList ::= Param ("," Param)*
type ParamList struct {
	Params []Param
}

func (pl *ParamList) String() string {
	var sb strings.Builder

	for i, p := range pl.Params {
		sb.WriteString(p.String())
		if i < len(pl.Params)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
