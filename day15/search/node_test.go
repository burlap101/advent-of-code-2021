package search

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestChildren(t *testing.T) {
	
	n := Node{
		Location: MatrixCoords{0, 0},
		ParentNode: nil,
		PathCost: 0,
	}
	riskmap := ExampleRiskMap()
	children := n.Children(riskmap)
	t.Run("Length", func(t *testing.T) {
		if len(children) != 2 {
			t.Errorf("children found = %d; expected 2", len(children))
		}
	})

	for i, child := range children {
		
		t.Run(fmt.Sprintf("Child %d correct", i+1), func(t *testing.T) {
			if child.PathCost != 1 {
				t.Errorf("child.Pathcost = %d; expected 1", child.PathCost)
			}
			if *child.ParentNode != n {
				t.Errorf("child.ParentNode = %+v; expected %+v", *child.ParentNode, n)
			}
			expectedCoords := []MatrixCoords{{1, 0}, {0, 1}}
			if child.Location != expectedCoords[0] && child.Location != expectedCoords[1] {
				t.Errorf("child.Location = %+v; expected either %+v or %+v", child.Location, expectedCoords[0], expectedCoords[1])
			}
		})
	}

}

func ExampleRiskMap() RiskMap {
	result := make(RiskMap, 0)
	dat := `1163751742
					1381373672
					2136511328
					3694931569
					7463417111
					1319128137
					1359912421
					3125421639
					1293138521
					2311944581`
	
	for m, row := range strings.Split(dat, "\n") {
		result = append(result, make([]uint8, 0))
		row = strings.TrimSpace(row)
		for _, val := range row {
			num, err := strconv.Atoi(string(val))
			if err != nil {
				panic(err)
			}
			result[m] = append(result[m], uint8(num))
		}
	}
	return result
}