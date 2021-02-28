package object

// Func represents function type
type Func func(args ...Item) Item

// Type represents Item Type
type Type string

// Item Types
const (
	NilType     Type = "nil"
	ErrorType   Type = "error"
	PholderType Type = "?"
	VarrefType  Type = "$"

	ItemType     Type = "item"
	SequenceType Type = "sequence"
	EmptySeqType Type = "empty-sequence"

	// function
	MapType   Type = "map"
	ArrayType Type = "array"
	FuncType  Type = "function"

	// node
	ErrorNodeType    Type = "0"
	TextNodeType     Type = "1"
	DocumentNodeType Type = "2"
	ElementNodeType  Type = "3"
	CommentNodeType  Type = "4"

	// atomic
	DoubleType  Type = "xs:double"
	DecimalType Type = "xs:decimal"
	IntegerType Type = "xs:integer"
	StringType  Type = "xs:string"
	BooleanType Type = "xs:boolean"

	// abstract
	NodeType      Type = "node"
	NumericType   Type = "xs:numeric"
	AnyAtomicType Type = "xs:anyAtomic"
)
