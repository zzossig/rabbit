package bif

import (
	"strconv"
	"strings"

	"github.com/zzossig/xpath/object"
)

func concat(args ...object.Item) object.Item {
	var sb strings.Builder
	for _, arg := range args {
		switch arg := arg.(type) {
		case *object.Integer:
			val := strconv.FormatInt(int64(arg.Value()), 10)
			sb.WriteString(val)
		case *object.Decimal:
			val := strconv.FormatFloat(arg.Value(), 'f', -1, 64)
			sb.WriteString(val)
		case *object.Double:
			val := strconv.FormatFloat(arg.Value(), 'f', -1, 64)
			sb.WriteString(val)
		default:
			sb.WriteString(arg.Inspect())
		}

	}
	return NewString(sb.String())
}

func upperCase(args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arg := args[0]
	if arg.Type() != object.StringType {
		return object.NIL
	}
	strItem := arg.(*object.String)
	return NewString(strings.ToUpper(strItem.Value()))
}

func lowerCase(args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arg := args[0]
	if arg.Type() != object.StringType {
		return object.NIL
	}
	strItem := arg.(*object.String)
	return NewString(strings.ToLower(strItem.Value()))
}
