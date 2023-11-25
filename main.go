package main

import (
	. "ece-ascii-dag/dag"
)

func main() {
	DAGtoText("chrome -> content\nchrome -> blink\nchrome -> base\n\ncontent -> blink\ncontent -> net\ncontent -> base\n\nblink -> v8\nblink -> CC\nblink -> WTF\nblink -> skia\nblink -> base\nblink -> net\n\nweblayer -> content\nweblayer -> chrome\nweblayer -> base\n\nnet -> base\nWTF -> base\nweblayer -> v8")
}
