package bif

import (
	"math"

	"github.com/zzossig/xpath/object"
)

func abs(ctx *object.Context, args ...object.Item) object.Item {
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

	return object.NIL
}
