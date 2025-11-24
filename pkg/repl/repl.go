package repl

import (
	"bufio"
	"fmt"
	"io"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
)

const PROMPT = ">> "

const PARSER_LOGO = `
 ____   __    ____   _____ 
|  _ \ / _\  |  _ \ / ____|
| |_) / /_\  | |_) | (___ 
|  __/ ____ |  _ < \___ \ 
| |  / /    | |_) |____) |
|_| /_/     |____/|_____/ 

Welcome to the Pars programming language!
Feel free to type in commands
`

// Start starts the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := evaluator.NewEnvironment()

	fmt.Fprintf(out, PARSER_LOGO)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			fmt.Fprintf(out, "Goodbye!\n")
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if errors := p.Errors(); len(errors) != 0 {
			printParserErrors(out, errors)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// printParserErrors prints parser errors
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
