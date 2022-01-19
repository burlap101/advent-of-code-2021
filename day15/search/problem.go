package search

import (
	"fmt"
	"sort"
)

type RiskMap [][]uint8

type Problem struct {
	frontier []NodeRank
	StateSpace RiskMap
	explored map[MatrixCoords]struct{}
}

type NodeRank struct {
	N Node
	Rank float32
}

func (baseRiskMap RiskMap) PathCoster(parentCost uint32, location MatrixCoords) uint32 {
	section := MatrixCoords{location.m / uint32(len(baseRiskMap)), location.n / uint32(len(baseRiskMap[0]))} // section refers to the 2D slice of the overall riskmap
	baseRisk := baseRiskMap[int(location.m) % len(baseRiskMap)][int(location.n) % len(baseRiskMap)]
	resultingRisk := (baseRisk + uint8(section.n) + uint8(section.m))
	if resultingRisk >= 10 {
		resultingRisk = resultingRisk % 10 + 1
	}
	return parentCost + uint32(resultingRisk)
}

func (p *Problem) GoalTest(n Node) bool {
	goalPosition := MatrixCoords{uint32(len(p.StateSpace)-1), uint32(len(p.StateSpace[0])-1)}
	return n.Location == goalPosition 
}

func (p * Problem) GoalTestPart2(n Node) bool {
	goalPosition := MatrixCoords{uint32(len(p.StateSpace)*5-1), uint32(len(p.StateSpace[0])*5-1)}
	return n.Location == goalPosition 
}

func StraightLineDistance(n Node, goalPosition MatrixCoords) float64 {
	// return math.Sqrt(math.Pow(float64(goalPosition.m-n.Location.m), 2) + math.Pow(float64(goalPosition.n-n.Location.n), 2))
	return float64(goalPosition.m-n.Location.m)+float64(goalPosition.n-n.Location.n)
}

func (p *Problem) Ranker(n Node) float32 {
	goalLocation := MatrixCoords{uint32(len(p.StateSpace)-1), uint32(len(p.StateSpace[0])-1)}
	return float32(StraightLineDistance(n, goalLocation)) + float32(n.PathCost)
}

func (p *Problem) RankerPart2(n Node) float32 {
	goalLocation := MatrixCoords{uint32(len(p.StateSpace)*5-1), uint32(len(p.StateSpace[0])*5-1)}
	return float32(StraightLineDistance(n, goalLocation)) + float32(n.PathCost)
}

func (p *Problem) Search() (uint32, error) {
	startNode := Node{
		Location: MatrixCoords{0,0}, 
		ParentNode: nil, 
		PathCost: 0,
	}
	p.explored = make(map[MatrixCoords]struct{})
	p.frontier = make([]NodeRank, 1)
	p.frontier[0] = NodeRank{N: startNode, Rank: p.Ranker(startNode)}
	for len(p.frontier) > 0 {
		fmt.Printf("Frontier size: %d; Explored size: %d                       \r", len(p.frontier), len(p.explored))
		node := p.frontier[0].N
		if p.GoalTest(node) {
			return node.PathCost, nil
		}
		p.explored[node.Location] = struct{}{}
		p.frontier = p.frontier[1:]
		for _, cnode := range node.Children(p.StateSpace) {
			if _, ok := p.explored[cnode.Location]; !ok {
				if i := sort.Search(len(p.frontier), func(i int) bool {return p.frontier[i].N.Location == cnode.Location}); i == len(p.frontier) {
					p.frontier = append(p.frontier, NodeRank{N: cnode, Rank: p.Ranker(cnode)})
					sort.Slice(p.frontier, func(i, j int) bool { return p.frontier[i].Rank < p.frontier[j].Rank })
				}
			}
		}
	}
	return 0, fmt.Errorf("no paths were returned")
}

func (p *Problem) SearchPart2() (uint32, MatrixCoords, error) {
	startNode := Node{
		Location: MatrixCoords{0,0}, 
		ParentNode: nil, 
		PathCost: 0,
	}
	p.explored = make(map[MatrixCoords]struct{})
	p.frontier = make([]NodeRank, 1)
	p.frontier[0] = NodeRank{N: startNode, Rank: p.Ranker(startNode)}
	for len(p.frontier) > 0 {
		fmt.Printf("Frontier size: %d; Explored size: %d                       \r", len(p.frontier), len(p.explored))
		node := p.frontier[0].N
		if p.GoalTestPart2(node) {
			return node.PathCost, node.Location, nil
		}
		p.explored[node.Location] = struct{}{}
		p.frontier = p.frontier[1:]
		for _, cnode := range node.Children(p.StateSpace) {
			if _, ok := p.explored[cnode.Location]; !ok {
				if i := sort.Search(len(p.frontier), func(i int) bool {return p.frontier[i].N.Location == cnode.Location}); i == len(p.frontier) {
					p.frontier = append(p.frontier, NodeRank{N: cnode, Rank: p.RankerPart2(cnode)})
					sort.Slice(p.frontier, func(i, j int) bool { return p.frontier[i].Rank < p.frontier[j].Rank })
				}
			}
		}
	}
	return 0, MatrixCoords{0,0}, fmt.Errorf("no paths were returned")
}

