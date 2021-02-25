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
		return evalInfixExpr(expr, ctx)
	case *ast.UnionExpr:
		return evalInfixExpr(expr, ctx)
	case *ast.IntersectExceptExpr:
		return evalInfixExpr(expr, ctx)
	case *ast.ComparisonExpr:
		return evalInfixExpr(expr, ctx)
	case *ast.OrExpr:
		return evalLogicalExpr(expr, ctx)
	case *ast.AndExpr:
		return evalLogicalExpr(expr, ctx)
	case *ast.SimpleMapExpr:
		return evalSimpleMapExpr(expr, ctx)
	case *ast.UnaryExpr:
		return evalPrefixExpr(expr, ctx)
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

func evalXPath(expr *ast.XPath, ctx *object.Context) object.Item {
	xpath := &object.Sequence{}

	for _, e := range expr.Exprs {
		item := Eval(e, ctx)

		switch item := item.(type) {
		case *object.Sequence:
			xpath.Items = append(xpath.Items, item.Items...)
		default:
			xpath.Items = append(xpath.Items, item)
		}
	}

	return xpath
}

func evalExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	switch expr := expr.(type) {
	case *ast.Expr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.ParenthesizedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.EnclosedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.Predicate:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	}
	return object.NIL
}
