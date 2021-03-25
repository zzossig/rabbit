package util

import (
	"regexp"
	"strings"
)

// IsChar checks
// Char ::= #x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
func IsChar(ch rune) bool {
	re := regexp.MustCompile(`^[\x09\x0A\x0D\x20-\x{D7FF}\x{E000}-\x{FFFD}\x{10000}-\x{10FFFF}]$`)
	return re.MatchString(string(ch))
}

// IsNameStartChar checks
// NameStartChar ::= ":" | [A-Z] | "_" | [a-z] | [#xC0-#xD6] | [#xD8-#xF6] | [#xF8-#x2FF] | [#x370-#x37D] | [#x37F-#x1FFF] | [#x200C-#x200D] | [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF] | [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
func IsNameStartChar(ch rune) bool {
	re := regexp.MustCompile(`^[:A-Za-z_\xC0-\xD6\xD8-\xF6\xF8-\x{2FF}\x{370}-\x{37D}\x{37F}-\x{1FFF}\x{200C}-\x{200D}\x{2070}-\x{218F}\x{2C00}-\x{2FEF}\x{3001}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFFD}\x{10000}-\x{EFFFF}]$`)
	return re.MatchString(string(ch))
}

// IsNameChar checks
// NameChar ::= NameStartChar | "-" | "." | [0-9] | #xB7 | [#x0300-#x036F] | [#x203F-#x2040]
func IsNameChar(ch rune) bool {
	if IsNameStartChar(ch) {
		return true
	}
	re := regexp.MustCompile(`^[-\.0-9\xB7\x{0300}-\x{036F}\x{203F}-\x{2040}]$`)
	return re.MatchString(string(ch))
}

// IsName checks
// Name ::= NameStartChar (NameChar)*
func IsName(str string) bool {
	if len(str) == 0 {
		return false
	}
	if !IsNameStartChar(rune(str[0])) {
		return false
	}
	for _, c := range str[1:] {
		if !IsNameChar(c) {
			return false
		}
	}
	return true
}

// IsNCName checks
// NCName ::= Name - (Char* ':' Char*)
func IsNCName(str string) bool {
	return IsName(str) && !strings.Contains(str, ":")
}

// IsPrefixedName checks
// PrefixedName ::= Prefix ':' LocalPart
func IsPrefixedName(str string) bool {
	ss := strings.Split(str, ":")
	if len(ss) != 2 {
		return false
	}
	for _, s := range ss {
		if !IsNCName(s) {
			return false
		}
	}
	return true
}

// IsQName checks
// QName ::= PrefixedName | UnprefixedName
// PrefixedName ::= Prefix ':' LocalPart
// UnprefixedName ::= LocalPart
// Prefix ::= NCName
// LocalPart ::= NCName
func IsQName(str string) bool {
	if IsNCName(str) {
		return true
	}
	return IsPrefixedName(str)
}

// IsURIQualifiedName checks
// URIQualifiedName ::= BracedURILiteral NCName
func IsURIQualifiedName(str string) bool {
	re := regexp.MustCompile(`^Q{[^{}]*}`)
	loc := re.FindStringIndex(str)
	if len(loc) <= 1 {
		return false
	}
	return IsNCName(str[loc[1]:])
}

// IsBracedURILiteral checks
// BracedURILiteral ::= "Q" "{" [^{}]* "}"
func IsBracedURILiteral(str string) bool {
	re := regexp.MustCompile(`^Q{[^{}]*}$`)
	return re.MatchString(str)
}

// IsEQName checks
// EQName ::= QName | URIQualifiedName
func IsEQName(str string) bool {
	return IsQName(str) || IsURIQualifiedName(str)
}

// IsWildcard checks
// Wildcard ::= "*" | (NCName ":*") | ("*:" NCName) | (BracedURILiteral "*")
func IsWildcard(str string) bool {
	if str == "*" {
		return true
	} else if strings.HasSuffix(str, ":*") && IsNCName(str[:len(str)-2]) {
		return true
	} else if strings.HasPrefix(str, "*:") && IsNCName(str[2:]) {
		return true
	} else if strings.HasSuffix(str, "*") && IsBracedURILiteral(str[:len(str)-1]) {
		return true
	}
	return false
}

// IsValueComp checks
// ValueComp ::= "eq" | "ne" | "lt" | "le" | "gt" | "ge"
func IsValueComp(str string) bool {
	re := regexp.MustCompile(`^(eq|ne|lt|le|gt|ge)$`)
	return re.MatchString(str)
}

// IsGeneralComp checks
// GeneralComp ::= "=" | "!=" | "<" | "<=" | ">" | ">="
func IsGeneralComp(str string) bool {
	re := regexp.MustCompile(`^(=|!=|<|<=|>|>=)$`)
	return re.MatchString(str)
}

// IsNodeComp checks
// NodeComp ::= "is" | "<<" | ">>"
func IsNodeComp(str string) bool {
	re := regexp.MustCompile(`^(is|<<|>>)$`)
	return re.MatchString(str)
}

// IsForwardAxis checks
// ForwardAxis ::= ("child" "::") | ("descendant" "::") | ("attribute" "::") | ("self" "::") | ("descendant-or-self" "::") | ("following-sibling" "::") | ("following" "::") | ("namespace" "::")
func IsForwardAxis(str string) bool {
	re := regexp.MustCompile(`^(child::|descendant::|attribute::|self::|descendant-or-self::|following-sibling::|following::|namespace::)$`)
	return re.MatchString(str)
}

// IsReverseAxis checks
// ReverseAxis ::= ("parent" "::") | ("ancestor" "::") | ("preceding-sibling" "::") | ("preceding" "::") | ("ancestor-or-self" "::")
func IsReverseAxis(str string) bool {
	re := regexp.MustCompile(`^(parent::|ancestor::|preceding-sibling::|preceding::|ancestor-or-self::)$`)
	return re.MatchString(str)
}

// IsNumber checks number. eg) .05e2
func IsNumber(str string) bool {
	// re := regexp.MustCompile(`^([1-9]{1}\d*|[0]?)(\.\d*)?(([e|E][+|-]?)?\d*)?$`)
	re := regexp.MustCompile(`^(\d*[\.])?\d+([eE][+-]?\d+)?$`)
	return re.MatchString(str)
}

// IsDigit checks Digits ::= [0-9]+
func IsDigit(str string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(str)
}

// IsOccurrenceIndicator checks
// OccurrenceIndicator ::= "?" | "*" | "+"
func IsOccurrenceIndicator(str string) bool {
	re := regexp.MustCompile(`^(\?|\*|\+)$`)
	return re.MatchString(str)
}

// CheckKindTest checks TypeID field
func CheckKindTest(str string) byte {
	switch str {
	case "document-node":
		return 1
	case "element":
		return 2
	case "attribute":
		return 3
	case "schema-element":
		return 4
	case "schema-attribute":
		return 5
	case "processing-instruction":
		return 6
	case "comment":
		return 7
	case "text":
		return 8
	case "namespace-node":
		return 9
	case "node":
		return 10
	default:
		return 0
	}
}

// CheckItemType checks TypeID field
func CheckItemType(str string) byte {
	if CheckKindTest(str) != 0 {
		return 1
	}

	switch str {
	case "item":
		return 2
	case "function":
		return 3
	case "map":
		return 4
	case "array":
		return 5
	case "(":
		return 7
	default:
		return 6
	}
}
