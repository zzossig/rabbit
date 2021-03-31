package bif

import (
	"math"
	"strings"

	"github.com/zzossig/rabbit/object"
)

func fnCount(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:count")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:count")
	}

	if !IsSeq(args[0]) {
		return NewInteger(1)
	}

	seq := args[0].(*object.Sequence)
	return NewInteger(len(seq.Items))
}

func fnAvg(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:avg")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:avg")
	}

	if IsSeqEmpty(args[0]) {
		return args[0]
	}

	sum := fnSum(nil, args[0])
	cnt := fnCount(nil, args[0])
	ty := sum.Type()

	dsum := CastType(sum, object.DecimalType)
	if IsError(dsum) {
		return dsum
	}
	dcnt := CastType(cnt, object.DecimalType)
	if IsError(dcnt) {
		return dcnt
	}

	dsumObj := dsum.(*object.Decimal)
	dcntObj := dcnt.(*object.Decimal)
	avg := NewDecimal(dsumObj.Value() / dcntObj.Value())

	return CastType(avg, ty)
}

func fnMax(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:max")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:max")
	}

	if IsSeqEmpty(args[0]) {
		return args[0]
	}
	if IsArrayEmpty(args[0]) {
		return NewSequence()
	}
	if !IsSeq(args[0]) && !IsArray(args[0]) {
		if IsNumeric(args[0]) {
			return args[0]
		} else {
			return NewError("cannot match item type with required type")
		}
	}

	src := &object.Sequence{}
	src.Items = UnwrapArr(args[0])

	switch {
	case IsNumeric(src.Items[0]):
		var max float64
		ty := src.Items[0].Type()

		for i, item := range src.Items {
			if !IsNumeric(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			d := CastType(item, object.DecimalType)
			if IsError(d) {
				return d
			}
			dObj := d.(*object.Decimal)

			if i == 0 {
				max = dObj.Value()
			} else {
				max = math.Max(dObj.Value(), max)
			}
		}

		maxObj := NewDecimal(max)
		return CastType(maxObj, ty)
	case IsString(src.Items[0]):
		var max string

		for i, item := range src.Items {
			if !IsString(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			str := item.(*object.String)
			if i == 0 {
				max = str.Value()
			} else {
				if strings.Compare(str.Value(), max) > 0 {
					max = str.Value()
				}
			}
		}

		return NewString(max)
	case IsBoolean(src.Items[0]):
		for _, item := range src.Items {
			if !IsBoolean(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			b := item.(*object.Boolean)
			if b.Value() {
				return NewBoolean(true)
			}
		}
		return NewBoolean(false)
	default:
		return NewError("unexpected argument type: %s", src.Items[0].Type())
	}
}

func fnMin(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:min")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:min")
	}

	if IsSeqEmpty(args[0]) {
		return args[0]
	}
	if IsArrayEmpty(args[0]) {
		return NewSequence()
	}
	if !IsSeq(args[0]) && !IsArray(args[0]) {
		if IsNumeric(args[0]) {
			return args[0]
		} else {
			return NewError("cannot match item type with required type")
		}
	}

	src := &object.Sequence{}
	src.Items = UnwrapArr(args[0])

	switch {
	case IsNumeric(src.Items[0]):
		var min float64
		ty := src.Items[0].Type()

		for i, item := range src.Items {
			if !IsNumeric(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			d := CastType(item, object.DecimalType)
			if IsError(d) {
				return d
			}
			dObj := d.(*object.Decimal)

			if i == 0 {
				min = dObj.Value()
			} else {
				min = math.Min(dObj.Value(), min)
			}
		}

		minObj := NewDecimal(min)
		return CastType(minObj, ty)
	case IsString(src.Items[0]):
		var min string

		for i, item := range src.Items {
			if !IsString(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			str := item.(*object.String)
			if i == 0 {
				min = str.Value()
			} else {
				if strings.Compare(str.Value(), min) < 0 {
					min = str.Value()
				}
			}
		}

		return NewString(min)
	case IsBoolean(src.Items[0]):
		for _, item := range src.Items {
			if !IsBoolean(item) {
				return NewError("unexpected argument type: %s", item.Type())
			}

			b := item.(*object.Boolean)
			if !b.Value() {
				return NewBoolean(false)
			}
		}
		return NewBoolean(true)
	default:
		return NewError("unexpected argument type: %s", src.Items[0].Type())
	}
}

func fnSum(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:sum")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:sum")
	}

	if IsSeqEmpty(args[0]) {
		return NewInteger(0)
	}
	if !IsSeq(args[0]) && !IsArray(args[0]) {
		if IsNumeric(args[0]) {
			return args[0]
		} else {
			return NewError("cannot match item type with required type")
		}
	}

	src := &object.Sequence{}
	src.Items = UnwrapArr(args[0])

	sum := 0.0
	ty := object.IntegerType

	for _, item := range src.Items {
		if !IsNumeric(item) {
			return NewError("cannot match item type with required type")
		}
		if (ty == object.IntegerType || ty == object.DecimalType) && item.Type() == object.DoubleType {
			ty = item.Type()
		} else if ty == object.IntegerType && (item.Type() == object.DecimalType || item.Type() == object.DoubleType) {
			ty = item.Type()
		}

		d := CastType(item, object.DecimalType)
		dObj := d.(*object.Decimal)
		sum += dObj.Value()
	}

	switch ty {
	case object.IntegerType:
		return NewInteger(int(sum))
	case object.DecimalType:
		return NewDecimal(sum)
	case object.DoubleType:
		return NewDouble(sum)
	}

	return NewInteger(0)
}
