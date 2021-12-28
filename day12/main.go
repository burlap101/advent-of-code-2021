package main

import (
	"fmt"
	"os"
	"strings"
	"sort"
)

type State struct {
	VertexSet map[string]*Vertex
	Paths map[string]Path
}

type Vertex struct {
	name string
	isBig bool
	nextHops map[Edge]*Vertex
}

type Edge struct {
	start string
	end string
}

type Path struct {
	steps []*Vertex
} 

func (v *Vertex) AddHop(e Edge, state *State) (int) {
	if _, ok := v.nextHops[e]; !ok {
		eV, ok2 := state.VertexSet[e.end]
		if !ok2 {
			eV = &Vertex{e.end, strings.ToUpper(e.end) == e.end, make(map[Edge]*Vertex)}
			state.AddVertexToSet(eV)
		}
		v.nextHops[e] = eV
	}
	return len(v.nextHops)
}

func (s *State) AddVertexToSet(v *Vertex) (int) {
	if _, ok := s.VertexSet[v.name]; !ok {
		s.VertexSet[v.name] = v
	}
	return len(s.VertexSet)
}

func processInputFile(filename string) (State, error) {
	state := State{make(map[string]*Vertex), make(map[string]Path)}
	dat, err := os.ReadFile(filename)
	
	if err != nil {
		return state, err
	}
	edges := make([]Edge, 0)
	for _, line := range strings.Split(string(dat), "\n") {
		caves := strings.Split(line, "-")
		edges = append(edges, Edge{caves[0], caves[1]})
	}
	
	for _, edge := range edges {
		vStart, ok := state.VertexSet[edge.start]
		if !ok {
			state.VertexSet[edge.start] = &Vertex{edge.start, strings.ToUpper(edge.start) == edge.start, make(map[Edge]*Vertex)}
			vStart = state.VertexSet[edge.start]
		}
		vEnd, ok2 := state.VertexSet[edge.end]; 
		
		if !ok2 {
			state.AddVertexToSet(&Vertex{edge.end, strings.ToUpper(edge.end)==edge.end, make(map[Edge]*Vertex)})
			vEnd = state.VertexSet[edge.end]
		}
		vStart.nextHops[edge] = vEnd
		if !(edge.start == "start" || edge.end == "end") {
			vEnd.AddHop(Edge{edge.end, edge.start}, &state)
		}
	}
	for _, v := range state.VertexSet {
		for edge := range v.nextHops {
			if edge.end == "start" {
				delete(v.nextHops, edge)
			}
		}
	}
	return state, nil
}

func (s State) String() string {
	result := ""
	for name, v := range s.VertexSet {
		result += fmt.Sprintf("Vertex %s {\n", name)
		result += fmt.Sprintf("  isBig: %v\n", v.isBig)
		result += "  nextHops: {"
		if len(v.nextHops) == 0 {
			result += "}\n"
		}
		for edge, nextHopV := range v.nextHops {
			result += fmt.Sprintf("\n    %+v: %s", edge, nextHopV.name)
		}
		if len(v.nextHops) > 0 {
			result += "\n  }\n"
		}
		result += "}\n"
	}
	result += "Paths: ["
	if len(s.Paths) == 0 {
		result += "]\n"
	}
	for pathStr := range s.Paths {
		result += "\n  " + pathStr
	}
	if len(s.Paths) > 0 {
		result += "\n]\n"
	}
	return result
}

func (s *State) FindPaths() (int, error) {
	startVertex := s.VertexSet["start"]
	for _, v := range startVertex.nextHops {
		path := Path{make([]*Vertex, 2)}
		path.steps[0] = startVertex
		path.steps[1] = v
		pathStr := "start," + v.name + ","
		s.Paths[pathStr] = path
		err := s.TraverseGraph(v, pathStr)
		if err != nil {
			return 0, err
		}
	}
	return len(s.Paths), nil 
}

func (s *State) TraverseGraph(currentVertex *Vertex, currentPathString string) (error) {
	currentPath := s.Paths[currentPathString]
	currentPath.SortSteps()
	for _, v := range currentVertex.nextHops {
		isHopOnPathAlready := strings.Contains(currentPathString, "," + v.name + ",")
		if v.isBig || (isHopOnPathAlready && !currentPath.LittleCaveVisitedTwice()) || !isHopOnPathAlready {
			newPathStr := currentPathString + v.name
			if v.name != "end" {
				newPathStr += ","
			}
			newPath := Path{make([]*Vertex, len(currentPath.steps))}
			copy(newPath.steps, currentPath.steps)
			newPath.steps = append(newPath.steps, v)
			s.Paths[newPathStr] = newPath
			if v.name != "end" {
				err := s.TraverseGraph(v, newPathStr)
				if err != nil {
					return err
				}
			}
		}
	}
	delete(s.Paths, currentPathString)
	return nil
}

func (p Path) LittleCaveVisitedTwice() bool {
	namesVisits := make(map[string]int)
	for _, v := range p.steps {
		if _, ok := namesVisits[v.name]; !ok {
			namesVisits[v.name] = 0
		}
		namesVisits[v.name] += 1
		if !v.isBig && namesVisits[v.name] == 2 {
			return true
		}
	}
	return false
}

func (p *Path) SortSteps() {
	sort.Slice(p.steps, func(i, j int) bool { return p.steps[i].name < p.steps[j].name})
}



func main() {
	fmt.Println("Advent of code day 12")
	state, err := processInputFile("input.txt")
	if err != nil {
		panic(err)
	}
	
	pathCount, err := state.FindPaths()
	if err != nil {
		panic(err)
	}
	fmt.Println(state)
	fmt.Println("Part 2 Answer:", pathCount)
}