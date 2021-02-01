package parser

import (
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
	token.IS:         EQ,
	token.EQ:         EQ,
	token.NE:         EQ,
	token.LT:         EQ,
	token.LE:         EQ,
	token.GT:         EQ,
	token.GE:         EQ,
	token.NEV:        EQ,
	token.LTV:        EQ,
	token.LEV:        EQ,
	token.GTV:        EQ,
	token.GEV:        EQ,
	token.DGT:        EQ,
	token.DLT:        EQ,
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
	xpath  *ast.XPath
	errors []error

	curToken  token.Token
	peekToken token.Token

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
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.DECIMAL] = p.parseDecimalLiteral
	p.prefixParseFns[token.DOUBLE] = p.parseDoubleLiteral
	p.prefixParseFns[token.STRING] = p.parseStringLiteral
	p.prefixParseFns[token.DOLLAR] = p.parseVariable
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpr
	p.prefixParseFns[token.PLUS] = p.parsePrefixExpr
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpr
	p.prefixParseFns[token.ARRAY] = p.parseCurlyArrayExpr
	p.prefixParseFns[token.LBRACKET] = p.parseSquareArrayExpr
	p.prefixParseFns[token.IF] = p.parseIfExpr
	p.prefixParseFns[token.FOR] = p.parseForExpr
	p.prefixParseFns[token.LET] = p.parseLetExpr
	p.prefixParseFns[token.QUESTION] = p.parseUnaryLookupExpr
	p.prefixParseFns[token.MAP] = p.parseMapExpr
	p.prefixParseFns[token.SOME] = p.parseQuantifiedExpr
	p.prefixParseFns[token.EVERY] = p.parseQuantifiedExpr
	p.prefixParseFns[token.FUNCTION] = p.parseInlineFunctionExpr
	p.prefixParseFns[token.DOT] = p.parseContextItemExpr
	p.prefixParseFns[token.SLASH] = p.parsePathExpr
	p.prefixParseFns[token.DSLASH] = p.parsePathExpr
	p.prefixParseFns[token.DDOT] = p.parseAxisStep
	p.prefixParseFns[token.AT] = p.parseAxisStep

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseAdditiveExpr
	p.infixParseFns[token.MINUS] = p.parseAdditiveExpr
	p.infixParseFns[token.ASTERISK] = p.parseMultiplicativeExpr
	p.infixParseFns[token.DIV] = p.parseMultiplicativeExpr
	p.infixParseFns[token.IDIV] = p.parseMultiplicativeExpr
	p.infixParseFns[token.MOD] = p.parseMultiplicativeExpr
	p.infixParseFns[token.ARROW] = p.parseArrowExpr
	p.infixParseFns[token.BANG] = p.parseBangExpr
	p.infixParseFns[token.IS] = p.parseComparisonExpr
	p.infixParseFns[token.EQ] = p.parseComparisonExpr
	p.infixParseFns[token.NE] = p.parseComparisonExpr
	p.infixParseFns[token.LT] = p.parseComparisonExpr
	p.infixParseFns[token.LE] = p.parseComparisonExpr
	p.infixParseFns[token.GT] = p.parseComparisonExpr
	p.infixParseFns[token.GE] = p.parseComparisonExpr
	p.infixParseFns[token.DGT] = p.parseComparisonExpr
	p.infixParseFns[token.DLT] = p.parseComparisonExpr
	p.infixParseFns[token.EQV] = p.parseComparisonExpr
	p.infixParseFns[token.NEV] = p.parseComparisonExpr
	p.infixParseFns[token.LTV] = p.parseComparisonExpr
	p.infixParseFns[token.LEV] = p.parseComparisonExpr
	p.infixParseFns[token.GTV] = p.parseComparisonExpr
	p.infixParseFns[token.GEV] = p.parseComparisonExpr
	p.infixParseFns[token.OR] = p.parseOrExpr
	p.infixParseFns[token.AND] = p.parseAndExpr
	p.infixParseFns[token.TO] = p.parseRangeExpr
	p.infixParseFns[token.UNION] = p.parseUnionExpr
	p.infixParseFns[token.VBAR] = p.parseUnionExpr
	p.infixParseFns[token.INTERSECT] = p.parseIntersectExceptExpr
	p.infixParseFns[token.EXCEPT] = p.parseIntersectExceptExpr
	p.infixParseFns[token.INSTANCE] = p.parseInstanceofExpr
	p.infixParseFns[token.CAST] = p.parseCastExpr
	p.infixParseFns[token.CASTABLE] = p.parseCastableExpr
	p.infixParseFns[token.TREAT] = p.parseTreatExpr
	p.infixParseFns[token.HASH] = p.parseNamedFunctionRef

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// cur t or t1 or t2 or ..
func (p *Parser) curTokenIs(t token.Type, ts ...token.Type) bool {
	if len(ts) > 0 {
		for _, tt := range ts {
			if p.curToken.Type == tt {
				return true
			}
		}
	}
	return p.curToken.Type == t
}

// peek t or t1 or t2 or ..
func (p *Parser) peekTokenIs(t token.Type, ts ...token.Type) bool {
	if len(ts) > 0 {
		for _, tt := range ts {
			if p.peekToken.Type == tt {
				return true
			}
		}
	}
	return p.peekToken.Type == t
}

// expect t or t1 or t2 or ...
func (p *Parser) expectPeek(t token.Type, ts ...token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else if len(ts) > 0 {
		for _, tt := range ts {
			if p.peekTokenIs(tt) {
				p.nextToken()
				return true
			}
		}
	}
	p.peekError(t)
	return false
}

func (p *Parser) expectCur(t token.Type) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// *must* used in a grouped expressions
func (p *Parser) hasComma() bool {
	input := p.l.Input()[p.l.Pos():]
	lCnt := 0
	rCnt := 0
	cCnt := 0

	for _, ch := range input {
		if ch == '(' {
			lCnt++
		}
		if ch == ')' {
			rCnt++
		}
		if ch == ',' {
			cCnt++
		}
		if rCnt == lCnt+1 {
			break
		}
	}

	if cCnt > 0 && rCnt == lCnt+1 {
		return true
	}

	return false
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
	p.xpath = xpath

	for !p.curTokenIs(token.EOF) {
		items := p.parseItem()
		if items != nil {
			xpath.Items = append(xpath.Items, items...)
		}
		p.nextToken()
	}
	return xpath
}
