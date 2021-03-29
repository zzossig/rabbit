package bif

import "github.com/zzossig/xpath/object"

func fnForEach(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:for-each")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: fn:for-each")
	}

	seq := &object.Sequence{}

	if IsSeq(args[0]) {
		seq = args[0].(*object.Sequence)
	} else {
		seq.Items = append(seq.Items, args[0])
	}

	result := &object.Sequence{}

	switch action := args[1].(type) {
	case *object.FuncNamed:
		f := *action.Func

		if action.Num != 1 {
			return NewError("wrong number of parameters. got=#%d, expected=#1", action.Num)
		}

		for _, item := range seq.Items {
			result.Items = append(result.Items, f(ctx, item))
		}
	case *object.FuncInline:
		if len(action.PL.Params) != 1 {
			return NewError("wrong number of parameters. got=%d, expected=1", len(action.PL.Params))
		}

		enclosedCtx := object.NewEnclosedContext(ctx)
		for _, item := range seq.Items {
			enclosedCtx.Set(action.PL.Params[0].Value(), item)
			a := action.Fn(action.Body, enclosedCtx)
			result.Items = append(result.Items, a)
		}
	case *object.FuncPartial:
		f := *action.Func

		if action.PCnt != 1 {
			return NewError("wrong number of placeholder. got=%d, expected=1", action.PCnt)
		}

		for _, item := range seq.Items {
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
					it, ok := ctx.Get(arg.Inspect())
					if !ok {
						return NewError("variable not defined: $%s", arg.Inspect())
					}
					a = append(a, it)
				case object.PholderType:
					a = append(a, item)
					pcnt++
				}
			}

			result.Items = append(result.Items, f(ctx, a...))
		}
	default:
		return NewError("not supported action type in fn:for-each. got %s", action.Type())
	}

	return result
}

func fnForEachPair(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 3 {
		return NewError("too few parameters for function call: fn:for-each-pair")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: fn:for-each-pair")
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

	var minLen int
	if len(seq[0].Items) > len(seq[1].Items) {
		minLen = len(seq[1].Items)
	} else {
		minLen = len(seq[0].Items)
	}

	result := &object.Sequence{}

	switch action := args[2].(type) {
	case *object.FuncNamed:
		f := *action.Func

		if action.Num != 2 {
			return NewError("wrong number of parameters. got=#%d, expected=#2", action.Num)
		}

		for i := 0; i < minLen; i++ {
			a := []object.Item{}
			for j := 0; j < action.Num; j++ {
				a = append(a, seq[j].Items[i])
			}

			result.Items = append(result.Items, f(ctx, a...))
		}
	case *object.FuncInline:
		if len(action.PL.Params) != 2 {
			return NewError("wrong number of parameters. got=%d, expected=2", len(action.PL.Params))
		}

		enclosedCtx := object.NewEnclosedContext(ctx)
		for i := 0; i < minLen; i++ {
			for j, param := range action.PL.Params {
				enclosedCtx.Set(param.Value(), seq[j].Items[i])
			}

			a := action.Fn(action.Body, enclosedCtx)
			result.Items = append(result.Items, a)
		}
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
					it, ok := ctx.Get(arg.Inspect())
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
		return NewError("not supported action type in for-each-pair. got %s", action.Type())
	}

	return result
}

func fnFilter(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: fn:filter")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: fn:filter")
	}

	seq := &object.Sequence{}

	if IsSeq(args[0]) {
		seq = args[0].(*object.Sequence)
	} else {
		seq.Items = append(seq.Items, args[0])
	}

	result := &object.Sequence{}

	switch action := args[1].(type) {
	case *object.FuncNamed:
		f := *action.Func

		if action.Num != 1 {
			return NewError("wrong number of parameters. got=#%d, expected=#1", action.Num)
		}

		for _, item := range seq.Items {
			b := f(ctx, item)
			if !IsBoolean(b) {
				return NewError("cannot match item type with required type")
			}
			bObj := b.(*object.Boolean)
			if bObj.Value() {
				result.Items = append(result.Items, item)
			}
		}
	case *object.FuncInline:
		if len(action.PL.Params) != 1 {
			return NewError("wrong number of parameters. got=%d, expected=1", len(action.PL.Params))
		}

		enclosedCtx := object.NewEnclosedContext(ctx)
		for _, item := range seq.Items {
			enclosedCtx.Set(action.PL.Params[0].Value(), item)
			a := action.Fn(action.Body, enclosedCtx)
			i := UnwrapSeq(a)
			if len(i) == 1 {
				if b, ok := i[0].(*object.Boolean); ok {
					if b.Value() {
						result.Items = append(result.Items, item)
					}
				} else {
					return NewError("cannot match item type with required type")
				}
			}
		}
	case *object.FuncPartial:
		f := *action.Func

		if action.PCnt != 1 {
			return NewError("wrong number of placeholder. got=%d, expected=1", action.PCnt)
		}

		for _, item := range seq.Items {
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
					it, ok := ctx.Get(arg.Inspect())
					if !ok {
						return NewError("variable not defined: $%s", arg.Inspect())
					}
					a = append(a, it)
				case object.PholderType:
					a = append(a, item)
					pcnt++
				}
			}

			b := f(ctx, a...)
			if !IsBoolean(b) {
				return NewError("cannot match item type with required type")
			}
			bObj := b.(*object.Boolean)
			if bObj.Value() {
				result.Items = append(result.Items, item)
			}
		}
	default:
		return NewError("not supported action type in fn:filter. got %s", action.Type())
	}

	return result
}
