package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

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
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
