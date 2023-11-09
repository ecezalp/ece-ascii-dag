package util

func GetDefaultBoxRuneMap() map[string]rune {
	var boxRuneMap map[string]rune
	boxRuneMap["topRightCorner"] = '┐'
	boxRuneMap["bottomRightCorner"] = '┘'
	boxRuneMap["topLeftCorner"] = '┌'
	boxRuneMap["bottomLeftCorner"] = '└'
	boxRuneMap["horizontalLine"] = '─'
	boxRuneMap["verticalLine"] = '│'
	return boxRuneMap
}
