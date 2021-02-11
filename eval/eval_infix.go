package eval

import (
	"math"
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalIdentifier(ident *ast.Identifier, env *object.Env) object.Item {
	if val, ok := env.Get(ident.EQName.Value()); ok {
		return val
	}

	return bif.NewError("identifier not found: " + ident.EQName.Value())
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
	case *ast.StringConcatExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.RangeExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.ComparisonExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.SimpleMapExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	default:
		return bif.NewError("%T is not an infix expression\n", expr)
	}

	if bif.IsError(left) {
		return left
	}

	if bif.IsError(right) {
		return right
	}

	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalInfixIntInt(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DecimalType:
		return evalInfixIntDecimal(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DoubleType:
		return evalInfixIntDouble(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.StringType:
		return evalInfixIntString(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.IntegerType:
		return evalInfixDecimalInt(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DecimalType:
		return evalInfixDecimalDecimal(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DoubleType:
		return evalInfixDecimalDouble(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.StringType:
		return evalInfixDecimalString(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.IntegerType:
		return evalInfixDoubleInt(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DecimalType:
		return evalInfixDoubleDecimal(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
		return evalInfixDoubleDouble(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.StringType:
		return evalInfixDoubleString(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.IntegerType:
		return evalInfixStringInt(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DecimalType:
		return evalInfixStringDecimal(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DoubleType:
		return evalInfixStringDouble(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return evalInfixStringString(op, left, right)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.TO:
		seq := &object.Sequence{}
		for i := leftVal; i <= rightVal; i++ {
			seq.Items = append(seq.Items, &object.Integer{Value: i})
		}
		return seq
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: float64(leftVal) == rightVal}
	case token.NE:
		return &object.Boolean{Value: float64(leftVal) != rightVal}
	case token.LT:
		return &object.Boolean{Value: float64(leftVal) < rightVal}
	case token.LE:
		return &object.Boolean{Value: float64(leftVal) <= rightVal}
	case token.GT:
		return &object.Boolean{Value: float64(leftVal) > rightVal}
	case token.GE:
		return &object.Boolean{Value: float64(leftVal) >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: float64(leftVal) == rightVal}
	case token.NEV:
		return &object.Boolean{Value: float64(leftVal) != rightVal}
	case token.LTV:
		return &object.Boolean{Value: float64(leftVal) < rightVal}
	case token.LEV:
		return &object.Boolean{Value: float64(leftVal) <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: float64(leftVal) > rightVal}
	case token.GEV:
		return &object.Boolean{Value: float64(leftVal) >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: float64(leftVal) == rightVal}
	case token.NE:
		return &object.Boolean{Value: float64(leftVal) != rightVal}
	case token.LT:
		return &object.Boolean{Value: float64(leftVal) < rightVal}
	case token.LE:
		return &object.Boolean{Value: float64(leftVal) <= rightVal}
	case token.GT:
		return &object.Boolean{Value: float64(leftVal) > rightVal}
	case token.GE:
		return &object.Boolean{Value: float64(leftVal) >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: float64(leftVal) == rightVal}
	case token.NEV:
		return &object.Boolean{Value: float64(leftVal) != rightVal}
	case token.LTV:
		return &object.Boolean{Value: float64(leftVal) < rightVal}
	case token.LEV:
		return &object.Boolean{Value: float64(leftVal) <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: float64(leftVal) > rightVal}
	case token.GEV:
		return &object.Boolean{Value: float64(leftVal) >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%d) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == float64(rightVal)}
	case token.NE:
		return &object.Boolean{Value: leftVal != float64(rightVal)}
	case token.LT:
		return &object.Boolean{Value: leftVal < float64(rightVal)}
	case token.LE:
		return &object.Boolean{Value: leftVal <= float64(rightVal)}
	case token.GT:
		return &object.Boolean{Value: leftVal > float64(rightVal)}
	case token.GE:
		return &object.Boolean{Value: leftVal >= float64(rightVal)}
	case token.EQV:
		return &object.Boolean{Value: leftVal == float64(rightVal)}
	case token.NEV:
		return &object.Boolean{Value: leftVal != float64(rightVal)}
	case token.LTV:
		return &object.Boolean{Value: leftVal < float64(rightVal)}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= float64(rightVal)}
	case token.GTV:
		return &object.Boolean{Value: leftVal > float64(rightVal)}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= float64(rightVal)}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == float64(rightVal)}
	case token.NE:
		return &object.Boolean{Value: leftVal != float64(rightVal)}
	case token.LT:
		return &object.Boolean{Value: leftVal < float64(rightVal)}
	case token.LE:
		return &object.Boolean{Value: leftVal <= float64(rightVal)}
	case token.GT:
		return &object.Boolean{Value: leftVal > float64(rightVal)}
	case token.GE:
		return &object.Boolean{Value: leftVal >= float64(rightVal)}
	case token.EQV:
		return &object.Boolean{Value: leftVal == float64(rightVal)}
	case token.NEV:
		return &object.Boolean{Value: leftVal != float64(rightVal)}
	case token.LTV:
		return &object.Boolean{Value: leftVal < float64(rightVal)}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= float64(rightVal)}
	case token.GTV:
		return &object.Boolean{Value: leftVal > float64(rightVal)}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= float64(rightVal)}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
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
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%f) with %s(%s)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%d)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LT:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GT:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GE:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.EQV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.NEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.LEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GTV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	case token.GEV:
		return bif.NewError("Cannot compare %s(%s) with %s(%f)", left.Type(), leftVal, right.Type(), rightVal)
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		return &object.String{Value: leftVal + rightVal}
	case token.EQ:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NE:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LT:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LE:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GT:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GE:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.EQV:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.NEV:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.LTV:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LEV:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.GTV:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GEV:
		return &object.Boolean{Value: leftVal >= rightVal}
	// case token.IS:
	// case token.DGT:
	// case token.DLT:
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}
