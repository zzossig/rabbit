package bif

import "github.com/zzossig/xpath/object"

func fnForEachPair(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) != 3 {
		return NewError("wrong number of arguments. got=%d, want=3", len(args))
	}

	var seq []object.Sequence

	if IsSeq(args[0]) {
		s := args[0].(*object.Sequence)
		seq = append(seq, *s)
	} else {
		s := object.Sequence{}
		s.Items = append(s.Items, args[0])
		seq = append(seq, s)
	}

	if IsSeq(args[1]) {
		s := args[1].(*object.Sequence)
		seq = append(seq, *s)
	} else {
		s := object.Sequence{}
		s.Items = append(s.Items, args[1])
		seq = append(seq, s)
	}

	action := args[2]
	var minLen int

	if len(seq[0].Items) > len(seq[1].Items) {
		minLen = len(seq[1].Items)
	} else {
		minLen = len(seq[0].Items)
	}

	var result object.Sequence

	switch action := action.(type) {
	case *object.FuncNamed:
	case *object.FuncInline:
	case *object.FuncPartial:
		f := *action.Func

		if action.PCnt != 2 {
			return NewError("wrong number of placeholder. got=%d, expected=2", action.PCnt)
		}

		for i := 0; i < minLen; i++ {
			pcnt := 0
			a := []object.Item{}

			for _, arg := range action.Args {
				switch arg.Type() {
				case object.IntegerType:
					fallthrough
				case object.DecimalType:
					fallthrough
				case object.DoubleType:
					fallthrough
				case object.StringType:
					a = append(a, arg)
				case object.VarrefType:
					it, ok := action.Context.Get(arg.Inspect())
					if !ok {
						return NewError("variable not defined: $%s", arg.Inspect())
					}
					a = append(a, it)
				case object.PholderType:
					a = append(a, seq[pcnt].Items[i])
					pcnt++
				}
			}

			result.Items = append(result.Items, f(ctx, a...))
		}

	default:
		return NewError("not supported type in for-each-pair, got %v", action.Type())
	}

	return &result
}
