package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/formatter"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
	"github.com/sambeau/parsley/pkg/repl"
)

// Version is set at compile time via -ldflags
var Version = "0.9.4"

var (
	// Display flags
	helpFlag       = flag.Bool("h", false, "Show help message")
	helpLongFlag   = flag.Bool("help", false, "Show help message")
	versionFlag    = flag.Bool("V", false, "Show version information")
	versionLongFlag = flag.Bool("version", false, "Show version information")
	prettyPrintFlag = flag.Bool("pp", false, "Pretty-print HTML output")
	prettyLongFlag  = flag.Bool("pretty", false, "Pretty-print HTML output")
	
	// Security flags
	restrictReadFlag     = flag.String("restrict-read", "", "Comma-separated read blacklist paths")
	noReadFlag           = flag.Bool("no-read", false, "Deny all file reads")
	allowWriteFlag       = flag.String("allow-write", "", "Comma-separated write whitelist paths")
	allowWriteAllFlag    = flag.Bool("allow-write-all", false, "Allow unrestricted writes")
	allowWriteAllShort   = flag.Bool("w", false, "Shorthand for --allow-write-all")
	allowExecuteFlag     = flag.String("allow-execute", "", "Comma-separated execute whitelist paths")
	allowExecuteAllFlag  = flag.Bool("allow-execute-all", false, "Allow unrestricted executes")
	allowExecuteAllShort = flag.Bool("x", false, "Shorthand for --allow-execute-all")
)

func main() {
	// Customize flag usage message
	flag.Usage = printHelp
	flag.Parse()

	// Check for help flag
	if *helpFlag || *helpLongFlag {
		printHelp()
		os.Exit(0)
	}

	// Check for version flag
	if *versionFlag || *versionLongFlag {
		fmt.Printf("pars version %s\n", Version)
		os.Exit(0)
	}

	// Get filename from remaining args
	args := flag.Args()
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	// Determine pretty print setting
	prettyPrint := *prettyPrintFlag || *prettyLongFlag

	if filename != "" {
		// File execution mode
		executeFile(filename, prettyPrint)
	} else {
		// REPL mode
		repl.Start(os.Stdin, os.Stdout, Version)
	}
}

func printHelp() {
	fmt.Printf(`pars - Parsley language interpreter version %s

Usage:
  pars [options] [file]

Display Options:
  -h, --help            Show this help message
  -V, --version         Show version information
  -pp, --pretty         Pretty-print HTML output with proper indentation

Security Options:
  --restrict-read=PATHS     Deny reading from comma-separated paths
  --no-read                 Deny all file reads
  --allow-write=PATHS       Allow writing to comma-separated paths
  --allow-write-all, -w     Allow unrestricted writes
  --allow-execute=PATHS     Allow executing scripts from paths
  --allow-execute-all, -x   Allow unrestricted script execution

Security Examples:
  pars -w script.pars                           # Allow all writes
  pars --allow-write=./output script.pars       # Allow writes to ./output only
  pars -x --allow-write=./data script.pars      # Allow all executes, writes to ./data
  pars --restrict-read=/etc script.pars         # Deny reads from /etc

Examples:
  pars                      Start interactive REPL
  pars script.pars          Execute a Parsley script
  pars -pp page.pars        Execute and pretty-print HTML output

For more information, visit: https://github.com/sambeau/parsley
`, Version)
}

// executeFile reads and executes a pars source file
func executeFile(filename string, prettyPrint bool) {
	// Build security policy (always create one to enable default restrictions)
	policy, err := buildSecurityPolicy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

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
	env.Security = policy
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
		output := evaluator.ObjectToPrintString(evaluated)

		// Apply HTML formatting if --pp flag is set
		if prettyPrint {
			output = formatter.FormatHTML(output)
		}

		fmt.Println(output)
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

// buildSecurityPolicy creates a SecurityPolicy from command-line flags
func buildSecurityPolicy() (*evaluator.SecurityPolicy, error) {
	policy := &evaluator.SecurityPolicy{
		NoRead:          *noReadFlag,
		AllowWriteAll:   *allowWriteAllFlag || *allowWriteAllShort,
		AllowExecuteAll: *allowExecuteAllFlag || *allowExecuteAllShort,
	}

	// Parse restrict list
	if *restrictReadFlag != "" {
		paths, err := parseAndResolvePaths(*restrictReadFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid --restrict-read: %s", err)
		}
		policy.RestrictRead = paths
	}

	// Parse allow lists
	if *allowWriteFlag != "" {
		paths, err := parseAndResolvePaths(*allowWriteFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid --allow-write: %s", err)
		}
		policy.AllowWrite = paths
	}

	if *allowExecuteFlag != "" {
		paths, err := parseAndResolvePaths(*allowExecuteFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid --allow-execute: %s", err)
		}
		policy.AllowExecute = paths
	}

	return policy, nil
}

// parseAndResolvePaths parses comma-separated paths and resolves them to absolute paths
func parseAndResolvePaths(pathList string) ([]string, error) {
	parts := strings.Split(pathList, ",")
	resolved := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Expand home directory
		if strings.HasPrefix(p, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("cannot expand ~: %s", err)
			}
			p = filepath.Join(home, p[2:])
		}

		// Convert to absolute path
		absPath, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("invalid path %s: %s", p, err)
		}

		// Clean path
		absPath = filepath.Clean(absPath)

		resolved = append(resolved, absPath)
	}

	return resolved, nil
}
