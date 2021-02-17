package bif

import (
	"github.com/zzossig/xpath/object"
)

func boolean(args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("%d arguments supplied, 1 expected", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Sequence:
		if len(arg.Items) == 0 {
			return object.FALSE
		}
		if arg.Items[0].Type() == object.NodeType {
			return object.TRUE
		}
		if len(arg.Items) == 1 {
			return boolean(arg.Items[0])
		}
	case *object.Boolean:
		return arg
	case *object.String:
		if len(arg.Value) == 0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Integer:
		if arg.Value == 0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Decimal:
		if arg.Value == 0.0 {
			return object.FALSE
		}
		return object.TRUE
	case *object.Double:
		if arg.Value == 0.0 {
			return object.FALSE
		}
		return object.TRUE
	}

	return NewError("[err:FORG0006]")
}
