package parser

import (
	"fmt"
	"os"
	"strings"
)

var traceLevel int = 0

func traceEnabled() bool {
	return os.Getenv("TRACE") == "1"
}

func identLevel() string {
	return strings.Repeat("  ", traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() {
	traceLevel += 1
}

func decIdent() {
	traceLevel -= 1
}

func trace(msg string) string {
	if !traceEnabled() {
		return ""
	}

	incIdent()
	tracePrint("↱ " + msg)
	return msg
}

func untrace(msg string) {
	if !traceEnabled() {
		return
	}

	tracePrint("↳ " + msg)
	decIdent()
}
