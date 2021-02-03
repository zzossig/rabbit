package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

// NodeTest ::= KindTest | NameTest
type NodeTest interface {
	nodeTest()
	String() string
}

// StepExpr ::= PostfixExpr | AxisStep
// TypeID ::= 	1						| 2
type StepExpr struct {
	ExprSingle
	TypeID byte
}

func (se *StepExpr) exprSingle() {}
func (se *StepExpr) String() string {
	return se.ExprSingle.String()
}

// PathExpr ::= ("/" RelativePathExpr?) | ("//" RelativePathExpr) | RelativePathExpr
type PathExpr struct {
	ExprSingle
	Token token.Token
}

func (pe *PathExpr) exprSingle() {}
func (pe *PathExpr) String() string {
	var sb strings.Builder

	sb.WriteString(pe.Token.Literal)
	sb.WriteString(pe.ExprSingle.String())

	return sb.String()
}

// RelativePathExpr ::= StepExpr (("/" | "//") StepExpr)*
type RelativePathExpr struct {
	LeftExpr  ExprSingle
	RightExpr ExprSingle
	Token     token.Token
}

func (rpe *RelativePathExpr) exprSingle() {}
func (rpe *RelativePathExpr) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(rpe.LeftExpr.String())
	sb.WriteString(" ")
	sb.WriteString(rpe.Token.Literal)
	sb.WriteString(" ")
	sb.WriteString(rpe.RightExpr.String())
	sb.WriteString(")")

	return sb.String()
}

// AxisStep ::= (ReverseStep | ForwardStep) PredicateList
// TypeID ::= 	1						 | 2
type AxisStep struct {
	ForwardStep
	ReverseStep
	PredicateList
	TypeID byte
}

func (as *AxisStep) exprSingle() {}
func (as *AxisStep) String() string {
	var sb strings.Builder

	switch as.TypeID {
	case 1:
		sb.WriteString(as.ReverseStep.String())
	case 2:
		sb.WriteString(as.ForwardStep.String())
	default:
		sb.WriteString("")
	}
	sb.WriteString(as.PredicateList.String())

	return sb.String()
}

// ForwardStep ::= (ForwardAxis NodeTest) | AbbrevForwardStep
// TypeID ::=			 1											| 2
type ForwardStep struct {
	ForwardAxis
	NodeTest
	AbbrevForwardStep
	TypeID byte
}

func (fs *ForwardStep) exprSingle() {}
func (fs *ForwardStep) String() string {
	var sb strings.Builder

	switch fs.TypeID {
	case 1:
		sb.WriteString(fs.ForwardAxis.Value())
		sb.WriteString(fs.NodeTest.String())
	case 2:
		sb.WriteString(fs.AbbrevForwardStep.String())
	default:
		sb.WriteString("")
	}

	return sb.String()
}

// ReverseStep ::= (ReverseAxis NodeTest) | AbbrevReverseStep
// TypeID ::=			 1											| 2
type ReverseStep struct {
	ReverseAxis
	NodeTest
	AbbrevReverseStep
	TypeID byte
}

func (rs *ReverseStep) exprSingle() {}
func (rs *ReverseStep) String() string {
	var sb strings.Builder

	switch rs.TypeID {
	case 1:
		sb.WriteString(rs.ReverseAxis.Value())
		sb.WriteString(rs.NodeTest.String())
	case 2:
		sb.WriteString(rs.AbbrevReverseStep.String())
	default:
		sb.WriteString("")
	}

	return sb.String()
}

// ForwardAxis ::= ("child" "::") | ("descendant" "::") | ("attribute" "::") | ("self" "::") | ("descendant-or-self" "::") | ("following-sibling" "::") | ("following" "::") | ("namespace" "::")
type ForwardAxis struct {
	value string
}

// Value is a getter for the value field
func (fa *ForwardAxis) Value() string {
	return fa.value
}

// SetValue is a setter for the value field
func (fa *ForwardAxis) SetValue(str string) {
	if util.IsForwardAxis(str) {
		fa.value = str
	} else {
		// TODO error
	}
}

// AbbrevForwardStep ::= "@"? NodeTest
type AbbrevForwardStep struct {
	Token token.Token
	NodeTest
}

func (afs *AbbrevForwardStep) exprSingle() {}
func (afs *AbbrevForwardStep) String() string {
	var sb strings.Builder

	sb.WriteString(afs.Token.Literal)
	sb.WriteString(afs.NodeTest.String())

	return sb.String()
}

// ReverseAxis ::= ("parent" "::") | ("ancestor" "::") | ("preceding-sibling" "::") | ("preceding" "::") | ("ancestor-or-self" "::")
type ReverseAxis struct {
	value string
}

// Value is a getter for the value field
func (ra *ReverseAxis) Value() string {
	return ra.value
}

// SetValue is a setter for the value field
func (ra *ReverseAxis) SetValue(str string) {
	if util.IsReverseAxis(str) {
		ra.value = str
	} else {
		// TODO error
	}
}

// AbbrevReverseStep ::= ".."
type AbbrevReverseStep struct {
	Token token.Token
}

func (ars *AbbrevReverseStep) exprSingle() {}
func (ars *AbbrevReverseStep) String() string {
	return ars.Token.Literal
}

// NameTest ::= EQName | Wildcard
// TypeID ::= 	1			 | 2
type NameTest struct {
	EQName
	Wildcard
	TypeID byte
}

func (nt *NameTest) nodeTest() {}
func (nt *NameTest) String() string {
	switch nt.TypeID {
	case 1:
		return nt.EQName.Value()
	case 2:
		return nt.Wildcard.String()
	default:
		return ""
	}
}

// Wildcard ::= "*" | (NCName ":*") | ("*:" NCName) | (BracedURILiteral "*")
// TypeID				1		| 2							| 3							| 4
type Wildcard struct {
	NCName
	BracedURILiteral
	TypeID byte
}

func (w *Wildcard) exprSingle() {}
func (w *Wildcard) nodeTest()   {}
func (w *Wildcard) String() string {
	var sb strings.Builder

	switch w.TypeID {
	case 1:
		sb.WriteString("*")
	case 2:
		sb.WriteString(w.NCName.Value())
		sb.WriteString(":")
		sb.WriteString("*")
	case 3:
		sb.WriteString("*")
		sb.WriteString(":")
		sb.WriteString(w.NCName.Value())
	case 4:
		sb.WriteString(w.BracedURILiteral.Value())
		sb.WriteString("*")
	default:
		sb.WriteString("")
	}

	return sb.String()
}
