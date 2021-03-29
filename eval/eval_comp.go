package eval

import (
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

func evalComparisonExpr(expr ast.ExprSingle, ctx *object.Context) object.Item {
	ce := expr.(*ast.ComparisonExpr)

	left := Eval(ce.LeftExpr, ctx)
	right := Eval(ce.RightExpr, ctx)
	op := ce.Token

	switch {
	case bif.IsSeq(left) && bif.IsNumeric(right):
		return compSeqNumber(op, left, right)
	case bif.IsSeq(left) && bif.IsString(right):
		return compSeqString(op, left, right)
	case bif.IsSeq(left) && bif.IsArray(right):
		return compSeqArray(op, left, right)
	case bif.IsSeq(left) && bif.IsSeq(right):
		return compSeqSeq(op, left, right, ctx)

	case bif.IsNode(left) && bif.IsString(right):
		return compNodeString(op, left, right, ctx)
	case bif.IsString(left) && bif.IsNode(right):
		return compStringNode(op, left, right, ctx)
	case bif.IsNode(left) && bif.IsNumeric(right):
		return compNodeNumber(op, left, right, ctx)
	case bif.IsNumeric(left) && bif.IsNode(right):
		return compNumberNode(op, left, right, ctx)
	case bif.IsNode(left) && bif.IsNode(right):
		return compNodeNode(op, left, right, ctx)

	case bif.IsNumeric(left) && bif.IsNumeric(right):
		return compNumberNumber(op, left, right)
	case bif.IsNumeric(left) && bif.IsSeq(right):
		return compNumberSeq(op, left, right)
	case bif.IsNumeric(left) && bif.IsArray(right):
		return compNumberArray(op, left, right)

	case bif.IsString(left) && bif.IsString(right):
		return compStringString(op, left, right)
	case bif.IsString(left) && bif.IsArray(right):
		return compStringArray(op, left, right)
	case bif.IsString(left) && bif.IsSeq(right):
		return compStringSeq(op, left, right)

	case bif.IsArray(left) && bif.IsNumeric(right):
		return compArrayNumber(op, left, right)
	case bif.IsArray(left) && bif.IsString(right):
		return compArrayString(op, left, right)
	case bif.IsArray(left) && bif.IsArray(right):
		return compArrayArray(op, left, right)
	case bif.IsArray(left) && bif.IsSeq(right):
		return compArraySeq(op, left, right)

	case bif.IsBoolean(left) && bif.IsBoolean(right):
		return compBool(op, left, right)
	}

	return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
}

func compNumberNumber(op token.Token, left, right object.Item) object.Item {
	switch op.Type {
	case token.EQ, token.EQV:
		builtin := bif.F["op:numeric-equal"]
		return builtin(nil, left, right)
	case token.NE, token.NEV:
		builtin := bif.F["op:numeric-equal"]
		b := builtin(nil, left, right)
		boolean := b.(*object.Boolean)
		return bif.NewBoolean(!boolean.Value())
	case token.LT, token.LTV:
		builtin := bif.F["op:numeric-less-than"]
		return builtin(nil, left, right)
	case token.GT, token.GTV:
		builtin := bif.F["op:numeric-greater-than"]
		return builtin(nil, left, right)
	case token.LE, token.LEV:
		builtin := bif.F["op:numeric-less-than"]
		b := builtin(nil, left, right)
		boolean := b.(*object.Boolean)

		if !boolean.Value() {
			builtin = bif.F["op:numeric-equal"]
			return builtin(nil, left, right)
		}

		return boolean
	case token.GE, token.GEV:
		builtin := bif.F["op:numeric-greater-than"]
		b := builtin(nil, left, right)
		boolean := b.(*object.Boolean)

		if !boolean.Value() {
			builtin = bif.F["op:numeric-equal"]
			return builtin(nil, left, right)
		}

		return boolean
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compNumberArray(op token.Token, left, right object.Item) object.Item {
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compNumberSeq(op token.Token, left, right object.Item) object.Item {
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.String).Value()
	rightVal := right.(*object.String).Value()

	switch op.Type {
	case token.EQ, token.EQV:
		return bif.NewBoolean(leftVal == rightVal)
	case token.NE, token.NEV:
		return bif.NewBoolean(leftVal != rightVal)
	case token.LT, token.LTV:
		return bif.NewBoolean(leftVal < rightVal)
	case token.LE, token.LEV:
		return bif.NewBoolean(leftVal <= rightVal)
	case token.GT, token.GTV:
		return bif.NewBoolean(leftVal > rightVal)
	case token.GE, token.GEV:
		return bif.NewBoolean(leftVal >= rightVal)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() == e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() != e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() < e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() <= e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() > e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() >= e.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringSeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() == e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.BaseNode); ok {
				if leftVal.Value() == e.Text() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.AttrNode); ok {
				if leftVal.Value() == e.Text() {
					return bif.NewBoolean(true)
				}
				continue
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() != e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.BaseNode); ok {
				if leftVal.Value() != e.Text() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.AttrNode); ok {
				if leftVal.Value() != e.Text() {
					return bif.NewBoolean(true)
				}
				continue
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() < e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() <= e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() > e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range rightVal.Items {
			if e, ok := item.(*object.String); ok {
				if leftVal.Value() >= e.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() == rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() != rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() < rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() <= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() > rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() >= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArraySeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() {
				return bif.NewBoolean(true)
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() == rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.BaseNode); ok {
				if e.Text() == rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.AttrNode); ok {
				if e.Text() == rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() != rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.BaseNode); ok {
				if e.Text() != rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if e, ok := item.(*object.AttrNode); ok {
				if e.Text() != rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() < rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() <= rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() > rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, item := range leftVal.Items {
			if e, ok := item.(*object.String); ok {
				if e.Value() >= rightVal.Value() {
					return bif.NewBoolean(true)
				}
				continue
			}
			if _, ok := item.(*object.BaseNode); ok {
				return bif.NewBoolean(true)
			}
			if _, ok := item.(*object.AttrNode); ok {
				return bif.NewBoolean(true)
			}

			return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqSeq(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ, token.EQV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.NE, token.NEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LT, token.LTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GT, token.GTV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() {
					return bif.NewBoolean(true)
				}
			}
		}
		return bif.NewBoolean(false)
	case token.IS, token.DGT, token.DLT:
		if len(leftVal.Items) != 1 {
			return bif.NewError("wrong number of items. got=%d, expected=1", len(leftVal.Items))
		}
		if len(rightVal.Items) != 1 {
			return bif.NewError("wrong number of items. got=%d, expected=1", len(rightVal.Items))
		}
		if !bif.IsNode(leftVal.Items[0]) || !bif.IsNode(rightVal.Items[0]) {
			return bif.NewError("node types expected. got=%s, %s", leftVal.Items[0].Type(), rightVal.Items[0].Type())
		}

		doc := ctx.Doc.(*object.BaseNode)
		if leftItem, ok := leftVal.Items[0].(*object.BaseNode); ok {
			switch rightItem := rightVal.Items[0].(type) {
			case *object.BaseNode:
				switch op.Type {
				case token.IS:
					return bif.NewBoolean(leftItem.Tree() == rightItem.Tree())
				case token.DGT:
					return bif.IsPrecede(rightItem, leftItem, doc)
				case token.DLT:
					return bif.IsPrecede(leftItem, rightItem, doc)
				}
			case *object.AttrNode:
				switch op.Type {
				case token.IS:
					return bif.NewBoolean(false)
				case token.DGT:
					return bif.IsPrecede(rightItem, leftItem, doc)
				case token.DLT:
					return bif.IsPrecede(leftItem, rightItem, doc)
				}
			}
		}

		if leftItem, ok := leftVal.Items[0].(*object.AttrNode); ok {
			switch rightItem := rightVal.Items[0].(type) {
			case *object.BaseNode:
				switch op.Type {
				case token.IS:
					return bif.NewBoolean(false)
				case token.DGT:
					return bif.IsPrecede(rightItem, leftItem, doc)
				case token.DLT:
					return bif.IsPrecede(leftItem, rightItem, doc)
				}
			case *object.AttrNode:
				switch op.Type {
				case token.IS:
					return bif.NewBoolean(leftItem.Tree() == rightItem.Tree() && leftItem.Key() == rightItem.Key())
				case token.DGT:
					return bif.IsPrecede(rightItem, leftItem, doc)
				case token.DLT:
					return bif.IsPrecede(leftItem, rightItem, doc)
				}
			}
		}

		return bif.NewBoolean(false)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compBool(op token.Token, left, right object.Item) object.Item {
	leftVal, ok := left.(*object.Boolean)
	if !ok {
		return bif.NewError("cannot compare Types: %s, %s", left.Type(), right.Type())
	}

	rightVal, ok := right.(*object.Boolean)
	if !ok {
		return bif.NewError("cannot compare Types: %s, %s", left.Type(), right.Type())
	}

	switch op.Type {
	case token.EQ, token.EQV:
		return bif.NewBoolean(leftVal.Value() == rightVal.Value())
	case token.NE, token.NEV:
		return bif.NewBoolean(leftVal.Value() != rightVal.Value())
	case token.GT, token.GTV:
		if leftVal.Value() && !rightVal.Value() {
			return bif.NewBoolean(true)
		}
		return bif.NewBoolean(false)
	case token.GE, token.GEV:
		if !leftVal.Value() && rightVal.Value() {
			return bif.NewBoolean(false)
		}
		return bif.NewBoolean(true)
	case token.LT, token.LTV:
		if !leftVal.Value() && rightVal.Value() {
			return bif.NewBoolean(true)
		}
		return bif.NewBoolean(false)
	case token.LE, token.LEV:
		if leftVal.Value() && !rightVal.Value() {
			return bif.NewBoolean(false)
		}
		return bif.NewBoolean(true)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compNodeString(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	rightVal := right.(*object.String)

	switch leftNode := left.(type) {
	case *object.BaseNode:
		switch op.Type {
		case token.EQ, token.EQV:
			if leftNode.Text() == rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.NE, token.NEV:
			if leftNode.Text() != rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.GT, token.GTV:
			if leftNode.Text() > rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.GE, token.GEV:
			if leftNode.Text() >= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.LT, token.LTV:
			if leftNode.Text() < rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.LE, token.LEV:
			if leftNode.Text() <= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.IS, token.DGT, token.DLT:
			return bif.NewError("cannot compare node with string: %s, %s", leftNode.Type(), rightVal.Type())
		}
	case *object.AttrNode:
		switch op.Type {
		case token.EQ, token.EQV:
			if leftNode.Text() == rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.NE, token.NEV:
			if leftNode.Text() != rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.GT, token.GTV:
			if leftNode.Text() > rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.GE, token.GEV:
			if leftNode.Text() >= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.LT, token.LTV:
			if leftNode.Text() < rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.LE, token.LEV:
			if leftNode.Text() <= rightVal.Value() {
				return bif.NewBoolean(true)
			}
		case token.IS, token.DGT, token.DLT:
			return bif.NewError("cannot compare node with string: %s, %s", leftNode.Type(), rightVal.Type())
		}
	}

	return bif.NewBoolean(false)
}

func compStringNode(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	leftVal := left.(*object.String)

	switch rightNode := right.(type) {
	case *object.BaseNode:
		switch op.Type {
		case token.EQ, token.EQV:
			if leftVal.Value() == rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.NE, token.NEV:
			if leftVal.Value() != rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.GT, token.GTV:
			if leftVal.Value() > rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.GE, token.GEV:
			if leftVal.Value() >= rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.LT, token.LTV:
			if leftVal.Value() < rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.LE, token.LEV:
			if leftVal.Value() <= rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.IS, token.DGT, token.DLT:
			return bif.NewError("cannot compare node with string: %s, %s", leftVal.Type(), rightNode.Type())
		}
	case *object.AttrNode:
		switch op.Type {
		case token.EQ, token.EQV:
			if leftVal.Value() == rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.NE, token.NEV:
			if leftVal.Value() != rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.GT, token.GTV:
			if leftVal.Value() > rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.GE, token.GEV:
			if leftVal.Value() >= rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.LT, token.LTV:
			if leftVal.Value() < rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.LE, token.LEV:
			if leftVal.Value() <= rightNode.Text() {
				return bif.NewBoolean(true)
			}
		case token.IS, token.DGT, token.DLT:
			return bif.NewError("cannot compare node with string: %s, %s", leftVal.Type(), rightNode.Type())
		}
	}

	return bif.NewBoolean(false)
}

func compNodeNumber(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	var leftVal string

	if leftNode, ok := left.(*object.BaseNode); ok {
		leftVal = leftNode.Text()
	}
	if leftNode, ok := left.(*object.AttrNode); ok {
		leftVal = leftNode.Text()
	}

	switch op.Type {
	case token.EQ, token.EQV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsEQ(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsEQ(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsEQ(bif.NewInteger(int(l)), rightVal)
			}
		}
	case token.NE, token.NEV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsNE(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsNE(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsNE(bif.NewInteger(int(l)), rightVal)
			}
		}
	case token.LT, token.LTV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsLT(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsLT(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsLT(bif.NewInteger(int(l)), rightVal)
			}
		}
	case token.LE, token.LEV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsLE(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsLE(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsLE(bif.NewInteger(int(l)), rightVal)
			}
		}
	case token.GT, token.GTV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsGT(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsGT(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsGT(bif.NewInteger(int(l)), rightVal)
			}
		}
	case token.GE, token.GEV:
		switch rightVal := right.(type) {
		case *object.Double:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsGE(bif.NewDouble(l), rightVal)
			}
		case *object.Decimal:
			if l, err := strconv.ParseFloat(leftVal, 64); err == nil {
				return bif.IsGE(bif.NewDecimal(l), rightVal)
			}
		case *object.Integer:
			if l, err := strconv.ParseInt(leftVal, 0, 64); err == nil {
				return bif.IsGE(bif.NewInteger(int(l)), rightVal)
			}
		}
	}

	return bif.NewError("cannot compare")
}

func compNumberNode(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	var rightVal string
	if rightNode, ok := right.(*object.BaseNode); ok {
		rightVal = rightNode.Text()
	}
	if rightNode, ok := right.(*object.AttrNode); ok {
		rightVal = rightNode.Text()
	}

	switch op.Type {
	case token.EQ, token.EQV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsEQ(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsEQ(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsEQ(leftVal, bif.NewInteger(int(r)))
			}
		}
	case token.NE, token.NEV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsNE(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsNE(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsNE(leftVal, bif.NewInteger(int(r)))
			}
		}
	case token.LT, token.LTV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsLT(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsLT(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsLT(leftVal, bif.NewInteger(int(r)))
			}
		}
	case token.LE, token.LEV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsLE(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsLE(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsLE(leftVal, bif.NewInteger(int(r)))
			}
		}
	case token.GT, token.GTV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsGT(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsGT(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsGT(leftVal, bif.NewInteger(int(r)))
			}
		}
	case token.GE, token.GEV:
		switch leftVal := left.(type) {
		case *object.Double:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsGE(leftVal, bif.NewDouble(r))
			}
		case *object.Decimal:
			if r, err := strconv.ParseFloat(rightVal, 64); err == nil {
				return bif.IsGE(leftVal, bif.NewDecimal(r))
			}
		case *object.Integer:
			if r, err := strconv.ParseInt(rightVal, 0, 64); err == nil {
				return bif.IsGE(leftVal, bif.NewInteger(int(r)))
			}
		}
	}

	return bif.NewError("cannot compare")
}

func compNodeNode(op token.Token, left, right object.Item, ctx *object.Context) object.Item {
	if leftNode, ok := left.(*object.BaseNode); ok {
		switch rightNode := right.(type) {
		case *object.BaseNode:
			switch op.Type {
			case token.EQ, token.EQV:
				if leftNode.Text() == rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.NE, token.NEV:
				if leftNode.Text() != rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GT, token.GTV:
				if leftNode.Text() > rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GE, token.GEV:
				if leftNode.Text() >= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LT, token.LTV:
				if leftNode.Text() < rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LE, token.LEV:
				if leftNode.Text() <= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.IS:
				return bif.NewBoolean(leftNode.Tree() == rightNode.Tree())
			case token.DGT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(rightNode, leftNode, doc)
			case token.DLT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(leftNode, rightNode, doc)
			}
		case *object.AttrNode:
			switch op.Type {
			case token.EQ, token.EQV:
				if leftNode.Text() == rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.NE, token.NEV:
				if leftNode.Text() != rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GT, token.GTV:
				if leftNode.Text() > rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GE, token.GEV:
				if leftNode.Text() >= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LT, token.LTV:
				if leftNode.Text() < rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LE, token.LEV:
				if leftNode.Text() <= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.IS:
				return bif.NewBoolean(false)
			case token.DGT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(rightNode, leftNode, doc)
			case token.DLT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(leftNode, rightNode, doc)
			}
		}
	}

	if leftNode, ok := left.(*object.AttrNode); ok {
		switch rightNode := right.(type) {
		case *object.BaseNode:
			switch op.Type {
			case token.EQ, token.EQV:
				if leftNode.Text() == rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.NE, token.NEV:
				if leftNode.Text() != rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GT, token.GTV:
				if leftNode.Text() > rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GE, token.GEV:
				if leftNode.Text() >= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LT, token.LTV:
				if leftNode.Text() < rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LE, token.LEV:
				if leftNode.Text() <= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.IS:
				return bif.NewBoolean(false)
			case token.DGT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(rightNode, leftNode, doc)
			case token.DLT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(leftNode, rightNode, doc)
			}
		case *object.AttrNode:
			switch op.Type {
			case token.EQ, token.EQV:
				if leftNode.Text() == rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.NE, token.NEV:
				if leftNode.Text() != rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GT, token.GTV:
				if leftNode.Text() > rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.GE, token.GEV:
				if leftNode.Text() >= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LT, token.LTV:
				if leftNode.Text() < rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.LE, token.LEV:
				if leftNode.Text() <= rightNode.Text() {
					return bif.NewBoolean(true)
				}
			case token.IS:
				return bif.NewBoolean(leftNode.Tree() == rightNode.Tree() && leftNode.Key() == rightNode.Key())
			case token.DGT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(rightNode, leftNode, doc)
			case token.DLT:
				doc := ctx.Doc.(*object.BaseNode)
				return bif.IsPrecede(leftNode, rightNode, doc)
			}
		}
	}

	return bif.NewBoolean(false)
}
