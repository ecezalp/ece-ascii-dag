package dag

import (
	. "ece-ascii-dag/screen"
	. "ece-ascii-dag/util"
	"sort"
	"strings"
)

type Node struct {
	DownwardClosure       Set
	DownwardNodeIds       Set
	DownwardNodeIdsSorted []int
	Height                int
	IsConnector           bool
	Layer                 int
	Padding               int
	Row                   int
	UpwardNodeIds         Set
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
	Inputs    []Set
	IsEnabled bool
	Outputs   []Set
	Runes     [][]rune
	Y         int
}

func NewAdapter() *Adapter {
	// Implement the Construct function here
	adapter := Adapter{}
	return &adapter
}

func (a *Adapter) Render(screen *Screen) {
	// Implement the Render function here
}

type Layer struct {
	Adapter Adapter
	Edges   []Edge
	NodeIds []int
}

type Context struct {
	Layers         []Layer
	NodeLabelIdMap map[string]int
	NodeLabels     []string // node names
	Nodes          []Node
}

func (c *Context) Process(input string) string {
	return ""
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
		DownwardClosure:       nil,
		DownwardNodeIds:       nil,
		DownwardNodeIdsSorted: nil,
		Height:                0,
		IsConnector:           false,
		Layer:                 0,
		Padding:               1,
		Row:                   0,
		UpwardNodeIds:         nil,
		UpwardNodeIdsSorted:   nil,
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
	c.Nodes = append(c.Nodes, Node{
		DownwardClosure:       nil,
		DownwardNodeIds:       nil,
		DownwardNodeIdsSorted: nil,
		Height:                0,
		IsConnector:           true,
		Layer:                 c.Nodes[aId].Layer + 1,
		Padding:               0,
		Row:                   0,
		UpwardNodeIds:         nil,
		UpwardNodeIdsSorted:   nil,
		Width:                 0,
		X:                     0,
		Y:                     0,
	})
	c.NodeLabels = append(c.NodeLabels, "connector")

	// insert connector node between nodeA and nodeB
	c.Nodes[aId].DownwardNodeIds.Remove(bId)
	c.Nodes[bId].UpwardNodeIds.Remove(aId)
	c.Nodes[aId].DownwardNodeIds.Add(cId)
	c.Nodes[cId].UpwardNodeIds.Add(aId)
	c.Nodes[cId].DownwardNodeIds.Add(bId)
	c.Nodes[bId].UpwardNodeIds.Add(cId)
}

func (c *Context) AddVertex(aLabel, bLabel string) {
	// place nodeB below nodeA
	c.Nodes[c.NodeLabelIdMap[aLabel]].DownwardNodeIds.Add(c.NodeLabelIdMap[bLabel])
	c.Nodes[c.NodeLabelIdMap[bLabel]].UpwardNodeIds.Add(c.NodeLabelIdMap[aLabel])
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
	for _, node := range c.Nodes {
		for upwardNodeId := range node.UpwardNodeIds {
			node.UpwardNodeIdsSorted = append(node.UpwardNodeIdsSorted, upwardNodeId)
		}
		sort.Slice(node.UpwardNodeIds, func(i, j int) bool {
			return c.Nodes[node.UpwardNodeIdsSorted[i]].Row < c.Nodes[node.UpwardNodeIdsSorted[j]].Row
		})
		for downwardNodeId := range node.DownwardNodeIds {
			node.DownwardNodeIdsSorted = append(node.DownwardNodeIdsSorted, downwardNodeId)
		}
		sort.Slice(node.DownwardNodeIds, func(i, j int) bool {
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
		for y := len(c.Layers) - 2; y > 0; y-- {
			currentLayer := &c.Layers[y]
			for _, upwardNodeId := range currentLayer.NodeIds {
				upwardNode := &c.Nodes[upwardNodeId]
				for downwardNodeId := range upwardNode.DownwardNodeIds {
					upwardNode.DownwardClosure[downwardNodeId] = true
					for nodeId := range c.Nodes[downwardNodeId].DownwardClosure {
						upwardNode.DownwardClosure[nodeId] = true
					}
				}
			}
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
	return
}

func (c *Context) Layout() {
	return
}

func (c *Context) LayoutNodeDoNotTouch() bool {
	return false
}

func (c *Context) LayoutEdgesDoNotTouch() bool {
	return false
}

func (c *Context) LayoutGrowNode() bool {
	return false
}

func (c *Context) LayoutShiftEdges() bool {
	return false
}

func (c *Context) LayoutShiftConnectorNode() bool {
	return false
}

func (c *Context) Render() string {
	return ""
}
