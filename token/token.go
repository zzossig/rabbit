package token

import "strings"

// Token represents lexical token
type Token struct {
	Type    Type
	Literal string
}

// Type represents Token Type
type Type string

// Token Types
const (
	ILLEGAL    Type = "-1"
	EOF        Type = "0"
	IDENT      Type = "identifier"
	STRING     Type = "string"
	INT        Type = "integer"
	DOUBLE     Type = "double"
	DECIMAL    Type = "decimal"
	SQUOTE     Type = "'"
	DQUOTE     Type = "\""
	PLUS       Type = "+"   // binary
	UPLUS      Type = "(+)" // unary
	MINUS      Type = "-"   // binary
	UMINUS     Type = "(-)" // unary
	ASTERISK   Type = "*"
	SLASH      Type = "/"
	DSLASH     Type = "//"
	LT         Type = "<"
	DLT        Type = "<<"
	LE         Type = "<="
	GT         Type = ">"
	DGT        Type = ">>"
	GE         Type = ">="
	EQ         Type = "="
	ARROW      Type = "=>"
	BANG       Type = "!"
	NE         Type = "!="
	AT         Type = "@"
	DOLLAR     Type = "$"
	HASH       Type = "#"
	DOT        Type = "."
	DDOT       Type = ".."
	COMMA      Type = ","
	COLON      Type = ":"
	DCOLON     Type = "::"
	ASSIGN     Type = ":="
	SCOLON     Type = ";"
	LPAREN     Type = "("
	RPAREN     Type = ")"
	LBRACE     Type = "{"
	RBRACE     Type = "}"
	LBRACKET   Type = "["
	RBRACKET   Type = "]"
	VBAR       Type = "|"
	DVBAR      Type = "||"
	QUESTION   Type = "?"
	UQUESTION  Type = "(?)" // unary
	LET        Type = "let"
	THEN       Type = "then"
	ELSE       Type = "else"
	RETURN     Type = "return"
	FOR        Type = "for"
	IN         Type = "in"
	SOME       Type = "some"
	EVERY      Type = "every"
	OF         Type = "of"
	AS         Type = "as"
	IS         Type = "is"
	EQV        Type = "eq"
	NEV        Type = "ne"
	LTV        Type = "lt"
	LEV        Type = "le"
	GTV        Type = "gt"
	GEV        Type = "ge"
	TO         Type = "to"
	OR         Type = "or"
	AND        Type = "and"
	DIV        Type = "div"
	IDIV       Type = "idiv"
	MOD        Type = "mod"
	UNION      Type = "union"
	INTERSECT  Type = "intersect"
	EXCEPT     Type = "except"
	INSTANCE   Type = "instance"
	TREAT      Type = "treat"
	CAST       Type = "cast"
	CASTABLE   Type = "castable"
	SATISFIES  Type = "satisfies"
	INSTANCEOF Type = "instance of"
	TREATAS    Type = "treat as"
	CASTABLEAS Type = "castable as"
	CASTAS     Type = "cast as"

	//Reserved Function
	ARRAY      Type = "array"
	ATTRIBUTE  Type = "attribute"
	COMMENT    Type = "comment"
	DNODE      Type = "docunemt-node"
	ELEMENT    Type = "element"
	ES         Type = "empty-sequence"
	FUNCTION   Type = "function"
	IF         Type = "if"
	ITEM       Type = "item"
	MAP        Type = "map"
	NSNODE     Type = "namespace-node"
	NODE       Type = "node"
	PI         Type = "processing-instruction"
	SA         Type = "schema-attribute"
	SE         Type = "schema-element"
	SWITCH     Type = "switch"
	TEXT       Type = "text"
	TYPESWITCH Type = "typeswitch"
)

var keywords = map[string]Type{
	"let":       LET,
	"then":      THEN,
	"else":      ELSE,
	"return":    RETURN,
	"every":     EVERY,
	"some":      SOME,
	"for":       FOR,
	"in":        IN,
	"of":        OF,
	"as":        AS,
	"is":        IS,
	"eq":        EQV,
	"ne":        NEV,
	"lt":        LTV,
	"le":        LEV,
	"gt":        GTV,
	"ge":        GEV,
	"to":        TO,
	"or":        OR,
	"and":       AND,
	"div":       DIV,
	"idiv":      IDIV,
	"mod":       MOD,
	"union":     UNION,
	"intersect": INTERSECT,
	"except":    EXCEPT,
	"instance":  INSTANCE,
	"treat":     TREAT,
	"cast":      CAST,
	"castable":  CASTABLE,
	"satisfies": SATISFIES,

	"array":                  ARRAY,
	"attribute":              ATTRIBUTE,
	"comment":                COMMENT,
	"docunemt-node":          DNODE,
	"element":                ELEMENT,
	"empty-sequence":         ES,
	"function":               FUNCTION,
	"if":                     IF,
	"item":                   ITEM,
	"map":                    MAP,
	"namespace-node":         NSNODE,
	"node":                   NODE,
	"processing-instruction": PI,
	"schema-attribute":       SA,
	"schema-element":         SE,
	"switch":                 SWITCH,
	"text":                   TEXT,
	"typeswitch":             TYPESWITCH,
}

// LookupIdent returns identifier token type
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// LookupNumber returns number token type
func LookupNumber(number string) Type {
	if strings.Contains(number, "e") || strings.Contains(number, "E") {
		return DOUBLE
	} else if strings.Contains(number, ".") {
		return DECIMAL
	}
	return INT
}
