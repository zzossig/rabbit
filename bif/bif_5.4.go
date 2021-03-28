package bif

import (
	"strings"

	"github.com/zzossig/xpath/object"
)

func fnConcat(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:concat")
	}

	seq := &object.Sequence{}
	seq.Items = append(seq.Items, args...)
	return fnStringJoin(ctx, seq)
}

func fnStringJoin(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) == 0 {
		return NewError("too few parameters for function call: fn:string-join")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: fn:string-join")
	}

	if len(args) == 1 {
		if IsSeq(args[0]) {
			var sb strings.Builder
			seq := args[0].(*object.Sequence)

			for _, item := range seq.Items {
				str := CastType(item, object.StringType)
				if IsError(str) {
					return str
				}
				strObj := str.(*object.String)
				sb.WriteString(strObj.Value())
			}

			return NewString(sb.String())
		} else {
			return CastType(args[0], object.StringType)
		}
	}

	if args[1].Type() != object.StringType {
		return NewError("cannot match item type with required type")
	}
	if !IsSeq(args[0]) {
		return CastType(args[0], object.StringType)
	}

	var elems []string
	sep := args[1].(*object.String)

	seq := args[0].(*object.Sequence)
	for _, item := range seq.Items {
		str := CastType(item, object.StringType)
		if IsError(str) {
			return str
		}
		strObj := str.(*object.String)
		elems = append(elems, strObj.Value())
	}

	return NewString(strings.Join(elems, sep.Value()))
}

func fnSubstring(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func fnStringLength(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func fnNormalizeSpace(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func fnUpperCase(ctx *object.Context, args ...object.Item) object.Item {
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
		return fnUpperCase(ctx, seq.Items[0])
	}

	return NewError("cannot match item type with required type")
}

func fnLowerCase(ctx *object.Context, args ...object.Item) object.Item {
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
		return fnLowerCase(ctx, seq.Items[0])
	}

	return NewError("cannot match item type with required type")
}
