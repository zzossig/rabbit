package eval

import (
	"fmt"
	"math"
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

// predefined
var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval ..
func Eval(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.XPath:
		return evalXPath(expr, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: expr.Value}
	case *ast.DecimalLiteral:
		return &object.Decimal{Value: expr.Value}
	case *ast.DoubleLiteral:
		return &object.Double{Value: expr.Value}
	case *ast.StringLiteral:
		return &object.String{Value: expr.Value}
	case *ast.Expr:
		return evalExpr(expr, env)
	case *ast.ParenthesizedExpr:
		return evalExpr(expr, env)
	case *ast.EnclosedExpr:
		return evalExpr(expr, env)
	case *ast.Predicate:
		return evalExpr(expr, env)
	case *ast.Identifier:
		return evalIdentifier(expr, env)
	case *ast.InlineFunctionExpr:
		return evalFunctionLiteral(expr, env)
	case *ast.NamedFunctionRef:
		return evalFunctionLiteral(expr, env)
	case *ast.FunctionCall:
		return evalFunctionCall(expr, env)
	case *ast.AdditiveExpr:
		return evalInfixExpr(expr, env)
	case *ast.MultiplicativeExpr:
		return evalInfixExpr(expr, env)
	case *ast.StringConcatExpr:
		return evalInfixExpr(expr, env)
	case *ast.RangeExpr:
		return evalInfixExpr(expr, env)
	case *ast.UnaryExpr:
		return evalPrefixExpr(expr, env)
	case *ast.SquareArrayConstructor:
		return evalArrayExpr(expr, env)
	case *ast.CurlyArrayConstructor:
		return evalArrayExpr(expr, env)
	}

	return nil
}

func evalXPath(expr *ast.XPath, env *object.Env) object.Item {
	xpath := &object.Sequence{}

	for _, e := range expr.Exprs {
		item := Eval(e, env)

		switch item := item.(type) {
		case *object.Sequence:
			xpath.Items = append(xpath.Items, item.Items...)
		default:
			xpath.Items = append(xpath.Items, item)
		}
	}

	return xpath
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(item object.Item) bool {
	if item != nil {
		return item.Type() == object.ErrorType
	}
	return false
}

func evalExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.Expr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.ParenthesizedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.EnclosedExpr:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	case *ast.Predicate:
		seq := &object.Sequence{}
		for _, e := range expr.Exprs {
			item := Eval(e, env)
			seq.Items = append(seq.Items, item)
		}
		return seq
	}
	return nil
}

func evalIdentifier(ident *ast.Identifier, env *object.Env) object.Item {
	if val, ok := env.Get(ident.EQName.Value()); ok {
		return val
	}

	return newError("identifier not found: " + ident.EQName.Value())
}

func evalFunctionLiteral(expr ast.ExprSingle, env *object.Env) object.Item {
	switch expr := expr.(type) {
	case *ast.NamedFunctionRef:
		return &object.FuncNamed{Name: expr.EQName.Value(), Num: expr.IntegerLiteral.Value}
	case *ast.InlineFunctionExpr:
		body := Eval(&expr.FunctionBody, env)
		return &object.FuncInline{Body: body, Params: &expr.ParamList, SType: &expr.SequenceType}
	}
	return nil
}

func evalFunctionCall(expr ast.ExprSingle, env *object.Env) object.Item {
	f := expr.(*ast.FunctionCall)

	builtin, ok := builtins[f.EQName.Value()]
	if !ok {
		return newError("function not found: " + f.EQName.Value())
	}

	enclosedEnv := object.NewEnclosedEnv(env)
	fc := &object.FuncCall{}
	fc.Env = enclosedEnv
	fc.Name = f.EQName.Value()
	fc.Func = &builtin

	var args []object.Item
	pcnt := 0
loop:
	for _, arg := range f.Args {
		switch arg.TypeID {
		case 0:
			break loop
		case 1:
			evaled := Eval(arg.ExprSingle, enclosedEnv)
			args = append(args, evaled)
		case 2:
			args = append(args, &object.Placeholder{})
			pcnt++
		default:
			break loop
		}
	}

	enclosedEnv.Args = append(enclosedEnv.Args, args...)

	if pcnt > 0 {
		return fc
	}

	return builtin(args...)
}

func evalInfixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	var left object.Item
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.AdditiveExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.MultiplicativeExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.StringConcatExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	case *ast.RangeExpr:
		left = Eval(expr.LeftExpr, env)
		right = Eval(expr.RightExpr, env)
		op = expr.Token
	default:
		return newError("%T is not an infix expression\n", expr)
	}

	if isError(left) {
		return left
	}

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalInfixIntInt(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DecimalType:
		return evalInfixIntDecimal(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.DoubleType:
		return evalInfixIntDouble(op, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.StringType:
		return evalInfixIntString(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.IntegerType:
		return evalInfixDecimalInt(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DecimalType:
		return evalInfixDecimalDecimal(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.DoubleType:
		return evalInfixDecimalDouble(op, left, right)
	case left.Type() == object.DecimalType && right.Type() == object.StringType:
		return evalInfixDecimalString(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.IntegerType:
		return evalInfixDoubleInt(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DecimalType:
		return evalInfixDoubleDecimal(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
		return evalInfixDoubleDouble(op, left, right)
	case left.Type() == object.DoubleType && right.Type() == object.StringType:
		return evalInfixDoubleString(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.IntegerType:
		return evalInfixStringInt(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DecimalType:
		return evalInfixStringDecimal(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.DoubleType:
		return evalInfixStringDouble(op, left, right)
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return evalInfixStringString(op, left, right)
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalPrefixExpr(expr ast.ExprSingle, env *object.Env) object.Item {
	var right object.Item
	var op token.Token

	switch expr := expr.(type) {
	case *ast.UnaryExpr:
		right = Eval(expr.ExprSingle, env)
		op = expr.Token
	default:
		return newError("%T is not an prefix expression\n", expr)
	}

	if isError(right) {
		return right
	}

	switch {
	case right.Type() == object.IntegerType:
		return evalPrefixInt(op, right)
	case right.Type() == object.DecimalType:
		return evalPrefixDecimal(op, right)
	case right.Type() == object.DoubleType:
		return evalPrefixDouble(op, right)
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalInfixIntInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: float64(leftVal) / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / float64(rightVal))}
	case token.MOD:
		return &object.Integer{Value: leftVal % rightVal}
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	case token.TO:
		seq := &object.Sequence{}
		for i := leftVal; i <= rightVal; i++ {
			seq.Items = append(seq.Items, &object.Integer{Value: i})
		}
		return seq
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: float64(leftVal) + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: float64(leftVal) - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: float64(leftVal) * rightVal}
	case token.DIV:
		return &object.Decimal{Value: float64(leftVal) / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(float64(leftVal), rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: float64(leftVal) + rightVal}
	case token.MINUS:
		return &object.Double{Value: float64(leftVal) - rightVal}
	case token.ASTERISK:
		return &object.Double{Value: float64(leftVal) * rightVal}
	case token.DIV:
		return &object.Double{Value: float64(leftVal) / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(float64(leftVal) / rightVal)}
	case token.MOD:
		return &object.Double{Value: math.Mod(float64(leftVal), rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixIntString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatInt(int64(leftVal), 10)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + float64(rightVal)}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - float64(rightVal)}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * float64(rightVal)}
	case token.DIV:
		return &object.Decimal{Value: leftVal / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / float64(rightVal))}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, float64(rightVal))}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDecimalString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Decimal).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: leftVal + float64(rightVal)}
	case token.MINUS:
		return &object.Double{Value: leftVal - float64(rightVal)}
	case token.ASTERISK:
		return &object.Double{Value: leftVal * float64(rightVal)}
	case token.DIV:
		return &object.Double{Value: leftVal / float64(rightVal)}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / float64(rightVal))}
	case token.MOD:
		return &object.Double{Value: math.Mod(leftVal, float64(rightVal))}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Decimal{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Decimal{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Decimal{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Decimal{Value: math.Mod(leftVal, rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Double{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Double{Value: leftVal * rightVal}
	case token.DIV:
		return &object.Double{Value: leftVal / rightVal}
	case token.IDIV:
		return &object.Integer{Value: int(leftVal / rightVal)}
	case token.MOD:
		return &object.Double{Value: math.Mod(leftVal, rightVal)}
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixDoubleString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.Double).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		leftVal := strconv.FormatFloat(leftVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringInt(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatInt(int64(rightVal), 10)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringDecimal(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringDouble(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Double).Value

	switch op.Type {
	case token.DVBAR:
		rightVal := strconv.FormatFloat(rightVal, 'f', -1, 64)
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalInfixStringString(op token.Token, left object.Item, right object.Item) object.Item {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op.Type {
	case token.DVBAR:
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("The operator '%s' is not defined for operands of type %s and %s\n", op.Literal, left.Type(), right.Type())
	}
}

func evalPrefixInt(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Integer).Value

	switch op.Type {
	case token.PLUS:
		return &object.Integer{Value: rightVal}
	case token.MINUS:
		return &object.Integer{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDecimal(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Decimal{Value: rightVal}
	case token.MINUS:
		return &object.Decimal{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
}

func evalPrefixDouble(op token.Token, right object.Item) object.Item {
	rightVal := right.(*object.Decimal).Value

	switch op.Type {
	case token.PLUS:
		return &object.Double{Value: rightVal}
	case token.MINUS:
		return &object.Double{Value: -1 * rightVal}
	default:
		return newError("The operator '%s' is not defined for operand of type %s\n", op.Literal, right.Type())
	}
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
