package main

import (
	"os"

	"pars/pkg/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
