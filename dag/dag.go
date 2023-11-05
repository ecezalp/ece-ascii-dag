package dag

type Node struct {
	columnIndex              int
	coordinateX              int
	coordinateY              int
	downwardNodeIds          map[string]int
	downwardNodeIdsSorted    []int
	isConnectorNode          bool
	nodeIdsInDownwardClosure map[string]int
	padding                  int
	rowIndex                 int
	upwardNodeIds            map[string]int
	upwardNodeIdsSorted      []int
	visualHeight             int
	visualWidth              int
}

type Edge struct {
	upwardNodeId   int
	downwardNodeId int
	coordinateX    int
	coordinateY    int
}

func (a Edge) isEqualTo(b Edge) bool {
	return a.upwardNodeId == b.upwardNodeId &&
		a.downwardNodeId == b.downwardNodeId
}

type Adapter struct {
	coordinateY      int
	inputNodeIds     []map[int]struct{}
	isEnabled        bool
	outputNodeIds    []map[int]struct{}
	visualCharacters [][]string
	visualHeight     int
}

func (a *Adapter) Construct() {
	// Implement the Construct function here
}

func (a *Adapter) Render(screen *Screen) {
	// Implement the Render function here
}
