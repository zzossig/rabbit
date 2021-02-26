package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
)

// Eval ..
func Eval(expr ast.ExprSingle, ctx *object.Context) object.Item {
	switch expr := expr.(type) {
	case *ast.XPath:
		return evalXPath(expr, ctx)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(expr, ctx)
	case *ast.DecimalLiteral:
		return evalDecimalLiteral(expr, ctx)
	case *ast.DoubleLiteral:
		return evalDoubleLiteral(expr, ctx)
	case *ast.StringLiteral:
		return evalStringLiteral(expr, ctx)
	case *ast.ContextItemExpr:
		return ctx.CItem
	case *ast.Expr:
		return evalExpr(expr, ctx)
	case *ast.ParenthesizedExpr:
		return evalExpr(expr, ctx)
	case *ast.EnclosedExpr:
		return evalExpr(expr, ctx)
	case *ast.Predicate:
		return evalExpr(expr, ctx)
	case *ast.Identifier:
		return evalIdentifier(expr, ctx)
	case *ast.InlineFunctionExpr:
		return evalFunctionLiteral(expr, ctx)
	case *ast.NamedFunctionRef:
		return evalFunctionLiteral(expr, ctx)
	case *ast.FunctionCall:
		return evalFunctionCall(expr, ctx)
	case *ast.VarRef:
		return evalVarRef(expr, ctx)
	case *ast.ArrowExpr:
		return evalArrowExpr(expr, ctx)
	case *ast.PostfixExpr:
		return evalPostfixExpr(expr, ctx)
	case *ast.AdditiveExpr:
		return evalAdditiveExpr(expr, ctx)
	case *ast.MultiplicativeExpr:
		return evalMultiplicativeExpr(expr, ctx)
	case *ast.StringConcatExpr:
		return evalStringConcatExpr(expr, ctx)
	case *ast.RangeExpr:
		return evalRangeExpr(expr, ctx)
	case *ast.ComparisonExpr:
		return evalComparisonExpr(expr, ctx)
	case *ast.UnionExpr:
		return evalUnionExpr(expr, ctx)
	case *ast.IntersectExceptExpr:
		return evalIntersectExceptExpr(expr, ctx)
	case *ast.OrExpr:
		return evalLogicalExpr(expr, ctx)
	case *ast.AndExpr:
		return evalLogicalExpr(expr, ctx)
	case *ast.SimpleMapExpr:
		return evalSimpleMapExpr(expr, ctx)
	case *ast.UnaryExpr:
		return evalUnaryExpr(expr, ctx)
	case *ast.SquareArrayConstructor:
		return evalArrayExpr(expr, ctx)
	case *ast.CurlyArrayConstructor:
		return evalArrayExpr(expr, ctx)
	case *ast.IfExpr:
		return evalIfExpr(expr, ctx)
	case *ast.ForExpr:
		return evalForExpr(expr, ctx)
	case *ast.LetExpr:
		return evalLetExpr(expr, ctx)
	case *ast.QuantifiedExpr:
		return evalQuantifiedExpr(expr, ctx)
	case *ast.MapConstructor:
		return evalMapExpr(expr, ctx)
	case *ast.UnaryLookup:
		return evalUnaryLookup(expr, ctx)
	}

	return object.NIL
}
