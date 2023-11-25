package dag

import (
	. "ece-ascii-dag/screen"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewAdapter(inputs []IntSet, outputs []IntSet) Adapter {
	solutionFound := false
	height := 3
	width := len(inputs)
	var adapterNodes []AdapterNode
	var adapterEdges []AdapterEdge
	rendering := make([][]rune, height)

	connectorLength := 0
	for _, inputItem := range inputs {
		inputSlice := inputItem.Slice()
		connectorLength = max(connectorLength, len(inputSlice))
	}
	for !solutionFound {
		adapterNodes = make([]AdapterNode, width*height*2)
		adapterEdges = make([]AdapterEdge, width*height*3)

		IndexFunc := func(x, y, layer int) int {
			return x + width*(y+height*layer)
		}

		ConnectFunc := func(edge *AdapterEdge, a *AdapterNode, b *AdapterNode, weight int) {
			edge.A = a
			edge.B = b
			edge.Weight = weight
			a.Edges = append(a.Edges, edge)
			b.Edges = append(b.Edges, edge)
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {

				// vertical
				if y != height-1 {
					ConnectFunc(
						&adapterEdges[IndexFunc(x, y+0, 0)],
						&adapterNodes[IndexFunc(x, y+0, 0)],
						&adapterNodes[IndexFunc(x, y+1, 0)],
						1,
					)
				}

				// horizontal
				if y >= 1 && y <= height-3 && x != width-1 {
					ConnectFunc(
						&adapterEdges[IndexFunc(x+0, y, 1)],
						&adapterNodes[IndexFunc(x+0, y, 1)],
						&adapterNodes[IndexFunc(x+1, y, 1)],
						1,
					)
				}

				// corners
				dy := height/2 - y
				ConnectFunc(
					&adapterEdges[IndexFunc(x, y, 2)],
					&adapterNodes[IndexFunc(x, y, 0)],
					&adapterNodes[IndexFunc(x, y, 1)],
					10+dy*dy,
				)
			}
		}

		// try a solution
		solutionFound = true

		// add path one by one
		for connectorId := 1; connectorId < connectorLength; connectorId++ {
			bigNumber := 1 << 15

			// clear previous costs
			for i, node := range adapterNodes {
				node.Visited = false
				node.Cost = bigNumber
				adapterNodes[i] = node
			}

			start := NewNodeSet()
			end := NewNodeSet()

			for xVal := 0; xVal < width; xVal++ {
				if inputs[xVal].Contains(connectorId) {
					start.Add(&adapterNodes[IndexFunc(xVal, 0, 0)])
				}
				if outputs[xVal].Contains(connectorId) {
					end.Add(&adapterNodes[IndexFunc(xVal, height-1, 0)])
				}

				var pending NodeAndCostHeap
				for node := range start {
					pending.Push(NodeAndCost{
						Node: node,
						Cost: 0,
					})
				}

				for len(pending) != 0 {
					item := pending.Pop()
					node := item.Node
					if node.Visited {
						continue
					}
					node.Visited = true
					node.Cost = item.Cost
					for _, edge := range node.Edges {
						oppositeNode := edge.A
						if edge.A == node {
							oppositeNode = edge.B
						}
						if oppositeNode.Visited {
							continue
						}
						if edge.Assigned != 0 {
							continue
						}
						pending.Push(
							NodeAndCost{
								Node: oppositeNode,
								Cost: node.Cost + edge.Weight,
							},
						)
					}
				}

				// Reconstruct the path from end to start.
				bestScore := bigNumber
				var currentNode *AdapterNode
				for endNode := range end {
					if bestScore >= endNode.Cost {
						bestScore = endNode.Cost
						currentNode = endNode
					}
				}

				// No path found.
				if bestScore == bigNumber {
					solutionFound = false
					continue
				}

				for !start.Contains(currentNode) {
					for _, edge := range currentNode.Edges {
						oppositeNode := edge.A
						if edge.A == currentNode {
							oppositeNode = edge.B
						}
						if currentNode.Cost == oppositeNode.Cost+edge.Weight {
							edge.Assigned = connectorId
							currentNode = oppositeNode
						}
					}
				}

				isAssignedFunc := func(x int, y int, layer int) bool {
					return adapterEdges[IndexFunc(x, y, layer)].Assigned != 0
				}

				for itemY := 0; itemY < height; itemY++ {
					for itemX := 0; itemX < width; itemX++ {
						if isAssignedFunc(itemX, itemY, 0) {
							adapterEdges[IndexFunc(itemX, itemY, 1)].Weight = 20
						}
						if isAssignedFunc(itemX, itemY, 1) {
							adapterEdges[IndexFunc(itemX, itemY, 0)].Weight = 20
						}
					}
				}

				if height > 30 {
					solutionFound = true
				}

				if !solutionFound {
					height++
					continue
				}

				for i := range rendering {
					rendering[i] = make([]rune, width)
					for j := range rendering[i] {
						rendering[i][j] = ' '
					}
				}

				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						if isAssignedFunc(x, y, 1) {
							rendering[y][x] = '─'
						}
						if isAssignedFunc(x, y, 0) {
							rendering[y][x] = '│'
						}
						if isAssignedFunc(x, y, 2) {
							if isAssignedFunc(x, y, 0) {
								if isAssignedFunc(x, y, 1) {
									rendering[y][x] = '┌'
								} else {
									rendering[y][x] = '┐'
								}
							} else {
								if isAssignedFunc(x, y, 1) {
									rendering[y][x] = '└'
								} else {
									rendering[y][x] = '┘'
								}
							}
						}
					}
				}
			}
		}
	}
	return Adapter{
		Inputs:    inputs,
		Outputs:   outputs,
		IsEnabled: false,
		Runes:     rendering,
		Height:    height,
	}
}

type Context struct {
	Layers         []Layer
	NodeLabelIdMap map[string]int
	NodeLabels     []string // node names
	Nodes          []Node
}

func DAGtoText(input string) string {
	var context Context
	context.NodeLabelIdMap = make(map[string]int)
	return context.Process(input)
}

func printStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// Make sure s is a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fmt.Printf("%s: %v\n", fieldName, field.Interface())
	}
}

func (c *Context) Process(input string) string {
	c.Parse(input)
	if len(c.Nodes) == 0 {
		return ""
	}
	if !c.TopoSort() {
		return "There are cycles"
	}

	c.Complete()
	c.AddToLayers()
	c.ResolveCrossingEdges()
	c.Layout()

	for _, node := range c.Nodes {
		printStructFields(node)
		fmt.Println()
	}

	return c.Render()
}

func (c *Context) Parse(input string) {
	for _, line := range strings.Split(input, "\n") {
		parts := strings.Split(line, "->")
		for i := range parts {
			name := strings.TrimSpace(parts[i])
			c.CreateNode(name)
			if i > 0 {
				c.AddVertex(parts[i-1], name)
			}
		}
	}
}

func (c *Context) CreateNode(label string) {
	// handle already created
	if c.NodeLabelIdMap[label] != 0 {
		return
	}

	// create empty node
	newNodeId := len(c.Nodes)
	c.Nodes = append(c.Nodes, Node{
		DownwardClosure:       NewIntSet(),
		DownwardNodeIds:       NewIntSet(),
		DownwardNodeIdsSorted: []int{},
		Height:                0,
		IsConnector:           false,
		Layer:                 0,
		Padding:               1,
		Row:                   0,
		UpwardNodeIds:         NewIntSet(),
		UpwardNodeIdsSorted:   []int{},
		Width:                 0,
		X:                     0,
		Y:                     0,
	})
	c.NodeLabelIdMap[label] = newNodeId
	c.NodeLabels = append(c.NodeLabels, label)
}

func (c *Context) AddConnector(aId, bId int) {
	// create connector node
	cId := len(c.Nodes)
	connectorNode := Node{
		DownwardClosure:       NewIntSet(),
		DownwardNodeIds:       NewIntSet(),
		DownwardNodeIdsSorted: []int{},
		Height:                0,
		IsConnector:           true,
		Layer:                 c.Nodes[aId].Layer + 1,
		Padding:               0,
		Row:                   0,
		UpwardNodeIds:         NewIntSet(),
		UpwardNodeIdsSorted:   []int{},
		Width:                 0,
		X:                     0,
		Y:                     0,
	}
	updatedNodes := append(c.Nodes, connectorNode)
	c.NodeLabels = append(c.NodeLabels, "connector")
	c.Nodes = updatedNodes

	nodeA := c.Nodes[aId]
	nodeB := c.Nodes[bId]
	nodeC := c.Nodes[cId]

	nodeA.DownwardNodeIds.Remove(bId)
	nodeA.DownwardNodeIds.Add(cId)
	c.Nodes[aId] = nodeA

	nodeB.UpwardNodeIds.Add(cId)
	nodeB.UpwardNodeIds.Remove(aId)
	c.Nodes[bId] = nodeB

	nodeC.UpwardNodeIds.Add(aId)
	nodeC.DownwardNodeIds.Add(bId)
	c.Nodes[cId] = nodeC
}

func (c *Context) AddVertex(aLabel, bLabel string) {
	nodeAId := c.NodeLabelIdMap[aLabel]
	nodeA := c.Nodes[nodeAId]

	nodeBId := c.NodeLabelIdMap[bLabel]
	nodeB := c.Nodes[nodeBId]

	nodeA.DownwardNodeIds.Add(nodeBId)
	nodeB.UpwardNodeIds.Add(nodeAId)

	c.Nodes[nodeAId] = nodeA
	c.Nodes[nodeBId] = nodeB
}

func (c *Context) TopoSort() (success bool) {
	i := 0
	isThereMoreWork := true
	for isThereMoreWork {
		i++
		isThereMoreWork = false
		for a := 0; a < len(c.Nodes); a++ {
			for b := range c.Nodes[a].DownwardNodeIds {
				if c.Nodes[b].Layer <= c.Nodes[a].Layer {
					c.Nodes[b].Layer = c.Nodes[a].Layer + 1
					isThereMoreWork = true
				}
			}
		}
		// detect cycles
		if i > len(c.Nodes)*len(c.Nodes) {
			return false
		}
	}
	return true
}

func (c *Context) Complete() {
	isThereMoreWork := true
	for isThereMoreWork {
		isThereMoreWork = false
		for a := 0; a < len(c.Nodes); a++ {
			for b := range c.Nodes[a].DownwardNodeIds {
				if c.Nodes[a].Layer+1 != c.Nodes[b].Layer {
					isThereMoreWork = true
					c.AddConnector(a, b)
					break
				}
			}
		}
	}
}

func (c *Context) AddToLayers() {
	//Compute the number of layers necessary.
	var lastLayer int
	for _, node := range c.Nodes {
		if node.Layer > lastLayer {
			lastLayer = node.Layer
		}
	}
	c.Layers = make([]Layer, lastLayer+1)

	// Put the elements in the layers.
	for i, node := range c.Nodes {
		c.Layers[node.Layer].NodeIds = append(c.Layers[node.Layer].NodeIds, i)
	}

	// optimize row order
	c.OptimizeRowOrder()

	// Precompute upward_sorted, downward_sorted.
	for i, node := range c.Nodes {
		for upwardNodeId := range node.UpwardNodeIds {
			node.UpwardNodeIdsSorted = append(node.UpwardNodeIdsSorted, upwardNodeId)
			c.Nodes[i] = node
		}
		sort.Slice(node.UpwardNodeIdsSorted, func(i, j int) bool {
			return c.Nodes[node.UpwardNodeIdsSorted[i]].Row < c.Nodes[node.UpwardNodeIdsSorted[j]].Row
		})
		for downwardNodeId := range node.DownwardNodeIds {
			node.DownwardNodeIdsSorted = append(node.DownwardNodeIdsSorted, downwardNodeId)
			c.Nodes[i] = node
		}
		sort.Slice(node.DownwardNodeIdsSorted, func(i, j int) bool {
			return c.Nodes[node.DownwardNodeIdsSorted[i]].Row < c.Nodes[node.DownwardNodeIdsSorted[j]].Row
		})
	}

	// Add the edges
	for _, layer := range c.Layers {
		for upwardNodeId := range layer.NodeIds {
			for downwardNodeId := range c.Nodes[upwardNodeId].DownwardNodeIdsSorted {
				layer.Edges = append(layer.Edges, Edge{
					UpwardNodeId:   upwardNodeId,
					DownwardNodeId: downwardNodeId,
					X:              0,
					Y:              0,
				})
			}
		}
	}
}

func (c *Context) OptimizeRowOrder() {
	computeDownwardClosure := func() {
		fmt.Println("HERE")
		fmt.Println(len(c.Layers))
		for y := len(c.Layers) - 2; y > 0; y-- {
			currentLayer := c.Layers[y]
			for _, upwardNodeId := range currentLayer.NodeIds {
				upwardNode := c.Nodes[upwardNodeId]
				for downwardNodeId := range upwardNode.DownwardNodeIds {
					upwardNode.DownwardClosure[downwardNodeId] = true
					for nodeId := range c.Nodes[downwardNodeId].DownwardClosure {
						upwardNode.DownwardClosure[nodeId] = true
					}
				}
				c.Nodes[upwardNodeId] = upwardNode
			}
			c.Layers[y] = currentLayer
		}
	}

	computeDownwardDistances := func(layerId int, layerWidth int) [][]int {
		distanceMatrix := make([][]int, layerWidth)
		for i := range distanceMatrix {
			distanceMatrix[i] = make([]int, layerWidth)
		}
		for a := 0; a < layerWidth; a++ {
			for b := 0; b < layerWidth; b++ {
				nodeA := &c.Nodes[c.Layers[layerId].NodeIds[a]]
				nodeB := &c.Nodes[c.Layers[layerId].NodeIds[b]]
				commonDownwardNodeIds := make([]int, 0)
				for downwardNodeIdForA := range nodeA.DownwardClosure {
					if nodeB.DownwardClosure[downwardNodeIdForA] {
						commonDownwardNodeIds = append(commonDownwardNodeIds, downwardNodeIdForA)
					}
				}
				for _, commonDownwardNodeId := range commonDownwardNodeIds {
					d := distanceMatrix[a][b]
					distanceMatrix[a][b] = min(d, c.Nodes[commonDownwardNodeId].Layer-c.Nodes[a].Layer)
				}
			}
		}
		return distanceMatrix
	}

	computeParentMean := func(layerId int, layerWidth int) []float64 {
		parentMeans := make([]float64, layerWidth)
		for a := 0; a < layerWidth; a++ {
			accumulatedRowSum := 0.0
			downwardNode := c.Nodes[c.Layers[layerId].NodeIds[a]]
			for b := range downwardNode.UpwardNodeIds {
				upwardNode := c.Nodes[b]
				accumulatedRowSum += float64(upwardNode.Row)
			}
			parentMeans[a] = accumulatedRowSum / (float64(len(downwardNode.UpwardNodeIds)) + 0.01)
		}
		return parentMeans
	}

	computePermutationSlice := func(layerWidth int) []int {
		permutation := make([]int, layerWidth)
		for i := 0; i < layerWidth; i++ {
			permutation[i] = i
		}
		return permutation
	}

	evaluateScore := func(layerWidth int, distanceMatrix [][]int, permutation []int, parentMeans []float64) float64 {
		score := 0.0
		for i := 0; i < layerWidth-1; i++ {
			score += float64(distanceMatrix[permutation[i]][permutation[i+1]])
		}
		for i := 0; i < layerWidth; i++ {
			d := float64(i) - parentMeans[permutation[i]]
			score += d * d * 15
		}
		return score
	}

	findLowestScore := func(score float64, layerWidth int, distanceMatrix [][]int, permutation []int, parentMeans []float64) {
		scoreLastLoop := 0.0
		for score != scoreLastLoop {
			scoreLastLoop = score
			for a := 0; a < layerWidth; a++ {
				for b := 0; b < layerWidth; b++ {
					permutation[a], permutation[b] = permutation[b], permutation[a]
					newScore := evaluateScore(layerWidth, distanceMatrix, permutation, parentMeans)
					if newScore < score {
						score = newScore
					} else {
						permutation[a], permutation[b] = permutation[b], permutation[a]
					}
				}
			}
		}
	}

	reorderNodeIds := func(layerId int, layerWidth int, permutation []int) {
		orderedNodeIds := make([]int, layerWidth)
		for i, p := range permutation {
			orderedNodeIds[i] = c.Layers[layerId].NodeIds[p]
		}
		c.Layers[layerId].NodeIds = orderedNodeIds
	}

	computeRows := func(layerId int, layerWidth int) {
		for i := 0; i < layerWidth; i++ {
			node := &c.Nodes[c.Layers[layerId].NodeIds[i]]
			node.Row = i
		}
	}

	computeDownwardClosure()
	for layerId, layer := range c.Layers {
		layerWidth := len(layer.NodeIds)
		distanceMatrix := computeDownwardDistances(layerId, layerWidth)
		parentMeans := computeParentMean(layerId, layerWidth)
		permutation := computePermutationSlice(layerWidth)
		score := evaluateScore(layerWidth, distanceMatrix, permutation, parentMeans)
		findLowestScore(score, layerWidth, distanceMatrix, permutation, parentMeans)
		evaluateScore(layerWidth, distanceMatrix, permutation, parentMeans)
		reorderNodeIds(layerId, layerWidth, permutation)
		computeRows(layerId, layerWidth)
	}
}

func (c *Context) ResolveCrossingEdges() {
	for _, layer := range c.Layers {
		upwardEdges := layer.Edges
		downwardEdges := layer.Edges

		sort.Slice(upwardEdges, func(i, j int) bool {
			iUpwardNodeRow := c.Nodes[upwardEdges[i].UpwardNodeId].Row
			jUpwardNodeRow := c.Nodes[upwardEdges[j].UpwardNodeId].Row
			iDownwardNodeRow := c.Nodes[upwardEdges[i].DownwardNodeId].Row
			jDownwardNodeRow := c.Nodes[upwardEdges[j].DownwardNodeId].Row

			return iUpwardNodeRow < jUpwardNodeRow ||
				(iUpwardNodeRow == jUpwardNodeRow && iDownwardNodeRow < jDownwardNodeRow)
		})

		sort.Slice(downwardEdges, func(i, j int) bool {
			iDownwardNodeRow := c.Nodes[downwardEdges[i].DownwardNodeId].Row
			jDownwardNodeRow := c.Nodes[downwardEdges[j].DownwardNodeId].Row
			iUpwardNodeRow := c.Nodes[downwardEdges[i].UpwardNodeId].Row
			jUpwardNodeRow := c.Nodes[downwardEdges[j].UpwardNodeId].Row

			return iDownwardNodeRow < jDownwardNodeRow ||
				(iDownwardNodeRow == jDownwardNodeRow && iUpwardNodeRow < jUpwardNodeRow)
		})

		for i := 0; i < len(upwardEdges); i++ {
			if !(upwardEdges[i] == downwardEdges[i]) {
				layer.Edges = []Edge{}
				layer.Adapter.IsEnabled = true
			}
		}
	}
}

func (c *Context) Layout() {
	width := 0
	for i, node := range c.Nodes {
		if node.IsConnector {
			width = 1
		} else {
			width = max(0, len(c.NodeLabels[i]))
			width = max(width, len(node.UpwardNodeIds))
			width = max(width, len(node.DownwardNodeIds))
			width += 2
		}
		node.Width = width
		c.Nodes[i] = node
	}

	for i := 0; i < 1000; i++ {
		if !c.isLayoutNodesNotTouching() {
			continue
		}
		if !c.isLayoutNodesNotTouching() {
			continue
		}
		if !c.LayoutGrowNode() {
			continue
		}
		if !c.LayoutShiftEdges() {
			continue
		}
		if !c.LayoutShiftConnectorNode() {
			continue
		}
		break
	}

	for y := 0; y < len(c.Layers)-1; y++ {
		upwardLayer := c.Layers[y]
		downwardLayer := c.Layers[y+1]
		if upwardLayer.Adapter.IsEnabled {
			continue
		}
		width := 0
		for _, nodeId := range upwardLayer.NodeIds {
			width = max(width, c.Nodes[nodeId].X+c.Nodes[nodeId].Width)
		}
		for _, nodeId := range downwardLayer.NodeIds {
			width = max(width, c.Nodes[nodeId].X+c.Nodes[nodeId].Width)
		}

		type Pair struct {
			Origin, Destination int
		}

		ids := make(map[Pair]int)
		getId := func(origin int, destination int) int {
			value := ids[Pair{Origin: origin, Destination: destination}]
			if value == 0 {
				value = len(ids)
			}
			return value
		}
		input := make([]IntSet, width)
		output := make([]IntSet, width)

		for _, nodeIdA := range upwardLayer.NodeIds {
			nodeA := c.Nodes[nodeIdA]
			for x := nodeA.X + nodeA.Padding; x < nodeA.X-nodeA.Padding+nodeA.Width; x++ {
				for downwardNodeId := range nodeA.DownwardNodeIds {
					input[x] = NewIntSet()
					input[x].Add(getId(nodeIdA, downwardNodeId))
				}
			}
		}

		for _, nodeIdB := range downwardLayer.NodeIds {
			nodeB := c.Nodes[nodeIdB]
			for x := nodeB.X + nodeB.Padding; x < nodeB.X-nodeB.Padding+nodeB.Width; x++ {
				for upwardNodeId := range nodeB.UpwardNodeIds {
					output[x] = NewIntSet()
					output[x].Add(getId(upwardNodeId, nodeIdB))
				}
			}
		}

		upwardLayer.Adapter = NewAdapter(input, output)
	}

	// y-axis: size and position.
	y := 0
	for _, layer := range c.Layers {
		for _, nodeId := range layer.NodeIds {
			node := c.Nodes[nodeId]
			node.Y = y
			node.Height = 3
			c.Nodes[nodeId] = node
		}

		for i, edge := range layer.Edges {
			edge.Y = y + 2
			layer.Edges[i] = edge

			if layer.Adapter.IsEnabled {
				layer.Adapter.Y = y + 2
				y += layer.Adapter.Height - 3
			}

			y += 3
		}
	}
}

func (c *Context) isLayoutNodesNotTouching() bool {
	isNotTouching := true
	for _, layer := range c.Layers {
		x := 0
		for _, nodeId := range layer.NodeIds {
			if c.Nodes[nodeId].X < x {
				isNotTouching = false
			}
			x = c.Nodes[nodeId].X + c.Nodes[nodeId].Width
		}
	}
	return isNotTouching
}

// LayoutGrowNode Grow the nodes to fit their edges
func (c *Context) LayoutGrowNode() bool {
	for _, layer := range c.Layers {
		for _, edge := range layer.Edges {
			upwardNode := c.Nodes[edge.UpwardNodeId]
			if upwardNode.X+upwardNode.Width-2 < edge.X && !upwardNode.IsConnector {
				upwardNode.Width = edge.X + 2 - upwardNode.X
				c.Nodes[edge.UpwardNodeId] = upwardNode
				return false
			}

			downwardNode := c.Nodes[edge.DownwardNodeId]
			if downwardNode.X+downwardNode.Width-2 < edge.X && !downwardNode.IsConnector {
				downwardNode.Width = edge.X + 2 - downwardNode.X
				c.Nodes[edge.DownwardNodeId] = downwardNode
				return false
			}
		}
	}

	return true
}

// LayoutShiftEdges Shift the edges to the right, so that they reach their nodes.
func (c *Context) LayoutShiftEdges() bool {
	for _, layer := range c.Layers {
		for _, edge := range layer.Edges {
			upwardNode := c.Nodes[edge.UpwardNodeId]
			downwardNode := c.Nodes[edge.DownwardNodeId]
			minX := max(upwardNode.X+upwardNode.Padding, downwardNode.X+downwardNode.Padding)
			if edge.X < minX {
				edge.X = minX
				return false
			}
		}
	}
	return true
}

func (c *Context) LayoutShiftConnectorNode() bool {
	for nodeId, node := range c.Nodes {
		if !node.IsConnector {
			continue
		}
		minX := 0
		for _, edge := range c.Layers[node.Layer-1].Edges {
			if edge.DownwardNodeId == nodeId {
				minX = max(minX, edge.X)
			}
		}
		for _, edge := range c.Layers[node.Layer].Edges {
			if edge.UpwardNodeId == nodeId {
				minX = max(minX, edge.X)
			}
		}
		if node.X < minX {
			node.X = minX
			c.Nodes[nodeId] = node
			return false
		}
	}
	return true
}

func (c *Context) Render() string {
	width := 0
	height := 0
	for _, node := range c.Nodes {
		width = max(width, node.X+node.Width)
		height = max(height, node.Y+node.Height)
	}

	s := NewScreen(width, height)
	for i, node := range c.Nodes {
		if node.IsConnector {
			if node.Width == 1 {
				s.PlaceVerticalLine(node.Y, node.Y+2, node.X, s.BoxStyle["verticalLine"])
			} else {
				s.PlaceBox(node.X, node.Y, node.Width, node.Height)
			}
		} else {
			fmt.Println(node.X, node.Y, node.Width, node.Height)

			s.PlaceBox(node.X, node.Y, node.Width, node.Height)
			s.PlaceWord(node.X+1, node.Y+1, []rune(c.NodeLabels[i]))
		}
	}

	for y := 0; y < len(c.Layers); y++ {
		layer := c.Layers[y]
		for _, edge := range layer.Edges {
			up := '┬'
			if c.Nodes[edge.UpwardNodeId].IsConnector {
				up = '│'
			}
			down := '▽'
			if c.Nodes[edge.DownwardNodeId].IsConnector {
				up = '│'
			}
			s.PlaceRune(edge.X, edge.Y, up)
			s.PlaceRune(edge.X, edge.Y+1, down)
		}
	}

	for y := 0; y < len(c.Layers); y++ {
		layer := c.Layers[y]
		if layer.Adapter.IsEnabled {
			layer.Adapter.Render(s)
		}
	}

	return s.String()
}

type Node struct {
	DownwardClosure       IntSet
	DownwardNodeIds       IntSet
	DownwardNodeIdsSorted []int
	Height                int
	IsConnector           bool
	Layer                 int
	Padding               int
	Row                   int
	UpwardNodeIds         IntSet
	UpwardNodeIdsSorted   []int
	Width                 int
	X                     int
	Y                     int
}

type Edge struct {
	DownwardNodeId int
	UpwardNodeId   int
	X              int
	Y              int
}

func (a *Edge) isEqualTo(b *Edge) bool {
	return a.UpwardNodeId == b.UpwardNodeId &&
		a.DownwardNodeId == b.DownwardNodeId
}

type Adapter struct {
	Height    int
	Inputs    []IntSet
	IsEnabled bool
	Outputs   []IntSet
	Runes     [][]rune
	Y         int
}

func (a *Adapter) Render(screen *Screen) {
	for dy := 0; dy < a.Height-1; dy++ {
		x := 0
		for _, value := range a.Runes[dy] {
			if value == ' ' {
				x++
				continue
			}

			currentRune := screen.ReadRune(x, a.Y+dy)
			if dy == 0 {
				if *currentRune == '─' {
					*currentRune = '┬'
					x++
					continue
				}
			}

			if dy == a.Height-2 {
				if *currentRune == '─' {
					*currentRune = '▽'
					x++
					continue
				}
			}

			*currentRune = value
			x++
			continue
		}
	}
}

type Layer struct {
	Adapter Adapter
	Edges   []Edge
	NodeIds []int
}

// Edge represents an edge between two nodes.
type AdapterEdge struct {
	A        *AdapterNode
	B        *AdapterNode
	Weight   int
	Assigned int
}

// Node represents a node in the graph.
type AdapterNode struct {
	Visited bool
	Cost    int
	Edges   []*AdapterEdge
}

type NodeAndCost struct {
	Node *AdapterNode
	Cost int
}

type IntSet map[int]bool

func NewIntSet() IntSet {
	return make(map[int]bool)
}

func (s IntSet) Add(item int) {
	s[item] = true
}

func (s IntSet) Contains(item int) bool {
	return s[item]
}

func (s IntSet) Remove(item int) {
	delete(s, item)
}

func (s IntSet) Size() int {
	return len(s)
}

func (s IntSet) Slice() []int {
	var intSlice []int
	for item := range s {
		intSlice = append(intSlice, item)
	}
	return intSlice
}

type NodeSet map[*AdapterNode]bool

func NewNodeSet() NodeSet {
	return make(NodeSet)
}

func (s NodeSet) Add(item *AdapterNode) {
	s[item] = true
}

func (s NodeSet) Contains(item *AdapterNode) bool {
	return s[item]
}

func (s NodeSet) Remove(item *AdapterNode) {
	delete(s, item)
}

func (s NodeSet) Size() int {
	return len(s)
}

// NodeAndCostHeap is a min heap of NodeAndCost.
type NodeAndCostHeap []NodeAndCost

//
//func (h NodeAndCostHeap) Len() int           { return len(h) }
//func (h NodeAndCostHeap) Less(i, j int) bool { return h[i].Cost < h[j].Cost }
//func (h NodeAndCostHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NodeAndCostHeap) Push(x NodeAndCost) {
	*h = append(*h, x)
}

func (h *NodeAndCostHeap) Pop() NodeAndCost {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
