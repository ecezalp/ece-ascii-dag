package dag

import (
	. "ece-ascii-dag/screen"
	. "ece-ascii-dag/util"
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
	return
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
	return
}

func (c *Context) AddToLayers() {
	return
}

func (c *Context) OptimizeRowOrder() {
	return
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
