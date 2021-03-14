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

	if pe.Token.Type == token.SLASH && pe.ExprSingle == nil {
		return ctx.Doc
	}
	if pe.ExprSingle == nil {
		return bif.NewError("unexpected end of xpath statement")
	}

	if pe.Token.Type == token.DSLASH {
		ctx.CAxis = "descendant-or-self::"
	} else {
		ctx.CAxis = "child::"
	}

	return Eval(pe.ExprSingle, ctx)
}

func evalRelativePathExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	rpe := expr.(*ast.RelativePathExpr)

	left := Eval(rpe.LeftExpr, ctx)
	if !bif.IsNode(left) && !bif.IsNodeSeq(left) {
		return bif.NewError("path expression cannot contain type %s except the last step", left.Type())
	}

	if rpe.Token.Type == token.DSLASH {
		ctx.CAxis = "descendant-or-self::"
	} else {
		ctx.CAxis = "child::"
	}

	return Eval(rpe.RightExpr, ctx)
}

func evalAxisStep(expr ast.ExprSingle, ctx *object.Context) object.Item {
	as := expr.(*ast.AxisStep)

	if _, ok := ctx.CItem.(*object.BaseNode); !ok {
		if ctx.Doc == nil {
			return bif.NewError("context node is undefined")
		}
		ctx.CItem = ctx.Doc
		ctx.CNode = []object.Node{ctx.Doc}
	}

	var nodes []object.Node
	switch as.TypeID {
	case 1: // ReverseStep
		switch as.ReverseStep.TypeID {
		case 1:
			switch as.ReverseAxis.Value() {
			case "parent::":
				for _, c := range ctx.CNode {
					if c.Parent() != nil && !bif.IsContainN(nodes, c.Parent()) {
						nodes = append(nodes, c.Parent())
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ReverseAxis.Value()
				return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
			case "ancestor::":
				for _, c := range ctx.CNode {
					for p := c.Parent(); p != nil; p = p.Parent() {
						if !bif.IsContainN(nodes, p) {
							nodes = append(nodes, p)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ReverseAxis.Value()
				return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
			case "preceding-sibling::":
				for _, c := range ctx.CNode {
					for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
						if !bif.IsContainN(nodes, s) {
							nodes = append(nodes, s)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ReverseAxis.Value()
				return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
			case "preceding::":
				for _, c := range ctx.CNode {
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

						if !bif.IsContainN(nodes, s) {
							nodes = append(nodes, s)
						}
						bif.WalkDesc(nodes, s)
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ReverseAxis.Value()
				return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
			case "ancestor-or-self::":
				for _, c := range ctx.CNode {
					nodes = append(nodes, c.Self())
					for p := c.Parent(); p != nil; p = p.Parent() {
						if !bif.IsContainN(nodes, p) {
							nodes = append(nodes, p)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ReverseAxis.Value()
				return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
			default:
				return bif.NewError("unexpected AxisStep expression")
			}
		case 2:
			for _, c := range ctx.CNode {
				if c.Parent() != nil && !bif.IsContainN(nodes, c.Parent()) {
					nodes = append(nodes, c.Parent())
				}
			}

			ctx.CNode = nodes
			ctx.CAxis = "parent::"
			return evalNodeTest(as.ReverseStep.NodeTest, &as.PredicateList, ctx)
		default:
			return bif.NewError("not supported axis: %s", as.ReverseAxis.Value())
		}
	case 2: // ForwardStep
		switch as.ForwardStep.TypeID {
		case 1:
			switch as.ForwardAxis.Value() {
			case "child::":
				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "descendant::":
				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "attribute::":
				for _, c := range ctx.CNode {
					if c, ok := c.(*object.BaseNode); ok {
						for _, a := range c.Attr() {
							n := &object.AttrNode{}
							n.SetAttr(&a)
							n.SetTree(c.Tree())
							nodes = append(nodes, n)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "self::":
				for _, c := range ctx.CNode {
					nodes = append(nodes, c.Self())
				}

				ctx.CNode = nodes
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "descendant-or-self::":
				for _, c := range ctx.CNode {
					nodes = append(nodes, c.Self())
					bif.WalkDesc(nodes, c)
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "following-sibling::":
				for _, c := range ctx.CNode {
					for s := c.NextSibling(); s != nil; s = s.NextSibling() {
						if !bif.IsContainN(nodes, s) {
							nodes = append(nodes, s)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "following::":
				for _, c := range ctx.CNode {
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

						if !bif.IsContainN(nodes, s) {
							nodes = append(nodes, s)
						}
						bif.WalkDesc(nodes, s)
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = as.ForwardAxis.Value()
				return evalNodeTest(as.ForwardStep.NodeTest, &as.PredicateList, ctx)
			case "namespace::":
				fallthrough
			default:
				return bif.NewError("not supported axis: %s", as.ForwardAxis.Value())
			}
		case 2:
			if as.AbbrevForwardStep.Token.Type == token.AT {
				for _, c := range ctx.CNode {
					if c, ok := c.(*object.BaseNode); ok {
						for _, a := range c.Attr() {
							n := &object.AttrNode{}
							n.SetAttr(&a)
							n.SetTree(c.Tree())
							nodes = append(nodes, n)
						}
					}
				}

				ctx.CNode = nodes
				ctx.CAxis = "attribute::"
			}

			return evalNodeTest(as.ForwardStep.AbbrevForwardStep.NodeTest, &as.PredicateList, ctx)
		default:
			return bif.NewError("unexpected AxisStep expression")
		}
	default:
		return bif.NewError("unexpected AxisStep expression")
	}
}

func evalNodeTest(test ast.NodeTest, plist *ast.PredicateList, ctx *object.Context) object.Item {
	var nodes []object.Node

	if t, ok := test.(*ast.KindTest); ok {
		switch ctx.CAxis {
		case "child::":
			nodes = kindTestChild(t, ctx)
		case "descendant::":
			nodes = kindTestDesc(t, ctx)
		case "attribute::":
			nodes = kindTestAttr(t, ctx)
		case "self::":
			nodes = kindTestSelf(t, ctx)
		case "descendant-or-self::":
			nodes = kindTestDescOrSelf(t, ctx)
		case "following-sibling::":
			nodes = kindTestFS(t, ctx)
		case "following::":
			nodes = kindTestFollowing(t, ctx)
		case "namespace::":
			nodes = kindTestNS(t, ctx)
		case "parent::":
			nodes = kindTestParent(t, ctx)
		case "ancestor::":
			nodes = kindTestAncestor(t, ctx)
		case "preceding-sibling::":
			nodes = kindTestPS(t, ctx)
		case "preceding::":
			nodes = kindTestPreceding(t, ctx)
		case "ancestor-or-self::":
			nodes = kindTestAncestorOrSelf(t, ctx)
		default:
			return bif.NewError("not supported axis: %s", ctx.CAxis)
		}
	}

	if t, ok := test.(*ast.NameTest); ok {
		switch ctx.CAxis {
		case "child::":
			nodes = nameTestChild(t, ctx)
		case "descendant::":
			nodes = nameTestDesc(t, ctx)
		case "attribute::":
			nodes = nameTestAttr(t, ctx)
		case "self::":
			nodes = nameTestSelf(t, ctx)
		case "descendant-or-self::":
			nodes = nameTestDescOrSelf(t, ctx)
		case "following-sibling::":
			nodes = nameTestFS(t, ctx)
		case "following::":
			nodes = nameTestFollowing(t, ctx)
		case "namespace::":
			nodes = nameTestNS(t, ctx)
		case "parent::":
			nodes = nameTestParent(t, ctx)
		case "ancestor::":
			nodes = nameTestAncestor(t, ctx)
		case "preceding-sibling::":
			nodes = nameTestPS(t, ctx)
		case "preceding::":
			nodes = nameTestPreceding(t, ctx)
		case "ancestor-or-self::":
			nodes = nameTestAncestorOrSelf(t, ctx)
		default:
			return bif.NewError("not supported axis: %s", ctx.CAxis)
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for i, node := range nodes {
		ctx.CPos = i + 1
		ctx.CItem = node

		if len(plist.PL) > 0 {
			if evalPredicateList(node, plist, ctx) {
				seq.Items = append(seq.Items, node)
			}
		} else {
			seq.Items = append(seq.Items, node)
		}
	}

	return seq
}

func evalPredicateList(node object.Node, plist *ast.PredicateList, ctx *object.Context) bool {
	for _, p := range plist.PL {
		p := Eval(&p.Expr, ctx)

		switch p := p.(type) {
		case *object.Boolean:
			return p.Value()
		case *object.Integer:
			if ctx.CPos == p.Value() {
				return true
			}
		case *object.Sequence:
			for _, item := range p.Items {
				switch item := item.(type) {
				case *object.Integer:
					if ctx.CPos == item.Value() {
						return true
					}
				case *object.BaseNode:
					if node.Type() != object.ElementNodeType && node.Type() != object.DocumentNodeType {
						return false
					}
					if item.Tree().Data == node.Tree().Data {
						return true
					}
				case *object.AttrNode:
					if node.Type() != object.ElementNodeType && node.Type() != object.DocumentNodeType {
						return false
					}
					for _, attr := range node.Tree().Attr {
						if attr.Key == item.Key() {
							return true
						}
					}
				}
			}
		default:
			builtin := bif.Builtins["fn:boolean"]
			bl := builtin(p)
			boolObj := bl.(*object.Boolean)
			return boolObj.Value()
		}
	}

	return false
}

func kindTestChild(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

Loop:
	for _, c := range ctx.CNode {
		for n := c.FirstChild(); n != nil; n = n.NextSibling() {
			switch t.TypeID {
			case 1:
				if c.Type() == object.DocumentNodeType {
					nodes = append(nodes, c)
					break Loop
				}
			case 2:
				if n.Type() == object.ElementNodeType {
					nodes = append(nodes, n)
				}
			case 3:
				if n.Type() == object.ElementNodeType {
					n := n.(*object.BaseNode)
					nodes = append(nodes, n.Attr()...)
				}
			case 7:
				if n.Type() == object.CommentNodeType {
					nodes = append(nodes, n)
				}
			case 8:
				if n.Type() == object.TextNodeType {
					nodes = append(nodes, n)
				}
			case 10:
				nodes = append(nodes, n)
				if n.Type() == object.ElementNodeType {
					n := n.(*object.BaseNode)
					nodes = append(nodes, n.Attr()...)
				}
			}
		}
	}

	return nodes
}

func kindTestDesc(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		bif.WalkDescKind(nodes, c, t.TypeID)
	}
	return nodes
}

func kindTestAttr(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		if c.Type() == object.ElementNodeType {
			c := c.(*object.BaseNode)
			nodes = append(nodes, c.Attr()...)
		}
	}
	return nodes
}

func kindTestSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestDescOrSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		nodes = append(nodes, c)
		bif.WalkDescKind(nodes, c, t.TypeID)
	}
	return nodes
}

func kindTestFS(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestFollowing(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestNS(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestParent(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestAncestor(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestPS(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestPreceding(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestAncestorOrSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestChild(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestDesc(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestAttr(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestDescOrSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestFS(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestFollowing(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestNS(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestParent(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestAncestor(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestPS(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestPreceding(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestAncestorOrSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}
