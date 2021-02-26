package eval

import (
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
	case bif.IsNumeric(left) && bif.IsNumeric(right):
		if op.Type == token.EQ || op.Type == token.EQV {
			builtin := bif.Builtins["op:numeric-equal"]
			return builtin(left, right)
		} else if op.Type == token.NE || op.Type == token.NEV {
			builtin := bif.Builtins["op:numeric-equal"]
			b := builtin(left, right)
			boolean := b.(*object.Boolean)

			return bif.NewBoolean(!boolean.Value())
		} else if op.Type == token.LT || op.Type == token.LTV {
			builtin := bif.Builtins["op:numeric-less-than"]
			return builtin(left, right)
		} else if op.Type == token.GT || op.Type == token.GTV {
			builtin := bif.Builtins["op:numeric-greater-than"]
			return builtin(left, right)
		} else if op.Type == token.LE || op.Type == token.LEV {
			builtin := bif.Builtins["op:numeric-less-than"]
			b := builtin(left, right)
			boolean := b.(*object.Boolean)

			if !boolean.Value() {
				builtin = bif.Builtins["op:numeric-equal"]
				return builtin(left, right)
			}
			return boolean
		} else if op.Type == token.GE || op.Type == token.GEV {
			builtin := bif.Builtins["op:numeric-greater-than"]
			b := builtin(left, right)
			boolean := b.(*object.Boolean)

			if !boolean.Value() {
				builtin = bif.Builtins["op:numeric-equal"]
				return builtin(left, right)
			}
			return boolean
		}
	case bif.IsNumeric(left) && bif.IsSeq(right):
		return compNumberSeq(op, left, right)
	case bif.IsNumeric(left) && right.Type() == object.ArrayType:
		return compNumberArray(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return compStringString(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.ArrayType:
		return compStringArray(op, left, right)
	case left.Type() == object.StringType && bif.IsSeq(right):
		return compStringSeq(op, left, right)
	case left.Type() == object.ArrayType && bif.IsNumeric(right):
		return compArrayNumber(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.StringType:
		return compArrayString(op, left, right)
	case left.Type() == object.ArrayType && right.Type() == object.ArrayType:
		return compArrayArray(op, left, right)
	case left.Type() == object.ArrayType && bif.IsSeq(right):
		return compArraySeq(op, left, right)
	case bif.IsSeq(left) && bif.IsNumeric(right):
		return compSeqNumber(op, left, right)
	case bif.IsSeq(left) && right.Type() == object.StringType:
		return compSeqString(op, left, right)
	case bif.IsSeq(left) && right.Type() == object.ArrayType:
		return compSeqArray(op, left, right)
	case bif.IsSeq(left) && bif.IsSeq(right):
		return compSeqSeq(op, left, right)
	case left.Type() == object.BooleanType && right.Type() == object.BooleanType:
		return compBool(op, left, right)
	}

	return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
}

func compNumberArray(op token.Token, left object.Item, right object.Item) object.Item {
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compNumberSeq(op token.Token, left object.Item, right object.Item) object.Item {
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e := bif.IsEQ(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e := bif.IsNE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e := bif.IsLT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e := bif.IsLE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e := bif.IsGT(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e := bif.IsGE(left, item)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value()
	rightVal := right.(*object.String).Value()

	switch op.Type {
	case token.EQ:
		return bif.NewBoolean(leftVal == rightVal)
	case token.NE:
		return bif.NewBoolean(leftVal != rightVal)
	case token.LT:
		return bif.NewBoolean(leftVal < rightVal)
	case token.LE:
		return bif.NewBoolean(leftVal <= rightVal)
	case token.GT:
		return bif.NewBoolean(leftVal > rightVal)
	case token.GE:
		return bif.NewBoolean(leftVal >= rightVal)
	case token.EQV:
		return bif.NewBoolean(leftVal == rightVal)
	case token.NEV:
		return bif.NewBoolean(leftVal != rightVal)
	case token.LTV:
		return bif.NewBoolean(leftVal < rightVal)
	case token.LEV:
		return bif.NewBoolean(leftVal <= rightVal)
	case token.GTV:
		return bif.NewBoolean(leftVal > rightVal)
	case token.GEV:
		return bif.NewBoolean(leftVal >= rightVal)
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringArray(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() == e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() != e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() < e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() <= e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() > e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() >= e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compStringSeq(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() == e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() != e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() < e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() <= e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() > e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range rightVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if leftVal.Value() >= e.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() == rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() != rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() < rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() <= rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() > rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() >= rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArraySeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compArrayArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqNumber(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e := bif.IsEQ(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e := bif.IsNE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e := bif.IsLT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e := bif.IsLE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e := bif.IsGT(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e := bif.IsGE(item, right)
			if bif.IsError(e) {
				return e
			}
			bl := e.(*object.Boolean)
			if bl.Value() == true {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqString(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.String)

	switch op.Type {
	case token.EQ:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() == rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.NE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() != rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() < rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.LE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() <= rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GT:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() > rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	case token.GE:
		for _, item := range leftVal.Items {
			e, ok := item.(*object.String)
			if !ok {
				return bif.NewError("Types %s and %s are not comparable.", leftVal.Type(), item.Type())
			}
			if e.Value() >= rightVal.Value() {
				return object.TRUE
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqArray(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Array)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compSeqSeq(op token.Token, left, right object.Item) object.Item {
	leftVal := left.(*object.Sequence)
	rightVal := right.(*object.Sequence)

	switch op.Type {
	case token.EQ:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsEQ(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.NE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsNE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.LE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsLE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GT:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGT(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	case token.GE:
		for _, li := range leftVal.Items {
			for _, ri := range rightVal.Items {
				e := bif.IsGE(li, ri)
				if bif.IsError(e) {
					return e
				}
				bl := e.(*object.Boolean)
				if bl.Value() == true {
					return object.TRUE
				}
			}
		}
		return object.FALSE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func compBool(op token.Token, left, right object.Item) object.Item {
	leftVal, ok := left.(*object.Boolean)
	if !ok {
		return bif.NewError("[XPTY0004] Types %s and %s are not comparable.", left.Type(), right.Type())
	}

	rightVal, ok := right.(*object.Boolean)
	if !ok {
		return bif.NewError("[XPTY0004] Types %s and %s are not comparable.", left.Type(), right.Type())
	}

	switch op.Type {
	case token.EQ:
		fallthrough
	case token.EQV:
		return bif.NewBoolean(leftVal.Value() == rightVal.Value())
	case token.NE:
		fallthrough
	case token.NEV:
		return bif.NewBoolean(leftVal.Value() != rightVal.Value())
	case token.GT:
		fallthrough
	case token.GTV:
		if leftVal.Value() && !rightVal.Value() {
			return object.TRUE
		}
		return object.FALSE
	case token.GE:
		fallthrough
	case token.GEV:
		if !leftVal.Value() && rightVal.Value() {
			return object.FALSE
		}
		return object.TRUE
	case token.LT:
		fallthrough
	case token.LTV:
		if !leftVal.Value() && rightVal.Value() {
			return object.TRUE
		}
		return object.FALSE
	case token.LE:
		fallthrough
	case token.LEV:
		if leftVal.Value() && !rightVal.Value() {
			return object.FALSE
		}
		return object.TRUE
	default:
		return bif.NewError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}
