package search

type Node struct {
	Location MatrixCoords
	ParentNode *Node
	PathCost uint32
}

type MatrixCoords struct {
	m, n uint32
}

func (n Node) Children(riskMap RiskMap) []Node {
	nodeLocations := make([]MatrixCoords, 0)
	if n.Location.m != 0 && int(n.Location.m) != len(riskMap)*5-1 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m+1, n.Location.n})
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m-1, n.Location.n})
	} else if n.Location.m == 0 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m+1, n.Location.n})
	} else if int(n.Location.m) == len(riskMap)*5-1 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m-1, n.Location.n})
	}
	
	if n.Location.n != 0 && int(n.Location.n) != len(riskMap[0])*5-1 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m, n.Location.n+1})
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m, n.Location.n-1})
	} else if n.Location.n == 0 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m, n.Location.n+1})
	} else if int(n.Location.n) == len(riskMap[0])*5-1 {
		nodeLocations = append(nodeLocations, MatrixCoords{n.Location.m, n.Location.n-1})
	}

	result := make([]Node, len(nodeLocations))
	for i, nodeLocation := range nodeLocations {
		result[i] = Node{
			Location: nodeLocation,
			ParentNode: &n,
			PathCost: riskMap.PathCoster(n.PathCost, nodeLocation),
		}
	}
	return result
}


