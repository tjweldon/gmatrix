package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
	"tjweldon/gmatrix/src"
)

var debugOut string

var Charset []rune

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

	debugOut += string(drawing.GetCharset())

	drawing.Draw()

	fmt.Printf("%s", debugOut)
	fmt.Print("\n")
}
