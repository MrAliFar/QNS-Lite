package path

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/request"
)

const (
	MODIFIED_GREEDY = "modified greedy"
	Q_PASS          = "Q-Pass"
	Q_CAST          = "Q-CAST"
	SLMP            = "SLMP"
)

type Path []*graph.Node

// The PathFinder interface captures the path finding algorithm abstraction.
type PathFinder interface {
	Build(graph.Topology)
	Clear()
	Find(src, dest *graph.Node) Path
	//GetPath() path
	next(dest *graph.Node) *graph.Node
	add(*graph.Node)
	//getNetwork() *graph.Topology
}

func BuildPathFinder(algorithm string, network graph.Topology) PathFinder {
	if algorithm == MODIFIED_GREEDY {
		var mg *modifiedGreedy = new(modifiedGreedy)
		mg.Build(network)
		return mg
	}
	fmt.Println("path.go: Warning! algorithm not recognized!")
	return nil
}

func PF(network graph.Topology, reqs []*request.Request, algorithm string, hasContention bool) {
	pf := BuildPathFinder(algorithm, network)
	//paths := make([]Path, len(reqs))
	for i, req := range reqs {
		//fmt.Println("PF - request number", i)
		for j := 1; j <= config.GetConfig().GetAggressiveness(); j++ {
			//fmt.Println("PF - path number", j)
			//fmt.Println("New Path", PathToNode(pf.Find(req.Src, req.Dest)))
			//fmt.Println("Before", reqs[i].Paths)
			if j == 1 {
				nodes := PathToNode(pf.Find(req.Src, req.Dest))
				//reqs[i].Paths[j] = PathToNode(pf.Find(req.Src, req.Dest))
				//fmt.Println("Found Path is", nodes)
				//fmt.Println("Slot to copy is", reqs[i].Paths[j-1])
				reqs[i].Paths[j-1] = make([]*graph.Node, len(nodes))
				copy(reqs[i].Paths[j-1], nodes)
				//fmt.Println("Copied Path is", reqs[i].Paths[j-1])
			} else {
				reqs[i].Paths = append(reqs[i].Paths, PathToNode(pf.Find(req.Src, req.Dest)))
			}
			//fmt.Println("The new path is", reqs[i].Paths[len(reqs[i].Paths)-1])
			for _, node := range reqs[i].Paths[len(reqs[i].Paths)-1] {
				//fmt.Println(node.ID)
			}
			//fmt.Println(j, "Appended!")
			if !hasContention {
				////////////////////////////// COMPLETE THIS!!!!!!!!!!!!!!!!!!!!!!
				//fmt.Println("En route to PathToLinks", reqs[i].Paths[len(reqs[i].Paths)-1])
				graph.Prune(PathToLinks(reqs[i].Paths[len(reqs[i].Paths)-1], network))
				//fmt.Println(j, "Pruned!")
			}
			//paths[i] = pf.Find(req.Src, req.Dest)
			pf.Clear()
		}
	}
	graph.Deprune(network)
	//return paths
}

//func batchFind()

//func copyNetwork(network graph.Topology) copied graph.Topology

// TODO: Check this!
func PathToNode(path Path) []*graph.Node {
	if path == nil {
		return nil
	}
	//fmt.Println("Inside PathToNodes. Input:", path)
	nodes := make([]*graph.Node, len(path))
	copy(nodes, path)
	//fmt.Println("Inside PathToNodes. Output:", nodes)
	return nodes
}

func PathToLinks(path Path, network graph.Topology) []*graph.Link {
	links := make([]*graph.Link, 0)
	i := 0
	//fmt.Println("Inside PathToLinks. Length of path:", len(path), "Path is:", path)
	//for _, node := range path {
	//	fmt.Println("NODE", node.ID)
	//}
	for i <= len(path)-2 {
		links = append(links, network.GetLinkBetween(path[i], path[i+1]))
		i += 1
	}
	return links
}
