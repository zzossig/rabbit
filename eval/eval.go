package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
)

// predefined
var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval ..
func Eval(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.XPath:
		return evalXPath(expr, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: expr.Value}
	case *ast.DecimalLiteral:
		return &object.Decimal{Value: expr.Value}
	case *ast.DoubleLiteral:
		return &object.Double{Value: expr.Value}
	case *ast.StringLiteral:
		return &object.String{Value: expr.Value}
	case *ast.Expr:
		return evalExpr(expr, env)
	case *ast.ParenthesizedExpr:
		return evalExpr(expr, env)
	case *ast.EnclosedExpr:
		return evalExpr(expr, env)
	case *ast.Predicate:
		return evalExpr(expr, env)
	case *ast.Identifier:
		return evalIdentifier(expr, env)
	case *ast.InlineFunctionExpr:
		return evalFunctionLiteral(expr, env)
	case *ast.NamedFunctionRef:
		return evalFunctionLiteral(expr, env)
	case *ast.FunctionCall:
		return evalFunctionCall(expr, env)
	case *ast.VarRef:
		return evalVarRef(expr, env)
	case *ast.ArrowExpr:
		return evalArrowExpr(expr, env)
	case *ast.PostfixExpr:
		return evalPostfixExpr(expr, env)
	case *ast.AdditiveExpr:
		return evalInfixExpr(expr, env)
	case *ast.MultiplicativeExpr:
		return evalInfixExpr(expr, env)
	case *ast.StringConcatExpr:
		return evalInfixExpr(expr, env)
	case *ast.RangeExpr:
		return evalInfixExpr(expr, env)
	case *ast.UnaryExpr:
		return evalPrefixExpr(expr, env)
	case *ast.SquareArrayConstructor:
		return evalArrayExpr(expr, env)
	case *ast.CurlyArrayConstructor:
		return evalArrayExpr(expr, env)
	}

	return nil
}

func evalXPath(expr *ast.XPath, env *object.Env) object.Item {
	xpath := &object.Sequence{}

	for _, e := range expr.Exprs {
		item := Eval(e, env)

		switch item := item.(type) {
		case *object.Sequence:
			xpath.Items = append(xpath.Items, item.Items...)
		default:
			xpath.Items = append(xpath.Items, item)
		}
	}

	return xpath
}

func evalExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.Expr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.ParenthesizedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.EnclosedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.Predicate:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	}
	return nil
}
