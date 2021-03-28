package bif

import (
	"math"

	"github.com/zzossig/xpath/object"
)

func mathPI(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 0 {
		return NewError("too many parameters for function call: math:pi")
	}
	return NewDouble(math.Pi)
}

func mathExp(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:exp. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}
	dbl := d.(*object.Double)
	dbl.SetValue(math.Exp(dbl.Value()))
	return dbl
}

func mathExp2(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:exp2. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}
	dbl := d.(*object.Double)
	dbl.SetValue(math.Exp2(dbl.Value()))
	return dbl
}

func mathLog(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:log. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}
	dbl := d.(*object.Double)
	dbl.SetValue(math.Log(dbl.Value()))
	return dbl
}

func mathLog2(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:log2. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}
	dbl := d.(*object.Double)
	dbl.SetValue(math.Log2(dbl.Value()))
	return dbl
}

func mathLog10(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:log10. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Log10(dbl.Value()))
	return dbl
}

func mathPow(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 2 {
		return NewError("wrong number of parameters for function call: math:pow. got=%d, want=2", len(args))
	}
	if !IsNumeric(args[0]) || !IsNumeric(args[1]) {
		return NewError("cannot match item type with required type")
	}

	d1 := CastType(args[0], object.DoubleType)
	d2 := CastType(args[1], object.DoubleType)
	if IsError(d1) || IsError(d2) {
		return NewError("cannot match item type with required type")
	}

	dbl1 := d1.(*object.Double)
	dbl2 := d2.(*object.Double)
	return NewDouble(math.Pow(dbl1.Value(), dbl2.Value()))
}

func mathSqrt(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:sqrt. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Sqrt(dbl.Value()))
	return dbl
}

func mathSin(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:sin. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Sin(dbl.Value()))
	return dbl
}

func mathCos(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:cos. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Cos(dbl.Value()))
	return dbl
}

func mathTan(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:tan. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Tan(dbl.Value()))
	return dbl
}

func mathAsin(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:asin. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Asin(dbl.Value()))
	return dbl
}

func mathAcos(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:acos. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Acos(dbl.Value()))
	return dbl
}

func mathAtan(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of parameters for function call: math:atan. got=%d, want=1", len(args))
	}
	if !IsNumeric(args[0]) {
		return NewError("cannot match item type with required type")
	}

	d := CastType(args[0], object.DoubleType)
	if IsError(d) {
		return NewError("cannot match item type with required type")
	}

	dbl := d.(*object.Double)
	dbl.SetValue(math.Atan(dbl.Value()))
	return dbl
}

func mathAtan2(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 2 {
		return NewError("wrong number of parameters for function call: math:atan2. got=%d, want=2", len(args))
	}
	if !IsNumeric(args[0]) || !IsNumeric(args[1]) {
		return NewError("cannot match item type with required type")
	}

	d1 := CastType(args[0], object.DoubleType)
	d2 := CastType(args[1], object.DoubleType)
	if IsError(d1) || IsError(d2) {
		return NewError("cannot match item type with required type")
	}

	dbl1 := d1.(*object.Double)
	dbl2 := d2.(*object.Double)
	return NewDouble(math.Atan2(dbl1.Value(), dbl2.Value()))
}
