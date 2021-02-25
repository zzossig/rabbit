package bif

import (
	"strings"

	"github.com/zzossig/xpath/object"
)

func concat(args ...object.Item) object.Item {
	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(arg.Inspect())
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
