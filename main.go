package main

import (
	"fmt"
	"os"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
	"pars/pkg/repl"
)

const VERSION = "0.1.0"

func main() {
	// Check if a filename argument was provided
	if len(os.Args) > 1 {
		arg := os.Args[1]

		// Check for version flag
		if arg == "-V" || arg == "--version" {
			fmt.Printf("pars version %s\n", VERSION)
			os.Exit(0)
		}

		// File execution mode
		executeFile(arg)
	} else {
		// REPL mode
		repl.Start(os.Stdin, os.Stdout)
	}
}

// executeFile reads and executes a pars source file
func executeFile(filename string) {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	// Create lexer and parser
	l := lexer.New(string(content))
	p := parser.New(l)

	// Parse the program
	program := p.ParseProgram()
	if errors := p.Errors(); len(errors) != 0 {
		fmt.Fprintf(os.Stderr, "Parser errors in '%s':\n", filename)
		for _, msg := range errors {
			fmt.Fprintf(os.Stderr, "\t%s\n", msg)
		}
		os.Exit(1)
	}

	// Evaluate the program
	env := evaluator.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	// Print result if not nil
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}
