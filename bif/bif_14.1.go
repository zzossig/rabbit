package bif

import "github.com/zzossig/rabbit/object"

func fnEmpty(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:empty")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:empty")
	}

	if IsSeqEmpty(args[0]) {
		return NewBoolean(true)
	}
	return NewBoolean(false)
}

func fnExists(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:exists")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:exists")
	}

	if IsSeqEmpty(args[0]) {
		return NewBoolean(false)
	}
	return NewBoolean(true)
}

func fnRemove(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:remove")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: fn:remove")
	}

	if !IsSeq(args[0]) {
		pos, ok := args[1].(*object.Integer)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		if pos.Value() == 1 {
			return NewSequence()
		}
		return args[0]
	}

	target, ok1 := args[0].(*object.Sequence)
	pos, ok2 := args[1].(*object.Integer)
	if !ok1 || !ok2 {
		return NewError("cannot match item type with required type")
	}
	if pos.Value() == 0 {
		return target
	}

	idx := pos.Value() - 1
	if idx > len(target.Items)-1 || idx < 0 {
		return target
	}

	target.Items = append(target.Items[:idx], target.Items[idx+1:]...)
	return target
}

func fnHead(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:head")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:head")
	}

	if IsSeqEmpty(args[0]) || !IsSeq(args[0]) {
		return args[0]
	}

	seq := args[0].(*object.Sequence)
	return NewSequence(seq.Items[0])
}

func fnTail(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:tail")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:tail")
	}

	if IsSeqEmpty(args[0]) || !IsSeq(args[0]) {
		return NewSequence()
	}

	seq := args[0].(*object.Sequence)
	return NewSequence(seq.Items[1:]...)
}

func fnInsertBefore(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 3 {
		return NewError("too few parameters for function call: fn:insert-before")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:insert-before")
	}

	position, ok := args[1].(*object.Integer)
	if !ok {
		return NewError("cannot match item type with required type")
	}
	idx := position.Value()
	if idx != 0 {
		idx -= 1
	}

	target := &object.Sequence{}
	inserts := &object.Sequence{}

	if IsSeq(args[0]) && !IsSeqEmpty(args[0]) {
		seq := args[0].(*object.Sequence)
		target = seq
	} else if !IsSeqEmpty(args[0]) {
		target.Items = append(target.Items, args[0])
	}

	if IsSeq(args[2]) && !IsSeqEmpty(args[2]) {
		seq := args[2].(*object.Sequence)
		inserts = seq
	} else if !IsSeqEmpty(args[2]) {
		inserts.Items = append(inserts.Items, args[2])
	}

	target.Items = append(target.Items[:idx], append(inserts.Items, target.Items[idx:]...)...)
	return target
}

func fnReverse(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:reverse")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:reverse")
	}

	if IsSeqEmpty(args[0]) || !IsSeq(args[0]) {
		return args[0]
	}

	seq := args[0].(*object.Sequence)
	for i := len(seq.Items)/2 - 1; i >= 0; i-- {
		opp := len(seq.Items) - 1 - i
		seq.Items[i], seq.Items[opp] = seq.Items[opp], seq.Items[i]
	}
	return seq
}

func fnSubsequence(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:subsequence")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:subsequence")
	}

	sourceSeq := &object.Sequence{}

	if IsSeq(args[0]) {
		sourceSeq = args[0].(*object.Sequence)
	} else {
		sourceSeq.Items = append(sourceSeq.Items, args[0])
	}

	d := fnRound(nil, args[1])
	if IsError(d) {
		return d
	}

	loc := CastType(d, object.IntegerType)
	if IsError(loc) {
		return loc
	}
	startingLoc := loc.(*object.Integer)

	idx := startingLoc.Value()
	if idx < 0 && len(args) == 2 {
		return sourceSeq
	}
	if idx != 0 {
		idx -= 1
	}

	if idx > len(sourceSeq.Items)-1 {
		return NewSequence()
	}

	if len(args) == 2 {
		return NewSequence(sourceSeq.Items[idx:]...)
	}

	l := fnRound(nil, args[2])
	if IsError(l) {
		return l
	}
	li := CastType(l, object.IntegerType)
	if IsError(li) {
		return li
	}

	length := li.(*object.Integer)
	if length.Value() <= 0 {
		return NewSequence()
	}

	last := idx + length.Value()
	if startingLoc.Value() == 0 {
		last -= 1
	}
	if idx < 0 {
		idx = 0
	}
	if last > len(sourceSeq.Items)-1 {
		return NewSequence(sourceSeq.Items[idx:]...)
	}
	return NewSequence(sourceSeq.Items[idx:last]...)
}
