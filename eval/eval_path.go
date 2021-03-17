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
		var nodes []object.Node
		nodes = bif.WalkDescKind(nodes, ctx.Doc, 10)
		ctx.CNode = append(ctx.CNode, nodes...)
		ctx.CAxis = "child::"
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
		var nodes []object.Node
		for _, c := range ctx.CNode {
			nodes = append(nodes, c)
			nodes = bif.WalkDescKind(nodes, c, 10)
		}

		ctx.CNode = nodes
		ctx.CAxis = "child::"
	} else {
		ctx.CAxis = "child::"
	}

	return Eval(rpe.RightExpr, ctx)
}

func evalAxisStep(expr ast.ExprSingle, ctx *object.Context) object.Item {
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

func evalWildcard(expr ast.ExprSingle, ctx *object.Context) object.Item {
	w := expr.(*ast.Wildcard)

	var nodes []object.Node
	switch w.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if c.Type() == object.ElementNodeType {
				nodes = bif.AppendNode(nodes, c)
			}
		}
	}

	ctx.CNode = nodes
	ctx.CSize = len(nodes)

	seq := &object.Sequence{}
	for _, n := range nodes {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if c.Type() == object.ElementNodeType {
				seq.Items = append(seq.Items, c)
			}
		}
	}

	return seq
}

func kindTestChild(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

Loop:
	for _, c := range ctx.CNode {
		for n := c.FirstChild(); n != nil; n = n.NextSibling() {
			if t.TypeID == 1 && c.Type() == object.DocumentNodeType {
				nodes = append(nodes, c)
				break Loop
			}
			nodes = bif.AppendKind(nodes, n, t.TypeID)
		}
	}

	return nodes
}

func kindTestDesc(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		nodes = bif.WalkDescKind(nodes, c, t.TypeID)
	}
	return nodes
}

func kindTestAttr(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		if c.Type() == object.ElementNodeType {
			c := c.(*object.BaseNode)
			for _, a := range c.Attr() {
				nodes = bif.AppendNode(nodes, a)
			}
		}
	}
	return nodes
}

func kindTestSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		nodes = bif.AppendKind(nodes, c, t.TypeID)
	}
	return nodes
}

func kindTestDescOrSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		nodes = bif.AppendKind(nodes, c, t.TypeID)
		nodes = bif.WalkDescKind(nodes, c, t.TypeID)
	}
	return nodes
}

func kindTestFS(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		for s := c.NextSibling(); s != nil; s = s.NextSibling() {
			if bif.IsKindMatch(s, t.TypeID) {
				nodes = bif.AppendNode(nodes, s)
			}
		}
	}
	return nodes
}

func kindTestFollowing(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
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

			if bif.IsKindMatch(s, t.TypeID) {
				nodes = bif.AppendNode(nodes, s)
			}
			nodes = bif.WalkDescKind(nodes, s, t.TypeID)
		}
	}
	return nodes
}

func kindTestNS(t *ast.KindTest, ctx *object.Context) []object.Node {
	return nil
}

func kindTestParent(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		if c.Parent() != nil && bif.IsKindMatch(c.Parent(), t.TypeID) {
			nodes = bif.AppendNode(nodes, c.Parent())
		}
	}
	return nodes
}

func kindTestAncestor(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		for p := c.Parent(); p != nil; p = p.Parent() {
			if bif.IsKindMatch(p, t.TypeID) {
				nodes = bif.AppendNode(nodes, p)
			}
		}
	}
	return nodes
}

func kindTestPS(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
			if bif.IsKindMatch(s, t.TypeID) {
				nodes = bif.AppendNode(nodes, s)
			}
		}
	}
	return nodes
}

func kindTestPreceding(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
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

			if bif.IsKindMatch(s, t.TypeID) {
				nodes = bif.AppendNode(nodes, s)
			}
			nodes = bif.WalkDescKind(nodes, s, t.TypeID)
		}
	}
	return nodes
}

func kindTestAncestorOrSelf(t *ast.KindTest, ctx *object.Context) []object.Node {
	var nodes []object.Node
	for _, c := range ctx.CNode {
		if bif.IsKindMatch(c, t.TypeID) {
			nodes = bif.AppendNode(nodes, c)
		}
		for p := c.Parent(); p != nil; p = p.Parent() {
			if bif.IsKindMatch(p, t.TypeID) {
				nodes = bif.AppendNode(nodes, p)
			}
		}
	}
	return nodes
}

func nameTestChild(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			for n := c.FirstChild(); n != nil; n = n.NextSibling() {
				if n.Type() == object.ElementNodeType && t.EQName.Value() == n.Tree().Data {
					nodes = append(nodes, n)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				for n := c.FirstChild(); n != nil; n = n.NextSibling() {
					nodes = append(nodes, n)
				}
			}
		}
	}

	return nodes
}

func nameTestDesc(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	for _, c := range ctx.CNode {
		nodes = bif.WalkDescName(nodes, c, t)
	}

	return nodes
}

func nameTestAttr(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if c.Type() == object.ElementNodeType {
				c := c.(*object.BaseNode)
				for _, a := range c.Attr() {
					a := a.(*object.AttrNode)
					if a.Key() == t.EQName.Value() {
						nodes = bif.AppendNode(nodes, a)
					}
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				if c.Type() == object.ElementNodeType {
					c := c.(*object.BaseNode)
					for _, a := range c.Attr() {
						nodes = bif.AppendNode(nodes, a)
					}
				}
			}
		}
	}

	return nodes
}

func nameTestSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if t.EQName.Value() == c.Tree().Data {
				nodes = append(nodes, c)
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			nodes = append(nodes, ctx.CNode...)
		}
	}

	return nodes
}

func nameTestDescOrSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	for _, c := range ctx.CNode {
		nodes = bif.WalkDescName(nodes, c, t)
	}

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if t.EQName.Value() == c.Tree().Data {
				nodes = bif.AppendNode(nodes, c)
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				nodes = bif.AppendNode(nodes, c)
			}
		}
	}

	return nodes
}

func nameTestFS(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			for s := c.NextSibling(); s != nil; s = s.NextSibling() {
				if t.EQName.Value() == s.Tree().Data {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				for s := c.NextSibling(); s != nil; s = s.NextSibling() {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	}

	return nodes
}

func nameTestFollowing(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
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

				if t.EQName.Value() == s.Tree().Data {
					nodes = bif.AppendNode(nodes, s)
				}
				nodes = bif.WalkDescName(nodes, s, t)
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
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

					nodes = bif.AppendNode(nodes, s)
					nodes = bif.WalkDescName(nodes, s, t)
				}
			}
		}
	}

	return nodes
}

func nameTestNS(t *ast.NameTest, ctx *object.Context) []object.Node {
	return nil
}

func nameTestParent(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if c.Parent() != nil && t.EQName.Value() == c.Parent().Tree().Data {
				nodes = bif.AppendNode(nodes, c.Parent())
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				if c.Parent() != nil {
					nodes = bif.AppendNode(nodes, c.Parent())
				}
			}
		}
	}

	return nodes
}

func nameTestAncestor(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			for p := c.Parent(); p != nil; p = p.Parent() {
				if t.EQName.Value() == p.Tree().Data {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				for p := c.Parent(); p != nil; p = p.Parent() {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	}

	return nodes
}

func nameTestPS(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
				if t.EQName.Value() == s.Tree().Data {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				for s := c.PrevSibling(); s != nil; s = s.PrevSibling() {
					nodes = bif.AppendNode(nodes, s)
				}
			}
		}
	}

	return nodes
}

func nameTestPreceding(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
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
				c = s

				if t.EQName.Value() == s.Tree().Data {
					nodes = bif.AppendNode(nodes, s)
				}
				nodes = bif.WalkDescName(nodes, s, t)
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
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
					c = s

					nodes = bif.AppendNode(nodes, s)
					nodes = bif.WalkDescName(nodes, s, t)
				}
			}
		}
	}

	return nodes
}

func nameTestAncestorOrSelf(t *ast.NameTest, ctx *object.Context) []object.Node {
	var nodes []object.Node

	switch t.TypeID {
	case 1:
		for _, c := range ctx.CNode {
			if t.EQName.Value() == c.Tree().Data {
				nodes = bif.AppendNode(nodes, c)
			}
			for p := c.Parent(); p != nil; p = p.Parent() {
				if t.EQName.Value() == p.Tree().Data {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	case 2:
		switch t.Wildcard.TypeID {
		case 1:
			for _, c := range ctx.CNode {
				nodes = bif.AppendNode(nodes, c)
				for p := c.Parent(); p != nil; p = p.Parent() {
					nodes = bif.AppendNode(nodes, p)
				}
			}
		}
	}

	return nodes
}
