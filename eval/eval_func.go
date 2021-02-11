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
		return bif.NewError("function not found: " + f.EQName.Value())
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
		return &object.Nil{}
	case 1:
		return Eval(arg.ExprSingle, env)
	case 2:
		return &object.Placeholder{}
	default:
		return &object.Nil{}
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
			bl := builtin(ev).(*object.Boolean)
			if bl.Value {
				return &object.Sequence{Items: src}
			}
			return &object.Sequence{}
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
