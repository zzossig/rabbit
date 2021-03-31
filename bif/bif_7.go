package bif

import "github.com/zzossig/xpath/object"

func fnTrue(ctx *object.Context, args ...object.Item) object.Item {
	return NewBoolean(true)
}

func fnFalse(ctx *object.Context, args ...object.Item) object.Item {
	return NewBoolean(false)
}

func fnBoolean(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:boolean")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:boolean")
	}

	switch arg := args[0].(type) {
	case *object.Sequence:
		if len(arg.Items) == 0 {
			return NewBoolean(false)
		}
		if IsNode(arg.Items[0]) {
			return NewBoolean(true)
		}
		if len(arg.Items) == 1 {
			return fnBoolean(nil, arg.Items[0])
		}
		if len(arg.Items) > 1 {
			return NewError("too many items supplied for effective boolean value")
		}
	case *object.Boolean:
		return arg
	case *object.String:
		if len(arg.Value()) == 0 {
			return NewBoolean(false)
		}
		return NewBoolean(true)
	case *object.Integer:
		if arg.Value() == 0 {
			return NewBoolean(false)
		}
		return NewBoolean(true)
	case *object.Decimal:
		if arg.Value() == 0.0 {
			return NewBoolean(false)
		}
		return NewBoolean(true)
	case *object.Double:
		if arg.Value() == 0.0 {
			return NewBoolean(false)
		}
		return NewBoolean(true)
	case *object.BaseNode:
		return NewBoolean(true)
	case *object.AttrNode:
		return NewBoolean(true)
	}

	return NewError("unexpected argument type: %s", args[0].Type())
}

func fnNot(ctx *object.Context, args ...object.Item) object.Item {
	item := fnBoolean(nil, args...)
	if IsError(item) {
		return item
	}

	boolObj := item.(*object.Boolean)
	boolObj.SetValue(!boolObj.Value())
	return boolObj
}

func opBooleanEqual(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: op:boolean-equal")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: op:boolean-equal")
	}
	if !IsBoolean(args[0]) || !IsBoolean(args[1]) {
		return NewError("cannot match item type with required type")
	}

	arg1 := args[0].(*object.Boolean)
	arg2 := args[1].(*object.Boolean)

	return NewBoolean(arg1.Value() == arg2.Value())
}

func opBooleanLessThan(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: op:boolean-equal")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: op:boolean-equal")
	}
	if !IsBoolean(args[0]) || !IsBoolean(args[1]) {
		return NewError("cannot match item type with required type")
	}

	arg1 := args[0].(*object.Boolean)
	arg2 := args[1].(*object.Boolean)

	return NewBoolean(!arg1.Value() && arg2.Value())
}

func opBooleanGreaterThan(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: op:boolean-equal")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: op:boolean-equal")
	}
	if !IsBoolean(args[0]) || !IsBoolean(args[1]) {
		return NewError("cannot match item type with required type")
	}

	arg1 := args[0].(*object.Boolean)
	arg2 := args[1].(*object.Boolean)

	return opBooleanLessThan(nil, arg2, arg1)
}
