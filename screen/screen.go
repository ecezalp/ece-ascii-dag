package screen

import (
	. "ece-ascii-dag/util"
	"fmt"
)

type Screen struct {
	width    int
	height   int
	runes    [][]rune
	BoxStyle map[string]rune
}

func (s *Screen) ReadRune(x, y int) *rune {
	return &s.runes[y][x]
}

func (s *Screen) PlaceRune(x, y int, rune rune) {
	s.runes[y][x] = rune
}

func (s *Screen) PlaceWord(x, y int, word []rune) {
	for i, rune := range word {
		s.PlaceRune(x+i, y, rune)
	}
}

func (s *Screen) PlaceHorizontalLine(left, right, y int, rune rune) {
	for x := left; x <= right; x++ {
		s.PlaceRune(x, y, rune)
	}
}

func (s *Screen) PlaceVerticalLine(top, bottom, x int, rune rune) {
	for y := top; y <= bottom; y++ {
		s.PlaceRune(x, y, rune)
	}
}

func (s *Screen) PlaceBox(x, y, width, height int) {
	// determine box coordinates
	boxStartY := y
	boxEndY := y + height - 1
	boxStartX := x
	boxEndX := x + width - 1

	// place corners
	s.PlaceRune(boxStartX, boxStartY, s.BoxStyle["topLeftCorner"])
	s.PlaceRune(boxEndX, boxStartY, s.BoxStyle["topLeftCorner"])
	s.PlaceRune(boxStartX, boxEndY, s.BoxStyle["bottomLeftCorner"])
	s.PlaceRune(boxEndX, boxEndY, s.BoxStyle["bottomRightCorner"])

	// place lines
	s.PlaceHorizontalLine(boxStartX, boxEndX, boxStartY, s.BoxStyle["horizontalLine"])
	s.PlaceHorizontalLine(boxStartX, boxEndX, boxEndY, s.BoxStyle["horizontalLine"])
	s.PlaceVerticalLine(boxStartY, boxEndY, boxStartX, s.BoxStyle["verticalLine"])
	s.PlaceVerticalLine(boxStartY, boxEndY, boxEndY, s.BoxStyle["verticalLine"])
}

func (s *Screen) PlaceTextBox(x, y int, word []rune) {
	s.PlaceBox(x, y, len(word)+2, 3)
	s.PlaceWord(x+1, y+1, word)
}

func (s *Screen) String() string {
	var stringVal string
	for _, runeLine := range s.runes {
		stringVal += string(runeLine) + "\n"
	}
	return stringVal
}

func NewScreen(width int, height int) *Screen {
	fmt.Printf("Screen Width %v", width)
	fmt.Printf("Screen Height %v", height)

	runes := make([][]rune, width)
	for i := range runes {
		runes[i] = make([]rune, height)
	}

	return &Screen{
		width:    width,
		height:   height,
		runes:    runes,
		BoxStyle: GetDefaultBoxRuneMap(),
	}
}
