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
			bif.NewError("built-in function not found: %s", f.Name)
		}
		if len(args) != f.Num {
			bif.NewError("wrong number of argument. got=%d, want=%d", len(args), f.Num)
		}

		return builtin(args...)
	case *object.Array:
	case *object.Sequence:
	default:
		bif.NewError("Cannot match item type with required type")
	}
	return nil
}
