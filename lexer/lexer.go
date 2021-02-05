package lexer

import (
	"strings"
	"unicode"

	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
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

// PeekSpace checks if next char is space or not
func (l *Lexer) PeekSpace() bool {
	return unicode.IsSpace(rune(l.ch))
}

// Remaining returns not yet parsed input
func (l *Lexer) Remaining() string {
	if l.fPos >= len(l.input) {
		return l.input
	}
	return l.input[l.pos:]
}

// NextToken returns next token by reading the input characters
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipSpace()

	switch l.ch {
	case '"', '\'':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}
	case '+', '-':
		tt, literal := l.readAdditive()
		tok.Literal = literal
		tok.Type = tt
		return tok
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
		tok = token.Token{Type: token.DOLLAR, Literal: "$"}
	case '#':
		tok = token.Token{Type: token.HASH, Literal: "#"}
	case '?':
		tok = token.Token{Type: token.QUESTION, Literal: "?"}
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
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)
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

// read char skip space
func (l *Lexer) readCharSS() {
	l.readChar()

	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.fPos >= len(l.input) {
		return 0
	}
	return l.input[l.fPos]
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
	l.readChar()
	for l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-' || util.IsDigit(string(l.ch)) {
		if l.ch == '+' || l.ch == '-' {
			break
		}
		if (l.ch == 'e' || l.ch == 'E') && (l.peekChar() == '+' || l.peekChar() == '-') {
			l.readChar()
			l.readChar()
		} else {
			l.readChar()
		}
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for unicode.IsLetter(rune(l.ch)) || l.ch == '-' || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readAdditive() (token.Type, string) {
	minusCnt := 0
	if l.ch == '-' {
		minusCnt++
	}
	l.readCharSS()

	for l.ch == '+' || l.ch == '-' {
		if l.ch == '-' {
			minusCnt++
		}
		l.readCharSS()
	}

	if minusCnt%2 == 1 {
		return token.MINUS, "-"
	}
	return token.PLUS, "+"
}
