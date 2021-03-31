package bif

import (
	"math"

	"github.com/zzossig/rabbit/object"
)

func fnAbs(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:abs")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:abs")
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
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:ceiling")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:ceiling")
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
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:floor")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:floor")
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
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:round")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:round")
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
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:round-half-to-even")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:round-half-to-even")
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
