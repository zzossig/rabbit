package bif

import "github.com/zzossig/xpath/object"

func boolean(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("%d arguments supplied, 1 expected", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Sequence:
		if len(arg.Items) == 0 {
			return object.FALSE
		}
		if IsNode(arg.Items[0]) {
			return object.TRUE
		}
		if len(arg.Items) == 1 {
			return boolean(ctx, arg.Items[0])
		}
	case *object.Boolean:
		return arg
	case *object.String:
		if len(arg.Value()) == 0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Integer:
		if arg.Value() == 0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Decimal:
		if arg.Value() == 0.0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Double:
		if arg.Value() == 0.0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.BaseNode:
		return object.TRUE
	case *object.AttrNode:
		return object.TRUE
	}

	return object.FALSE
}
