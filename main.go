package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"time"
	"tjweldon/gmatrix/src"
)

var debugOut string

var Charset []rune

var args struct {
	DumpCharset string `arg:"--dump-src"`
}

func main() {
	arg.MustParse(&args)

	// Dump the src and exit
	if args.DumpCharset != "" {
		src.DumpCharset(args.DumpCharset)
		os.Exit(0)
	}

	Charset = src.GetCharset()
	debugOut += string(Charset)

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.Clear()
	_, bg, _ := tcell.StyleDefault.Decompose()
	style := tcell.StyleDefault.Foreground(bg)
	s.SetStyle(style)
	s.Show()
	w, h := s.Size()
	maxLen := w * h
	for j, cell := range Charset {
		if j > maxLen {
			break
		}
		x := j%w + 1
		y := j / w

		s.SetContent(x, y, cell, nil, style.Foreground(tcell.NewHexColor(0x22aa33)))
	}

	s.Show()
	time.Sleep(time.Duration(0) * time.Second)
	s.Fini()

	fmt.Printf("%s", debugOut)
	fmt.Print("\n")
}
