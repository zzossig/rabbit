package parser

import (
	"fmt"

	"github.com/zzossig/xpath/token"
)

// Errors returns []error
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
