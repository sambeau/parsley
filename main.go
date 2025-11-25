package main

import (
	"fmt"
	"os"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
	"pars/pkg/repl"
)

const VERSION = "0.3.0"

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

	// Create lexer and parser with filename
	l := lexer.NewWithFilename(string(content), filename)
	p := parser.New(l)

	// Parse the program
	program := p.ParseProgram()
	if errors := p.Errors(); len(errors) != 0 {
		printErrors(filename, string(content), errors)
		os.Exit(1)
	}

	// Evaluate the program
	env := evaluator.NewEnvironment()
	env.Filename = filename
	evaluated := evaluator.Eval(program, env)

	// Check for evaluation errors
	if evaluated != nil && evaluated.Type() == evaluator.ERROR_OBJ {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filename, evaluated.Inspect())
		os.Exit(1)
	}

	// Print result if not nil
	if evaluated != nil && evaluated.Type() != evaluator.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
	}
}

// printErrors prints formatted error messages with context
func printErrors(filename string, source string, errors []string) {
	fmt.Fprintf(os.Stderr, "Error in '%s':\n", filename)
	lines := splitLines(source)

	for _, msg := range errors {
		fmt.Fprintf(os.Stderr, "  %s\n", msg)

		// Try to extract line number and column, then show context
		var lineNum, colNum int
		if n, _ := fmt.Sscanf(msg, "line %d, column %d", &lineNum, &colNum); n == 2 && lineNum > 0 && lineNum <= len(lines) {
			// Show the problematic line
			fmt.Fprintf(os.Stderr, "    %s\n", lines[lineNum-1])
			// Show pointer to the error position
			if colNum > 0 {
				pointer := ""
				for i := 1; i < colNum; i++ {
					pointer += " "
				}
				pointer += "^"
				fmt.Fprintf(os.Stderr, "    %s\n", pointer)
			}
		} else if _, err := fmt.Sscanf(msg, "line %d", &lineNum); err == nil && lineNum > 0 && lineNum <= len(lines) {
			// Fallback: show line without pointer if only line number available
			fmt.Fprintf(os.Stderr, "    %s\n", lines[lineNum-1])
		}
	}
}

// splitLines splits source code into lines
func splitLines(s string) []string {
	lines := []string{}
	line := ""
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(ch)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
