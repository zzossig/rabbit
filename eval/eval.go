package eval

import (
	"fmt"
	"math"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
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
	case *ast.StringLiteral:
		return &object.String{Value: expr.Value}
	case *ast.AdditiveExpr:
		return evalInfixExpr(expr, env)
	case *ast.MultiplicativeExpr:
		return evalInfixExpr(expr, env)
	case *ast.UnaryExpr:
		return evalPrefixExpr(expr, env)
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
		case *object.Array:
			xpath.Items = append(xpath.Items, item.Items...)
		default:
			xpath.Items = append(xpath.Items, item)
		}
	}

	return xpath
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(item object.Item) bool {
	if item != nil {
		return item.Type() == object.ErrorType
	}
	return false
}

func evalInfixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	var left object.Item
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.AdditiveExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.MultiplicativeExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	default:
		return newError("%T is not an infix expression\n", expr)
	}

	if isError(left) {
		return left
	}

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalInfixIntInt(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DecimalType:
		return evalInfixIntDecimal(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DoubleType:
		return evalInfixIntDouble(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.IntegerType:
		return evalInfixDecimalInt(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DecimalType:
		return evalInfixDecimalDecimal(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DoubleType:
		return evalInfixDecimalDouble(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.IntegerType:
		return evalInfixDoubleInt(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DecimalType:
		return evalInfixDoubleDecimal(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
		return evalInfixDoubleDouble(op, left, right)
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalPrefixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.UnaryExpr:
		right = Eval(expr.ExprSingle, env)
		op = expr.Token
	default:
		return newError("%T is not an prefix expression\n", expr)
	}

	if isError(right) {
		return right
	}

	switch {
	case right.Type() == object.IntegerType:
		return evalPrefixInt(op, right)
	case right.Type() == object.DecimalType:
		return evalPrefixDecimal(op, right)
	case right.Type() == object.DoubleType:
		return evalPrefixDouble(op, right)
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalInfixIntInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: float64(leftVal) / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / float64(rightVal))}
	case token.MOD:
		return &object.Integer{Value: leftVal % rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: float64(leftVal) + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: float64(leftVal) - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: float64(leftVal) * rightVal}
	case token.DIV:
		return &object.Decimal{Value: float64(leftVal) / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(float64(leftVal), rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: float64(leftVal) + rightVal}
	case token.MINUS:
		return &object.Double{Value: float64(leftVal) - rightVal}
	case token.ASTERISK:
		return &object.Double{Value: float64(leftVal) * rightVal}
	case token.DIV:
		return &object.Double{Value: float64(leftVal) / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / rightVal)}
	case token.MOD:
		return &object.Double{Value: math.Mod(float64(leftVal), rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + float64(rightVal)}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - float64(rightVal)}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * float64(rightVal)}
	case token.DIV:
		return &object.Decimal{Value: leftVal / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / float64(rightVal))}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, float64(rightVal))}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: leftVal + float64(rightVal)}
	case token.MINUS:
		return &object.Double{Value: leftVal - float64(rightVal)}
	case token.ASTERISK:
		return &object.Double{Value: leftVal * float64(rightVal)}
	case token.DIV:
		return &object.Double{Value: leftVal / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / float64(rightVal))}
	case token.MOD:
		return &object.Double{Value: math.Mod(leftVal, float64(rightVal))}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Double{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Double{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Double{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Double{Value: math.Mod(leftVal, rightVal)}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalPrefixInt(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Integer{Value: rightVal}
	case token.MINUS:
		return &object.Integer{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDecimal(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: rightVal}
	case token.MINUS:
		return &object.Decimal{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDouble(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: rightVal}
	case token.MINUS:
		return &object.Double{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}
