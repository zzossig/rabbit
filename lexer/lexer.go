package lexer

import (
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
		tok = token.Token{Type: token.DOLLAR, Literal: "$"}
	case '.':
		// TODO (.5) = (0.5)
		if l.peekChar() == '.' {
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
			tokType := token.LookupIdent(tok.Literal)

			if l.peekChar() == ':' && tok.Type != token.IDENT {

			}
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
	for {
		l.readChar()
		if l.ch == 0 || (ch == '"' && l.ch == '"') || (ch == '\'' && l.ch == '\'') {
			break
		}
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for unicode.IsDigit(rune(l.ch)) || l.ch == 'e' || l.ch == 'E' || l.ch == '.' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for unicode.IsLetter(rune(l.ch)) || l.ch == '-' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}
