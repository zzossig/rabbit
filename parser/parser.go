package parser

import (
	"fmt"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/lexer"
	"github.com/zzossig/xpath/token"
)

// Precedence Order
const (
	LOWEST = iota
	COMMA
	FOR
	OR
	AND
	EQ
	DVBAR
	TO
	SUM
	DIV
	UNION
	INTERSECT
	INSTANCEOF
	TREATAS
	CASTABLEAS
	CASTAS
	ARROW
	UNARY
	BANG
	SLASH
	PREDICATE
	LOOKUP
)

var precedences = map[token.Type]int{
	token.COMMA:      COMMA,
	token.FOR:        FOR,
	token.LET:        FOR,
	token.SOME:       FOR,
	token.EVERY:      FOR,
	token.IF:         FOR,
	token.OR:         OR,
	token.AND:        AND,
	token.EQ:         EQ,
	token.DVBAR:      DVBAR,
	token.TO:         TO,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.ASTERISK:   DIV,
	token.DIV:        DIV,
	token.IDIV:       DIV,
	token.MOD:        DIV,
	token.UNION:      UNION,
	token.VBAR:       UNION,
	token.INTERSECT:  INTERSECT,
	token.EXCEPT:     INTERSECT,
	token.INSTANCEOF: INSTANCEOF,
	token.TREATAS:    TREATAS,
	token.CASTABLEAS: CASTABLEAS,
	token.CASTAS:     CASTAS,
	token.ARROW:      ARROW,
	token.UPLUS:      UNARY,
	token.UMINUS:     UNARY,
	token.BANG:       BANG,
	token.SLASH:      SLASH,
	token.DSLASH:     SLASH,
	token.LBRACKET:   PREDICATE,
	token.RBRACKET:   PREDICATE,
	token.QUESTION:   PREDICATE,
	token.UQUESTION:  LOOKUP,
}

type (
	prefixParseFn func() ast.ExprSingle
	infixParseFn  func(ast.ExprSingle) ast.ExprSingle
)

// Parser object
type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  token.Token
	peekToken token.Token
	lastToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns parser object
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpr
	p.prefixParseFns[token.PLUS] = p.parsePrefixExpr
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpr

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseAdditiveExpr
	p.infixParseFns[token.MINUS] = p.parseAdditiveExpr
	p.infixParseFns[token.ASTERISK] = p.parseMultiplicativeExpr
	p.infixParseFns[token.DIV] = p.parseMultiplicativeExpr
	p.infixParseFns[token.IDIV] = p.parseMultiplicativeExpr
	p.infixParseFns[token.MOD] = p.parseMultiplicativeExpr

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Errors returns []error
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseXPath generates ast tree
func (p *Parser) ParseXPath() *ast.XPath {
	xpath := &ast.XPath{}
	xpath.Items = []ast.Item{}

	for !p.curTokenIs(token.EOF) {
		items := p.parseItem()
		if items != nil {
			xpath.Items = append(xpath.Items, items...)
		}
		p.nextToken()
	}
	return xpath
}
