package bif

import (
	"strings"

	"github.com/zzossig/rabbit/object"
)

func fnCodepointsToString(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:codepoints-to-string")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:codepoints-to-string")
	}

	if IsSeqEmpty(args[0]) || IsArrayEmpty(args[0]) {
		return NewString("")
	}

	seq := &object.Sequence{}

	switch arg := args[0].(type) {
	case *object.Sequence:
		seq.Items = arg.Items
	case *object.Array:
		seq.Items = arg.Items
	case *object.Integer:
		if arg.Value() == 0 {
			return NewError("not allowed value in fn:codepoints-to-string: 0")
		}
		seq.Items = append(seq.Items, arg)
	default:
		seq.Items = append(seq.Items, arg)
	}

	var sb strings.Builder

	for _, item := range seq.Items {
		switch item := item.(type) {
		case *object.Integer:
			sb.WriteRune(rune(item.Value()))
		default:
			return NewError("cannot match item type with required type")
		}
	}

	return NewString(sb.String())
}

func fnStringToCodepoints(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:string-to-codepoints")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:string-to-codepoints")
	}

	if IsSeqEmpty(args[0]) || IsArrayEmpty(args[0]) {
		return NewSequence()
	}

	var str string

	switch arg := args[0].(type) {
	case *object.Sequence:
		if len(arg.Items) > 1 {
			return NewError("too many items in the sequence")
		}
		return fnStringToCodepoints(ctx, arg.Items[0])
	case *object.Array:
		if len(arg.Items) > 1 {
			return NewError("too many items in the array")
		}
		return fnStringToCodepoints(ctx, arg.Items[0])
	case *object.String:
		str = arg.Value()
	default:
		return NewError("cannot match item type with required type")
	}

	result := &object.Sequence{}

	for _, r := range str {
		result.Items = append(result.Items, NewInteger(int(r)))
	}

	return result
}
