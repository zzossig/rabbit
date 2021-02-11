package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalIntegerLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	il := expr.(*ast.IntegerLiteral)
	return &object.Integer{Value: il.Value}
}

func evalDecimalLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	dl := expr.(*ast.DecimalLiteral)
	return &object.Decimal{Value: dl.Value}
}

func evalDoubleLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	dl := expr.(*ast.DoubleLiteral)
	return &object.Double{Value: dl.Value}
}

func evalStringLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	sl := expr.(*ast.StringLiteral)
	return &object.String{Value: sl.Value}
}

func evalPrefixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.UnaryExpr:
		right = Eval(expr.ExprSingle, env)
		op = expr.Token
	default:
		return bif.NewError("%T is not an prefix expression\n", expr)
	}

	if bif.IsError(right) {
		return right
	}

	switch {
	case right.Type() == object.IntegerType:
		return evalPrefixInt(op, right, env)
	case right.Type() == object.DecimalType:
		return evalPrefixDecimal(op, right, env)
	case right.Type() == object.DoubleType:
		return evalPrefixDouble(op, right, env)
	case right.Type() == object.NilType:
		return &object.Nil{}
	default:
		return bif.NewError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixInt(op token.Token, right object.Item, env *object.Env) object.Item {
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Integer{Value: rightVal}
	case token.MINUS:
		return &object.Integer{Value: -1 * rightVal}
	default:
		return bif.NewError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDecimal(op token.Token, right object.Item, env *object.Env) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: rightVal}
	case token.MINUS:
		return &object.Decimal{Value: -1 * rightVal}
	default:
		return bif.NewError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDouble(op token.Token, right object.Item, env *object.Env) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: rightVal}
	case token.MINUS:
		return &object.Decimal{Value: -1 * rightVal}
	default:
		return bif.NewError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalArrayExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	array := &object.Array{}
	var exprs []ast.ExprSingle

	switch expr := expr.(type) {
	case *ast.SquareArrayConstructor:
		exprs = expr.Exprs
	case *ast.CurlyArrayConstructor:
		exprs = expr.EnclosedExpr.Exprs
	}

	for _, e := range exprs {
		item := Eval(e, env)
		array.Items = append(array.Items, item)
	}

	return array
}

func evalIfExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	ie := expr.(*ast.IfExpr)
	builtin := bif.Builtins["boolean"]

	testE := Eval(ie.TestExpr, env)
	bl, ok := builtin(testE).(*object.Boolean)
	if !ok {
		return builtin(testE)
	}

	if bl.Value {
		return Eval(ie.ThenExpr, env)
	}
	return Eval(ie.ElseExpr, env)
}

func evalForExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	fe := expr.(*ast.ForExpr)
	var items []object.Item

	if len(fe.Bindings) > 1 {
		b := fe.Bindings[0]
		bval := Eval(b.ExprSingle, env)

		nfe := &ast.ForExpr{ExprSingle: fe.ExprSingle}
		nfe.Bindings = fe.Bindings[1:]

		switch bval := bval.(type) {
		case *object.Sequence:
			for _, item := range bval.Items {
				env.Set(b.VarName.Value(), item)
				e := evalForExpr(nfe, env).(*object.Sequence)
				items = append(items, e.Items...)
			}
		default:
			env.Set(b.VarName.Value(), bval)
			e := evalForExpr(nfe, env).(*object.Sequence)
			items = append(items, e.Items...)
		}

		return &object.Sequence{Items: items}
	}

	b := fe.Bindings[0]
	bval := Eval(b.ExprSingle, env)

	switch bval := bval.(type) {
	case *object.Sequence:
		for _, item := range bval.Items {
			env.Set(b.VarName.Value(), item)
			e := Eval(fe.ExprSingle, env)
			items = append(items, e)
		}
	default:
		env.Set(b.VarName.Value(), bval)
		e := Eval(fe.ExprSingle, env)
		items = append(items, e)
	}

	return &object.Sequence{Items: items}
}

func evalLetExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	le := expr.(*ast.LetExpr)

	for _, b := range le.Bindings {
		bval := Eval(b.ExprSingle, env)
		env.Set(b.VarName.Value(), bval)
	}

	return Eval(le.ExprSingle, env)
}
