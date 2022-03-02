package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
	"tjweldon/gmatrix/src"
)

var someOtherDebugOut string

var args struct {
	DumpCharset string `arg:"--dump-charset"`
}

func main() {
	arg.MustParse(&args)

	// Dump the charset and exit
	if args.DumpCharset != "" {
		drawing.DumpCharset(args.DumpCharset)
		os.Exit(0)
	}

	someOtherDebugOut += string(drawing.GetCharset())

	drawing.Draw()

	fmt.Printf("%s", someOtherDebugOut)
	fmt.Print("\n")
}
