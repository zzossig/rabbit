package ast

import "github.com/zzossig/xpath/token"

// ItemType ::= KindTest | ("item" "(" ")") | FunctionTest | MapTest | ArrayTest | AtomicOrUnionType | ParenthesizedItemType
type ItemType interface {
	itemType()
}

// KindTest ::= DocumentTest | ElementTest | AttributeTest | SchemaElementTest | SchemaAttributeTest | PITest | CommentTest | TextTest | NamespaceNodeTest | AnyKindTest
type KindTest interface {
	ItemType
	kindTest()
}

// FunctionTest ::= AnyFunctionTest | TypedFunctionTest
type FunctionTest interface {
	ItemType
	functionTest()
}

// MapTest ::= AnyMapTest | TypedMapTest
type MapTest interface {
	ItemType
	mapTest()
}

// ArrayTest ::= AnyArrayTest | TypedArrayTest
type ArrayTest interface {
	ItemType
	arrayTest()
}

// SequenceType ::= ("empty-sequence" "(" ")") | (ItemType OccurrenceIndicator?)
type SequenceType struct {
	ItemType
	OccurrenceIndicator
}

// OccurrenceIndicator ::= "?" | "*" | "+"
type OccurrenceIndicator struct {
	Token token.Token
}

// ParenthesizedItemType ::= "(" ItemType ")"
type ParenthesizedItemType struct {
	ItemType
}

// DocumentTest ::= "document-node" "(" (ElementTest | SchemaElementTest)? ")"
type DocumentTest struct {
	KindTest
}

func (dt *DocumentTest) itemType() {}
func (dt *DocumentTest) kindTest() {}

// ElementTest ::= "element" "(" (ElementNameOrWildcard ("," TypeName "?"?)?)? ")"
type ElementTest struct {
	ElementNameOrWildcard
	TypeName
}

func (et *ElementTest) itemType() {}
func (et *ElementTest) kindTest() {}

// AttributeTest ::= "attribute" "(" (AttribNameOrWildcard ("," TypeName)?)? ")"
type AttributeTest struct {
	AttribNameOrWildcard
	TypeName
}

func (at *AttributeTest) itemType() {}
func (at *AttributeTest) kindTest() {}

// SchemaElementTest ::= "schema-element" "(" ElementDeclaration ")"
type SchemaElementTest struct {
	ElementDeclaration
}

func (set *SchemaElementTest) itemType() {}
func (set *SchemaElementTest) kindTest() {}

// SchemaAttributeTest ::= "schema-attribute" "(" AttributeDeclaration ")"
type SchemaAttributeTest struct {
	AttributeDeclaration
}

func (sat *SchemaAttributeTest) itemType() {}
func (sat *SchemaAttributeTest) kindTest() {}

// PITest ::= "processing-instruction" "(" (NCName | StringLiteral)? ")"
type PITest struct {
	Name NCName
	StringLiteral
}

func (pit *PITest) itemType() {}
func (pit *PITest) kindTest() {}

// CommentTest ::= "comment" "(" ")"
type CommentTest struct{}

func (ct *CommentTest) itemType() {}
func (ct *CommentTest) kindTest() {}

// NamespaceNodeTest ::= "namespace-node" "(" ")"
type NamespaceNodeTest struct{}

func (nnt *NamespaceNodeTest) itemType() {}
func (nnt *NamespaceNodeTest) kindTest() {}

// TextTest ::= "text" "(" ")"
type TextTest struct{}

func (tt *TextTest) itemType() {}
func (tt *TextTest) kindTest() {}

// AnyKindTest ::= "node" "(" ")"
type AnyKindTest struct{}

func (akt *AnyKindTest) itemType() {}
func (akt *AnyKindTest) kindTest() {}

// AtomicOrUnionType ::= EQName
type AtomicOrUnionType = EQName

// ElementDeclaration ::= ElementName
type ElementDeclaration = ElementName

// ElementName ::= EQName
type ElementName = EQName

// ElementNameOrWildcard ::= ElementName | "*"
type ElementNameOrWildcard struct {
	ElementName
}

// AttributeDeclaration ::= AttributeName
type AttributeDeclaration = AttributeName

// AttribNameOrWildcard ::= AttributeName | "*"
type AttribNameOrWildcard struct {
	AttributeName
}

// AttributeName ::= EQName
type AttributeName = EQName

// TypeName ::= EQName
type TypeName = EQName

// AnyFunctionTest ::= "function" "(" "*" ")"
type AnyFunctionTest struct{}

func (aft *AnyFunctionTest) itemType()     {}
func (aft *AnyFunctionTest) functionTest() {}

// TypedFunctionTest ::= "function" "(" (SequenceType ("," SequenceType)*)? ")" "as" SequenceType
type TypedFunctionTest struct {
	ParamSTypes []SequenceType
	AsSType     SequenceType
}

func (tft *TypedFunctionTest) itemType()     {}
func (tft *TypedFunctionTest) functionTest() {}

// AnyMapTest ::= "map" "(" "*" ")"
type AnyMapTest struct{}

func (amt *AnyMapTest) itemType() {}
func (amt *AnyMapTest) mapTest()  {}

// TypedMapTest ::= "map" "(" AtomicOrUnionType "," SequenceType ")"
type TypedMapTest struct {
	AtomicOrUnionType
	SequenceType
}

func (tmt *TypedMapTest) itemType() {}
func (tmt *TypedMapTest) mapTest()  {}

// AnyArrayTest ::= "array" "(" "*" ")"
type AnyArrayTest struct{}

func (aat *AnyArrayTest) itemType()  {}
func (aat *AnyArrayTest) arrayTest() {}

// TypedArrayTest ::= "array" "(" SequenceType ")"
type TypedArrayTest struct{}

func (tat *TypedArrayTest) itemType()  {}
func (tat *TypedArrayTest) arrayTest() {}
