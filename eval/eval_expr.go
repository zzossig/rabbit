package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalXPath(expr *ast.XPath, ctx *object.Context) object.Item {
	xpath := &object.Sequence{}

	for _, e := range expr.Exprs {
		item := Eval(e, ctx)

		switch item := item.(type) {
		case *object.Sequence:
			for _, it := range item.Items {
				if bif.IsSeq(it) {
					xpath.Items = append(xpath.Items, bif.UnwrapSeq(it)...)
				} else {
					xpath.Items = append(xpath.Items, it)
				}
			}
		default:
			xpath.Items = append(xpath.Items, item)
		}
	}

	return xpath
}

func evalExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	switch expr := expr.(type) {
	case *ast.Expr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.ParenthesizedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.EnclosedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.Predicate:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, ctx)
			seq.Items = append(seq.Items, item)
		}
		return seq
	}
	return object.NIL
}

func evalIntegerLiteral(expr ast.ExprSingle, ctx *object.Context) object.Item {
	il := expr.(*ast.IntegerLiteral)
	return bif.NewInteger(il.Value)
}

func evalDecimalLiteral(expr ast.ExprSingle, ctx *object.Context) object.Item {
	dl := expr.(*ast.DecimalLiteral)
	return bif.NewDecimal(dl.Value)
}

func evalDoubleLiteral(expr ast.ExprSingle, ctx *object.Context) object.Item {
	dl := expr.(*ast.DoubleLiteral)
	return bif.NewDouble(dl.Value)
}

func evalStringLiteral(expr ast.ExprSingle, ctx *object.Context) object.Item {
	sl := expr.(*ast.StringLiteral)
	return bif.NewString(sl.Value)
}

func evalUnaryExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ue := expr.(*ast.UnaryExpr)

	right := Eval(ue.ExprSingle, ctx)
	op := ue.Token

	var funcName string
	if op.Type == token.UPLUS {
		funcName = "op:numeric-unary-plus"
	} else {
		funcName = "op:numeric-unary-minus"
	}

	builtin := bif.F[funcName]

	return builtin(nil, right)
}

func evalIfExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ie := expr.(*ast.IfExpr)
	builtin := bif.F["fn:boolean"]

	testE := Eval(ie.TestExpr, ctx)
	bl := builtin(nil, testE)
	boolObj := bl.(*object.Boolean)

	if boolObj.Value() {
		return Eval(ie.ThenExpr, ctx)
	}
	return Eval(ie.ElseExpr, ctx)
}

func evalForExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	fe := expr.(*ast.ForExpr)
	enclosedCtx := object.NewEnclosedContext(ctx)
	var items []object.Item

	if len(fe.Bindings) > 1 {
		b := fe.Bindings[0]
		bval := Eval(b.ExprSingle, ctx)

		nfe := &ast.ForExpr{ExprSingle: fe.ExprSingle}
		nfe.Bindings = fe.Bindings[1:]

		switch bval := bval.(type) {
		case *object.Sequence:
			for _, item := range bval.Items {
				enclosedCtx.Set(b.VarName.Value(), item)
				e := evalForExpr(nfe, enclosedCtx).(*object.Sequence)
				items = append(items, e.Items...)
			}
		default:
			enclosedCtx.Set(b.VarName.Value(), bval)
			e := evalForExpr(nfe, enclosedCtx).(*object.Sequence)
			items = append(items, e.Items...)
		}

		return &object.Sequence{Items: items}
	}

	b := fe.Bindings[0]
	bval := Eval(b.ExprSingle, enclosedCtx)

	switch bval := bval.(type) {
	case *object.Sequence:
		for _, item := range bval.Items {
			enclosedCtx.Set(b.VarName.Value(), item)
			e := Eval(fe.ExprSingle, enclosedCtx)
			items = append(items, e)
		}
	default:
		enclosedCtx.Set(b.VarName.Value(), bval)
		e := Eval(fe.ExprSingle, enclosedCtx)
		items = append(items, e)
	}

	return &object.Sequence{Items: items}
}

func evalLetExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	le := expr.(*ast.LetExpr)
	enclosedCtx := object.NewEnclosedContext(ctx)

	for _, b := range le.Bindings {
		bval := Eval(b.ExprSingle, enclosedCtx)
		enclosedCtx.Set(b.VarName.Value(), bval)
	}

	return Eval(le.ExprSingle, enclosedCtx)
}

func evalQuantifiedExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	qe := expr.(*ast.QuantifiedExpr)
	enclosedCtx := object.NewEnclosedContext(ctx)

	if len(qe.Bindings) > 1 {
		b := qe.Bindings[0]
		bval := Eval(b.ExprSingle, enclosedCtx)

		nqe := &ast.QuantifiedExpr{ExprSingle: qe.ExprSingle, Token: qe.Token}
		nqe.Bindings = qe.Bindings[1:]

		switch bval := bval.(type) {
		case *object.Sequence:
			for _, item := range bval.Items {
				enclosedCtx.Set(b.VarName.Value(), item)
				e := evalQuantifiedExpr(nqe, enclosedCtx).(*object.Boolean)

				if qe.Token.Type == token.EVERY && !e.Value() {
					return bif.NewBoolean(false)
				}
				if qe.Token.Type == token.SOME && e.Value() {
					return bif.NewBoolean(true)
				}
			}
		default:
			enclosedCtx.Set(b.VarName.Value(), bval)
			e := evalQuantifiedExpr(nqe, enclosedCtx).(*object.Boolean)

			if qe.Token.Type == token.EVERY && !e.Value() {
				return bif.NewBoolean(false)
			}
			if qe.Token.Type == token.SOME && e.Value() {
				return bif.NewBoolean(true)
			}
		}
	}

	b := qe.Bindings[0]
	bval := Eval(b.ExprSingle, enclosedCtx)

	switch bval := bval.(type) {
	case *object.Sequence:
		for _, item := range bval.Items {
			enclosedCtx.Set(b.VarName.Value(), item)
			e, ok := Eval(qe.ExprSingle, enclosedCtx).(*object.Boolean)

			if !ok {
				builtin := bif.F["fn:boolean"]
				bl := builtin(nil, e)

				boolObj := bl.(*object.Boolean)
				if qe.Token.Type == token.EVERY && !boolObj.Value() {
					return bif.NewBoolean(false)
				}
				if qe.Token.Type == token.SOME && boolObj.Value() {
					return bif.NewBoolean(true)
				}
			}

			if qe.Token.Type == token.EVERY && !e.Value() {
				return bif.NewBoolean(false)
			}
			if qe.Token.Type == token.SOME && e.Value() {
				return bif.NewBoolean(true)
			}
		}
	default:
		enclosedCtx.Set(b.VarName.Value(), bval)
		e, ok := Eval(qe.ExprSingle, enclosedCtx).(*object.Boolean)

		if !ok {
			builtin := bif.F["fn:boolean"]
			bl := builtin(nil, e)

			boolObj := bl.(*object.Boolean)
			if qe.Token.Type == token.EVERY && !boolObj.Value() {
				return bif.NewBoolean(false)
			}
			if qe.Token.Type == token.SOME && boolObj.Value() {
				return bif.NewBoolean(true)
			}
		}

		if qe.Token.Type == token.EVERY && !e.Value() {
			return bif.NewBoolean(false)
		}
		if qe.Token.Type == token.SOME && e.Value() {
			return bif.NewBoolean(true)
		}
	}

	if qe.Token.Type == token.EVERY {
		return bif.NewBoolean(true)
	}
	return bif.NewBoolean(false)
}

func evalAdditiveExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ae := expr.(*ast.AdditiveExpr)

	left := Eval(ae.LeftExpr, ctx)
	right := Eval(ae.RightExpr, ctx)
	op := ae.Token

	var funcName string
	if op.Type == token.PLUS {
		funcName = "op:numeric-add"
	} else {
		funcName = "op:numeric-subtract"
	}

	builtin := bif.F[funcName]

	return builtin(nil, left, right)
}

func evalMultiplicativeExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	me := expr.(*ast.MultiplicativeExpr)

	left := Eval(me.LeftExpr, ctx)
	right := Eval(me.RightExpr, ctx)
	op := me.Token

	var funcName string
	if op.Type == token.ASTERISK {
		funcName = "op:numeric-multiply"
	} else if op.Type == token.DIV {
		funcName = "op:numeric-divide"
	} else if op.Type == token.IDIV {
		funcName = "op:numeric-integer-divide"
	} else {
		funcName = "op:numeric-mod"
	}

	builtin := bif.F[funcName]

	return builtin(nil, left, right)
}

func evalStringConcatExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	sce := expr.(*ast.StringConcatExpr)

	left := Eval(sce.LeftExpr, ctx)
	right := Eval(sce.RightExpr, ctx)

	funcName := "fn:concat"
	builtin := bif.F[funcName]

	return builtin(nil, left, right)
}

func evalRangeExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	re := expr.(*ast.RangeExpr)

	l := Eval(re.LeftExpr, ctx)
	r := Eval(re.RightExpr, ctx)

	left, ok := l.(*object.Integer)
	if !ok {
		return bif.NewError("not allowed type in RangeExpr: %s", left.Type())
	}

	right, ok := r.(*object.Integer)
	if !ok {
		return bif.NewError("not allowed type in RangeExpr: %s", right.Type())
	}

	seq := &object.Sequence{}
	for i := left.Value(); i <= right.Value(); i++ {
		seq.Items = append(seq.Items, bif.NewInteger(i))
	}
	return seq
}

func evalLogicalExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	var left object.Item
	var right object.Item
	var op token.Token

	builtin := bif.F["fn:boolean"]

	switch expr := expr.(type) {
	case *ast.AndExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	case *ast.OrExpr:
		left = Eval(expr.LeftExpr, ctx)
		right = Eval(expr.RightExpr, ctx)
		op = expr.Token
	}

	l := builtin(nil, left)
	r := builtin(nil, right)

	leftBool := l.(*object.Boolean)
	rightBool := r.(*object.Boolean)

	switch op.Type {
	case token.AND:
		return bif.NewBoolean(leftBool.Value() && rightBool.Value())
	case token.OR:
		return bif.NewBoolean(leftBool.Value() || rightBool.Value())
	default:
		return object.NIL
	}
}

func evalUnionExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ue := expr.(*ast.UnionExpr)
	cnode := ctx.CNode

	left := Eval(ue.LeftExpr, ctx)
	ctx.CNode = cnode
	right := Eval(ue.RightExpr, ctx)

	var nodes []object.Node

	switch {
	case bif.IsSeq(left) && bif.IsSeq(right):
		lseq := left.(*object.Sequence)
		rseq := right.(*object.Sequence)

		for _, item := range lseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in UnionExpr: %s", item.Type())
		}

		for _, item := range rseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in UnionExpr: %s", item.Type())
		}
	case bif.IsSeq(left) && bif.IsNode(right):
		lseq := left.(*object.Sequence)

		for _, item := range lseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in UnionExpr: %s", item.Type())
		}

		switch right := right.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, right)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, right)
		}

	case bif.IsNode(left) && bif.IsSeq(right):
		rseq := right.(*object.Sequence)

		switch left := left.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, left)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, left)
		}

		for _, item := range rseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in UnionExpr: %s", item.Type())
		}
	case bif.IsNode(left) && bif.IsNode(right):
		switch left := left.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, left)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, left)
		}

		switch right := right.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, right)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, right)
		}
	default:
		return bif.NewError("not allowed types in UnionExpr: %s, %s", left.Type(), right.Type())
	}

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func evalIntersectExceptExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	iee := expr.(*ast.IntersectExceptExpr)

	left := Eval(iee.LeftExpr, ctx)
	right := Eval(iee.RightExpr, ctx)

	var nodes []object.Node
	var inodes []object.Node
	var enodes []object.Node

	switch {
	case bif.IsSeq(left) && bif.IsSeq(right):
		lseq := left.(*object.Sequence)
		rseq := right.(*object.Sequence)

		for _, item := range lseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in IntersectExceptExpr: %s", item.Type())
		}

		for _, item := range rseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				if bif.IsContainN(nodes, item) {
					inodes = append(inodes, item)
				}
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				if bif.IsContainN(nodes, item) {
					inodes = append(inodes, item)
				}
				continue
			}
			return bif.NewError("not allowed type in IntersectExceptExpr: %s", item.Type())
		}
	case bif.IsSeq(left) && bif.IsNode(right):
		lseq := left.(*object.Sequence)

		for _, item := range lseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				nodes = bif.AppendNode(nodes, item)
				continue
			}
			return bif.NewError("not allowed type in IntersectExceptExpr: %s", item.Type())
		}

		switch right := right.(type) {
		case *object.BaseNode:
			if bif.IsContainN(nodes, right) {
				inodes = append(inodes, right)
			}
		case *object.AttrNode:
			if bif.IsContainN(nodes, right) {
				inodes = append(inodes, right)
			}
		}
	case bif.IsNode(left) && bif.IsSeq(right):
		rseq := right.(*object.Sequence)

		switch left := left.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, left)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, left)
		}

		for _, item := range rseq.Items {
			if item, ok := item.(*object.BaseNode); ok {
				if bif.IsContainN(nodes, item) {
					inodes = append(inodes, item)
				}
				continue
			}
			if item, ok := item.(*object.AttrNode); ok {
				if bif.IsContainN(nodes, item) {
					inodes = append(inodes, item)
				}
				continue
			}
			return bif.NewError("not allowed type in IntersectExceptExpr: %s", item.Type())
		}
	case bif.IsNode(left) && bif.IsNode(right):
		switch left := left.(type) {
		case *object.BaseNode:
			nodes = bif.AppendNode(nodes, left)
		case *object.AttrNode:
			nodes = bif.AppendNode(nodes, left)
		}

		switch right := right.(type) {
		case *object.BaseNode:
			if bif.IsContainN(nodes, right) {
				inodes = append(inodes, right)
			}
		case *object.AttrNode:
			if bif.IsContainN(nodes, right) {
				inodes = append(inodes, right)
			}
		}
	default:
		return bif.NewError("not allowed types in IntersectExceptExpr: %s, %s", left.Type(), right.Type())
	}

	seq := &object.Sequence{}

	if iee.Token.Type == token.INTERSECT {
		for _, node := range inodes {
			seq.Items = append(seq.Items, node)
		}
		return seq
	}

	for _, n := range nodes {
		if n, ok := n.(*object.BaseNode); ok {
			if !bif.IsContainN(inodes, n) {
				enodes = append(enodes, n)
			}
			continue
		}
		if n, ok := n.(*object.AttrNode); ok {
			if !bif.IsContainN(inodes, n) {
				enodes = append(enodes, n)
			}
			continue
		}
	}

	for _, node := range enodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func evalInstanceofExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ie := expr.(*ast.InstanceofExpr)
	item := Eval(ie.ExprSingle, ctx)

	return bif.IsTypeMatch(item, &ie.SequenceType)
}

func evalCastExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ce := expr.(*ast.CastExpr)
	item := Eval(ce.ExprSingle, ctx)

	var ty object.Type
	switch ce.SingleType.Value() {
	case "xs:double":
		ty = object.DoubleType
	case "xs:decimal":
		ty = object.DecimalType
	case "xs:integer":
		ty = object.IntegerType
	case "xs:string":
		ty = object.StringType
	case "xs:boolean":
		ty = object.BooleanType
	}

	return bif.CastType(item, ty)
}

func evalCastableExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ce := expr.(*ast.CastableExpr)
	item := Eval(ce.ExprSingle, ctx)

	var ty object.Type
	switch ce.SingleType.Value() {
	case "xs:double":
		ty = object.DoubleType
	case "xs:decimal":
		ty = object.DecimalType
	case "xs:integer":
		ty = object.IntegerType
	case "xs:string":
		ty = object.StringType
	case "xs:boolean":
		ty = object.BooleanType
	}

	return bif.IsCastable(item, ty)
}

func evalTreatExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	return bif.NewError("treat expression not supported")
}
