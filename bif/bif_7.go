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
			return &object.Boolean{Value: false}
		}
		if arg.Items[0].Type() == object.NodeType {
			return &object.Boolean{Value: true}
		}
	case *object.Boolean:
		return arg
	case *object.String:
		if len(arg.Value) == 0 {
			return &object.Boolean{Value: false}
		}
		return &object.Boolean{Value: true}
	case *object.Integer:
		if arg.Value == 0 {
			return &object.Boolean{Value: false}
		}
		return &object.Boolean{Value: true}
	case *object.Decimal:
		if arg.Value == 0.0 {
			return &object.Boolean{Value: false}
		}
		return &object.Boolean{Value: true}
	case *object.Double:
		if arg.Value == 0.0 {
			return &object.Boolean{Value: false}
		}
		return &object.Boolean{Value: true}
	}

	return NewError("[err:FORG0006]")
}
