package lexer

import (
	"strings"
	"unicode"

	"github.com/zzossig/xpath/token"
)

// Lexer reads input string one by one
type Lexer struct {
	input string // user input
	pos   int    // current position within input
	fPos  int    // following position
	ch    byte   // current char under examination
}

// New returns Lexer pointer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns next token by reading the input characters
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipSpace()

	switch l.ch {
	case '"', '\'':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: "+"}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: "-"}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: "*"}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token.Token{Type: token.DSLASH, Literal: "//"}
		} else {
			tok = token.Token{Type: token.SLASH, Literal: "/"}
		}
	case '<':
		if l.peekChar() == '<' {
			l.readChar()
			tok = token.Token{Type: token.DLT, Literal: "<<"}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: "<="}
		} else {
			tok = token.Token{Type: token.LT, Literal: "<"}
		}
	case '>':
		if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.DGT, Literal: ">>"}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: ">="}
		} else {
			tok = token.Token{Type: token.GT, Literal: ">"}
		}
	case '=':
		if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "=>"}
		} else {
			tok = token.Token{Type: token.EQ, Literal: "="}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NE, Literal: "!="}
		} else {
			tok = token.Token{Type: token.BANG, Literal: "!"}
		}
	case '@':
		tok = token.Token{Type: token.AT, Literal: "@"}
	case '$':
		l.readChar()
		_, literal := l.readIdent()
		tok.Literal = "$" + literal
		tok.Type = token.VAR
		return tok
	case '.':
		if unicode.IsNumber(rune(l.peekChar())) {
			l.readChar()
			tok.Literal = "." + l.readNumber()
			tok.Type = token.LookupNumber(tok.Literal)
			return tok
		} else if l.peekChar() == '.' {
			l.readChar()
			tok = token.Token{Type: token.DDOT, Literal: ".."}
		} else {
			tok = token.Token{Type: token.DOT, Literal: "."}
		}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: ","}
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			tok = token.Token{Type: token.DCOLON, Literal: "::"}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.ASSIGN, Literal: ":="}
		} else {
			tok = token.Token{Type: token.COLON, Literal: ":"}
		}
	case ';':
		tok = token.Token{Type: token.SCOLON, Literal: ";"}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: ")"}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: "{"}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: "}"}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: "["}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: "]"}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.DVBAR, Literal: "||"}
		} else {
			tok = token.Token{Type: token.VBAR, Literal: "|"}
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if unicode.IsLetter(rune(l.ch)) {
			initPos, literal := l.readIdent()
			tok.Literal = literal
			tok.Type = l.readIdentType(tok.Literal, initPos)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Literal = l.readNumber()
			tok.Type = token.LookupNumber(tok.Literal)
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.fPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.fPos]
	}
	l.pos = l.fPos
	l.fPos++
}

func (l *Lexer) peekChar() byte {
	if l.fPos >= len(l.input) {
		return 0
	}
	return l.input[l.fPos]
}

// peek char skip space
func (l *Lexer) peekCharSS() byte {
	var ch byte

	for i := 0; i < len(l.input); i++ {
		if l.fPos+i >= len(l.input) {
			ch = 0
			break
		}
		if unicode.IsSpace(rune(l.input[l.fPos+i])) {
			continue
		}
		ch = l.input[l.fPos+i]
		break
	}

	return ch
}

func (l *Lexer) skipSpace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	ch := l.ch
	pos := l.pos + 1
	cnt := 0 // count escape string [" or ']
	for {
		l.readChar()
		if l.ch == 0 || (ch == '"' && l.ch == '"' && l.peekChar() != '"') || (ch == '\'' && l.ch == '\'' && l.peekChar() != '\'') {
			break
		}
		if (ch == '"' && l.ch == '"' && l.peekChar() == '"') || (ch == '\'' && l.ch == '\'' && l.peekChar() == '\'') {
			cnt++
		}
	}

	if cnt%2 != 0 {
		// TODO error occur err:XPST0003 - allowed: '''', notAllowed: '''''
	}

	if ch == '"' {
		return strings.ReplaceAll(l.input[pos:l.pos], "\"\"", "\"")
	}
	return strings.ReplaceAll(l.input[pos:l.pos], "''", "'")
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for unicode.IsDigit(rune(l.ch)) || l.ch == 'e' || l.ch == 'E' || l.ch == '.' || l.ch == '+' || l.ch == '-' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() (int, string) {
	pos := l.pos
	for unicode.IsLetter(rune(l.ch)) || l.ch == '-' || l.ch == '_' {
		l.readChar()
	}
	return pos, l.input[pos:l.pos]
}

func (l *Lexer) readIdentType(literal string, initPos int) token.Type {
	if l.ch == ':' {
		if l.peekChar() == ':' {
			return token.AXIS
		}
		return token.NS
	} else if l.ch == '(' || l.peekCharSS() == '(' {
		if token.IsBIF(literal) {
			return token.BIF
		} else if token.IsXType(literal) {
			return token.XTYPEF
		} else if literal == "function" {
			return token.FUNCTION
		}
	} else if initPos != 0 &&
		(l.input[initPos-1] == '/' ||
			l.input[initPos-1] == '@' ||
			l.input[initPos-1] == '$') {
		return token.IDENT
	} else if initPos != 0 && l.input[initPos-1] == ':' { // xs:, math:, fn:, ...
		if token.IsXType(literal) {
			return token.XTYPE
		}
	}
	return token.LookupIdent(literal)
}
