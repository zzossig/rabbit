package eval

import (
	"math"
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalIdentifier(ident *ast.Identifier, ctx *object.Context) object.Item {
	if val, ok := ctx.Get(ident.EQName.Value()); ok {
		return val
	}

	return bif.NewError("identifier not found: " + ident.EQName.Value())
}

func evalInfixExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	var left object.Item
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.AdditiveExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.MultiplicativeExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.StringConcatExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.RangeExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.ComparisonExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.SimpleMapExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.UnionExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.IntersectExceptExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
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
	case left.Type() == object.IntegerType && right.Type() == object.ArrayType:
		return evalInfixNumberArray(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.SequenceType:
		return evalInfixNumberSeq(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.IntegerType:
		return evalInfixDecimalInt(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DecimalType:
		return evalInfixDecimalDecimal(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DoubleType:
		return evalInfixDecimalDouble(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.StringType:
		return evalInfixDecimalString(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.ArrayType:
		return evalInfixNumberArray(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.SequenceType:
		return evalInfixNumberSeq(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.IntegerType:
		return evalInfixDoubleInt(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DecimalType:
		return evalInfixDoubleDecimal(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
		return evalInfixDoubleDouble(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.StringType:
		return evalInfixDoubleString(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.ArrayType:
		return evalInfixNumberArray(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.SequenceType:
		return evalInfixNumberSeq(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.IntegerType:
		return evalInfixStringInt(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DecimalType:
		return evalInfixStringDecimal(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DoubleType:
		return evalInfixStringDouble(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return evalInfixStringString(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.ArrayType:
		return evalInfixStringArray(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.SequenceType:
		return evalInfixStringSeq(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.IntegerType:
		return evalInfixArrayNumber(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.DecimalType:
		return evalInfixArrayNumber(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.DoubleType:
		return evalInfixArrayNumber(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.StringType:
		return evalInfixArrayString(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.ArrayType:
		return evalInfixArrayArray(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.SequenceType:
		return evalInfixArraySeq(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.IntegerType:
		return evalInfixSeqNumber(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.DecimalType:
		return evalInfixSeqNumber(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.DoubleType:
		return evalInfixSeqNumber(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.StringType:
		return evalInfixSeqString(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.ArrayType:
		return evalInfixSeqArray(op, left, right)
	case left.Type() == object.SequenceType && right.Type() == object.SequenceType:
		return evalInfixSeqSeq(op, left, right)
	case left.Type() == object.BooleanType || right.Type() == object.BooleanType:
		return evalInfixBool(op, left, right)
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
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixNumberArray(op token.Token, left object.Item, right object.Item) object.Item {
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixNumberSeq(op token.Token, left object.Item, right object.Item) object.Item {
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
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
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringArray(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value == e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value != e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value < e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value <= e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value > e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value >= e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringSeq(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value == e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value != e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value < e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value <= e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value > e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value >= e.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixSeqString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value == rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value != rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value < rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value <= rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value > rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value >= rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixArrayNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixArrayString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value == rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value != rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value < rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value <= rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value > rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value >= rightVal.Value {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixArrayArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixSeqArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixArraySeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixSeqSeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixSeqNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalLogicalExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	var left object.Item
	var right object.Item
	var op token.Token

	builtin := bif.Builtins["boolean"]

	switch expr := expr.(type) {
	case *ast.AndExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.OrExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	}

	l := builtin(left)
	if bif.IsError(l) {
		return l
	}
	r := builtin(right)
	if bif.IsError(r) {
		return r
	}

	leftBool := l.(*object.Boolean)
	rightBool := r.(*object.Boolean)

	switch op.Type {
	case token.AND:
		return &object.Boolean{Value: leftBool.Value && rightBool.Value}
	case token.OR:
		return &object.Boolean{Value: leftBool.Value || rightBool.Value}
	default:
		return object.NIL
	}
}

func evalInfixBool(op token.Token, left, right object.Item) object.Item {
	leftVal, ok := left.(*object.Boolean)
	if !ok {
		return bif.NewError("[XPTY0004] Types %s and %s are not comparable.", left.Type(), right.Type())
	}

	rightVal, ok := right.(*object.Boolean)
	if !ok {
		return bif.NewError("[XPTY0004] Types %s and %s are not comparable.", left.Type(), right.Type())
	}

	switch op.Type {
	case token.EQ:
		fallthrough
	case token.EQV:
		return &object.Boolean{Value: leftVal.Value == rightVal.Value}
	case token.NE:
		fallthrough
	case token.NEV:
		return &object.Boolean{Value: leftVal.Value != rightVal.Value}
	case token.GT:
		fallthrough
	case token.GTV:
		if leftVal.Value && !rightVal.Value {
			return object.TRUE
		}
		return object.FALSE
	case token.GE:
		fallthrough
	case token.GEV:
		if !leftVal.Value && rightVal.Value {
			return object.FALSE
		}
		return object.TRUE
	case token.LT:
		fallthrough
	case token.LTV:
		if !leftVal.Value && rightVal.Value {
			return object.TRUE
		}
		return object.FALSE
	case token.LE:
		fallthrough
	case token.LEV:
		if leftVal.Value && !rightVal.Value {
			return object.FALSE
		}
		return object.TRUE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}
