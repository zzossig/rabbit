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
	NodeType    Type = "node"
	DocType     Type = "document"
	ElemType    Type = "element"
	AttrType    Type = "attribute"
	PIType      Type = "processing-instruction"
	CommentType Type = "comment"
	NSNodeType  Type = "namespace-node"
	TextType    Type = "text"

	// atomic
	DoubleType  Type = "xs:double"
	DecimalType Type = "xs:decimal"
	IntegerType Type = "xs:integer"
	StringType  Type = "xs:string"
	BooleanType Type = "xs:boolean"

	// abstract
	NumericType   Type = "xs:numeric"
	AnyAtomicType Type = "xs:anyAtomic"
)
