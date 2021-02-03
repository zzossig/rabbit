package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// PAL := Predicate | ArgumentList | Lookup **custom**
type PAL interface {
	pal()
	String() string
}

// PostfixExpr ::= PrimaryExpr (Predicate | ArgumentList | Lookup)*
type PostfixExpr struct {
	ExprSingle
	Pals []PAL
}

func (pe *PostfixExpr) exprSingle() {}
func (pe *PostfixExpr) String() string {
	var sb strings.Builder

	sb.WriteString(pe.ExprSingle.String())
	for _, p := range pe.Pals {
		sb.WriteString(p.String())
	}

	return sb.String()
}

// Predicate ::= "[" Expr "]"
type Predicate struct {
	Expr
}

func (p *Predicate) pal() {}
func (p *Predicate) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	sb.WriteString(p.Expr.String())
	sb.WriteString("]")

	return sb.String()
}

// PredicateList ::= Predicate*
type PredicateList struct {
	PL []Predicate
}

func (pl *PredicateList) String() string {
	var sb strings.Builder

	for _, p := range pl.PL {
		sb.WriteString(p.String())
	}

	return sb.String()
}

// Lookup ::= "?" KeySpecifier
type Lookup struct {
	Token token.Token // token.QUESTION
	KeySpecifier
}

func (l *Lookup) pal() {}
func (l *Lookup) String() string {
	var sb strings.Builder

	sb.WriteString("?")
	sb.WriteString(l.KeySpecifier.String())

	return sb.String()
}
