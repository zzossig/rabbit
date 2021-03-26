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
	if item, ok := arg.(*object.String); ok {
		return NewString(strings.ToUpper(item.Value()))
	}
	if item, ok := arg.(*object.BaseNode); ok {
		return NewString(strings.ToUpper(item.Text()))
	}
	if item, ok := arg.(*object.AttrNode); ok {
		return NewString(strings.ToUpper(item.Text()))
	}
	if seq, ok := arg.(*object.Sequence); ok {
		if len(seq.Items) != 1 {
			return NewError("wrong number of sequence items. got=%d, want=1", len(args))
		}
		return upperCase(seq.Items[0])
	}

	return NewError("cannot match item type with required type")
}

func lowerCase(args ...object.Item) object.Item {
	if len(args) != 1 {
		return NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arg := args[0]
	if item, ok := arg.(*object.String); ok {
		return NewString(strings.ToLower(item.Value()))
	}
	if item, ok := arg.(*object.BaseNode); ok {
		return NewString(strings.ToLower(item.Text()))
	}
	if item, ok := arg.(*object.AttrNode); ok {
		return NewString(strings.ToLower(item.Text()))
	}
	if seq, ok := arg.(*object.Sequence); ok {
		if len(seq.Items) != 1 {
			return NewError("wrong number of sequence items. got=%d, want=1", len(args))
		}
		return lowerCase(seq.Items[0])
	}

	return NewError("cannot match item type with required type")
}
