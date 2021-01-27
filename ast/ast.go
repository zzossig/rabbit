package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// XPath ::= Expr
type XPath struct {
	Items []Item
}

func (x *XPath) String() string {
	var sb strings.Builder

	if len(x.Items) > 1 {
		sb.WriteString("(")
		for i, s := range x.Items {
			sb.WriteString(s.String())
			if i != len(x.Items)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
	} else if len(x.Items) == 1 {
		sb.WriteString(x.Items[0].String())
	}

	return sb.String()
}

// Item **custom**
type Item interface {
	item()
	String() string
}

// ExprItem **custom**
type ExprItem struct {
	Token      token.Token // the first token of the expression
	Expression ExprSingle
}

func (ei *ExprItem) item() {}
func (ei *ExprItem) String() string {
	return ei.Expression.String()
}

// ExprSingle ::= ForExpr | LetExpr | QuantifiedExpr | IfExpr | OrExpr
type ExprSingle interface {
	exprSingle()
	String() string
}
