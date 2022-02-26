package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"strings"
	"time"
	"unicode"
)

func main() {
	lower, upper := unicode.Braille.R16[0].Lo, unicode.Braille.R16[0].Hi

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
	var content string
	for i := lower; i <= upper; i++ {
		content += strings.Trim(string(rune(i)), " ")
	}
	for j, cell := range []rune(content) {
		if j > maxLen {
			break
		}
		x := j%w + 1
		y := j / w

		s.SetContent(x, y, cell, nil, style.Foreground(tcell.NewHexColor(0x22aa33)))
	}
	s.Show()
	time.Sleep(time.Duration(3) * time.Second)
	s.Fini()

	fmt.Println(content)
}
