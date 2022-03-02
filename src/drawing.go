package drawing

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

var seeded = false

var defaultStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack)
var bodyStyle = defaultStyle.Foreground(tcell.NewHexColor(0x119922))
var headStyle = defaultStyle.Foreground(tcell.NewHexColor(0x88ff99))
var sparkStyle = defaultStyle.Foreground(tcell.NewHexColor(0xaaffcc))

// Cursor is a screen location
type Cursor struct{ X, Y int }

func (c Cursor) Above(n int) Cursor {
	return Cursor{X: c.X, Y: c.Y - n}
}

func (c Cursor) Below(n int) Cursor {
	return Cursor{X: c.X, Y: c.Y + n}
}

// Layout represents the state of the screen
type Layout struct {
	W, H   int
	Cols   []*Column
	Sparks []*Spark
}

func NewLayout(s tcell.Screen) *Layout {
	var layout Layout
	layout.Sync(s)
	dims = &layout

	layout.Cols = make([]*Column, layout.W)
	for i := 0; i < layout.W; i++ {
		layout.Cols[i] = newColumn(i)
	}

	layout.Sparks = make([]*Spark, len(layout.Cols)/10)
	for i := 0; i < len(layout.Sparks); i++ {
		layout.Sparks[i] = newSpark(60)
	}
	return &layout
}

func (d *Layout) RandomCursor() Cursor {
	return Cursor{
		X: rand.Intn(d.W),
		Y: rand.Intn(d.H),
	}
}

func (d *Layout) Sync(s tcell.Screen) {
	d.W, d.H = s.Size()
}

func (d Layout) TotalCells() int {
	return d.W * d.H
}

func (d Layout) VContains(cursor Cursor) bool {
	return cursor.Y >= 0 && cursor.Y < d.H
}

func (d Layout) IsOutOfBoundsBelow(cursor Cursor) bool {
	return cursor.Y > d.H
}

func (d *Layout) Draw(s tcell.Screen) {
	for _, col := range d.Cols {
		col.Draw(s)
	}
	for _, spark := range d.Sparks {
		spark.Draw(s)
	}
}

func (d *Layout) Update() {
	for _, col := range d.Cols {
		col.Update()
	}
	for _, spark := range d.Sparks {
		spark.Update()
	}
}

func (d *Layout) SetStyleAt(s tcell.Screen, cursor Cursor, newStyle tcell.Style) {
	_, _, oldStyle, _ := s.GetContent(cursor.X, cursor.Y)
	if oldStyle != newStyle && d.Contains(cursor) {
		cursorContent := d.getRuneAt(cursor)
		s.SetContent(cursor.X, cursor.Y, cursorContent, nil, newStyle)
	}
}

func (d *Layout) getRuneAt(cursor Cursor) rune {
	cursorContent := d.Cols[cursor.X].Content[cursor.Y]
	return cursorContent
}

func (d *Layout) setRuneAt(glyph rune, cursor Cursor) {
	if d.Contains(cursor) {
		d.Cols[cursor.X].Content[cursor.Y] = glyph
	}
}

func (d *Layout) Contains(c Cursor) bool {
	return c.X >= 0 && c.X < d.W && c.Y >= 0 && c.Y < d.H
}

// dims is the shared module global layout
var dims *Layout

// Column represents a single character wide vertical
// line of cells
type Column struct {
	Index     int
	Content   []rune
	Raindrops []*Raindrop
}

func newColumn(index int) *Column {
	seedRand()
	randomN := SelectRandomN(dims.H)
	col := Column{Index: index, Content: randomN}

	col.Raindrops = make([]*Raindrop, 1)
	col.Raindrops[0] = newRaindrop(index)
	col.Raindrops[0].Progress.Y = rand.Intn(dims.H / 2)
	return &col
}

func (c *Column) Draw(s tcell.Screen) {
	for y, r := range c.Content {
		s.SetContent(c.Index, y, r, nil, defaultStyle)
	}
	for _, drop := range c.Raindrops {
		drop.Draw(s)
	}
}

func (c *Column) Update() {
	for _, drop := range c.Raindrops {
		drop.Update()
		if dims.IsOutOfBoundsBelow(drop.TailCursor()) {
			c.Raindrops = c.PopDrop()
		} else if dims.VContains(drop.Progress) {
			newContent := c.Content
			newContent[drop.Progress.Y] = SelectRand()
			c.Content = newContent
		}
	}
	if len(c.Raindrops) < 2 && c.Raindrops[0].TailCursor().Y > 15 {
		c.Raindrops = c.PushDrop()
	}
}

func (c *Column) PushDrop() []*Raindrop {
	return append([]*Raindrop{newRaindrop(c.Index)}, c.Raindrops...)
}

func (c *Column) PopDrop() []*Raindrop {
	return c.Raindrops[:len(c.Raindrops)-1]
}

// Raindrop represents the moving bright green vertical area
// of styling
type Raindrop struct {
	Length   int
	Progress Cursor
}

func newRaindrop(x int) *Raindrop {
	minLength := 4
	maxLength := 3 * dims.H / 4
	length := rand.Intn(maxLength-minLength) + minLength
	return &Raindrop{
		Length:   length,
		Progress: Cursor{X: x, Y: 0},
	}
}

func (r Raindrop) TailCursor() Cursor {
	return r.Progress.Above(r.Length)
}

func (r *Raindrop) Draw(s tcell.Screen) {
	if dims.VContains(r.TailCursor()) || dims.VContains(r.Progress) {
		dims.SetStyleAt(s, r.Progress, headStyle)
		for i := 0; i < r.Length; i++ {
			if r.Progress.Y-i == 0 {
				break
			}
			dims.SetStyleAt(s, r.Progress.Above(i+1), bodyStyle)
		}
		dims.SetStyleAt(s, r.TailCursor(), defaultStyle)
	}
}

func (r *Raindrop) Update() {
	r.Progress = r.Progress.Below(1)
}

func (r *Raindrop) ParentColumn() *Column {
	return dims.Cols[r.Progress.X]
}

type Spark struct {
	Position      Cursor
	Age, Lifespan int
}

func newSpark(lifespan int) *Spark {
	position := dims.RandomCursor()
	return &Spark{
		Position: position,
		Age:      0,
		Lifespan: lifespan,
	}
}

func (s *Spark) Update() {
	s.Age++
	if s.isDead() {
		s.Position = dims.RandomCursor()
		s.Age = 0
	}
	dims.setRuneAt(SelectRand(), s.Position)
}

func (s *Spark) Draw(screen tcell.Screen) {
	dims.SetStyleAt(screen, s.Position, sparkStyle)
}

func (s *Spark) isDead() bool {
	isDead := s.Age > s.Lifespan
	return isDead
}

// Draw is the public interface of the module.
// Hey homes, call this to enter the matrix, g.
func Draw() {
	s := Setup()
	dims = NewLayout(s)

	s.SetStyle(defaultStyle)
	s.Show()

	quit := make(chan struct{})
	defer func() {
		quit <- struct{}{}
		close(quit)
	}()
	events := make(chan tcell.Event)
	go func(q chan struct{}, e chan tcell.Event) {
		for {
			select {
			case <-q:
				close(events)
				return
			default:
				e <- s.PollEvent()
			}
		}
	}(quit, events)

	i := 0
	//debugOut = fmt.Sprintf("LAYOUT: %v\n", *dims)

	defer func() {
		s.Fini()
		os.Exit(0)
	}()
	for {
		time.Sleep(time.Second / 30)
		dims.Sync(s)
		dims.Draw(s)

		s.SetContent(dims.W-2, dims.H-2, rune(i+0x20), nil, defaultStyle.Foreground(tcell.ColorDarkOrange))

		select {
		case ev := <-events:
			switch et := ev.(type) {
			case *tcell.EventKey:
				if et.Key() == tcell.KeyESC || et.Key() == tcell.KeyCtrlC || et.Rune() == 'Q' {
					return
				} else {
					s.SetContent(dims.W/2, dims.H/2, et.Rune(), nil, defaultStyle.Background(tcell.ColorDarkOrange))
				}
			}
		default:
			// Nothing
		}

		s.Show()
		dims.Update()
		i = (i + 1) % 0xff
	}
}

// Setup is called in draw, fahgetabaahtit
func Setup() tcell.Screen {
	seedRand()
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.Clear()
	return s
}

func seedRand() {
	if !seeded {
		rand.Seed(time.Now().UnixNano())
		seeded = true
	}
}
