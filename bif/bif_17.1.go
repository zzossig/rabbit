package bif

import "github.com/zzossig/xpath/object"

func mapSize(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: map:size")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: map:size")
	}

	var m *object.Map

	switch arg := args[0].(type) {
	case *object.Map:
		m = arg
	case *object.Sequence:
		items := UnwrapSeq(arg)
		if len(items) != 1 {
			NewError("wrong number of arguments. got=%d, expected=1", len(items))
		}

		item, ok := items[0].(*object.Map)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		m = item
	default:
		return NewError("cannot match item type with required type")
	}

	return NewInteger(len(m.Pairs))
}

func mapKeys(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: map:keys")
	}
	if len(args) > 1 {
		return NewError("too many parameters for function call: map:keys")
	}

	var m *object.Map

	switch arg := args[0].(type) {
	case *object.Map:
		m = arg
	case *object.Sequence:
		items := UnwrapSeq(arg)
		if len(items) != 1 {
			NewError("wrong number of arguments. got=%d, expected=1", len(items))
		}

		item, ok := items[0].(*object.Map)
		if !ok {
			return NewError("cannot match item type with required type")
		}
		m = item
	default:
		return NewError("cannot match item type with required type")
	}

	keys := &object.Sequence{}

	for _, pair := range m.Pairs {
		keys.Items = append(keys.Items, pair.Key)
	}
	return keys
}

func mapContains(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: map:contains")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:contains")
	}

	if !IsMap(args[0]) || !IsAnyAtomic(args[1]) {
		return NewError("cannot match item type with required type")
	}

	m := args[0].(*object.Map)

	for _, pair := range m.Pairs {
		switch key := pair.Key.(type) {
		case *object.Integer:
			if arg1, ok := args[1].(*object.Integer); ok {
				if key.Value() == arg1.Value() {
					return NewBoolean(true)
				}
			}
		case *object.Decimal:
			if arg1, ok := args[1].(*object.Decimal); ok {
				if key.Value() == arg1.Value() {
					return NewBoolean(true)
				}
			}
		case *object.Double:
			if arg1, ok := args[1].(*object.Double); ok {
				if key.Value() == arg1.Value() {
					return NewBoolean(true)
				}
			}
		case *object.String:
			if arg1, ok := args[1].(*object.String); ok {
				if key.Value() == arg1.Value() {
					return NewBoolean(true)
				}
			}
		case *object.Boolean:
			if arg1, ok := args[1].(*object.Boolean); ok {
				if key.Value() == arg1.Value() {
					return NewBoolean(true)
				}
			}
		}
	}

	return NewBoolean(false)
}

func mapGet(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: map:get")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:get")
	}

	if !IsMap(args[0]) || !IsAnyAtomic(args[1]) {
		return NewError("cannot match item type with required type")
	}

	m := args[0].(*object.Map)

	switch arg1 := args[1].(type) {
	case *object.Integer:
		key := arg1.HashKey()
		if pair, ok := m.Pairs[key]; ok {
			return pair.Value
		}
	case *object.Decimal:
		key := arg1.HashKey()
		if pair, ok := m.Pairs[key]; ok {
			return pair.Value
		}
	case *object.Double:
		key := arg1.HashKey()
		if pair, ok := m.Pairs[key]; ok {
			return pair.Value
		}
	case *object.String:
		key := arg1.HashKey()
		if pair, ok := m.Pairs[key]; ok {
			return pair.Value
		}
	case *object.Boolean:
		key := arg1.HashKey()
		if pair, ok := m.Pairs[key]; ok {
			return pair.Value
		}
	}
	return NewSequence()
}

func mapPut(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 3 {
		return NewError("too few parameters for function call: map:put")
	}
	if len(args) > 3 {
		return NewError("too many parameters for function call: map:put")
	}

	if !IsMap(args[0]) || !IsAnyAtomic(args[1]) {
		return NewError("cannot match item type with required type")
	}

	m := args[0].(*object.Map)
	hashed := args[1].(object.Hasher).HashKey()
	m.Pairs[hashed] = object.Pair{Key: args[1], Value: args[2]}

	return m
}

func mapEntry(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: map:entry")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:entry")
	}

	if !IsAnyAtomic(args[0]) {
		return NewError("cannot match item type with required type")
	}

	pairs := make(map[object.HashKey]object.Pair)
	hashed := args[0].(object.Hasher).HashKey()
	pairs[hashed] = object.Pair{Key: args[0], Value: args[1]}

	return &object.Map{Pairs: pairs}
}

func mapRemove(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: map:remove")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:remove")
	}

	if !IsMap(args[0]) {
		return NewError("cannot match item type with required type")
	}
	if !IsAnyAtomic(args[1]) && !IsSeq(args[1]) {
		return NewError("cannot match item type with required type")
	}

	m := args[0].(*object.Map)

	if IsSeq(args[1]) {
		seq := args[1].(*object.Sequence)
		items := UnwrapSeq(seq)
		for _, item := range items {
			hashed := item.(object.Hasher).HashKey()
			delete(m.Pairs, hashed)
		}
	} else {
		hashed := args[1].(object.Hasher).HashKey()
		delete(m.Pairs, hashed)
	}

	return m
}

func mapMerge(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 1 {
		return NewError("too few parameters for function call: map:merge")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:merge")
	}

	if !IsSeq(args[0]) && !IsMap(args[0]) {
		return NewError("cannot match item type with required type")
	}
	if IsSeqEmpty(args[0]) {
		return &object.Map{}
	}

	var seq *object.Sequence
	if IsMap(args[0]) {
		return args[0]
	}
	seq = args[0].(*object.Sequence)

	var option string
	if len(args) == 2 {
		m := args[1].(*object.Map)
		dup, ok := m.Pairs[NewString("duplicates").HashKey()]

		if !ok {
			option = "use-first"
		} else {
			if v, ok := dup.Value.(*object.String); ok {
				option = v.Value()
			}
		}
	} else {
		option = "use-first"
	}

	var src []*object.Map
	result := &object.Map{}
	result.Pairs = make(map[object.HashKey]object.Pair)

	for _, m := range seq.Items {
		if m.Type() != object.MapType {
			return NewError("cannot match item type with required type")
		}
		m := m.(*object.Map)
		src = append(src, m)
	}

	switch option {
	case "reject":
		for _, m := range src {
			for key, pair := range m.Pairs {
				if _, ok := result.Pairs[key]; ok {
					return NewError("duplicate keys are rejected")
				} else {
					result.Pairs[key] = pair
				}
			}
		}
	case "use-first":
		for _, m := range src {
			for key, pair := range m.Pairs {
				if _, ok := result.Pairs[key]; ok {
					continue
				} else {
					result.Pairs[key] = pair
				}
			}
		}
	case "use-any":
		fallthrough
	case "use-last":
		for _, m := range src {
			for key, pair := range m.Pairs {
				if _, ok := result.Pairs[key]; ok {
					result.Pairs[key] = pair
				} else {
					result.Pairs[key] = pair
				}
			}
		}
	case "combine":
		for _, m := range src {
			for key, pair := range m.Pairs {
				if p, ok := result.Pairs[key]; ok {
					oldPairVal := p.Value
					newPairVal := pair.Value
					Val := NewSequence(oldPairVal, newPairVal)
					result.Pairs[key] = object.Pair{Key: pair.Key, Value: Val}
				} else {
					result.Pairs[key] = pair
				}
			}
		}
	default:
		return NewError("invalid duplicates option: %q", option)
	}

	return result
}

func mapForEach(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) < 2 {
		return NewError("too few parameters for function call: map:for-each")
	}
	if len(args) > 2 {
		return NewError("too many parameters for function call: map:for-each")
	}

	if !IsMap(args[0]) {
		return NewError("cannot match item type with required type")
	}

	m := args[0].(*object.Map)
	result := &object.Sequence{}

	switch action := args[1].(type) {
	case *object.FuncNamed:
		f := *action.Func

		if action.Num != 2 {
			return NewError("wrong number of parameters. got=#%d, expected=#2", action.Num)
		}

		for _, pair := range m.Pairs {
			result.Items = append(result.Items, f(ctx, pair.Key, pair.Value))
		}
	case *object.FuncInline:
		if len(action.PL.Params) != 2 {
			return NewError("wrong number of parameters. got=%d, expected=2", len(action.PL.Params))
		}

		enclosedCtx := object.NewEnclosedContext(ctx)
		for _, pair := range m.Pairs {
			enclosedCtx.Set(action.PL.Params[0].Value(), pair.Key)
			enclosedCtx.Set(action.PL.Params[1].Value(), pair.Value)
			a := action.Fn(action.Body, enclosedCtx)
			result.Items = append(result.Items, a)
		}
	case *object.FuncPartial:
		f := *action.Func

		if action.PCnt != 2 {
			return NewError("wrong number of placeholder. got=%d, expected=2", action.PCnt)
		}

		for _, pair := range m.Pairs {
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
					if pcnt > 0 {
						a = append(a, pair.Value)
					} else {
						a = append(a, pair.Key)
					}

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
