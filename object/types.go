package object

// Func represents function type
type Func func(args ...Item) Item

// Type represents Item Type
type Type string

// Value represents any value
type Value interface{}

// Item Types
const (
	NilType     Type = "nil"
	ErrorType   Type = "error"
	PholderType Type = "?"

	NodeType     Type = "node"
	MapType      Type = "map"
	ArrayType    Type = "array"
	SequenceType Type = "sequence"

	FuncCallType   Type = "functionC"
	FuncNamedType  Type = "functionN"
	FuncInlineType Type = "functionI"

	ByteType          Type = "xs:byte"
	ShortType         Type = "xs:short"
	IntType           Type = "xs:int"
	LongType          Type = "xs:long"
	IntegerType       Type = "xs:integer"
	DecimalType       Type = "xs:decimal"
	DoubleType        Type = "xs:double"
	BooleanType       Type = "xs:boolean"
	StringType        Type = "xs:string"
	UntypedAtomicType Type = "xs:untypedAtomic"
)
