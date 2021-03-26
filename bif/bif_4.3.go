package bif

import "github.com/zzossig/xpath/object"

func numericEqual(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal == rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(float64(leftVal) == rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(float64(leftVal) == rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal == float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal == rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal == rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal == float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal == rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal == rightVal)
	}

	return object.NIL
}

func numericLessThan(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal < rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(float64(leftVal) < rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(float64(leftVal) < rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal < float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal < rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal < rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal < float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal < rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal < rightVal)
	}

	return object.NIL
}

func numericGreaterThan(ctx *object.Context, args ...object.Item) object.Item {
	arg1 := args[0]
	arg2 := args[1]

	switch {
	case arg1.Type() == object.IntegerType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal > rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(float64(leftVal) > rightVal)
	case arg1.Type() == object.IntegerType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Integer).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(float64(leftVal) > rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal > float64(rightVal))
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal > rightVal)
	case arg1.Type() == object.DecimalType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Decimal).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal > rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.IntegerType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Integer).Value()
		return NewBoolean(leftVal > float64(rightVal))
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DecimalType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Decimal).Value()
		return NewBoolean(leftVal > rightVal)
	case arg1.Type() == object.DoubleType && arg2.Type() == object.DoubleType:
		leftVal := arg1.(*object.Double).Value()
		rightVal := arg2.(*object.Double).Value()
		return NewBoolean(leftVal > rightVal)
	}

	return object.NIL
}
