package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

const PROMPT = ">> "
const CONTINUATION_PROMPT = ".. "

const PARSER_LOGO = `
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ `

// Parsley keywords and builtins for tab completion
var completionWords = []string{
	// Keywords
	"let", "if", "else", "for", "in", "fn", "return", "export", "import",
	// Builtins - I/O
	"log", "logLine", "file", "dir", "JSON", "CSV", "MD", "SVG", "HTML",
	"text", "lines", "bytes", "SFTP", "Fetch", "SQL",
	// Builtins - Collections
	"len", "keys", "values", "type", "sort", "reverse", "join",
	// Builtins - Strings
	"split", "trim", "upper", "lower", "contains", "startsWith", "endsWith",
	"replace", "match", "test",
	// Builtins - Math
	"abs", "floor", "ceil", "round", "sqrt", "pow", "sin", "cos", "tan",
	"min", "max", "sum",
	// Builtins - DateTime
	"now", "date", "time", "duration", "format", "parse",
	// Builtins - Other
	"range", "glob", "toString",
	// Common values
	"true", "false", "null",
}

// Start starts the REPL with line editing, history, and tab completion
func Start(in io.Reader, out io.Writer, version string) {
	line := liner.NewLiner()
	defer line.Close()

	// Enable Ctrl+C to abort current line
	line.SetCtrlCAborts(true)

	// Set up tab completion
	line.SetCompleter(func(line string) []string {
		return filterCompletions(line)
	})

	// Load command history from file
	historyFile := filepath.Join(os.TempDir(), ".parsley_history")
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// Save history on exit
	defer func() {
		if f, err := os.Create(historyFile); err == nil {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	env := evaluator.NewEnvironment()

	fmt.Fprintf(out, "%s", PARSER_LOGO)
	fmt.Fprintln(out, "v", version)
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Type 'exit' or Ctrl+D to quit")
	fmt.Fprintln(out, "Use Tab for completion, ↑↓ for history")
	fmt.Fprintln(out, "")

	var inputBuffer strings.Builder
	currentPrompt := PROMPT

	for {
		input, err := line.Prompt(currentPrompt)
		if err != nil {
			// Ctrl+D or Ctrl+C
			if err == liner.ErrPromptAborted {
				// Ctrl+C - clear any buffered input and return to main prompt
				if inputBuffer.Len() > 0 {
					fmt.Fprintln(out, "^C (cleared)")
				} else {
					fmt.Fprintln(out, "^C")
				}
				inputBuffer.Reset()
				currentPrompt = PROMPT
				continue
			}
			if err == io.EOF {
				// Ctrl+D - exit
				fmt.Fprintln(out, "\nGoodbye!")
				return
			}
			fmt.Fprintf(out, "Error reading input: %v\n", err)
			continue
		}

		// Check for exit command
		trimmed := strings.TrimSpace(input)
		if inputBuffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			fmt.Fprintln(out, "Goodbye!")
			return
		}

		// Skip empty lines when no input buffered
		if inputBuffer.Len() == 0 && trimmed == "" {
			continue
		}

		// Add to input buffer
		if inputBuffer.Len() > 0 {
			inputBuffer.WriteString("\n")
		}
		inputBuffer.WriteString(input)

		// Check if input is complete (no unclosed braces/brackets)
		fullInput := inputBuffer.String()
		if needsMoreInput(fullInput) {
			// Continue multi-line input
			currentPrompt = CONTINUATION_PROMPT
			continue
		}

		// Input is complete - parse and evaluate
		currentPrompt = PROMPT

		// Add complete input to history
		if trimmed != "" {
			line.AppendHistory(fullInput)
		}

		// Parse and evaluate
		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if errors := p.Errors(); len(errors) != 0 {
			printParserErrors(out, errors)
			inputBuffer.Reset()
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		// Clear buffer for next input
		inputBuffer.Reset()
	}
}

// filterCompletions returns completion suggestions based on current input
func filterCompletions(line string) []string {
	// Get the last word being typed
	words := strings.Fields(line)
	if len(words) == 0 {
		return completionWords
	}

	lastWord := words[len(words)-1]

	// If line ends with space, no completion
	if len(line) > 0 && line[len(line)-1] == ' ' {
		return nil
	}

	var matches []string
	for _, word := range completionWords {
		if strings.HasPrefix(word, lastWord) {
			matches = append(matches, word)
		}
	}
	return matches
}

// needsMoreInput checks if the input has unclosed braces, brackets, or parentheses
func needsMoreInput(input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return false
	}

	braceCount := 0
	bracketCount := 0
	parenCount := 0
	inString := false
	escapeNext := false

	for i, ch := range input {
		if escapeNext {
			escapeNext = false
			continue
		}

		if ch == '\\' {
			escapeNext = true
			continue
		}

		// Track string state to ignore braces inside strings
		if ch == '"' && (i == 0 || input[i-1] != '\\') {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch ch {
		case '{':
			braceCount++
		case '}':
			braceCount--
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		case '(':
			parenCount++
		case ')':
			parenCount--
		}
	}

	// Need more input if any are unclosed
	return braceCount > 0 || bracketCount > 0 || parenCount > 0
}

// printParserErrors prints parser errors
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
