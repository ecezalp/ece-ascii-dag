package main

import (
	. "ece-ascii-dag/screen"
	"fmt"
)

func main() {
	screen := NewScreen(40, 30)
	fmt.Sprintf("%v", screen)
}
