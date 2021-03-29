package bif

import "github.com/zzossig/xpath/object"

func xsInteger(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: xs:integer")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: xs:integer")
	}
	return CastType(args[0], object.IntegerType)
}

func xsDecimal(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: xs:decimal")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: xs:decimal")
	}
	return CastType(args[0], object.DecimalType)
}

func xsDouble(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: xs:double")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: xs:double")
	}
	return CastType(args[0], object.DoubleType)
}

func xsString(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: xs:string")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: xs:string")
	}
	return CastType(args[0], object.StringType)
}

func xsBoolean(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: xs:boolean")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: xs:boolean")
	}
	return CastType(args[0], object.BooleanType)
}
