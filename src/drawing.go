package drawing

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

var seeded = false
var defaultBg, _, _ = tcell.StyleDefault.Decompose()
var defaultStyle = tcell.StyleDefault.Foreground(defaultBg)
var bodyStyle = defaultStyle.Foreground(tcell.NewHexColor(0x22aa33))
var headStyle = defaultStyle.Foreground(tcell.NewHexColor(0x66ff77))

type Cursor struct{ X, Y int }

func (c Cursor) Above(n int) Cursor {
	return Cursor{X: c.X, Y: c.Y - n}
}

func (c Cursor) Below(n int) Cursor {
	return Cursor{X: c.X, Y: c.Y + n}
}

type Layout struct {
	W, H int
	Cols []*Column
}

func (d *Layout) Sync(s tcell.Screen) {
	d.W, d.H = s.Size()
}

func (d Layout) TotalCells() int {
	return d.W * d.H
}

func (d Layout) HContains(cursor Cursor) bool {
	return cursor.Y >= 0 && cursor.Y < d.H
}

func (d Layout) IsOutOfBoundsBelow(cursor Cursor) bool {
	return cursor.Y > d.H
}

func MakeLayout(s tcell.Screen) *Layout {
	var layout Layout
	layout.Sync(s)

	layout.Cols = make([]*Column, layout.W)
	for i := 0; i < layout.H; i++ {
		layout.Cols[i] = MakeColumn(i)
	}

	return &layout
}

func (d Layout) Draw(s tcell.Screen) {
	for _, col := range d.Cols {
		col.Draw(s)
	}
}

func (d *Layout) Update() {
	for _, col := range d.Cols {
		col.Update()
	}
}

var dims *Layout

func SetStyleAt(s tcell.Screen, cursor Cursor, newStyle tcell.Style) {
	glyph, _, oldStyle, _ := s.GetContent(cursor.X, cursor.Y)
	if oldStyle != newStyle {
		s.SetContent(cursor.X, cursor.Y, glyph, nil, newStyle)
	}
}

type Raindrop struct {
	Length   int
	Progress Cursor
}

func (r Raindrop) TailCursor() Cursor {
	return r.Progress.Above(r.Length)
}

func (r Raindrop) Draw(s tcell.Screen) {
	dims.Sync(s)
	if dims.HContains(r.Progress) {
		SetStyleAt(s, r.Progress, headStyle)
		for i := 0; i < r.Length; i++ {
			SetStyleAt(s, r.Progress.Above(i), bodyStyle)
		}
		SetStyleAt(s, r.TailCursor(), defaultStyle)
	}
}

func (r *Raindrop) Update() {
	r.Progress = r.Progress.Below(1)
}

func MakeRaindrop(x int) *Raindrop {
	minLength := 4
	maxLength := 3 * dims.H / 4
	length := rand.Intn(maxLength-minLength) + minLength
	return &Raindrop{
		Length:   length,
		Progress: Cursor{X: x, Y: 0},
	}
}

type Column struct {
	Index     int
	Content   []rune
	Raindrops []*Raindrop
}

func MakeColumn(index int) *Column {
	seedRand()
	col := &Column{Index: index, Content: SelectRandomN(dims.H)}

	col.Raindrops[0] = MakeRaindrop(index)
	return col
}

func (c Column) Draw(s tcell.Screen) {
	for _, drop := range c.Raindrops {
		drop.Draw(s)
	}
}

func (c *Column) Update() {
	for _, drop := range c.Raindrops {
		drop.Update()
		if dims.IsOutOfBoundsBelow(drop.TailCursor()) {
			c.Raindrops = c.PopDrop()
		}
	}
	if len(c.Raindrops) < 2 && c.Raindrops[0].TailCursor().Y > 10 {
		c.Raindrops = c.PushDrop()
	}
}

func (c *Column) PushDrop() []*Raindrop {
	return append([]*Raindrop{MakeRaindrop(c.Index)}, c.Raindrops...)
}

func (c *Column) PopDrop() []*Raindrop {
	return c.Raindrops[:len(c.Raindrops)-1]
}

func Draw() {
	s := setup()
	dims = MakeLayout(s)

	s.SetStyle(defaultStyle)
	s.Show()
	defer func() {
		s.Fini()
		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second / 30)
		dims.Sync(s)
		dims.Draw(s)
		s.Show()
		ev := s.PollEvent()
		switch et := ev.(type) {
		case *tcell.EventKey:
			if et.Key() == tcell.KeyESC || et.Key() == tcell.KeyCtrlC || et.Rune() == 'Q' {
				return
			}
		}
		dims.Update()
	}
}

func setup() tcell.Screen {
	seedRand()
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	dims.Sync(s)
	s.Clear()
	return s
}

func seedRand() {
	if !seeded {
		rand.Seed(time.Now().UnixNano())
		seeded = true
	}
}
