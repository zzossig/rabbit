package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/zzossig/xpath/context"
	"github.com/zzossig/xpath/eval"
	"github.com/zzossig/xpath/lexer"
	"github.com/zzossig/xpath/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := context.NewContext()
	env.NewReaderFile("text.txt", true)

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

		evaled := eval.Eval(xpath, env)
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
