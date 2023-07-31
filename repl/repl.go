package repl

import (
	"bufio"
	"fmt"
	"fungo/lexer"
	"fungo/parser"
	"io"
)

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, "#: ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		parser := parser.New(lexer.New(line))
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		io.WriteString(out, program.String()+"\n")
	}
}
