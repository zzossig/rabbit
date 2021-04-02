package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/zzossig/rabbit/eval"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, ctx *object.Context) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		xpath := p.ParseXPath()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaled := eval.Eval(xpath, ctx)
		if evaled != nil {
			io.WriteString(out, evaled.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const RABBIT_FACE = `
   __     __
  /_/|   |\_\  
   |U|___|U|
   |       |
   | ,   , |
  (  = Y =  )
   |   '	 |
  /|       |\
  \| |   | |/
 (_|_|___|_|_)
   '"'   '"'
`

func printParserErrors(out io.Writer, errors []error) {
	io.WriteString(out, RABBIT_FACE)
	io.WriteString(out, "Woops! We ran into some rabbit business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Error()+"\n")
	}
}
