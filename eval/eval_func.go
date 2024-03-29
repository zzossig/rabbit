package eval

import (
	"github.com/zzossig/rabbit/ast"
	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/util"
)

func evalFunctionLiteral(expr ast.ExprSingle, ctx *object.Context) object.Item {
	switch expr := expr.(type) {
	case *ast.NamedFunctionRef:
		if expr.EQName.Prefix() == "" {
			expr.EQName.SetPrefix("fn")
		}

		builtin, ok := bif.F[expr.EQName.Value()]
		if !ok {
			return bif.NewError("function not found: %s", expr.EQName.Value())
		}

		return &object.FuncNamed{Name: expr.EQName, Num: expr.IntegerLiteral.Value, Func: &builtin}
	case *ast.InlineFunctionExpr:
		fi := &object.FuncInline{Body: &expr.FunctionBody, PL: &expr.ParamList}
		fi.Fn = Eval
		return fi
	}
	return bif.NewError("unexpected xpath expression. %#v", expr)
}

func evalFunctionCall(expr ast.ExprSingle, ctx *object.Context) object.Item {
	fc := expr.(*ast.FunctionCall)

	if ctxFunc, ok := ctx.Get(fc.EQName.Value()); ok {
		args := evalArgumentList(fc.Args, ctx)
		return evalDynamicFunctionCall(ctxFunc, args, ctx)
	}

	if fc.EQName.Prefix() == "" {
		fc.EQName.SetPrefix("fn")
	}

	builtin, ok := bif.F[fc.EQName.Value()]
	if !ok {
		return bif.NewError("function not found: %s", fc.EQName.Value())
	}

	pcnt := 0
	args := evalArgumentList(fc.Args, ctx)

	for _, arg := range args {
		if _, ok := arg.(*object.Placeholder); ok {
			pcnt++
		}
	}

	if pcnt > 0 {
		fp := &object.FuncPartial{}
		fp.Func = &builtin
		fp.Name = fc.EQName
		fp.Args = args
		fp.PCnt = pcnt
		fp.Context = ctx

		return fp
	}

	return builtin(ctx, args...)
}

func evalVarRef(expr ast.ExprSingle, ctx *object.Context) object.Item {
	vr := expr.(*ast.VarRef)

	if v, ok := ctx.Get(vr.VarName.Value()); ok {
		var nodes []object.Node

		if bif.IsNodeSeq(v) {
			v := v.(*object.Sequence)

			for _, item := range v.Items {
				if item, ok := item.(*object.BaseNode); ok {
					nodes = append(nodes, item)
				}
				if item, ok := item.(*object.AttrNode); ok {
					nodes = append(nodes, item)
				}
			}
		} else if bif.IsNode(v) {
			if item, ok := v.(*object.BaseNode); ok {
				nodes = append(nodes, item)
			}
			if item, ok := v.(*object.AttrNode); ok {
				nodes = append(nodes, item)
			}
		}

		if len(nodes) > 0 {
			ctx.CNode = nodes
			ctx.CSize = len(nodes)
			ctx.CAxis = "child::"
		}

		return v
	}

	return &object.Varref{Name: vr.VarName}
}

func evalArgument(arg ast.Argument, ctx *object.Context) object.Item {
	switch arg.TypeID {
	case 1:
		return Eval(arg.ExprSingle, ctx)
	case 2:
		return &object.Placeholder{}
	default:
		return bif.NewError("unexpected argument expression. %#v", arg)
	}
}

func evalArgumentList(args []ast.Argument, ctx *object.Context) []object.Item {
	var items []object.Item

	for _, arg := range args {
		item := evalArgument(arg, ctx)
		items = append(items, item)
	}

	return items
}

func evalPredicate(it object.Item, pred *ast.Predicate, ctx *object.Context) object.Item {
	var src []object.Item

	switch it := it.(type) {
	case *object.Sequence:
		if util.IsReverseAxis(ctx.CAxis) {
			for i := len(it.Items) - 1; i >= 0; i-- {
				src = append(src, it.Items[i])
			}
		} else {
			src = append(src, it.Items...)
		}
	default:
		src = append(src, it)
	}
	ctx.CSize = len(src)

	var items []object.Item
	for i, s := range src {
		ctx.CItem = s
		ctx.CPos = i + 1

		evaled := Eval(&pred.Expr, ctx).(*object.Sequence)
		if len(evaled.Items) != 1 {
			return bif.NewError("wrong number of argument. got=%d, want=1", len(evaled.Items))
		}
		if bif.IsError(evaled.Items[0]) {
			return evaled
		}

		switch ev := evaled.Items[0].(type) {
		case *object.Integer:
			if ev.Value()-1 == i {
				items = append(items, s)
			}
		case *object.Decimal:
			if ev.Value()-1 == float64(i) {
				items = append(items, s)
			}
		case *object.Double:
			if ev.Value()-1 == float64(i) {
				items = append(items, s)
			}
		case *object.String:
			builtin := bif.F["fn:boolean"]
			bl := builtin(nil, ev)

			boolObj := bl.(*object.Boolean)
			if boolObj.Value() {
				items = append(items, s)
			}
		case *object.Boolean:
			if ev.Value() {
				items = append(items, s)
			}
		}
	}

	return &object.Sequence{Items: items}
}

func evalLookup(it object.Item, lu *ast.Lookup, ctx *object.Context) object.Item {
	seq := &object.Sequence{}

	switch it := it.(type) {
	case *object.Array:
		switch lu.KeySpecifier.TypeID {
		case 1:
			return bif.NewError("cannot convert xs:string to xs:integer: %s.", lu.NCName.Value())
		case 2:
			if lu.IntegerLiteral.Value == 0 || lu.IntegerLiteral.Value > len(it.Items) {
				return bif.NewError("array index %d out of bounds (1..%d)", lu.IntegerLiteral.Value, len(it.Items))
			}
			return it.Items[lu.IntegerLiteral.Value-1]
		case 3:
			evaled := Eval(&lu.ParenthesizedExpr, ctx)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				if i, ok := item.(*object.Integer); ok {
					if i.Value() == 0 || i.Value() > len(it.Items) {
						return bif.NewError("array index %d out of bounds (1..%d)", i.Value(), len(it.Items))
					}
					seq.Items = append(seq.Items, it.Items[i.Value()-1])
				}
			}
		case 4:
			seq.Items = append(seq.Items, it.Items...)
		}
	case *object.Map:
		switch lu.KeySpecifier.TypeID {
		case 1:
			key := bif.NewString(lu.NCName.Value())
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 2:
			key := bif.NewInteger(lu.IntegerLiteral.Value)
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 3:
			evaled := Eval(&lu.ParenthesizedExpr, ctx)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				if key, ok := item.(object.Hasher); ok {
					if pair, ok := it.Pairs[key.HashKey()]; ok {
						seq.Items = append(seq.Items, pair.Value)
					}
				}
			}
		case 4:
			for _, pair := range it.Pairs {
				seq.Items = append(seq.Items, pair.Value)
			}
		}
	case *object.Sequence:
		for _, item := range it.Items {
			evaled := evalLookup(item, lu, ctx)
			seq.Items = append(seq.Items, evaled)
		}
	default:
		return bif.NewError("input of lookup operator is not a map or array: %v.", it)
	}

	return seq
}

func evalArrowExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ae := expr.(*ast.ArrowExpr)
	bindings := ae.Bindings

	e := Eval(ae.ExprSingle, ctx)
	args := []object.Item{e}
	var result object.Item

	for i, b := range bindings {
		switch b.TypeID {
		case 1:
			if b.EQName.Prefix() == "" {
				b.EQName.SetPrefix("fn")
			}
			builtin, ok := bif.F[b.EQName.Value()]
			if !ok {
				bif.NewError("function not defined: %s", b.EQName.Value())
			}

			evaled := evalArgumentList(b.Args, ctx)
			args = append(args, evaled...)
			result = builtin(ctx, args...)
			if i < len(bindings)-1 {
				args = []object.Item{result}
			}
		case 2:
			ctxItem, ok := ctx.Get(b.VarName.Value())
			if !ok {
				bif.NewError("function not defined: %s", b.VarRef.String())
			}
			ctxFunc, ok := ctxItem.(*object.FuncInline)
			if !ok {
				bif.NewError("function not defined: %s", b.VarRef.String())
			}

			evaled := evalArgumentList(b.Args, ctx)
			args = append(args, evaled...)
			if len(ctxFunc.PL.Params) != len(args) {
				return bif.NewError("wrong number of argument. got=%d, want=%d", len(args), len(ctxFunc.PL.Params))
			}

			enclosedCtx := object.NewEnclosedContext(ctx)
			for i, param := range ctxFunc.PL.Params {
				enclosedCtx.Set(param.Value(), args[i])
			}

			result = Eval(&ctxFunc.Body.Expr, enclosedCtx)
			if i < len(bindings)-1 {
				args = []object.Item{result}
			}
		case 3:
			result = Eval(&b.ParenthesizedExpr, ctx)
			if i < len(bindings)-1 {
				args = []object.Item{result}
			}
		}
	}

	return result
}

func evalPostfixExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	pe := expr.(*ast.PostfixExpr)
	evaled := Eval(pe.ExprSingle, ctx)

	for _, pal := range pe.Pals {
		switch pal := pal.(type) {
		case *ast.Predicate:
			evaled = evalPredicate(evaled, pal, ctx)
		case *ast.ArgumentList:
			args := evalArgumentList(pal.Args, ctx)
			evaled = evalDynamicFunctionCall(evaled, args, ctx)
		case *ast.Lookup:
			evaled = evalLookup(evaled, pal, ctx)
		}
	}

	return evaled
}

func evalDynamicFunctionCall(f object.Item, args []object.Item, ctx *object.Context) object.Item {
	switch f := f.(type) {
	case *object.FuncInline:
		if len(f.PL.Params) != len(args) {
			return bif.NewError("wrong number of argument. got=%d, want=%d", len(args), len(f.PL.Params))
		}

		enclosedCtx := object.NewEnclosedContext(ctx)
		for i, param := range f.PL.Params {
			enclosedCtx.Set(param.EQName.Value(), args[i])
		}

		return Eval(&f.Body.Expr, enclosedCtx)
	case *object.FuncNamed:
		if f.Name.Prefix() == "" {
			f.Name.SetPrefix("fn")
		}
		builtin, ok := bif.F[f.Name.Value()]
		if !ok {
			return bif.NewError("function not found: %s", f.Name.Value())
		}
		if len(args) != f.Num {
			return bif.NewError("wrong number of argument. got=%d, want=%d", len(args), f.Num)
		}

		return builtin(ctx, args...)
	case *object.Array:
		if len(args) != 1 {
			return bif.NewError("wrong number of argument. got=%d, want=1", len(args))
		}

		index, ok := args[0].(*object.Integer)
		if !ok {
			return bif.NewError("dynamic function call on array should have integer argument")
		}
		if index.Value() == 0 || index.Value() > len(f.Items) {
			return bif.NewError("Index out of range: size(%d)", len(f.Items))
		}
		return f.Items[index.Value()-1]
	case *object.Map:
		if len(args) != 1 {
			return bif.NewError("wrong number of argument. got=%d, want=1", len(args))
		}

		h, ok := args[0].(object.Hasher)
		if !ok {
			return bif.NewError("dynamic function call on map should have atomic argument")
		}

		key := h.HashKey()
		pair, ok := f.Pairs[key]
		if !ok {
			return bif.NewSequence()
		}
		return pair.Value
	default:
		bif.NewError("cannot match item type with required type")
	}
	return nil
}

func evalSimpleMapExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	sme := expr.(*ast.SimpleMapExpr)
	left := Eval(sme.LeftExpr, ctx)

	var cItems []object.Item
	switch left := left.(type) {
	case *object.Sequence:
		cItems = append(cItems, left.Items...)
	default:
		cItems = append(cItems, left)
	}

	var items []object.Item
	for _, c := range cItems {
		ctx.CItem = c
		right := Eval(sme.RightExpr, ctx)
		items = append(items, right)
	}

	return &object.Sequence{Items: items}
}

func evalArrayExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	array := &object.Array{}
	var exprs []ast.ExprSingle

	switch expr := expr.(type) {
	case *ast.SquareArrayConstructor:
		exprs = expr.Exprs
		for _, e := range exprs {
			item := Eval(e, ctx)
			array.Items = append(array.Items, item)
		}
	case *ast.CurlyArrayConstructor:
		exprs = expr.EnclosedExpr.Exprs
		for _, e := range exprs {
			item := Eval(e, ctx)
			if bif.IsSeq(item) {
				array.Items = append(array.Items, bif.UnwrapSeq(item)...)
			} else {
				array.Items = append(array.Items, item)
			}
		}
	}

	return array
}

func evalMapExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	mc := expr.(*ast.MapConstructor)
	pairs := make(map[object.HashKey]object.Pair)

	for _, entry := range mc.Entries {
		key := Eval(entry.MapKeyExpr.ExprSingle, ctx)

		hashKey, ok := key.(object.Hasher)
		if !ok {
			return bif.NewError("unusable as hash key: %s", key.Type())
		}

		value := Eval(entry.MapValueExpr.ExprSingle, ctx)

		hashed := hashKey.HashKey()
		pairs[hashed] = object.Pair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

func evalUnaryLookup(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ul := expr.(*ast.UnaryLookup)
	seq := &object.Sequence{}

	switch it := ctx.CItem.(type) {
	case *object.Array:
		switch ul.KeySpecifier.TypeID {
		case 1:
			return bif.NewError("NCName not supported in unary lookup")
		case 2:
			if ul.IntegerLiteral.Value == 0 || ul.IntegerLiteral.Value > len(it.Items) {
				return bif.NewError("array index %d out of bounds (1..%d)", ul.IntegerLiteral.Value, len(it.Items))
			}
			return it.Items[ul.IntegerLiteral.Value-1]
		case 3:
			evaled := Eval(&ul.ParenthesizedExpr, ctx)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				i, ok := item.(*object.Integer)
				if !ok {
					return bif.NewError("cannot convert %s to xs:integer", i.Type())
				}
				if i.Value() == 0 || i.Value() > len(it.Items) {
					return bif.NewError("array index %d out of bounds (1..%d)", ul.IntegerLiteral.Value, len(it.Items))
				}
				seq.Items = append(seq.Items, it.Items[i.Value()-1])
			}
		case 4:
			seq.Items = append(seq.Items, it.Items...)
		}
	case *object.Map:
		switch ul.KeySpecifier.TypeID {
		case 1:
			key := bif.NewString(ul.NCName.Value())
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 2:
			key := bif.NewInteger(ul.IntegerLiteral.Value)
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 3:
			evaled := Eval(&ul.ParenthesizedExpr, ctx)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				if key, ok := item.(object.Hasher); ok {
					if pair, ok := it.Pairs[key.HashKey()]; ok {
						seq.Items = append(seq.Items, pair.Value)
					}
				}
			}
		case 4:
			for _, pair := range it.Pairs {
				seq.Items = append(seq.Items, pair.Value)
			}
		}
	default:
		return bif.NewError("context item is not a map or an array")
	}

	return seq
}
