package eval

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalPathExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	if ctx.Doc == nil {
		return bif.NewError("context node is undefined")
	}

	pe := expr.(*ast.PathExpr)
	ctx.CNode = []object.Node{ctx.Doc}
	ctx.CItem = ctx.Doc

	if pe.Token.Type == token.DSLASH {
		nodes := []object.Node{ctx.Doc}
		var err object.Item

		nodes, err = walkDescKind(nodes, ctx.Doc, 10, nil, nil, ctx)
		if err != nil {
			return err
		}

		ctx.CNode = nodes
		ctx.CAxis = "child::"
	} else {
		ctx.CAxis = "child::"
	}

	e := Eval(pe.ExprSingle, ctx)

	if bif.IsAnyAtomic(e) || bif.IsAnyFunc(e) {
		seq := &object.Sequence{}
		for i := 0; i < len(ctx.CNode); i++ {
			seq.Items = append(seq.Items, e)
		}
		return seq
	}

	return e
}

func evalRelativePathExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	rpe := expr.(*ast.RelativePathExpr)

	left := Eval(rpe.LeftExpr, ctx)
	if !bif.IsNode(left) && !bif.IsNodeSeq(left) {
		return bif.NewError("not a valid xpath expression")
	}

	if rpe.Token.Type == token.DSLASH {
		var nodes []object.Node
		var err object.Item

		for _, c := range ctx.CNode {
			nodes = append(nodes, c)
			nodes, err = walkDescKind(nodes, c, 10, nil, nil, ctx)
			if err != nil {
				return err
			}
		}

		ctx.CNode = nodes
		ctx.CAxis = "child::"
	} else {
		ctx.CAxis = "child::"
	}

	e := Eval(rpe.RightExpr, ctx)

	if bif.IsAnyAtomic(e) || bif.IsAnyFunc(e) {
		seq := &object.Sequence{}
		for i := 0; i < len(ctx.CNode); i++ {
			seq.Items = append(seq.Items, e)
		}
		return seq
	}

	return e
}

func evalAxisStep(expr ast.ExprSingle, ctx *object.Context) object.Item {
	if ctx.Doc == nil {
		return bif.NewError("context node is undefined")
	}
	if ctx.CNode == nil {
		ctx.CNode = []object.Node{ctx.Doc}
		ctx.CItem = ctx.Doc
	}

	as := expr.(*ast.AxisStep)

	switch as.TypeID {
	case 1: // ReverseStep

		switch as.ReverseStep.TypeID {
		case 1:
			ctx.CAxis = as.ReverseAxis.Value()
			return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
		case 2:
			ctx.CAxis = "parent::"
			return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
		default:
			return bif.NewError("not supported axis: %s", as.ReverseAxis.Value())
		}
	case 2: // ForwardStep
		switch as.ForwardStep.TypeID {
		case 1:
			if as.ForwardAxis.Value() == "" {
				ctx.CAxis = "child::"
			} else {
				ctx.CAxis = as.ForwardAxis.Value()
			}

			return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
		case 2:
			if as.AbbrevForwardStep.Token.Type == token.AT {
				ctx.CAxis = "attribute::"
			} else {
				ctx.CAxis = "child::"
			}
			if as.AbbrevForwardStep.NodeTest == nil {
				if ctx.CItem.Type() == object.DocumentNodeType {
					return ctx.Doc
				} else {
					return bif.NewError("not a valid xpath expression")
				}
			}

			return evalNodeTest(as.ForwardStep.AbbrevForwardStep.NodeTest, &as.PredicateList, ctx)
		default:
			return bif.NewError("not supported axis: %s", as.ForwardAxis.Value())
		}
	default:
		return bif.NewError("unexpected AxisStep expression")
	}
}

func evalNodeTest(test ast.NodeTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	if t, ok := test.(*ast.KindTest); ok {
		switch ctx.CAxis {
		case "child::":
			return kindTestChild(t, plist, ctx)
		case "descendant::":
			return kindTestDesc(t, plist, ctx)
		case "attribute::":
			return kindTestAttr(t, plist, ctx)
		case "self::":
			return kindTestSelf(t, plist, ctx)
		case "descendant-or-self::":
			return kindTestDescOrSelf(t, plist, ctx)
		case "following-sibling::":
			return kindTestFS(t, plist, ctx)
		case "following::":
			return kindTestFollowing(t, plist, ctx)
		case "namespace::":
			return kindTestNS(t, plist, ctx)
		case "parent::":
			return kindTestParent(t, plist, ctx)
		case "ancestor::":
			return kindTestAncestor(t, plist, ctx)
		case "preceding-sibling::":
			return kindTestPS(t, plist, ctx)
		case "preceding::":
			return kindTestPreceding(t, plist, ctx)
		case "ancestor-or-self::":
			return kindTestAncestorOrSelf(t, plist, ctx)
		default:
			return bif.NewError("not supported axis: %s", ctx.CAxis)
		}
	}

	if t, ok := test.(*ast.NameTest); ok {
		switch ctx.CAxis {
		case "child::":
			return nameTestChild(t, plist, ctx)
		case "descendant::":
			return nameTestDesc(t, plist, ctx)
		case "attribute::":
			return nameTestAttr(t, plist, ctx)
		case "self::":
			return nameTestSelf(t, plist, ctx)
		case "descendant-or-self::":
			return nameTestDescOrSelf(t, plist, ctx)
		case "following-sibling::":
			return nameTestFS(t, plist, ctx)
		case "following::":
			return nameTestFollowing(t, plist, ctx)
		case "namespace::":
			return nameTestNS(t, plist, ctx)
		case "parent::":
			return nameTestParent(t, plist, ctx)
		case "ancestor::":
			return nameTestAncestor(t, plist, ctx)
		case "preceding-sibling::":
			return nameTestPS(t, plist, ctx)
		case "preceding::":
			return nameTestPreceding(t, plist, ctx)
		case "ancestor-or-self::":
			return nameTestAncestorOrSelf(t, plist, ctx)
		default:
			return bif.NewError("not supported axis: %s", ctx.CAxis)
		}
	}

	return object.NIL
}

// ii param is used when len(plist.PL.Params) > 1
func evalPredicateList(plist *ast.PredicateList, ii *int, ctx *object.Context) object.Item {
	result := bif.NewBoolean(false)
	cnode := ctx.CNode
	focus := bif.CopyFocus(ctx)

	for i, p := range plist.PL {
		if len(p.Exprs) == 0 {
			return bif.NewError("not a valid predicate expression")
		}
		if len(p.Exprs) > 1 {
			return bif.NewError("too many items in predicate expression")
		}

		ctx.CNode = cnode
		bif.ReplaceFocus(ctx, focus)

		p := Eval(&p.Expr, ctx)
		seq := p.(*object.Sequence)

		if len(seq.Items) == 0 {
			return bif.NewBoolean(false)
		}
		if len(seq.Items) > 1 {
			return bif.NewError("too many items in predicate expression")
		}

		switch item := seq.Items[0].(type) {
		case *object.Boolean:
			if item.Value() {
				*ii++
				result = bif.NewBoolean(true)
			} else {
				return bif.NewBoolean(false)
			}
		case *object.Integer:
			if i > 0 {
				if *ii == item.Value() {
					result = bif.NewBoolean(true)
				} else {
					return bif.NewBoolean(false)
				}
			} else {
				if ctx.CPos == item.Value() {
					*ii++
					result = bif.NewBoolean(true)
				} else {
					return bif.NewBoolean(false)
				}
			}
		default:
			builtin := bif.F["fn:boolean"]
			boolObj := builtin(nil, item).(*object.Boolean)

			if boolObj.Value() {
				*ii++
				result = bif.NewBoolean(true)
			} else {
				return bif.NewBoolean(false)
			}
		}
	}

	return result
}

func evalWildcard(expr ast.ExprSingle, ctx *object.Context) object.Item {
	w := expr.(*ast.Wildcard)
	seq := &object.Sequence{}

	switch w.TypeID {
	case 1:
		var base []object.Node
		var nodes []object.Node

		for _, c := range ctx.CNode {
			if c.Type() == object.ElementNodeType || c.Type() == object.DocumentNodeType {
				base = bif.AppendNode(base, c)
			}
		}

		for _, n := range base {
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Type() == object.ElementNodeType {
					nodes = append(nodes, c)
				}
			}
		}

		ctx.CNode = nodes
		ctx.CSize = len(nodes)

		for _, node := range nodes {
			seq.Items = append(seq.Items, node)
		}
	}

	return seq
}

func kindTestChild(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

Loop:
	for _, c := range ctx.CNode {
		if t.TypeID == 3 && c.Type() == object.ElementNodeType {
			c := c.(*object.BaseNode)
			j := 0

			for _, a := range c.Attr() {
				j++
				ctx.CPos = j
				ctx.CItem = a
				ctx.CNode = []object.Node{a}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, a)
					}
				} else {
					nodes = bif.AppendNode(nodes, a)
				}
			}
		}

		i := 0
		for n := c.FirstChild(); n != nil; n = n.NextSibling() {
			if t.TypeID == 1 && c.Type() == object.DocumentNodeType {
				i++
				ctx.CPos = i
				ctx.CItem = c
				ctx.CNode = []object.Node{c}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = append(nodes, c)
					}
				} else {
					nodes = append(nodes, c)
				}
				break Loop
			} else if t.TypeID != 3 && bif.IsKindMatch(n, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = n
				ctx.CNode = []object.Node{n}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, n)
					}
				} else {
					nodes = bif.AppendNode(nodes, n)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestDesc(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	for _, c := range ctx.CNode {
		nodes, err = walkDescKind(nodes, c, t.TypeID, &ii, plist, ctx)
		if err != nil {
			return err
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestAttr(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		if c.Type() == object.ElementNodeType {
			c := c.(*object.BaseNode)
			j := 0

			for _, a := range c.Attr() {
				j++
				ctx.CPos = j
				ctx.CItem = a
				ctx.CNode = []object.Node{a}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, a)
					}
				} else {
					nodes = bif.AppendNode(nodes, a)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestSelf(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		if bif.IsKindMatch(c, t.TypeID) {
			i++
			ctx.CPos = i
			ctx.CItem = c
			ctx.CNode = []object.Node{c}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, &ii, ctx)
				if bif.IsError(pred) {
					return pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c)
				}
			} else {
				nodes = bif.AppendNode(nodes, c)
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestDescOrSelf(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		if bif.IsKindMatch(c, t.TypeID) {
			i++
			ctx.CPos = i
			ctx.CItem = c
			ctx.CNode = []object.Node{c}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, &ii, ctx)
				if bif.IsError(pred) {
					return pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c)
				}
			} else {
				nodes = bif.AppendNode(nodes, c)
			}
		}

		nodes, err = walkDescKind(nodes, c, t.TypeID, &ii, plist, ctx)
		if err != nil {
			return err
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestFS(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		for s := c.NextSibling(); s != nil; s = s.NextSibling() {
			i := 0

			if bif.IsKindMatch(s, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = s
				ctx.CNode = []object.Node{s}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, s)
					}
				} else {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestFollowing(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	for _, c := range ctx.CNode {
		i := 0
		for {
			s := c.NextSibling()
			if s == nil {
				f := c.Parent()
				if f == nil {
					break
				}
				s = f.NextSibling()
				if s == nil {
					break
				}
			}
			c = s

			if bif.IsKindMatch(s, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = s
				ctx.CNode = []object.Node{s}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, s)
					}
				} else {
					nodes = bif.AppendNode(nodes, s)
				}
			}

			nodes, err = walkDescKind(nodes, s, t.TypeID, &ii, plist, ctx)
			if err != nil {
				return err
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestNS(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	return nil
}

func kindTestParent(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		if c.Parent() != nil && bif.IsKindMatch(c.Parent(), t.TypeID) {
			i++
			ctx.CPos = i
			ctx.CItem = c.Parent()
			ctx.CNode = []object.Node{c.Parent()}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, &ii, ctx)
				if bif.IsError(pred) {
					return pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c.Parent())
				}
			} else {
				nodes = bif.AppendNode(nodes, c.Parent())
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestAncestor(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		for p := c.Parent(); p != nil; p = p.Parent() {
			if bif.IsKindMatch(p, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = p
				ctx.CNode = []object.Node{p}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, p)
					}
				} else {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestPS(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
			if bif.IsKindMatch(s, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = s
				ctx.CNode = []object.Node{s}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, s)
					}
				} else {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestPreceding(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		for {
			s := c.PrevSibling()
			if s == nil {
				p := c.Parent()
				if p == nil {
					break
				}
				s = p.PrevSibling()
				if s == nil {
					break
				}
			}
			c = s

			nodes, err = walkPrevKind(nodes, s, t.TypeID, &i, &ii, plist, ctx)

			if bif.IsKindMatch(s, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = s
				ctx.CNode = []object.Node{s}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, s)
					}
				} else {
					nodes = bif.AppendNode(nodes, s)
				}
			}

			if err != nil {
				return err
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func kindTestAncestorOrSelf(t *ast.KindTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	for _, c := range ctx.CNode {
		i := 0

		if bif.IsKindMatch(c, t.TypeID) {
			i++
			ctx.CPos = i
			ctx.CItem = c
			ctx.CNode = []object.Node{c}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, &ii, ctx)
				if bif.IsError(pred) {
					return pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c)
				}
			} else {
				nodes = bif.AppendNode(nodes, c)
			}
		}

		for p := c.Parent(); p != nil; p = p.Parent() {
			if bif.IsKindMatch(p, t.TypeID) {
				i++
				ctx.CPos = i
				ctx.CItem = p
				ctx.CNode = []object.Node{p}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, p)
					}
				} else {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestChild(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			for n := c.FirstChild(); n != nil; n = n.NextSibling() {
				if n.Type() == object.ElementNodeType &&
					t.EQName.Value() == n.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = n
					ctx.CNode = []object.Node{n}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = append(nodes, n)
						}
					} else {
						nodes = append(nodes, n)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			var base []object.Node

			for _, c := range ctx.CNode {
				if c.Type() == object.ElementNodeType || c.Type() == object.DocumentNodeType {
					base = bif.AppendNode(base, c)
				}
			}

			for _, c := range base {
				i := 0
				for n := c.FirstChild(); n != nil; n = n.NextSibling() {
					i++
					ctx.CPos = i
					ctx.CItem = n
					ctx.CNode = []object.Node{n}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() && n.Type() == object.ElementNodeType {
							nodes = append(nodes, n)
						}
					} else {
						if n.Type() == object.ElementNodeType {
							nodes = append(nodes, n)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestDesc(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	for _, c := range ctx.CNode {
		nodes, err = walkDescName(nodes, c, t, &ii, plist, ctx)
		if err != nil {
			return err
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestAttr(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			if c.Type() == object.ElementNodeType {
				c := c.(*object.BaseNode)
				for _, a := range c.Attr() {
					a := a.(*object.AttrNode)
					if a.Key() == t.EQName.Value() {
						i++
						ctx.CPos = i
						ctx.CItem = a
						ctx.CNode = []object.Node{a}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, a)
							}
						} else {
							nodes = bif.AppendNode(nodes, a)
						}
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				if c.Type() == object.ElementNodeType {
					c := c.(*object.BaseNode)
					for _, a := range c.Attr() {
						i++
						ctx.CPos = i
						ctx.CItem = a
						ctx.CNode = []object.Node{a}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, a)
							}
						} else {
							nodes = bif.AppendNode(nodes, a)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestSelf(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			if c.Type() == object.ElementNodeType &&
				t.EQName.Value() == c.Tree().Data {
				i++
				ctx.CPos = i
				ctx.CItem = c
				ctx.CNode = []object.Node{c}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = append(nodes, c)
					}
				} else {
					nodes = append(nodes, c)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				if c.Type() == object.ElementNodeType {
					i++
					ctx.CPos = i
					ctx.CItem = c
					ctx.CNode = []object.Node{c}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = append(nodes, c)
						}
					} else {
						nodes = append(nodes, c)
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestDescOrSelf(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int
	cnode := ctx.CNode

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			if c.Type() == object.ElementNodeType &&
				t.EQName.Value() == c.Tree().Data {
				i++
				ctx.CPos = i
				ctx.CItem = c
				ctx.CNode = []object.Node{c}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, c)
					}
				} else {
					nodes = bif.AppendNode(nodes, c)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				if c.Type() == object.ElementNodeType {
					i++
					ctx.CPos = i
					ctx.CItem = c
					ctx.CNode = []object.Node{c}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						nodes = bif.AppendNode(nodes, c)
					}
				}
			}
		}
	}

	for _, c := range cnode {
		nodes, err = walkDescName(nodes, c, t, &ii, plist, ctx)
		if err != nil {
			return err
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestFS(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			for s := c.NextSibling(); s != nil; s = s.NextSibling() {
				if s.Type() == object.ElementNodeType &&
					t.EQName.Value() == s.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = s
					ctx.CNode = []object.Node{s}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, s)
						}
					} else {
						nodes = bif.AppendNode(nodes, s)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				for s := c.NextSibling(); s != nil; s = s.NextSibling() {
					if s.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = s
						ctx.CNode = []object.Node{s}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, s)
							}
						} else {
							nodes = bif.AppendNode(nodes, s)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestFollowing(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			for {
				s := c.NextSibling()
				if s == nil {
					p := c.Parent()
					if p == nil {
						break
					}
					s = p.NextSibling()
					if s == nil {
						break
					}
				}
				c = s

				if s.Type() == object.ElementNodeType &&
					t.EQName.Value() == s.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = s
					ctx.CNode = []object.Node{s}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, s)
						}
					} else {
						nodes = bif.AppendNode(nodes, s)
					}
				}

				nodes, err = walkDescName(nodes, s, t, &ii, plist, ctx)
				if err != nil {
					return err
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				for {
					s := c.NextSibling()
					if s == nil {
						p := c.Parent()
						if p == nil {
							break
						}
						s = p.NextSibling()
						if s == nil {
							break
						}
					}
					c = s

					if s.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = s
						ctx.CNode = []object.Node{s}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, s)
							}
						} else {
							nodes = bif.AppendNode(nodes, s)
						}
					}

					nodes, err = walkDescName(nodes, s, t, &ii, plist, ctx)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestNS(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	return nil
}

func nameTestParent(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			if c.Parent() != nil &&
				c.Parent().Type() == object.ElementNodeType &&
				t.EQName.Value() == c.Parent().Tree().Data {
				i++
				ctx.CPos = i
				ctx.CItem = c.Parent()
				ctx.CNode = []object.Node{c.Parent()}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, c.Parent())
					}
				} else {
					nodes = bif.AppendNode(nodes, c.Parent())
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				if c.Parent() != nil &&
					c.Type() == object.ElementNodeType {
					i++
					ctx.CPos = i
					ctx.CItem = c.Parent()
					ctx.CNode = []object.Node{c.Parent()}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, c.Parent())
						}
					} else {
						nodes = bif.AppendNode(nodes, c.Parent())
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestAncestor(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			for p := c.Parent(); p != nil; p = p.Parent() {
				if p.Type() == object.ElementNodeType &&
					t.EQName.Value() == p.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = p
					ctx.CNode = []object.Node{p}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, p)
						}
					} else {
						nodes = bif.AppendNode(nodes, p)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				for p := c.Parent(); p != nil; p = p.Parent() {
					if p.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = p
						ctx.CNode = []object.Node{p}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, p)
							}
						} else {
							nodes = bif.AppendNode(nodes, p)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestPS(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
				if s.Type() == object.ElementNodeType &&
					t.EQName.Value() == s.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = s
					ctx.CNode = []object.Node{s}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, s)
						}
					} else {
						nodes = bif.AppendNode(nodes, s)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
					if s.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = s
						ctx.CNode = []object.Node{s}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, s)
							}
						} else {
							nodes = bif.AppendNode(nodes, s)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestPreceding(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var err object.Item
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0

			for {
				s := c.PrevSibling()
				if s == nil {
					p := c.Parent()
					if p == nil {
						break
					}
					s = p.PrevSibling()
					if s == nil {
						break
					}
				}
				c = s

				nodes, err = walkPrevName(nodes, s, t, &i, &ii, plist, ctx)

				if t.EQName.Value() == s.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = s
					ctx.CNode = []object.Node{s}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, s)
						}
					} else {
						nodes = bif.AppendNode(nodes, s)
					}
				}

				if err != nil {
					return err
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0

				for {
					s := c.PrevSibling()
					if s == nil {
						p := c.Parent()
						if p == nil {
							break
						}
						s = p.PrevSibling()
						if s == nil {
							break
						}
					}
					c = s

					nodes, err = walkPrevName(nodes, s, t, &i, &ii, plist, ctx)

					if s.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = s
						ctx.CNode = []object.Node{s}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, s)
							}
						} else {
							nodes = bif.AppendNode(nodes, s)
						}
					}

					if err != nil {
						return err
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

func nameTestAncestorOrSelf(t *ast.NameTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node
	var ii int

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			i := 0
			if c.Type() == object.ElementNodeType &&
				t.EQName.Value() == c.Tree().Data {
				i++
				ctx.CPos = i
				ctx.CItem = c
				ctx.CNode = []object.Node{c}

				if plist != nil && len(plist.PL) > 0 {
					pred := evalPredicateList(plist, &ii, ctx)
					if bif.IsError(pred) {
						return pred
					}

					boolObj := pred.(*object.Boolean)
					if boolObj.Value() {
						nodes = bif.AppendNode(nodes, c)
					}
				} else {
					nodes = bif.AppendNode(nodes, c)
				}
			}
			for p := c.Parent(); p != nil; p = p.Parent() {
				if p.Type() == object.ElementNodeType &&
					t.EQName.Value() == p.Tree().Data {
					i++
					ctx.CPos = i
					ctx.CItem = p
					ctx.CNode = []object.Node{p}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, p)
						}
					} else {
						nodes = bif.AppendNode(nodes, p)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				i := 0
				if c.Type() == object.ElementNodeType {
					i++
					ctx.CPos = i
					ctx.CItem = c
					ctx.CNode = []object.Node{c}

					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, &ii, ctx)
						if bif.IsError(pred) {
							return pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						nodes = bif.AppendNode(nodes, c)
					}
				}

				for p := c.Parent(); p != nil; p = p.Parent() {
					if p.Type() == object.ElementNodeType {
						i++
						ctx.CPos = i
						ctx.CItem = p
						ctx.CNode = []object.Node{p}

						if plist != nil && len(plist.PL) > 0 {
							pred := evalPredicateList(plist, &ii, ctx)
							if bif.IsError(pred) {
								return pred
							}

							boolObj := pred.(*object.Boolean)
							if boolObj.Value() {
								nodes = bif.AppendNode(nodes, p)
							}
						} else {
							nodes = bif.AppendNode(nodes, p)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, node := range nodes {
		seq.Items = append(seq.Items, node)
	}

	return seq
}

// be careful using walkDescKind bacause this function changes the ctx.CNode
func walkDescKind(nodes []object.Node, n object.Node, typeID byte, ii *int, plist *ast.PredicateList, ctx *object.Context) ([]object.Node, object.Item) {
	var err object.Item

	i := 0
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if bif.IsKindMatch(c, typeID) {
			i++
			ctx.CPos = i
			ctx.CItem = c
			ctx.CNode = []object.Node{c}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, ii, ctx)
				if bif.IsError(pred) {
					return nodes, pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c)
				}
			} else {
				nodes = bif.AppendNode(nodes, c)
			}
		}

		if c.FirstChild() != nil {
			nodes, err = walkDescKind(nodes, c, typeID, ii, plist, ctx)
			if err != nil {
				return nodes, err
			}
		}
	}
	return nodes, nil
}

func walkPrevKind(nodes []object.Node, n object.Node, typeID byte, pos, ii *int, plist *ast.PredicateList, ctx *object.Context) ([]object.Node, object.Item) {
	var err object.Item

	for c := n.LastChild(); c != nil; c = c.PrevSibling() {
		if bif.IsKindMatch(c, typeID) {
			*pos++
			ctx.CPos = *pos
			ctx.CItem = c
			ctx.CNode = []object.Node{c}

			if plist != nil && len(plist.PL) > 0 {
				pred := evalPredicateList(plist, ii, ctx)
				if bif.IsError(pred) {
					return nodes, pred
				}

				boolObj := pred.(*object.Boolean)
				if boolObj.Value() {
					nodes = bif.AppendNode(nodes, c)
				}
			} else {
				nodes = bif.AppendNode(nodes, c)
			}
		}

		if c.LastChild() != nil {
			nodes, err = walkPrevKind(nodes, c, typeID, pos, ii, plist, ctx)
			if err != nil {
				return nodes, err
			}
		}
	}
	return nodes, nil
}

// be careful using walkDescName bacause this function changes the ctx.CNode
func walkDescName(nodes []object.Node, n object.Node, t *ast.NameTest, ii *int, plist *ast.PredicateList, ctx *object.Context) ([]object.Node, object.Item) {
	var err object.Item

	i := 0
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		i++
		ctx.CPos = i
		ctx.CItem = c
		ctx.CNode = []object.Node{c}

		if c.Type() == object.ElementNodeType {
			switch t.TypeID {
			case 1:
				if c.Tree().Data == t.EQName.Value() {
					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, ii, ctx)
						if bif.IsError(pred) {
							return nodes, pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						nodes = bif.AppendNode(nodes, c)
					}
				}
			case 2:
				switch t.Wildcard.TypeID {
				case 1:
					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, ii, ctx)
						if bif.IsError(pred) {
							return nodes, pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() && c.Type() == object.ElementNodeType {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						if c.Type() == object.ElementNodeType {
							nodes = bif.AppendNode(nodes, c)
						}
					}
				}
			}
		}

		if c.FirstChild() != nil {
			nodes, err = walkDescName(nodes, c, t, ii, plist, ctx)
			if err != nil {
				return nodes, err
			}
		}
	}
	return nodes, nil
}

func walkPrevName(nodes []object.Node, n object.Node, t *ast.NameTest, pos, ii *int, plist *ast.PredicateList, ctx *object.Context) ([]object.Node, object.Item) {
	var err object.Item

	for c := n.LastChild(); c != nil; c = c.PrevSibling() {
		*pos++
		ctx.CPos = *pos
		ctx.CItem = c
		ctx.CNode = []object.Node{c}

		if c.Type() == object.ElementNodeType {
			switch t.TypeID {
			case 1:
				if c.Tree().Data == t.EQName.Value() {
					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, ii, ctx)
						if bif.IsError(pred) {
							return nodes, pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						nodes = bif.AppendNode(nodes, c)
					}
				}
			case 2:
				switch t.Wildcard.TypeID {
				case 1:
					if plist != nil && len(plist.PL) > 0 {
						pred := evalPredicateList(plist, ii, ctx)
						if bif.IsError(pred) {
							return nodes, pred
						}

						boolObj := pred.(*object.Boolean)
						if boolObj.Value() && c.Type() == object.ElementNodeType {
							nodes = bif.AppendNode(nodes, c)
						}
					} else {
						if c.Type() == object.ElementNodeType {
							nodes = bif.AppendNode(nodes, c)
						}
					}
				}
			}
		}

		if c.LastChild() != nil {
			nodes, err = walkPrevName(nodes, c, t, pos, ii, plist, ctx)
			if err != nil {
				return nodes, err
			}
		}
	}
	return nodes, nil
}
