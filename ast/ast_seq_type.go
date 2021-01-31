package ast

import (
	"strings"

	"github.com/zzossig/xpath/token"
)

// ItemType ::= KindTest | ItemTest | FunctionTest | MapTest | ArrayTest | AtomicOrUnionType | ParenthesizedItemType
// TypeID ::= 	1				 | 2				| 3						 | 4			 | 5				 | 6								 | 7
type ItemType struct {
	NodeTest
	TypeID byte
}

func (it *ItemType) nodeTest() {}
func (it *ItemType) String() string {
	return it.NodeTest.String()
}

// KindTest ::= DocumentTest | ElementTest | AttributeTest | SchemaElementTest | SchemaAttributeTest | PITest | CommentTest | TextTest | NamespaceNodeTest | AnyKindTest
// TypeID ::=		1						 | 2					 | 3						 | 4								 | 5									 | 6			| 7						| 8				 | 9								 | 10
type KindTest struct {
	NodeTest
	TypeID byte
}

func (kt *KindTest) nodeTest() {}
func (kt *KindTest) String() string {
	return kt.NodeTest.String()
}

// ItemTest ::= ("item" "(" ")")
type ItemTest struct{}

func (it *ItemTest) nodeTest() {}
func (it *ItemTest) String() string {
	return "item()"
}

// FunctionTest ::= AnyFunctionTest | TypedFunctionTest
// TypeID ::= 			1								| 2
type FunctionTest struct {
	NodeTest
	TypeID byte
}

func (ft *FunctionTest) nodeTest() {}
func (ft *FunctionTest) String() string {
	return ft.NodeTest.String()
}

// MapTest ::= AnyMapTest | TypedMapTest
// TypeID ::=  1					| 2
type MapTest struct {
	NodeTest
	TypeID byte
}

func (mt *MapTest) nodeTest() {}
func (mt *MapTest) String() string {
	return mt.NodeTest.String()
}

// ArrayTest ::= AnyArrayTest | TypedArrayTest
// TypeID ::= 	 1 						| 2
type ArrayTest struct {
	NodeTest
	TypeID byte
}

func (at *ArrayTest) nodeTest() {}
func (at *ArrayTest) String() string {
	return at.NodeTest.String()
}

// SequenceType ::= ("empty-sequence" "(" ")") | (ItemType OccurrenceIndicator?)
// TypeID ::= 			1													 | 2
type SequenceType struct {
	NodeTest
	OccurrenceIndicator
	TypeID byte
}

func (st *SequenceType) String() string {
	var sb strings.Builder

	switch st.TypeID {
	case 1:
		sb.WriteString("empty-sequence()")
	case 2:
		sb.WriteString(st.NodeTest.String())
		sb.WriteString(st.OccurrenceIndicator.String())
	default:
		sb.WriteString("")
	}

	return sb.String()
}

// OccurrenceIndicator ::= "?" | "*" | "+"
type OccurrenceIndicator struct {
	Token token.Token
}

func (oi *OccurrenceIndicator) String() string {
	return oi.Token.Literal
}

// ParenthesizedItemType ::= "(" ItemType ")"
type ParenthesizedItemType struct {
	NodeTest
}

func (pit *ParenthesizedItemType) nodeTest() {}
func (pit *ParenthesizedItemType) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pit.NodeTest.String())
	sb.WriteString(")")

	return sb.String()
}

// DocumentTest ::= "document-node" "(" (ElementTest | SchemaElementTest)? ")"
type DocumentTest struct {
	NodeTest
}

func (dt *DocumentTest) nodeTest() {}
func (dt *DocumentTest) String() string {
	var sb strings.Builder

	sb.WriteString("document-node(")
	if dt.NodeTest != nil {
		sb.WriteString(dt.NodeTest.String())
	}
	sb.WriteString(")")

	return sb.String()
}

// ElementTest ::= "element" "(" (ElementNameOrWildcard ("," TypeName "?"?)?)? ")"
type ElementTest struct {
	ElementNameOrWildcard
	TypeName
	Token token.Token // token.QUESTION
}

func (et *ElementTest) nodeTest() {}
func (et *ElementTest) String() string {
	var sb strings.Builder

	sb.WriteString("element(")
	if et.ElementNameOrWildcard.String() != "" {
		sb.WriteString(et.ElementNameOrWildcard.String())
	}
	if et.TypeName.Value() != "" {
		sb.WriteString(", ")
		sb.WriteString(et.TypeName.Value())
		if et.Token.Literal != "" {
			sb.WriteString(et.Token.Literal)
		}
	}
	sb.WriteString(")")

	return sb.String()
}

// AttributeTest ::= "attribute" "(" (AttribNameOrWildcard ("," TypeName)?)? ")"
type AttributeTest struct {
	AttribNameOrWildcard
	TypeName
}

func (at *AttributeTest) nodeTest() {}
func (at *AttributeTest) String() string {
	var sb strings.Builder

	sb.WriteString("attribute(")
	if at.AttributeName.Value() != "" {
		sb.WriteString(at.AttributeName.Value())
	} else if at.WC != "" {
		sb.WriteString(at.WC)
	}
	if at.TypeName.Value() != "" {
		sb.WriteString(", ")
		sb.WriteString(at.TypeName.Value())
	}
	sb.WriteString(")")

	return sb.String()
}

// SchemaElementTest ::= "schema-element" "(" ElementDeclaration ")"
type SchemaElementTest struct {
	ElementDeclaration
}

func (set *SchemaElementTest) nodeTest() {}
func (set *SchemaElementTest) String() string {
	var sb strings.Builder

	sb.WriteString("schema-element(")
	sb.WriteString(set.ElementDeclaration.Value())
	sb.WriteString(")")

	return sb.String()
}

// SchemaAttributeTest ::= "schema-attribute" "(" AttributeDeclaration ")"
type SchemaAttributeTest struct {
	AttributeDeclaration
}

func (sat *SchemaAttributeTest) nodeTest() {}
func (sat *SchemaAttributeTest) String() string {
	var sb strings.Builder

	sb.WriteString("schema-attribute(")
	sb.WriteString(sat.AttributeDeclaration.Value())
	sb.WriteString(")")

	return sb.String()
}

// PITest ::= "processing-instruction" "(" (NCName | StringLiteral)? ")"
type PITest struct {
	NCName
	StringLiteral
}

func (pit *PITest) nodeTest() {}
func (pit *PITest) String() string {
	var sb strings.Builder

	sb.WriteString("processing-instruction(")
	if pit.NCName.Value() != "" {
		sb.WriteString(pit.NCName.Value())
	} else if pit.StringLiteral.String() != "" {
		sb.WriteString(pit.StringLiteral.String())
	}
	sb.WriteString(")")

	return sb.String()
}

// CommentTest ::= "comment" "(" ")"
type CommentTest struct{}

func (ct *CommentTest) nodeTest() {}
func (ct *CommentTest) String() string {
	return "comment()"
}

// NamespaceNodeTest ::= "namespace-node" "(" ")"
type NamespaceNodeTest struct{}

func (nnt *NamespaceNodeTest) nodeTest() {}
func (nnt *NamespaceNodeTest) String() string {
	return "namespace-node()"
}

// TextTest ::= "text" "(" ")"
type TextTest struct{}

func (tt *TextTest) nodeTest() {}
func (tt *TextTest) String() string {
	return "text()"
}

// AnyKindTest ::= "node" "(" ")"
type AnyKindTest struct{}

func (akt *AnyKindTest) nodeTest() {}
func (akt *AnyKindTest) String() string {
	return "node()"
}

// AtomicOrUnionType ::= EQName
type AtomicOrUnionType struct {
	EQName
}

func (aout *AtomicOrUnionType) nodeTest() {}
func (aout *AtomicOrUnionType) String() string {
	return aout.EQName.Value()
}

// ElementDeclaration ::= ElementName
type ElementDeclaration = ElementName

// ElementName ::= EQName
type ElementName = EQName

// ElementNameOrWildcard ::= ElementName | "*"
type ElementNameOrWildcard struct {
	ElementName
	WC string // "*"
}

func (eow *ElementNameOrWildcard) String() string {
	if eow.WC != "" {
		return eow.WC
	}
	return eow.ElementName.Value()
}

// AttributeDeclaration ::= AttributeName
type AttributeDeclaration = AttributeName

// AttribNameOrWildcard ::= AttributeName | "*"
type AttribNameOrWildcard struct {
	AttributeName
	WC string
}

func (aow *AttribNameOrWildcard) String() string {
	if aow.WC != "" {
		return aow.WC
	}
	return aow.AttributeName.Value()
}

// AttributeName ::= EQName
type AttributeName = EQName

// TypeName ::= EQName
type TypeName = EQName

// AnyFunctionTest ::= "function" "(" "*" ")"
type AnyFunctionTest struct{}

func (aft *AnyFunctionTest) nodeTest() {}
func (aft *AnyFunctionTest) String() string {
	return "function(*)"
}

// TypedFunctionTest ::= "function" "(" (SequenceType ("," SequenceType)*)? ")" "as" SequenceType
type TypedFunctionTest struct {
	ParamSTypes []SequenceType
	AsSType     SequenceType
}

func (tft *TypedFunctionTest) nodeTest() {}
func (tft *TypedFunctionTest) String() string {
	var sb strings.Builder

	sb.WriteString("function(")
	for i, param := range tft.ParamSTypes {
		sb.WriteString(param.String())
		if i < len(tft.ParamSTypes)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString("as")
	sb.WriteString(" ")
	sb.WriteString(tft.AsSType.String())

	return sb.String()
}

// AnyMapTest ::= "map" "(" "*" ")"
type AnyMapTest struct{}

func (amt *AnyMapTest) nodeTest() {}
func (amt *AnyMapTest) String() string {
	return "map(*)"
}

// TypedMapTest ::= "map" "(" AtomicOrUnionType "," SequenceType ")"
type TypedMapTest struct {
	AtomicOrUnionType
	SequenceType
}

func (tmt *TypedMapTest) nodeTest() {}
func (tmt *TypedMapTest) String() string {
	var sb strings.Builder

	sb.WriteString("map(")
	sb.WriteString(tmt.AtomicOrUnionType.Value())
	sb.WriteString(tmt.SequenceType.String())
	sb.WriteString(")")

	return sb.String()
}

// AnyArrayTest ::= "array" "(" "*" ")"
type AnyArrayTest struct{}

func (aat *AnyArrayTest) nodeTest() {}
func (aat *AnyArrayTest) String() string {
	return "array(*)"
}

// TypedArrayTest ::= "array" "(" SequenceType ")"
type TypedArrayTest struct {
	SequenceType
}

func (tat *TypedArrayTest) nodeTest() {}
func (tat *TypedArrayTest) String() string {
	var sb strings.Builder

	sb.WriteString("array(")
	sb.WriteString(tat.SequenceType.String())
	sb.WriteString(")")

	return sb.String()
}

// InstanceofExpr ::= TreatExpr ( "instance" "of" SequenceType )?
type InstanceofExpr struct {
	ExprSingle
	SequenceType
}

func (ie *InstanceofExpr) exprSingle() {}
func (ie *InstanceofExpr) String() string {
	var sb strings.Builder

	sb.WriteString(ie.ExprSingle.String())
	if ie.SequenceType.TypeID != 0 {
		sb.WriteString(" ")
		sb.WriteString("instance of")
		sb.WriteString(" ")
		sb.WriteString(ie.SequenceType.String())
	}

	return sb.String()
}

// CastExpr ::= ArrowExpr ( "cast" "as" SingleType )?
type CastExpr struct {
	ExprSingle
	SingleType
}

func (ce *CastExpr) exprSingle() {}
func (ce *CastExpr) String() string {
	var sb strings.Builder

	sb.WriteString(ce.ExprSingle.String())
	if ce.SingleType.SimpleTypeName.Value() != "" {
		sb.WriteString(" ")
		sb.WriteString("cast as")
		sb.WriteString(" ")
		sb.WriteString(ce.SingleType.SimpleTypeName.Value())
	}

	return sb.String()
}

// SingleType ::= SimpleTypeName "?"?
type SingleType struct {
	SimpleTypeName
	Token token.Token
}

func (st *SingleType) String() string {
	var sb strings.Builder

	sb.WriteString(st.SimpleTypeName.Value())
	sb.WriteString(st.Token.Literal)

	return sb.String()
}

// CastableExpr ::= CastExpr ( "castable" "as" SingleType )?
type CastableExpr struct {
	ExprSingle
	SingleType
}

func (ce *CastableExpr) exprSingle() {}
func (ce *CastableExpr) String() string {
	var sb strings.Builder

	sb.WriteString(ce.ExprSingle.String())
	if ce.SingleType.SimpleTypeName.Value() != "" {
		sb.WriteString(" ")
		sb.WriteString("castable as")
		sb.WriteString(" ")
		sb.WriteString(ce.SingleType.SimpleTypeName.Value())
	}

	return sb.String()
}

// TreatExpr ::= CastableExpr ( "treat" "as" SequenceType )?
type TreatExpr struct {
	ExprSingle
	SequenceType
}

func (te *TreatExpr) exprSingle() {}
func (te *TreatExpr) String() string {
	var sb strings.Builder

	sb.WriteString(te.ExprSingle.String())
	if te.SequenceType.TypeID != 0 {
		sb.WriteString(" ")
		sb.WriteString("treat as")
		sb.WriteString(" ")
		sb.WriteString(te.SequenceType.String())
	}

	return sb.String()
}
