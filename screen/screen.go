package screen

func GetBoxRuneMap() map[string]rune {
	var boxRuneMap map[string]rune
	boxRuneMap["topRightCorner"] = '┐'
	boxRuneMap["bottomRightCorner"] = '┘'
	boxRuneMap["topLeftCorner"] = '┌'
	boxRuneMap["bottomLeftCorner"] = '└'
	boxRuneMap["horizontalLine"] = '─'
	boxRuneMap["verticalLine"] = '│'
	return boxRuneMap
}

type Screen struct {
	width  int
	height int
	runes  [][]rune
}

func (s *Screen) PlaceChar(x, y int, rune rune) {
	s.runes[y][x] = rune
}

func (s *Screen) PlaceWord(x, y int, word []rune) {
	for i, rune := range word {
		s.PlaceChar(x+i, y, rune)
	}
}

func (s *Screen) PlaceBox(x, y, width, height int) {
	boxRuneMap := GetBoxRuneMap()

	// determine box coordinates
	boxStartY := y
	boxEndY := y + height - 1
	boxStartX := x
	boxEndX := x + width - 1

	// place corners
	s.PlaceChar(boxStartX, boxStartY, boxRuneMap["topLeftCorner"])
	s.PlaceChar(boxEndX, boxStartY, boxRuneMap["topLeftCorner"])
	s.PlaceChar(boxStartX, boxEndY, boxRuneMap["bottomLeftCorner"])
	s.PlaceChar(boxEndX, boxEndY, boxRuneMap["bottomRightCorner"])

	// place lines
	for xLine := 1; xLine < width-1; xLine++ {
		s.PlaceChar(boxStartX+xLine, boxStartY, boxRuneMap["horizontalLine"])
		s.PlaceChar(boxStartX+xLine, boxEndY, boxRuneMap["horizontalLine"])
	}
	for yLine := 1; yLine < height-1; yLine++ {
		s.PlaceChar(boxStartX, boxStartY+yLine, boxRuneMap["verticalLine"])
		s.PlaceChar(boxEndX, boxStartY+yLine, boxRuneMap["verticalLine"])
	}
}

func (s *Screen) PlaceTextBox(x, y int, word []rune) {
	s.PlaceWord(x+1, y+1, word)
	s.PlaceBox(x, y, len(word)+2, 3)
}

func NewScreen(width int, height int) *Screen {
	runes := make([][]rune, width)
	for i := range runes {
		runes[i] = make([]rune, height)
	}

	return &Screen{
		width:  width,
		height: height,
		runes:  runes,
	}
}
