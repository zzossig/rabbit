package bif

import (
	"math"

	"github.com/zzossig/xpath/object"
)

func fnAbs(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		v := arg.Value()
		if v < 0 {
			arg.SetValue(-v)
		}
		return arg
	case *object.Decimal:
		arg.SetValue(math.Abs(arg.Value()))
		return arg
	case *object.Double:
		arg.SetValue(math.Abs(arg.Value()))
		return arg
	}

	return NewError("cannot match item type with required type")
}

func fnCeiling(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Decimal:
		arg.SetValue(math.Ceil(arg.Value()))
		return arg
	case *object.Double:
		arg.SetValue(math.Ceil(arg.Value()))
		return arg
	}

	return NewError("cannot match item type with required type")
}

func fnFloor(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Decimal:
		arg.SetValue(math.Floor(arg.Value()))
		return arg
	case *object.Double:
		arg.SetValue(math.Floor(arg.Value()))
		return arg
	}

	return NewError("cannot match item type with required type")
}

func fnRound(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Decimal:
		arg.SetValue(math.Round(arg.Value()))
		return arg
	case *object.Double:
		arg.SetValue(math.Round(arg.Value()))
		return arg
	}

	return NewError("cannot match item type with required type")
}

// round-half-to-even
func fnRoundHTE(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Decimal:
		arg.SetValue(math.RoundToEven(arg.Value()))
		return arg
	case *object.Double:
		arg.SetValue(math.RoundToEven(arg.Value()))
		return arg
	}

	return NewError("cannot match item type with required type")
}
