package util

import (
	"regexp"
	"strings"
)

// IsChar checks
// Char ::= #x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
func IsChar(ch rune) bool {
	re := regexp.MustCompile(`\x9|\xA|\xD|[\x20-\x{D7FF}]|[\x{E000}-\x{FFFD}]|[\x{10000}-\x{10FFFF}]`)
	return re.MatchString(string(ch))
}

// IsNameStartChar checks
// NameStartChar ::= ":" | [A-Z] | "_" | [a-z] | [#xC0-#xD6] | [#xD8-#xF6] | [#xF8-#x2FF] | [#x370-#x37D] | [#x37F-#x1FFF] | [#x200C-#x200D] | [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF] | [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
func IsNameStartChar(ch rune) bool {
	re := regexp.MustCompile(`:|[A-Za-z]|_|[\xC0-\xD6]|[\xD8-\xF6]|[\xF8-\x{2FF}]|[\x{370}-\x{37D}]|[\x{37F}-\x{1FFF}]|[\x{200C}-\x{200D}]|[\x{2070}-\x{218F}]|[\x{2C00}-\x{2FEF}]|[\x{3001}-\x{D7FF}]|[\x{F900}-\x{FDCF}]|[\x{FDF0}-\x{FFFD}]|[\x{10000}-\x{EFFFF}]`)
	return re.MatchString(string(ch))
}

// IsNameChar checks
// NameChar ::= NameStartChar | "-" | "." | [0-9] | #xB7 | [#x0300-#x036F] | [#x203F-#x2040]
func IsNameChar(ch rune) bool {
	if IsNameStartChar(ch) {
		return true
	}
	re := regexp.MustCompile(`-|.|[0-9]|\xB7|[\x{0300}-\x{036F}]|[\x{203F}-\x{2040}]`)
	return re.MatchString(string(ch))
}

// IsName checks
// Name ::= NameStartChar (NameChar)*
func IsName(str string) bool {
	if len(str) == 0 {
		return false
	}
	for i, c := range str {
		if i == 0 && !IsNameStartChar(c) {
			return false
		}
		if i != 0 && !IsNameChar(c) {
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

// IsQName checks
// QName ::= PrefixedName | UnprefixedName
// PrefixedName ::= Prefix ':' LocalPart
// UnprefixedName ::= LocalPart
// Prefix ::= NCName
// LocalPart ::= NCName
func IsQName(str string) bool {
	ss := strings.Split(str, ":")
	if len(ss) > 2 {
		return false
	}
	for _, s := range ss {
		if !IsNCName(s) {
			return false
		}
	}
	return true
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
