package bif

import (
	"math"

	"github.com/zzossig/xpath/object"
)

func numericAdd(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(leftVal + rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(float64(leftVal) + rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(float64(leftVal) + rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(leftVal + float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(leftVal + rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal + rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDouble(leftVal + float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDouble(leftVal + rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal + rightVal)
	}

	return object.NIL
}

func numericSubtract(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(leftVal - rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(float64(leftVal) - rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(float64(leftVal) - rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(leftVal - float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(leftVal - rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal - rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDouble(leftVal - float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDouble(leftVal - rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal - rightVal)
	}

	return object.NIL
}

func numericMultiply(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(leftVal * rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(float64(leftVal) * rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(float64(leftVal) * rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(leftVal * float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(leftVal * rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal * rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDouble(leftVal * float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDouble(leftVal * rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal * rightVal)
	}

	return object.NIL
}

func numericDivide(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(float64(leftVal) / float64(rightVal))
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(float64(leftVal) / rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(float64(leftVal) / rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(leftVal / float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(leftVal / rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal / rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDouble(leftVal / float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDouble(leftVal / rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(leftVal / rightVal)
	}

	return object.NIL
}

func numericIntegerDivide(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(int(float64(leftVal) / float64(rightVal)))
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewInteger(int(float64(leftVal) / rightVal))
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewInteger(int(float64(leftVal) / rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(int(leftVal / float64(rightVal)))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewInteger(int(leftVal / rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewInteger(int(leftVal / rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(int(leftVal / float64(rightVal)))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewInteger(int(leftVal / rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewInteger(int(leftVal / rightVal))
	}

	return object.NIL
}

func numericMod(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewInteger(leftVal % rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(math.Mod(float64(leftVal), rightVal))
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(math.Mod(float64(leftVal), rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDecimal(math.Mod(leftVal, float64(rightVal)))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(math.Mod(leftVal, rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDecimal(math.Mod(leftVal, rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewDouble(math.Mod(leftVal, float64(rightVal)))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewDecimal(math.Mod(leftVal, rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewDouble(math.Mod(leftVal, rightVal))
	}

	return object.NIL
}

func numericUnaryPlus(ctx *object.Context, args ...object.Item) object.Item {
	arg := args[0]

	switch {
	case arg.Type() == object.IntegerType:
		rightVal := arg.(*object.Integer).Value()
		return NewInteger(rightVal)
	case arg.Type() == object.DecimalType:
		rightVal := arg.(*object.Decimal).Value()
		return NewDecimal(rightVal)
	case arg.Type() == object.DoubleType:
		rightVal := arg.(*object.Double).Value()
		return NewDouble(rightVal)
	}

	return object.NIL
}

func numericUnaryMinus(ctx *object.Context, args ...object.Item) object.Item {
	arg := args[0]

	switch {
	case arg.Type() == object.IntegerType:
		rightVal := arg.(*object.Integer).Value()
		return NewInteger(-1 * rightVal)
	case arg.Type() == object.DecimalType:
		rightVal := arg.(*object.Decimal).Value()
		return NewDecimal(-1 * rightVal)
	case arg.Type() == object.DoubleType:
		rightVal := arg.(*object.Double).Value()
		return NewDouble(-1 * rightVal)
	}

	return object.NIL
}
