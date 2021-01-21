package ast

// OrExpr ::= AndExpr ( "or" AndExpr )*
type OrExpr struct {
	Exprs []ExprSingle
}

func (oe *OrExpr) exprSingle() {}
func (oe *OrExpr) argument()   {}

// AndExpr ::= ComparisonExpr ( "and" ComparisonExpr )*
type AndExpr struct {
	Exprs []ExprSingle
}

func (ae *AndExpr) exprSingle() {}
func (ae *AndExpr) argument()   {}
