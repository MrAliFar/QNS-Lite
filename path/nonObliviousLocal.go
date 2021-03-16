package path

import (
	"fmt"

	"example.com/graph"
)

type nonObliviousLocal struct {
	network    graph.Topology
	path       Path
	isFinished bool
	curr       *graph.Node
	//src        *graph.Node
	//dest       *graph.Node
}

func (nol *nonObliviousLocal) Build(network graph.Topology) {
	nol.path = make([]*graph.Node, 1)
	nol.path[0] = new(graph.Node)
	nol.curr = new(graph.Node)
	nol.isFinished = false
	nol.network = network
	//mg.path[0] = mg.src
	//mg.curr = mg.src
}

func (nol *nonObliviousLocal) Clear() {
	nol.path = nil
	nol.isFinished = false
	nol.path = make([]*graph.Node, 1)
	nol.path[0] = new(graph.Node)
	//nol.curr = nil
}

func (nol *nonObliviousLocal) Find(src, dest *graph.Node) Path {
	nol.curr = src
	nol.add(src)
	cntr := 0
	for !nol.curr.IsEqual(dest) {
		//fmt.Println("Inside Find - Network Size is:", mg.network.GetSize())
		//fmt.Println("Inside Find - Counter Threshold is:", mg.network.GetSize()^2)
		if cntr >= nol.network.GetSize()*nol.network.GetSize() {
			fmt.Println("Inside Find. Counter overflow.")
			return nil
		}
		cntr = cntr + 1
		next := nol.next(dest)
		//fmt.Println("Next found", next, "CNTR", cntr)
		if next.Memory == 0 {
			//fmt.Println(mg.path)
			return nil
		}
		nol.add(nol.next(dest))
		nol.curr = nol.path[len(nol.path)-1]
	}
	//fmt.Println("Found Path - inside find: ", mg.path)
	return nol.path
}

func (nol *nonObliviousLocal) next(dest *graph.Node) *graph.Node {
	neighbors, neighIsNil := nol.network.GetNeighbors(nol.curr)
	//fmt.Println("Inside next - The neighbors are:", neighbors)
	if neighIsNil {
		return nil
	}
	optimumNode := neighbors[0]
	choices := make([]*graph.Node, 1)
	choices[0] = optimumNode
	return nil
}

func (nol *nonObliviousLocal) add(n *graph.Node) {
	//fmt.Println(nol.path)
	if nol.path[0].Memory == 0 {
		nol.path[0] = n
		//copy(nol.path[0], n)
		//fmt.Println("HERE!!!")
	} else {
		nol.path = append(nol.path, n)
	}
	//fmt.Println("Inside add. Input:", n)
	//fmt.Println("Inside add. Output:", nol.path)
}
