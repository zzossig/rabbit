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
	PLUS       Type = "+"
	MINUS      Type = "-"
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
	XSCHEMA    Type = "xs"
	XFUNC      Type = "fn"
	XMAP       Type = "map"
	XARRAY     Type = "array"
	XMATH      Type = "math"
	FUNC       Type = "function"
	LET        Type = "let"
	IF         Type = "if"
	THEN       Type = "then"
	ELSE       Type = "else"
	RETURN     Type = "return"
	FOR        Type = "for"
	SOME       Type = "some"
	EVERY      Type = "every"
	IS         Type = "is"
	TO         Type = "to"
	OR         Type = "or"
	AND        Type = "and"
	DIV        Type = "div"
	IDIV       Type = "idiv"
	MOD        Type = "mod"
	UNION      Type = "union"
	INTERSECT  Type = "intersect"
	EXCEPT     Type = "except"
	INSTANCEOF Type = "instance of"
	TREATAS    Type = "treat as"
	CASTAS     Type = "cast as"
	CASTABLEAS Type = "castable as"
)

var keywords = map[string]Type{
	"xs":          XSCHEMA,
	"fn":          XFUNC,
	"map":         XMAP,
	"array":       XARRAY,
	"math":        XMATH,
	"function":    FUNC,
	"let":         LET,
	"if":          IF,
	"then":        THEN,
	"else":        ELSE,
	"return":      RETURN,
	"for":         FOR,
	"some":        SOME,
	"every":       EVERY,
	"is":          IS,
	"to":          TO,
	"or":          OR,
	"and":         AND,
	"div":         DIV,
	"idiv":        IDIV,
	"mod":         MOD,
	"union":       UNION,
	"intersect":   INTERSECT,
	"except":      EXCEPT,
	"instance of": INSTANCEOF,
	"treat as":    TREATAS,
	"cast as":     CASTAS,
	"castable as": CASTABLEAS,
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
