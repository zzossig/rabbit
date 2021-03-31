package bif

import "github.com/zzossig/rabbit/object"

func opNumericEqual(ctx *object.Context, args ...object.Item) object.Item {
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

	return NewError("cannot eqaul types: %s, %s", arg1.Type(), arg2.Type())
}

func opNumericLessThan(ctx *object.Context, args ...object.Item) object.Item {
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

	return NewError("cannot less than types: %s, %s", arg1.Type(), arg2.Type())
}

func opNumericGreaterThan(ctx *object.Context, args ...object.Item) object.Item {
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

	return NewError("cannot greater than types: %s, %s", arg1.Type(), arg2.Type())
}
