package object

import "github.com/zzossig/xpath/ast"

// Func represents function type
type Func func(ctx *Context, args ...Item) Item

// Ev represents Eval function
type Ev func(expr ast.ExprSingle, ctx *Context) Item

// Type represents Item Type
type Type string

// Item Types
const (
	NilType      Type = "nil"
	ErrorType    Type = "error"
	PholderType  Type = "?"
	VarrefType   Type = "$"
	SequenceType Type = "sequence"

	// function
	MapType   Type = "map"
	ArrayType Type = "array"
	FuncType  Type = "function"

	// node
	ErrorNodeType     Type = "0"
	TextNodeType      Type = "1"
	DocumentNodeType  Type = "2"
	ElementNodeType   Type = "3"
	CommentNodeType   Type = "4"
	DoctypeNodeType   Type = "5"
	RawNodeType       Type = "6"
	AttributeNodeType Type = "7"

	// atomic
	DoubleType  Type = "xs:double"
	DecimalType Type = "xs:decimal"
	IntegerType Type = "xs:integer"
	StringType  Type = "xs:string"
	BooleanType Type = "xs:boolean"
)
