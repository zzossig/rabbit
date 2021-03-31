package bif

import (
	"strings"

	"github.com/zzossig/rabbit/object"
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
				if IsSeqEmpty(item) {
					continue
				}
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
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:substring")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:substring")
	}

	if len(args) == 2 {
		d := fnRound(nil, args[1])
		if IsError(d) {
			return NewError("cannot match item type with required type")
		}

		start := CastType(d, object.DoubleType)
		str := CastType(args[0], object.StringType)
		if IsError(str) {
			return str
		}

		startObj := start.(*object.Double)
		strObj := str.(*object.String)
		if int(startObj.Value()) > len(strObj.Value()) {
			return NewString("")
		}

		startIdx := int(startObj.Value())
		if startIdx != 0 {
			startIdx -= 1
		}
		if startIdx > len(strObj.Value()) {
			return NewString("")
		}
		return NewString(strObj.Value()[startIdx:])
	}

	d1 := fnRound(nil, args[1])
	d2 := fnRound(nil, args[2])
	if IsError(d1) || IsError(d2) {
		return NewError("cannot match item type with required type")
	}
	if IsSeqEmpty(args[0]) {
		return NewString("")
	}

	start := CastType(d1, object.DoubleType)
	length := CastType(d2, object.DoubleType)
	str := CastType(args[0], object.StringType)
	if IsError(str) {
		return str
	}

	startObj := start.(*object.Double)
	lengthObj := length.(*object.Double)
	strObj := str.(*object.String)

	if int(startObj.Value()) > len(strObj.Value()) {
		return NewString("")
	}

	startIdx := int(startObj.Value())
	lengthVal := int(lengthObj.Value())
	sum := startIdx + lengthVal - 1
	if startIdx != 0 {
		startIdx -= 1
	}
	if startIdx < 0 {
		startIdx = 0
	}

	if sum > len(strObj.Value()) {
		if startIdx > len(strObj.Value()) {
			return NewString("")
		}
		return NewString(strObj.Value()[startIdx:])
	}
	if startIdx > sum {
		return NewString("")
	}
	return NewString(strObj.Value()[startIdx:sum])
}

func fnStringLength(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:string-length")
	}

	if len(args) == 1 {
		if IsSeqEmpty(args[0]) {
			return NewInteger(0)
		}

		str, ok := args[0].(*object.String)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		return NewInteger(len(str.Value()))
	}

	seq := fnString(ctx)
	if IsError(seq) {
		return seq
	}

	result := &object.Sequence{}
	seqObj := seq.(*object.Sequence)
	for _, item := range seqObj.Items {
		item, ok := item.(*object.String)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		result.Items = append(result.Items, NewInteger(len(item.Value())))
	}

	return result
}

func fnNormalizeSpace(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("wrong number of arguments. got=%d, expected=0 or 1", len(args))
	}

	if len(args) == 1 {
		str := CastType(args[0], object.StringType)
		strObj := str.(*object.String)
		return NewString(strings.TrimSpace(strObj.Value()))
	}

	seq := &object.Sequence{}

	if len(ctx.CNode) > 0 {
		for _, n := range ctx.CNode {
			texts := collectText(nil, n)
			str := combineTextsNormalize(texts)
			seq.Items = append(seq.Items, NewString(str))
		}
		return seq
	}

	if ctx.Doc != nil {
		texts := collectText(nil, ctx.Doc)
		str := combineTextsNormalize(texts)
		seq.Items = append(seq.Items, NewString(str))
		return seq
	}

	return NewError("context node is not defined")
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
