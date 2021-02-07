package eval

import (
	"math"
	"strings"

	"github.com/zzossig/xpath/object"
)

var builtins = map[string]object.Func{
	"abs": func(args ...object.Item) object.Item {
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
	"concat": func(args ...object.Item) object.Item {
		var sb strings.Builder
		for _, arg := range args {
			sb.WriteString(arg.Inspect())
		}
		return &object.String{Value: sb.String()}
	},
	"for-each-pair": func(args ...object.Item) object.Item {
		if len(args) != 3 {
			return newError("wrong number of arguments. got=%d, want=3", len(args))
		}

		var seq []object.Sequence

		if isSeq(args[0]) {
			s := args[0].(*object.Sequence)
			seq = append(seq, *s)
		} else {
			s := object.Sequence{}
			s.Items = append(s.Items, args[0])
			seq = append(seq, s)
		}

		if isSeq(args[1]) {
			s := args[1].(*object.Sequence)
			seq = append(seq, *s)
		} else {
			s := object.Sequence{}
			s.Items = append(s.Items, args[1])
			seq = append(seq, s)
		}

		action := args[2]
		var minLen int

		if len(seq[0].Items) > len(seq[1].Items) {
			minLen = len(seq[1].Items)
		} else {
			minLen = len(seq[0].Items)
		}

		var result object.Sequence

		switch action := action.(type) {
		case *object.FuncNamed:
		case *object.FuncInline:
		case *object.FuncCall:
			f := *action.Func

			for i := 0; i < minLen; i++ {
				pcnt := 0
				a := []object.Item{}

				for _, arg := range action.Args {
					switch arg.Type() {
					case object.IntegerType:
						fallthrough
					case object.DecimalType:
						fallthrough
					case object.DoubleType:
						fallthrough
					case object.StringType:
						a = append(a, arg)
					case object.PholderType:
						if len(args)-1 <= pcnt {
							return newError("too many arguments")
						}
						ph := arg.(*object.Placeholder)
						ph.Value = seq[pcnt].Items[i]
						pcnt++

						a = append(a, ph)
					}
				}

				if pcnt != 0 && pcnt < len(action.Env.Args)-1 {
					return newError("too few arguments")
				}

				result.Items = append(result.Items, f(a...))
			}

		default:
			return newError("not supported type in for-each-pair, got %v", action.Type())
		}

		return &result
	},
}

func isSeq(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}

func isPlaceholder(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}
