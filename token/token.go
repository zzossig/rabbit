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
	NS         Type = "namespace"
	BIF        Type = "built-in function"
	AXIS       Type = "axis"
	IDENT      Type = "identifier"
	XTYPE      Type = "xs:type"
	XTYPEF     Type = "xs:type()"
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
	UQUESTION  Type = "(?)" // unary lookup
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
	"for":       FOR,
	"in":        IN,
	"some":      SOME,
	"every":     EVERY,
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

// built-in functions
var bifs = map[string]string{
	// node-test
	"item":                   "item()",
	"node":                   "node()",
	"attribute":              "attribute()",
	"comment":                "comment()",
	"document":               "document()",
	"element":                "element()",
	"schema-element":         "schema-element()",
	"processing-instruction": "processing-instruction()",
	"text":                   "text()",
	"document-node":          "document-node()",

	// func
	"abs":                               "abs()",
	"acos":                              "acos()",
	"add-dayTimeDurations":              "add-dayTimeDurations()",
	"add-dayTimeDuration-to-date":       "add-dayTimeDuration-to-date()",
	"add-dayTimeDuration-to-dateTime":   "add-dayTimeDuration-to-dateTime()",
	"add-dayTimeDuration-to-time":       "add-dayTimeDuration-to-time()",
	"add-yearMonthDurations":            "add-yearMonthDurations()",
	"add-yearMonthDuration-to-date":     "add-yearMonthDuration-to-date()",
	"add-yearMonthDuration-to-dateTime": "add-yearMonthDuration-to-dateTime()",
	"adjust-dateTime-to-timezone":       "adjust-dateTime-to-timezone()",
	"adjust-date-to-timezone":           "adjust-date-to-timezone()",
	"adjust-time-to-timezone":           "adjust-time-to-timezone()",
	"analyze-string":                    "analyze-string()",
	"asin":                              "asin()",
	"atan":                              "atan()",
	"atan2":                             "atan2()",
	"available-environment-variables":   "available-environment-variables()",
	"avg":                               "avg()",
	"base64Binary-equal":                "base64Binary-equal()",
	"base-uri":                          "base-uri()",
	"boolean":                           "boolean()",
	"boolean-equal":                     "boolean-equal()",
	"boolean-greater-than":              "boolean-greater-than()",
	"boolean-less-than":                 "boolean-less-than()",
	"ceiling":                           "ceiling()",
	"codepoint-equal":                   "codepoint-equal()",
	"codepoints-to-string":              "codepoints-to-string()",
	"collection":                        "collection()",
	"compare":                           "compare()",
	"concat":                            "concat()",
	"concatenate":                       "concatenate()",
	"contains":                          "contains()",
	"cos":                               "cos()",
	"count":                             "count()",
	"current-date":                      "current-date()",
	"current-dateTime":                  "current-dateTime()",
	"current-time":                      "current-time()",
	"data":                              "data()",
	"date-equal":                        "date-equal()",
	"date-greater-than":                 "date-greater-than()",
	"date-less-than":                    "date-less-than()",
	"dateTime":                          "dateTime()",
	"dateTime-equal":                    "dateTime-equal()",
	"dateTime-greater-than":             "dateTime-greater-than()",
	"dateTime-less-than":                "dateTime-less-than()",
	"day-from-date":                     "day-from-date()",
	"day-from-dateTime":                 "day-from-dateTime()",
	"days-from-duration":                "days-from-duration()",
	"dayTimeDuration-greater-than":      "dayTimeDuration-greater-than()",
	"dayTimeDuration-less-than":         "dayTimeDuration-less-than()",
	"deep-equal":                        "deep-equal()",
	"default-collation":                 "default-collation()",
	"distinct-values":                   "distinct-values()",
	"divide-dayTimeDuration":            "divide-dayTimeDuration()",
	"divide-dayTimeDuration-by-dayTimeDuration":     "divide-dayTimeDuration-by-dayTimeDuration()",
	"divide-yearMonthDuration":                      "divide-yearMonthDuration()",
	"divide-yearMonthDuration-by-yearMonthDuration": "divide-yearMonthDuration-by-yearMonthDuration()",
	"doc":                                    "doc()",
	"doc-available":                          "doc-available()",
	"document-uri":                           "document-uri()",
	"duration-equal":                         "duration-equal()",
	"element-with-id":                        "element-with-id()",
	"empty":                                  "empty()",
	"encode-for-uri":                         "encode-for-uri()",
	"ends-with":                              "ends-with()",
	"environment-variable":                   "environment-variable()",
	"error":                                  "error()",
	"escape-html-uri":                        "escape-html-uri()",
	"exactly-one":                            "exactly-one()",
	"except":                                 "except()",
	"exists":                                 "exists()",
	"exp":                                    "exp()",
	"exp10":                                  "exp10()",
	"false":                                  "false()",
	"filter":                                 "filter()",
	"floor":                                  "floor()",
	"fold-left":                              "fold-left()",
	"fold-right":                             "fold-right()",
	"for-each":                               "for-each()",
	"for-each-pair":                          "for-each-pair()",
	"format-date":                            "format-date()",
	"format-dateTime":                        "format-dateTime()",
	"format-integer":                         "format-integer()",
	"format-number":                          "format-number()",
	"format-time":                            "format-time()",
	"function-arity":                         "function-arity()",
	"function-lookup":                        "function-lookup()",
	"function-name":                          "function-name()",
	"gDay-equal":                             "gDay-equal()",
	"generate-id":                            "generate-id()",
	"gMonthDay-equal":                        "gMonthDay-equal()",
	"gMonth-equal":                           "gMonth-equal()",
	"gYear-equal":                            "gYear-equal()",
	"gYearMonth-equal":                       "gYearMonth-equal()",
	"has-children":                           "has-children()",
	"head":                                   "head()",
	"hexBinary-equal":                        "hexBinary-equal()",
	"hours-from-dateTime":                    "hours-from-dateTime()",
	"hours-from-duration":                    "hours-from-duration()",
	"hours-from-time":                        "hours-from-time()",
	"id":                                     "id()",
	"idref":                                  "idref()",
	"implicit-timezone":                      "implicit-timezone()",
	"index-of":                               "index-of()",
	"innermost":                              "innermost()",
	"in-scope-prefixes":                      "in-scope-prefixes()",
	"insert-before":                          "insert-before()",
	"intersect":                              "intersect()",
	"iri-to-uri":                             "iri-to-uri()",
	"is-same-node":                           "is-same-node()",
	"lang":                                   "lang()",
	"last":                                   "last()",
	"local-name":                             "local-name()",
	"local-name-from-QName":                  "local-name-from-QName()",
	"log":                                    "log()",
	"log10":                                  "log10()",
	"lower-case":                             "lower-case()",
	"matches":                                "matches()",
	"max":                                    "max()",
	"min":                                    "min()",
	"minutes-from-dateTime":                  "minutes-from-dateTime()",
	"minutes-from-duration":                  "minutes-from-duration()",
	"minutes-from-time":                      "minutes-from-time()",
	"month-from-date":                        "month-from-date()",
	"month-from-dateTime":                    "month-from-dateTime()",
	"months-from-duration":                   "months-from-duration()",
	"multiply-dayTimeDuration":               "multiply-dayTimeDuration()",
	"multiply-yearMonthDuration":             "multiply-yearMonthDuration()",
	"name":                                   "name()",
	"namespace-uri":                          "namespace-uri()",
	"namespace-uri-for-prefix":               "namespace-uri-for-prefix()",
	"namespace-uri-from-QName":               "namespace-uri-from-QName()",
	"nilled":                                 "nilled()",
	"node-after":                             "node-after()",
	"node-before":                            "node-before()",
	"node-name":                              "node-name()",
	"normalize-space":                        "normalize-space()",
	"normalize-unicode":                      "normalize-unicode()",
	"not":                                    "not()",
	"NOTATION-equal":                         "NOTATION-equal()",
	"number":                                 "number()",
	"numeric-add":                            "numeric-add()",
	"numeric-divide":                         "numeric-divide()",
	"numeric-equal":                          "numeric-equal()",
	"numeric-greater-than":                   "numeric-greater-than()",
	"numeric-integer-divide":                 "numeric-integer-divide()",
	"numeric-less-than":                      "numeric-less-than()",
	"numeric-mod":                            "numeric-mod()",
	"numeric-multiply":                       "numeric-multiply()",
	"numeric-subtract":                       "numeric-subtract()",
	"numeric-unary-minus":                    "numeric-unary-minus()",
	"numeric-unary-plus":                     "numeric-unary-plus()",
	"one-or-more":                            "one-or-more()",
	"outermost":                              "outermost()",
	"parse-xml":                              "parse-xml()",
	"parse-xml-fragment":                     "parse-xml-fragment()",
	"path":                                   "path()",
	"pi":                                     "pi()",
	"position":                               "position()",
	"pow":                                    "pow()",
	"prefix-from-QName":                      "prefix-from-QName()",
	"QName":                                  "QName()",
	"QName-equal":                            "QName-equal()",
	"remove":                                 "remove()",
	"replace":                                "replace()",
	"resolve-QName":                          "resolve-QName()",
	"resolve-uri":                            "resolve-uri()",
	"reverse":                                "reverse()",
	"root":                                   "root()",
	"round":                                  "round()",
	"round-half-to-even":                     "round-half-to-even()",
	"seconds-from-dateTime":                  "seconds-from-dateTime()",
	"seconds-from-duration":                  "seconds-from-duration()",
	"seconds-from-time":                      "seconds-from-time()",
	"serialize":                              "serialize()",
	"sin":                                    "sin()",
	"sqrt":                                   "sqrt()",
	"starts-with":                            "starts-with()",
	"static-base-uri":                        "static-base-uri()",
	"string":                                 "string()",
	"string-join":                            "string-join()",
	"string-length":                          "string-length()",
	"string-to-codepoints":                   "string-to-codepoints()",
	"subsequence":                            "subsequence()",
	"substring":                              "substring()",
	"substring-after":                        "substring-after()",
	"substring-before":                       "substring-before()",
	"subtract-dates":                         "subtract-dates()",
	"subtract-dateTimes":                     "subtract-dateTimes()",
	"subtract-dayTimeDuration-from-date":     "subtract-dayTimeDuration-from-date()",
	"subtract-dayTimeDuration-from-dateTime": "subtract-dayTimeDuration-from-dateTime()",
	"subtract-dayTimeDuration-from-time":     "subtract-dayTimeDuration-from-time()",
	"subtract-dayTimeDurations":              "subtract-dayTimeDurations()",
	"subtract-times":                         "subtract-times()",
	"subtract-yearMonthDuration-from-date":   "subtract-yearMonthDuration-from-date()",
	"subtract-yearMonthDuration-from-dateTime": "subtract-yearMonthDuration-from-dateTime()",
	"subtract-yearMonthDurations":              "subtract-yearMonthDurations()",
	"sum":                                      "sum()",
	"tail":                                     "tail()",
	"tan":                                      "tan()",
	"time-equal":                               "time-equal()",
	"time-greater-than":                        "time-greater-than()",
	"time-less-than":                           "time-less-than()",
	"timezone-from-date":                       "timezone-from-date()",
	"timezone-from-dateTime":                   "timezone-from-dateTime()",
	"timezone-from-time":                       "timezone-from-time()",
	"to":                                       "to()",
	"tokenize":                                 "tokenize()",
	"trace":                                    "trace()",
	"translate":                                "translate()",
	"true":                                     "true()",
	"union":                                    "union()",
	"unordered":                                "unordered()",
	"unparsed-text":                            "unparsed-text()",
	"unparsed-text-available":                  "unparsed-text-available()",
	"unparsed-text-lines":                      "unparsed-text-lines()",
	"upper-case":                               "upper-case()",
	"uri-collection":                           "uri-collection()",
	"year-from-date":                           "year-from-date()",
	"year-from-dateTime":                       "year-from-dateTime()",
	"yearMonthDuration-greater-than":           "yearMonthDuration-greater-than()",
	"yearMonthDuration-less-than":              "yearMonthDuration-less-than()",
	"years-from-duration":                      "years-from-duration()",
	"zero-or-one":                              "zero-or-one()",
}

var xtypes = map[string]string{
	"anyAtomicType":      "xs:anyAtomicType",
	"untypedAtomic":      "xs:untypedAtomic",
	"dateTime":           "xs:dateTime",
	"dateTimeStamp":      "xs:dateTimeStamp",
	"date":               "xs:date",
	"time":               "xs:time",
	"duration":           "xs:duration",
	"yearMonthDuration":  "xs:yearMonthDuration",
	"dayTimeDuration":    "xs:dayTimeDuration",
	"float":              "xs:float",
	"double":             "xs:double",
	"decimal":            "xs:decimal",
	"integer":            "xs:integer",
	"nonPositiveInteger": "xs:nonPositiveInteger",
	"negativeInteger":    "xs:negativeInteger",
	"long":               "xs:long",
	"int":                "xs:int",
	"short":              "xs:short",
	"byte":               "xs:byte",
	"nonNegativeInteger": "xs:nonNegativeInteger",
	"unsignedLong":       "xs:unsignedLong",
	"unsignedInt":        "xs:unsignedInt",
	"unsignedShort":      "xs:unsignedShort",
	"unsignedByte":       "xs:unsignedByte",
	"positiveInteger":    "xs:positiveInteger",
	"gYearMonth":         "xs:gYearMonth",
	"gYear":              "xs:gYear",
	"gMonthDay":          "xs:gMonthDay",
	"gDay":               "xs:gDay",
	"gMonth":             "xs:gMonth",
	"string":             "xs:string",
	"normalizedString":   "xs:normalizedString",
	"token":              "xs:token",
	"language":           "xs:language",
	"NMTOKEN":            "xs:NMTOKEN",
	"Name":               "xs:Name",
	"NCName":             "xs:NCName",
	"ID":                 "xs:ID",
	"IDREF":              "xs:IDREF",
	"ENTITY":             "xs:ENTITY",
	"boolean":            "xs:boolean",
	"base64Binary":       "xs:base64Binary",
	"hexBinary":          "xs:hexBinary",
	"anyURI":             "xs:anyURI",
	"QName":              "xs:QName",
	"NOTATION":           "xs:NOTATION",
}

// IsBIF checks if (ident string) is a built-in function or not
func IsBIF(ident string) bool {
	if _, ok := bifs[ident]; ok {
		return true
	}
	return false
}

// IsXType checks if (ident string) is a xs:type or not
func IsXType(ident string) bool {
	if _, ok := xtypes[ident]; ok {
		return true
	}
	return false
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
