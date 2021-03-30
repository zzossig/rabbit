package bif

import "github.com/zzossig/xpath/object"

func arrSize(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:size")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:size")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)
	return NewInteger(len(arr.Items))
}

func arrGet(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: arr:get")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: arr:get")
	}

	if !IsArray(args[0]) || args[1].Type() != object.IntegerType {
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)
	pos := args[1].(*object.Integer)
	idx := pos.Value() - 1

	if idx > len(arr.Items)-1 || idx < 0 {
		return NewError("index out of range")
	}

	return arr.Items[idx]
}

func arrPut(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 3 {
		return NewError("too few parameters for function call: arr:put")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: arr:put")
	}

	if !IsArray(args[0]) || args[1].Type() != object.IntegerType {
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)
	pos := args[1].(*object.Integer)
	mem := args[2]
	idx := pos.Value() - 1

	if idx > len(arr.Items)-1 || idx < 0 {
		return NewError("index out of range")
	}

	arr.Items[idx] = mem
	return arr
}

func arrAppend(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: arr:append")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: arr:append")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)
	arr.Items = append(arr.Items, args[1])

	return arr
}

func arrSubarray(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: arr:subarray")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: arr:subarray")
	}

	var arr *object.Array
	var start int
	var length int

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}
	arr = args[0].(*object.Array)

	if len(args) == 3 {
		arg1, ok1 := args[1].(*object.Integer)
		arg2, ok2 := args[2].(*object.Integer)
		if !ok1 || !ok2 {
			return NewError("cannot match item type with required type")
		}
		start = arg1.Value() - 1
		length = arg2.Value()
	} else {
		arg1, ok := args[1].(*object.Integer)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		start = arg1.Value() - 1
		length = len(arr.Items) - start
	}

	if start < 0 || length < 0 || start+length > len(arr.Items) {
		return NewError("index out of range")
	}

	arr.Items = arr.Items[start : start+length]

	return arr
}

func arrRemove(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: arr:remove")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: arr:remove")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}

	var positions []int
	switch arg := args[1].(type) {
	case *object.Sequence:
		for _, item := range arg.Items {
			if i, ok := item.(*object.Integer); ok {
				positions = append(positions, i.Value())
			} else {
				return NewError("cannot match item type with required type")
			}
		}
	case *object.Integer:
		positions = append(positions, arg.Value())
	default:
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)

	for _, pos := range positions {
		idx := pos - 1
		if idx < 0 || idx > len(arr.Items)-1 {
			return NewError("index out of range")
		}

		arr.Items[idx] = nil
	}

	result := &object.Array{}

	for _, item := range arr.Items {
		if item != nil {
			result.Items = append(result.Items, item)
		}
	}

	return result
}

func arrInsertBefore(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 3 {
		return NewError("too few parameters for function call: arr:insert-before")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: arr:insert-before")
	}

	if !IsArray(args[0]) || args[1].Type() != object.IntegerType {
		return NewError("cannot match item type with required type")
	}

	arr := args[0].(*object.Array)
	pos := args[1].(*object.Integer)
	mem := args[2]
	idx := pos.Value() - 1

	if idx > len(arr.Items) || idx < 0 {
		return NewError("index out of range")
	}

	arr.Items = append(arr.Items[:idx], append([]object.Item{mem}, arr.Items[idx:]...)...)
	return arr
}

func arrHead(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:head")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:head")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}
	if IsArrayEmpty(args[0]) {
		return NewError("index our of range")
	}

	arr := args[0].(*object.Array)
	return arr.Items[0]
}

func arrTail(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:tail")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:tail")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}
	if IsArrayEmpty(args[0]) {
		return NewError("index our of range")
	}

	return arrRemove(nil, args[0], NewInteger(1))
}

func arrReverse(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:reverse")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:reverse")
	}

	if !IsArray(args[0]) {
		return NewError("cannot match item type with required type")
	}
	if IsArrayEmpty(args[0]) {
		return args[0]
	}

	arr := args[0].(*object.Array)
	for i := len(arr.Items)/2 - 1; i >= 0; i-- {
		opp := len(arr.Items) - 1 - i
		arr.Items[i], arr.Items[opp] = arr.Items[opp], arr.Items[i]
	}

	return arr
}

func arrJoin(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:join")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:join")
	}

	if IsArray(args[0]) {
		return args[0]
	}
	if !IsSeq(args[0]) {
		return NewError("cannot match item type with required type")
	}

	result := &object.Array{}

	if IsSeqEmpty(args[0]) {
		return result
	}

	seq := args[0].(*object.Sequence)
	for _, item := range seq.Items {
		if arr, ok := item.(*object.Array); ok {
			result.Items = append(result.Items, arr.Items...)
		} else {
			return NewError("cannot match item type with required type")
		}
	}

	return result
}

func arrFlatten(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: arr:flatten")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: arr:flatten")
	}

	if IsSeqEmpty(args[0]) || IsArrayEmpty(args[0]) {
		return NewSequence()
	}

	items := UnwrapArr(args[0])
	return &object.Sequence{Items: items}
}

func arrFilter(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func arrSort(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func arrForEach(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}

func arrForEachPair(ctx *object.Context, args ...object.Item) object.Item {
	return nil
}
