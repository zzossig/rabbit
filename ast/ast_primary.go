package ast

import (
	"github.com/zzossig/xpath/token"
)

// PrimaryExpr ::= Literal | VarRef | ParenthesizedExpr | ContextItemExpr | FunctionCall | FunctionItemExpr | MapConstructor | ArrayConstructor | UnaryLookup
type PrimaryExpr interface {
	ExprSingle
	primaryExpr()
}

// Literal ::= NumericLiteral | StringLiteral
type Literal interface {
	PrimaryExpr
	literal()
}

// NumericLiteral ::= IntegerLiteral | DecimalLiteral | DoubleLiteral
type NumericLiteral interface {
	Literal
	numericLiteral()
}

// Argument ::= ExprSingle | ArgumentPlaceholder
type Argument interface {
	argument()
}

// FunctionItemExpr ::= NamedFunctionRef | InlineFunctionExpr
type FunctionItemExpr interface {
	PrimaryExpr
	functionItemExpr()
}

// IntegerLiteral ::= Digits
// Digits ::= [0-9]+
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) exprSingle()     {}
func (il *IntegerLiteral) argument()       {}
func (il *IntegerLiteral) primaryExpr()    {}
func (il *IntegerLiteral) literal()        {}
func (il *IntegerLiteral) numericLiteral() {}

// DecimalLiteral ::= ("." Digits) | (Digits "." [0-9]*)
type DecimalLiteral struct {
	Value float64
}

func (dl *DecimalLiteral) exprSingle()     {}
func (dl *DecimalLiteral) argument()       {}
func (dl *DecimalLiteral) primaryExpr()    {}
func (dl *DecimalLiteral) literal()        {}
func (dl *DecimalLiteral) numericLiteral() {}

// DoubleLiteral ::= (("." Digits) | (Digits ("." [0-9]*)?)) [eE] [+-]? Digits
type DoubleLiteral struct {
	Value float64
}

func (dl *DoubleLiteral) exprSingle()     {}
func (dl *DoubleLiteral) argument()       {}
func (dl *DoubleLiteral) primaryExpr()    {}
func (dl *DoubleLiteral) literal()        {}
func (dl *DoubleLiteral) numericLiteral() {}

// StringLiteral ::= ('"' (EscapeQuot | [^"])* '"') | ("'" (EscapeApos | [^'])*
// EscapeQuot ::= '""'
// EscapeApos ::= "''"
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) exprSingle()  {}
func (sl *StringLiteral) argument()    {}
func (sl *StringLiteral) primaryExpr() {}
func (sl *StringLiteral) literal()     {}

// VarRef ::= "$" VarName
type VarRef struct {
	VarName VarName
}

func (vr *VarRef) exprSingle()  {}
func (vr *VarRef) argument()    {}
func (vr *VarRef) primaryExpr() {}

// VarName ::= EQName
type VarName EQName

// ParenthesizedExpr ::= "(" Expr? ")"
type ParenthesizedExpr struct {
	Exprs []ExprSingle
}

func (pe *ParenthesizedExpr) exprSingle()  {}
func (pe *ParenthesizedExpr) argument()    {}
func (pe *ParenthesizedExpr) primaryExpr() {}

// EnclosedExpr ::= "{" Expr? "}"
type EnclosedExpr struct {
	Exprs []ExprSingle
}

func (ee *EnclosedExpr) exprSingle()  {}
func (ee *EnclosedExpr) argument()    {}
func (ee *EnclosedExpr) primaryExpr() {}

// ContextItemExpr ::= "."
type ContextItemExpr struct {
	Token token.Token // token.DOT
}

func (cie *ContextItemExpr) exprSingle()  {}
func (cie *ContextItemExpr) argument()    {}
func (cie *ContextItemExpr) primaryExpr() {}

// FunctionCall ::= EQName ArgumentList
type FunctionCall struct {
	Name EQName
	ArgumentList
}

func (fc *FunctionCall) exprSingle()  {}
func (fc *FunctionCall) argument()    {}
func (fc *FunctionCall) primaryExpr() {}

// ArgumentList ::= "(" (Argument ("," Argument)*)? ")"
type ArgumentList struct {
	Args []Argument
}

// ArgumentPlaceholder ::= "?"
type ArgumentPlaceholder struct {
	Token token.Token // token.QUESTION
}

func (ap *ArgumentPlaceholder) argument() {}

// NamedFunctionRef ::= EQName "#" IntegerLiteral
type NamedFunctionRef struct {
	Name EQName
	IntegerLiteral
}

func (nfr *NamedFunctionRef) exprSingle()       {}
func (nfr *NamedFunctionRef) argument()         {}
func (nfr *NamedFunctionRef) functionItemExpr() {}

// InlineFunctionExpr ::= "function" "(" ParamList? ")" ("as" SequenceType)? FunctionBody
type InlineFunctionExpr struct {
	ParamList
	SequenceType
	FunctionBody
}

func (ifr *InlineFunctionExpr) exprSingle()       {}
func (ifr *InlineFunctionExpr) argument()         {}
func (ifr *InlineFunctionExpr) functionItemExpr() {}

// FunctionBody ::= EnclosedExpr
type FunctionBody = EnclosedExpr

// TypeDeclaration ::= "as" SequenceType
type TypeDeclaration struct {
	SequenceType
}

// Param ::= "$" EQName TypeDeclaration?
type Param struct {
	Name EQName
	TypeDeclaration
}

// ParamList ::= Param ("," Param)*
type ParamList struct {
	Params []Param
}
