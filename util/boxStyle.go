package util

func GetDefaultBoxRuneMap() map[string]rune {
	boxRuneMap := make(map[string]rune)
	boxRuneMap["topRightCorner"] = '┐'
	boxRuneMap["bottomRightCorner"] = '┘'
	boxRuneMap["topLeftCorner"] = '┌'
	boxRuneMap["bottomLeftCorner"] = '└'
	boxRuneMap["horizontalLine"] = '─'
	boxRuneMap["verticalLine"] = '│'
	return boxRuneMap
}
