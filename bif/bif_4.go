package bif

import (
	"math"

	"github.com/zzossig/xpath/object"
)

func fnNumber(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("wrong number of arguments. got=%d, expected=0 or 1", len(args))
	}

	if len(args) == 1 {
		casted := CastType(args[0], object.DoubleType)
		if IsError(casted) {
			return NewDouble(math.NaN())
		}
		return casted
	}

	seq := &object.Sequence{}

	if len(ctx.CNode) > 0 {
		for _, n := range ctx.CNode {
			casted := CastType(NewString(n.Text()), object.DoubleType)
			if IsError(casted) {
				casted = NewDouble(math.NaN())
			}
			seq.Items = append(seq.Items, casted)
		}
		return seq
	}

	return NewError("context node is not defined")
}
