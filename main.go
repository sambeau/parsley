package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
	"github.com/sambeau/parsley/pkg/repl"
)

// Version is set at compile time via -ldflags
var Version = "dev"

func main() {
	// Check if a filename argument was provided
	if len(os.Args) > 1 {
		arg := os.Args[1]

		// Check for version flag
		if arg == "-V" || arg == "--version" {
			fmt.Printf("pars version %s\n", Version)
			os.Exit(0)
		}

		// File execution mode
		executeFile(arg)
	} else {
		// REPL mode
		repl.Start(os.Stdin, os.Stdout, Version)
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
		// Format runtime errors the same way as parse errors
		errObj, ok := evaluated.(*evaluator.Error)
		if ok && errObj.Line > 0 {
			// Error has position information
			printErrors(filename, string(content), []string{errObj.Inspect()})
		} else {
			// Error without position information (legacy format)
			fmt.Fprintf(os.Stderr, "%s: %s\n", filename, evaluated.Inspect())
		}
		os.Exit(1)
	}

	// Print result if not null and not an error
	if evaluated != nil && evaluated.Type() != evaluator.ERROR_OBJ && evaluated.Type() != evaluator.NULL_OBJ {
		fmt.Println(evaluator.ObjectToPrintString(evaluated))
	}
}

// printErrors prints formatted error messages with context
func printErrors(filename string, source string, errors []string) {
	fmt.Fprintf(os.Stderr, "Error in '%s':\n", filename)
	lines := strings.Split(source, "\n")

	for _, msg := range errors {
		fmt.Fprintf(os.Stderr, "  %s\n", msg)

		// Try to extract line number and column, then show context
		var lineNum, colNum int
		if n, _ := fmt.Sscanf(msg, "line %d, column %d", &lineNum, &colNum); n == 2 && lineNum > 0 && lineNum <= len(lines) {
			sourceLine := lines[lineNum-1]

			// Calculate how many columns to trim from the left
			trimCount := 0
			for i := 0; i < len(sourceLine); i++ {
				if sourceLine[i] == ' ' || sourceLine[i] == '\t' {
					if sourceLine[i] == '\t' {
						trimCount += 8
					} else {
						trimCount++
					}
				} else {
					break
				}
			}

			// Trim left whitespace from the source line
			trimmedLine := strings.TrimLeft(sourceLine, " \t")

			// Show the trimmed line with slight indentation
			fmt.Fprintf(os.Stderr, "    %s\n", trimmedLine)

			// Show pointer to the error position
			if colNum > 0 {
				// Calculate visual column accounting for tabs (8 spaces each) up to error position
				visualCol := 0
				for i := 0; i < colNum-1 && i < len(sourceLine); i++ {
					if sourceLine[i] == '\t' {
						visualCol += 8
					} else {
						visualCol++
					}
				}

				// Adjust pointer position by subtracting trimmed columns
				adjustedCol := visualCol - trimCount
				if adjustedCol < 0 {
					adjustedCol = 0
				}

				pointer := strings.Repeat(" ", adjustedCol) + "^"
				fmt.Fprintf(os.Stderr, "    %s\n", pointer)
			}
		} else if _, err := fmt.Sscanf(msg, "line %d", &lineNum); err == nil && lineNum > 0 && lineNum <= len(lines) {
			// Fallback: show line without pointer if only line number available
			trimmedLine := strings.TrimLeft(lines[lineNum-1], " \t")
			fmt.Fprintf(os.Stderr, "    %s\n", trimmedLine)
		}
	}
}
