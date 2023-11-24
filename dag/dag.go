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

func NewAdapter(inputs []Set, outputs []Set) Adapter {
	solutionFound := false
	height := 3
	width := len(inputs)
	var adapterNodes []AdapterNode
	var adapterEdges []AdapterEdge
	connectorLength := 0
	for _, input := range inputs {
		inputSlice := input.Slice()
		connectorLength = max(connectorLength, len(inputSlice))
	}
	for !solutionFound {
		adapterNodes = make([]AdapterNode, width*height*2)
		adapterEdges = make([]AdapterEdge, width*height*3)

		IndexFunc := func(x, y, layer int) int {
			return x + width*(y+height*layer)
		}

		ConnectFunc := func(edge AdapterEdge, nodeA AdapterNode, nodeB AdapterNode, weight int) {
			edge.A = &nodeA
			edge.B = &nodeB
			edge.Weight = weight
			nodeA.Edges = append(nodeA.Edges, &edge)
			nodeB.Edges = append(nodeB.Edges, &edge)
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// vertical
				if y != height-1 {
					ConnectFunc(
						adapterEdges[IndexFunc(x, y+0, 0)],
						adapterNodes[IndexFunc(x, y+0, 0)],
						adapterNodes[IndexFunc(x, y+1, 0)],
						1,
					)
				}
			}
		}
	}

	//
	//		// Horizontal:
	//		if (y >= 1 && y<= height - 3 && x != width-1) {
	//		connect(edges[index(x + 0, y, 1)],  //
	//		nodes[index(x + 0, y, 1)],  //
	//		nodes[index(x + 1, y, 1)],  //
	//		1);
	//	}
	//
	//		// Corners:
	//	{
	//		//int dx = width / 2 - x;
	//		int dy = height / 2 - y;
	//		connect(edges[index(x, y, 2)],  //
	//		nodes[index(x, y, 0)],  //
	//		nodes[index(x, y, 1)],  //
	//		10+dy*dy);
	//	}
	//	}
	//	}
	//
	//		// Assume a solution will be found. Otherwise, it will be reset to false.
	//		solution_found = true;
	//
	//		// Add path one by one.
	//		for(int connector = 1; connector <= connector_length; ++connector) {
	//
	//		int big_number = 1 << 15;
	//
	//		// Clear:
	//		for (auto& node : nodes) {
	//		node.visited = false;
	//		node.cost = big_number;
	//	}
	//
	//		std::set<Node*> start;
	//		std::set<Node*> end;
	//		for(int x = 0; x<width; ++x) {
	//		if (inputs[x].count(connector))
	//		start.insert(&nodes[index(x, 0, 0)]);
	//		if (outputs[x].count(connector))
	//		end.insert(&nodes[index(x, height - 1, 0)]);
	//	}
	//
	//	struct NodeAndCost {
	//		Node* node;
	//		int cost;
	//		bool operator<(const NodeAndCost& other) const {
	//		return cost > other.cost;
	//	}
	//	};
	//
	//		std::priority_queue<NodeAndCost> pending;
	//		for (auto& node : start) {
	//		pending.push({node, 0});
	//	}
	//
	//		while (pending.size() != 0) {
	//		auto element = pending.top();
	//		pending.pop();
	//		Node* node = element.node;
	//		if (node->visited)
	//		continue;
	//		node->visited = true;
	//		node->cost = element.cost;
	//		for (Edge* edge : node->edges) {
	//		Node* opposite = edge->a == node ? edge->b : edge->a;
	//		if (opposite->visited)
	//		continue;
	//		if (edge->assigned)
	//		continue;
	//		pending.push({opposite, node->cost + edge->weight});
	//	}
	//	}
	//
	//		// Reconstruct the path from end to start.
	//		int best_score = big_number;
	//		Node* current = nullptr;
	//		for(Node* node : end) {
	//		if (best_score >= node->cost) {
	//		best_score = node->cost;
	//		current = node;
	//	}
	//	}
	//		// No path found.
	//		if (best_score == big_number) {
	//		solution_found = false;
	//		continue;
	//	}
	//
	//		while (!start.count(current)) {
	//		for (Edge* edge : current->edges) {
	//		Node* opposite = edge->a == current ? edge->b : edge->a;
	//		if (current->cost == opposite->cost + edge->weight) {
	//		edge->assigned = connector;
	//		current = opposite;
	//	}
	//	}
	//	}
	//
	//		for (int y = 0; y < height; ++y) {
	//		for (int x = 0; x < width; ++x) {
	//		if (edges[index(x, y, 0)].assigned)
	//		edges[index(x, y, 1)].weight = 20;
	//		if (edges[index(x, y, 1)].assigned)
	//		edges[index(x, y, 0)].weight = 20;
	//	}
	//	}
	//	}
	//
	//		if (height > 30)
	//		solution_found = true;
	//
	//		auto assigned = [&](int x, int y, int layer) -> bool {
	//		return edges[index(x, y, layer)].assigned;
	//	};
	//
	//		if (!solution_found) {
	//		height++;
	//		continue;
	//	}
	//
	//		rendering = std::vector<std::vector<wchar_t>>(
	//		height, std::vector<wchar_t>(width, L' '));
	//		for (int y = 0; y < height; ++y) {
	//		for (int x = 0; x < width; ++x) {
	//		wchar_t& v = rendering[y][x];
	//		if (assigned(x, y, 1))
	//		v = L'─';
	//		if (assigned(x, y, 0))
	//		v = L'│';
	//		if (assigned(x, y, 2)) {
	//		if (assigned(x, y, 0))
	//		v = assigned(x, y, 1) ? L'┌' : L'┐';
	//		else
	//		v = assigned(x, y, 1) ? L'└' : L'┘';
	//	}
	//	}
	//	}
	//		return;
	//	}

	// Implement the Construct function here
	adapter := Adapter{
		Inputs:  inputs,
		Outputs: outputs,
	}
	return adapter
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
	// x-axis: minimal size to draw their content.
	for i, node := range c.Nodes {
		if node.IsConnector {
			node.Width = 1
		} else {
			node.Width = max(0, len(c.NodeLabels[i]))
			node.Width = max(node.Width, len(node.UpwardNodeIds))
			node.Width = max(node.Width, len(node.DownwardNodeIds))
			node.Width += 2
		}
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
		input := make([]Set, width)
		output := make([]Set, width)

		for _, nodeIdA := range upwardLayer.NodeIds {
			nodeA := c.Nodes[nodeIdA]
			for x := nodeA.X + nodeA.Padding; x < nodeA.X-nodeA.Padding+nodeA.Width; x++ {
				for downwardNodeId := range nodeA.DownwardNodeIds {
					input[x].Add(getId(nodeIdA, downwardNodeId))
				}
			}
		}

		for _, nodeIdB := range downwardLayer.NodeIds {
			nodeB := c.Nodes[nodeIdB]
			for x := nodeB.X + nodeB.Padding; x < nodeB.X-nodeB.Padding+nodeB.Width; x++ {
				for upwardNodeId := range nodeB.UpwardNodeIds {
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
		}

		for _, edge := range layer.Edges {
			edge.Y = y + 2

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
				return false
			}

			downwardNode := c.Nodes[edge.DownwardNodeId]
			if downwardNode.X+downwardNode.Width-2 < edge.X && !downwardNode.IsConnector {
				downwardNode.Width = edge.X + 2 - downwardNode.X
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
			return false
		}
	}
	return true
}

func (c *Context) Render() string {
	return ""
}
