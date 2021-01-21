package ast

import (
	"regexp"

	"github.com/zzossig/xpath/token"
)

// ComparisonExpr ::= StringConcatExpr ( (ValueComp | GeneralComp | NodeComp) StringConcatExpr )?
type ComparisonExpr struct {
	Token  token.Token
	SCExpr []ExprSingle
}

// IsValueComp checks
// ValueComp ::= "eq" | "ne" | "lt" | "le" | "gt" | "ge"
func (ce *ComparisonExpr) IsValueComp() bool {
	re := regexp.MustCompile(`^(eq|ne|lt|le|gt|ge)$`)
	return re.MatchString(ce.Token.Literal)
}

// IsGeneralComp checks
// GeneralComp ::= "=" | "!=" | "<" | "<=" | ">" | ">="
func (ce *ComparisonExpr) IsGeneralComp() bool {
	re := regexp.MustCompile(`^(=|!=|<|<=|>|>=)$`)
	return re.MatchString(ce.Token.Literal)
}

// IsNodeComp checks
// NodeComp ::= "is" | "<<" | ">>"
func (ce *ComparisonExpr) IsNodeComp() bool {
	re := regexp.MustCompile(`^(is||<<||>>)$`)
	return re.MatchString(ce.Token.Literal)
}
