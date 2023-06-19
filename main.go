package main

import (
	"fungo/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
