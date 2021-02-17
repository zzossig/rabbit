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
	return &object.String{Value: sb.String()}
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
	return &object.String{Value: strings.ToUpper(strItem.Value)}
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
	return &object.String{Value: strings.ToLower(strItem.Value)}
}
