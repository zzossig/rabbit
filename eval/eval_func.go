package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
)

func evalFunctionLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.NamedFunctionRef:
		return &object.FuncNamed{Name: expr.EQName.Value(), Num: expr.IntegerLiteral.Value}
	case *ast.InlineFunctionExpr:
		return &object.FuncInline{Body: &expr.FunctionBody, PL: &expr.ParamList, SType: &expr.SequenceType}
	}
	return nil
}

func evalFunctionCall(expr ast.ExprSingle, env *object.Env) object.Item {
	f := expr.(*ast.FunctionCall)

	builtin, ok := bif.Builtins[f.EQName.Value()]
	if !ok {
		envFunc, ok := env.Get(f.EQName.Value())
		if !ok {
			return bif.NewError("function not found: " + f.EQName.Value())
		}

		args, _ := evalArgumentList(f.Args, env)
		return evalDynamicFunctionCall(envFunc, args, env)
	}

	enclosedEnv := object.NewEnclosedEnv(env)
	fc := &object.FuncCall{}
	fc.Env = enclosedEnv
	fc.Name = f.EQName.Value()
	fc.Func = &builtin

	args, pcnt := evalArgumentList(f.Args, env)
	enclosedEnv.Args = append(enclosedEnv.Args, args...)

	if pcnt > 0 {
		return fc
	}

	return builtin(args...)
}

func evalVarRef(expr ast.ExprSingle, env *object.Env) object.Item {
	vr := expr.(*ast.VarRef)

	if v, ok := env.Get(vr.VarName.Value()); ok {
		return v
	}
	return bif.NewError("Undefined variable %s", vr.VarName.Value())
}

func evalArgument(arg ast.Argument, env *object.Env) object.Item {
	switch arg.TypeID {
	case 0:
		return NIL
	case 1:
		return Eval(arg.ExprSingle, env)
	case 2:
		return &object.Placeholder{}
	default:
		return NIL
	}
}

func evalArgumentList(args []ast.Argument, env *object.Env) ([]object.Item, int) {
	var items []object.Item
	pcnt := 0

	for _, arg := range args {
		item := evalArgument(arg, env)
		items = append(items, item)
		if item.Type() == object.PholderType {
			pcnt++
		}
	}

	return items, pcnt
}

func evalPredicate(it object.Item, pred *ast.Predicate, env *object.Env) object.Item {
	var src []object.Item

	switch it := it.(type) {
	case *object.Sequence:
		for _, item := range it.Items {
			src = append(src, item)
		}
	default:
		src = append(src, it)
	}

	var items []object.Item
	for i, s := range src {
		env.CItem = s

		evaled := Eval(&pred.Expr, env).(*object.Sequence)
		if len(evaled.Items) != 1 {
			return bif.NewError("Wrong number of argument. got=%d, want=1", len(evaled.Items))
		}

		switch ev := evaled.Items[0].(type) {
		case *object.Integer:
			if ev.Value-1 == i {
				items = append(items, s)
			}
		case *object.Decimal:
			if ev.Value-1 == float64(i) {
				items = append(items, s)
			}
		case *object.Double:
			if ev.Value-1 == float64(i) {
				items = append(items, s)
			}
		case *object.String:
			builtin := bif.Builtins["boolean"]
			bl := builtin(ev)
			if bif.IsError(bl) {
				return bl
			}

			boolObj := bl.(*object.Boolean)
			if boolObj.Value {
				items = append(items, s)
			}
		case *object.Boolean:
			if ev.Value {
				items = append(items, s)
			}
		case *object.Nil:
			break
		}
	}

	return &object.Sequence{Items: items}
}

func evalLookup(it object.Item, lu *ast.Lookup, env *object.Env) object.Item {
	seq := &object.Sequence{}

	switch it := it.(type) {
	case *object.Array:
		switch lu.KeySpecifier.TypeID {
		case 1:
			return bif.NewError("[XPTY0004] Cannot convert xs:string to xs:integer: %s.", lu.NCName.Value())
		case 2:
			if lu.IntegerLiteral.Value == 0 || lu.IntegerLiteral.Value > len(it.Items) {
				return bif.NewError("[FOAY0001] Array index %d out of bounds (1..%d)", lu.IntegerLiteral.Value, len(it.Items))
			}
			return it.Items[lu.IntegerLiteral.Value-1]
		case 3:
			evaled := Eval(&lu.ParenthesizedExpr, env)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				if i, ok := item.(*object.Integer); ok {
					if i.Value == 0 || i.Value > len(it.Items) {
						return bif.NewError("[FOAY0001] Array index %d out of bounds (1..%d)", i.Value, len(it.Items))
					}
					seq.Items = append(seq.Items, it.Items[i.Value-1])
				}
			}
		case 4:
			for _, item := range it.Items {
				seq.Items = append(seq.Items, item)
			}
		}
	case *object.Map:
		switch lu.KeySpecifier.TypeID {
		case 1:
			key := object.String{Value: lu.NCName.Value()}
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 2:
			key := object.Integer{Value: lu.IntegerLiteral.Value}
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 3:
			evaled := Eval(&lu.ParenthesizedExpr, env)
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
			evaled := evalLookup(item, lu, env)
			seq.Items = append(seq.Items, evaled)
		}
	default:
		return bif.NewError("[XPTY0004] Input of lookup operator is not a map or array: %v.", it)
	}

	return seq
}

func evalArrowExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	ae := expr.(*ast.ArrowExpr)
	bindings := ae.Bindings

	e := Eval(ae.ExprSingle, env)
	args := []object.Item{e}
	var result object.Item

	for i, b := range bindings {
		switch b.TypeID {
		case 1:
			builtin, ok := bif.Builtins[b.EQName.Value()]
			if !ok {
				bif.NewError("function not defined: %s", b.EQName.Value())
			}

			evaled, _ := evalArgumentList(b.Args, env)
			args = append(args, evaled...)
			result = builtin(args...)
			if i < len(bindings)-1 {
				args = []object.Item{result}
			}
		case 2:
			// TODO VarRef
		case 3:
			// TODO ParenthesizedExpr
		}
	}

	return result
}

func evalPostfixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	pe := expr.(*ast.PostfixExpr)
	evaled := Eval(pe.ExprSingle, env)

	for _, pal := range pe.Pals {
		switch pal := pal.(type) {
		case *ast.Predicate:
			evaled = evalPredicate(evaled, pal, env)
		case *ast.ArgumentList:
			args, _ := evalArgumentList(pal.Args, env)
			evaled = evalDynamicFunctionCall(evaled, args, env)
		case *ast.Lookup:
			evaled = evalLookup(evaled, pal, env)
		}
	}

	return evaled
}

func evalDynamicFunctionCall(f object.Item, args []object.Item, env *object.Env) object.Item {
	switch f := f.(type) {
	case *object.FuncInline:
		if len(f.PL.Params) != len(args) {
			return bif.NewError("wrong number of argument. got=%d, want=%d", len(args), len(f.PL.Params))
		}
		for i, param := range f.PL.Params {
			env.Set(param.EQName.Value(), args[i])
		}

		return Eval(&f.Body.Expr, env)
	case *object.FuncNamed:
		builtin, ok := bif.Builtins[f.Name]
		if !ok {
			return bif.NewError("built-in function not found: %s", f.Name)
		}
		if len(args) != f.Num {
			return bif.NewError("wrong number of argument. got=%d, want=%d", len(args), f.Num)
		}

		return builtin(args...)
	case *object.Array:
		if len(args) != 1 {
			return bif.NewError("wrong number of argument. got=%d, want=1", len(args))
		}

		index, ok := args[0].(*object.Integer)
		if !ok {
			return bif.NewError("dynamic function call on array should have integer argument")
		}
		if index.Value == 0 || index.Value > len(f.Items) {
			return bif.NewError("Index out of range: size(%d)", len(f.Items))
		}
		return f.Items[index.Value-1]
	default:
		bif.NewError("Cannot match item type with required type")
	}
	return nil
}

func evalSimpleMapExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	sme := expr.(*ast.SimpleMapExpr)
	left := Eval(sme.LeftExpr, env)

	var cItems []object.Item
	switch left := left.(type) {
	case *object.Sequence:
		for _, item := range left.Items {
			cItems = append(cItems, item)
		}
	default:
		cItems = append(cItems, left)
	}

	var items []object.Item
	for _, c := range cItems {
		env.CItem = c
		right := Eval(sme.RightExpr, env)
		items = append(items, right)
	}

	return &object.Sequence{Items: items}
}

func evalArrayExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	array := &object.Array{}
	var exprs []ast.ExprSingle

	switch expr := expr.(type) {
	case *ast.SquareArrayConstructor:
		exprs = expr.Exprs
	case *ast.CurlyArrayConstructor:
		exprs = expr.EnclosedExpr.Exprs
	}

	for _, e := range exprs {
		item := Eval(e, env)
		array.Items = append(array.Items, item)
	}

	return array
}

func evalMapExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	mc := expr.(*ast.MapConstructor)
	pairs := make(map[object.HashKey]object.Pair)

	for _, entry := range mc.Entries {
		key := Eval(entry.MapKeyExpr.ExprSingle, env)

		hashKey, ok := key.(object.Hasher)
		if !ok {
			return bif.NewError("unusable as hash key: %s", key.Type())
		}

		value := Eval(entry.MapValueExpr.ExprSingle, env)

		hashed := hashKey.HashKey()
		pairs[hashed] = object.Pair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

// KeySpecifier ::= NCName | IntegerLiteral | ParenthesizedExpr | "*"
// TypeID ::=				1			 | 2							| 3									| 4
func evalUnaryLookup(expr ast.ExprSingle, env *object.Env) object.Item {
	ul := expr.(*ast.UnaryLookup)
	seq := &object.Sequence{}

	switch it := env.CItem.(type) {
	case *object.Array:
		switch ul.KeySpecifier.TypeID {
		case 1:
			return bif.NewError("[err:XPTY0004] NCName not supported in unary lookup")
		case 2:
			if ul.IntegerLiteral.Value == 0 || ul.IntegerLiteral.Value > len(it.Items) {
				return bif.NewError("[FOAY0001] Array index %d out of bounds (1..%d)", ul.IntegerLiteral.Value, len(it.Items))
			}
			return it.Items[ul.IntegerLiteral.Value-1]
		case 3:
			evaled := Eval(&ul.ParenthesizedExpr, env)
			src := evaled.(*object.Sequence)

			for _, item := range src.Items {
				i, ok := item.(*object.Integer)
				if !ok {
					return bif.NewError("[XPTY0004] Cannot convert %s to xs:integer", i.Type())
				}
				if i.Value == 0 || i.Value > len(it.Items) {
					return bif.NewError("[FOAY0001] Array index %d out of bounds (1..%d)", ul.IntegerLiteral.Value, len(it.Items))
				}
				seq.Items = append(seq.Items, it.Items[i.Value-1])
			}
		case 4:
			for _, item := range it.Items {
				seq.Items = append(seq.Items, item)
			}
		}
	case *object.Map:
		switch ul.KeySpecifier.TypeID {
		case 1:
			key := object.String{Value: ul.NCName.Value()}
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 2:
			key := object.Integer{Value: ul.IntegerLiteral.Value}
			pair, ok := it.Pairs[key.HashKey()]
			if !ok {
				return seq
			}
			return pair.Value
		case 3:
			evaled := Eval(&ul.ParenthesizedExpr, env)
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
		return bif.NewError("[err:XPTY0004] context item is not a map or an array")
	}

	return seq
}
