package path

import (
	"fmt"

	"example.com/graph"
)

type modifiedGreedy struct {
	network    graph.Topology
	path       Path
	isFinished bool
	curr       *graph.Node
	//src        *graph.Node
	//dest       *graph.Node
}

func (mg *modifiedGreedy) Build(network graph.Topology) {
	mg.path = make([]*graph.Node, 1)
	mg.path[0] = new(graph.Node)
	mg.curr = new(graph.Node)
	mg.isFinished = false
	mg.network = network
	//mg.path[0] = mg.src
	//mg.curr = mg.src
}

// Clear flushes the path finder after it has found a path for a request.
func (mg *modifiedGreedy) Clear() {
	mg.path = nil
	mg.isFinished = false
	mg.path = make([]*graph.Node, 1)
	mg.path[0] = new(graph.Node)
	//mg.curr = nil
}

func (mg *modifiedGreedy) Find(src, dest *graph.Node) Path {
	mg.curr = src
	mg.add(src)
	cntr := 0
	for !mg.curr.IsEqual(dest) {
		//fmt.Println("Inside Find - Network Size is:", mg.network.GetSize())
		//fmt.Println("Inside Find - Counter Threshold is:", mg.network.GetSize()^2)
		if cntr >= mg.network.GetSize()*mg.network.GetSize() {
			fmt.Println("Inside Find. Counter overflow.")
			return nil
		}
		cntr = cntr + 1
		next := mg.next(dest)
		//fmt.Println("Next found", next, "CNTR", cntr)
		if next.Memory == 0 {
			//fmt.Println(mg.path)
			return nil
		}
		mg.add(mg.next(dest))
		mg.curr = mg.path[len(mg.path)-1]
	}
	//fmt.Println("Found Path - inside find: ", mg.path)
	return mg.path
}

func (mg *modifiedGreedy) next(dest *graph.Node) *graph.Node {
	neighbors := mg.network.GetNeighbors(mg.curr)
	//fmt.Println("Inside next - The neighbors are:", neighbors)
	if neighbors == nil {
		return nil
	}
	optimumNode := neighbors[0]
	choices := make([]*graph.Node, 1)
	choices[0] = optimumNode
	for _, node := range neighbors {
		if mg.network.Distance(node, dest, "hop") == mg.network.Distance(optimumNode, dest, "hop") {
			choices = append(choices, node)
		}
		if mg.network.Distance(node, dest, "hop") < mg.network.Distance(optimumNode, dest, "hop") {
			optimumNode = node
			choices = make([]*graph.Node, 1)
			choices[0] = optimumNode
		}
	}
	if len(choices) == 0 {
		return optimumNode
	} else {
		for _, node := range choices {
			if mg.network.GetLinkBetween(mg.curr, node).IsActivated() == true {
				return node
			}
		}
		return choices[0]
	}
}

func (mg *modifiedGreedy) add(n *graph.Node) {
	//fmt.Println(mg.path)
	if mg.path[0].Memory == 0 {
		mg.path[0] = n
		//copy(mg.path[0], n)
		//fmt.Println("HERE!!!")
	} else {
		mg.path = append(mg.path, n)
	}
	//fmt.Println("Inside add. Input:", n)
	//fmt.Println("Inside add. Output:", mg.path)
}
