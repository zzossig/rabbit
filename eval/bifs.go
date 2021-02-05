package eval

import (
	"math"

	"github.com/zzossig/xpath/object"
)

var builtins = map[string]*object.Builtin{
	"abs": {
		Name: "abs",
		Func: func(args ...object.Item) object.Item {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				v := arg.Value
				if v < 0 {
					v = -v
				}
				return &object.Integer{Value: v}
			case *object.Decimal:
				return &object.Decimal{Value: math.Abs(arg.Value)}
			case *object.Double:
				return &object.Double{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to `abs` not supported, got %s", args[0].Type())
			}
		},
	},
}
