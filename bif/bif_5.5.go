package bif

import (
	"strings"

	"github.com/zzossig/xpath/object"
)

func fnContains(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:contains")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:contains")
	}

	var s string
	var ss string

	switch args[0].Type() {
	case object.StringType:
		s = args[0].(*object.String).Value()
		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	case object.SequenceType:
		seq := args[0].(*object.Sequence)
		if len(seq.Items) == 0 {
			s = ""
		}
		if len(seq.Items) == 1 {
			seqItem, ok := seq.Items[0].(*object.String)
			if !ok {
				return NewError("cannot match item type with required type")
			}
			s = seqItem.Value()
		}

		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	default:
		return NewError("cannot match item type with required type")
	}

	isContain := strings.Contains(s, ss)
	return NewBoolean(isContain)
}

func fnStartsWith(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:starts-with")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:starts-with")
	}

	var s string
	var ss string

	switch args[0].Type() {
	case object.StringType:
		s = args[0].(*object.String).Value()
		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	case object.SequenceType:
		seq := args[0].(*object.Sequence)
		if len(seq.Items) == 0 {
			s = ""
		}
		if len(seq.Items) == 1 {
			seqItem, ok := seq.Items[0].(*object.String)
			if !ok {
				return NewError("cannot match item type with required type")
			}
			s = seqItem.Value()
		}

		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	default:
		return NewError("cannot match item type with required type")
	}

	isSW := strings.HasPrefix(s, ss)
	return NewBoolean(isSW)
}

func fnEndsWith(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:ends-with")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:ends-with")
	}

	var s string
	var ss string

	switch args[0].Type() {
	case object.StringType:
		s = args[0].(*object.String).Value()
		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	case object.SequenceType:
		seq := args[0].(*object.Sequence)
		if len(seq.Items) == 0 {
			s = ""
		}
		if len(seq.Items) == 1 {
			seqItem, ok := seq.Items[0].(*object.String)
			if !ok {
				return NewError("cannot match item type with required type")
			}
			s = seqItem.Value()
		}

		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	default:
		return NewError("cannot match item type with required type")
	}

	isEW := strings.HasSuffix(s, ss)
	return NewBoolean(isEW)
}
func fnSubstringBefore(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:substring-before")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:substring-before")
	}

	var s string
	var ss string

	switch args[0].Type() {
	case object.StringType:
		s = args[0].(*object.String).Value()
		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	case object.SequenceType:
		seq := args[0].(*object.Sequence)
		if len(seq.Items) == 0 {
			s = ""
		}
		if len(seq.Items) == 1 {
			seqItem, ok := seq.Items[0].(*object.String)
			if !ok {
				return NewError("cannot match item type with required type")
			}
			s = seqItem.Value()
		}

		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	default:
		return NewError("cannot match item type with required type")
	}

	str := strings.TrimRight(s, ss)
	if str == s {
		return NewString("")
	}
	return NewString(str)
}

func fnSubstringAfter(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:substring-after")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:substring-after")
	}

	var s string
	var ss string

	switch args[0].Type() {
	case object.StringType:
		s = args[0].(*object.String).Value()
		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	case object.SequenceType:
		seq := args[0].(*object.Sequence)
		if len(seq.Items) == 0 {
			s = ""
		}
		if len(seq.Items) == 1 {
			seqItem, ok := seq.Items[0].(*object.String)
			if !ok {
				return NewError("cannot match item type with required type")
			}
			s = seqItem.Value()
		}

		switch args[1].Type() {
		case object.StringType:
			ss = args[1].(*object.String).Value()
		case object.SequenceType:
			seq := args[1].(*object.Sequence)
			if len(seq.Items) == 0 {
				ss = ""
			}
			if len(seq.Items) == 1 {
				seqItem, ok := seq.Items[0].(*object.String)
				if !ok {
					return NewError("cannot match item type with required type")
				}
				ss = seqItem.Value()
			}
		default:
			return NewError("cannot match item type with required type")
		}
	default:
		return NewError("cannot match item type with required type")
	}

	str := strings.TrimLeft(s, ss)
	if str == s {
		return NewString("")
	}
	return NewString(str)
}
