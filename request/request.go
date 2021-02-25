package request

import (
	"errors"
	"fmt"
	"math/rand"

	"example.com/config"
	"example.com/graph"
)

type Request struct {
	Src            *graph.Node
	Dest           *graph.Node
	Paths          [][]*graph.Node
	Priority       int
	Position       int
	GenerationTime int
	ServingTime    int
	HasReached     bool
}

/////////////////////////// This package should be more general. Use an interface to avoid
/////////////////////////// many if-else blocks.

///////////////////////////////////////////// Change this function to receive requests and append to them.
func genRequests(N int, ids [][]int, priority []int, topology string, roundNum int) ([]*Request, error) {
	if topology != graph.GRID && topology != graph.RING {
		return nil, errors.New("request.genRequests: The requested topology is not valid!")
	}
	var reqs []*Request
	if topology == graph.GRID {
		var isSame bool
		var r [2]int
		reqs = make([]*Request, N)
		//s := make([]int, len(ids[0]))
		//d := make([]int, len(ids[0]))
		for i := 0; i < N; i++ {
			isSame = true
			for isSame == true {
				r[0] = rand.Intn(len(ids))
				r[1] = rand.Intn(len(ids))
				isSame = r[0] == r[1]
			}
			reqs[i] = new(Request)
			reqs[i].Src = new(graph.Node)
			reqs[i].Dest = new(graph.Node)
			reqs[i].Src.ID = make([]int, 2)
			reqs[i].Dest.ID = make([]int, 2)
			reqs[i].Src = graph.MakeNode(ids[r[0]], config.GetConfig().GetMemory())
			reqs[i].Dest = graph.MakeNode(ids[r[1]], config.GetConfig().GetMemory())
			reqs[i].Priority = priority[i]
			reqs[i].Position = 0
			reqs[i].GenerationTime = roundNum
			reqs[i].HasReached = false

			////////////// TODO: IMPORTANT!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			reqs[i].Paths = make([][]*graph.Node, 1)

			//s = rand.Intn(len(nodes))
			//isSame = true
			//for isSame == true {
			//	d = rand.Intn(len(nodes))
			//	isSame = (s == d)
			//}
			//reqs[i].Src = s
			//reqs[i].Dest = d
			//reqs[i].Priority = priority[i]
		}
	}
	if topology == graph.RING {
		//reqs = make([]*Request, N)
		fmt.Println("Caution! Input is ring, with no code.")
		return nil, nil
	}
	return reqs, nil
}

// TODO: Handle the priorities
func RG(N int, ids [][]int, priority []int, topology string, roundNum int) ([]*Request, error) {
	return genRequests(N, ids, priority, topology, roundNum)
}

// GatherRemainingRequests() gathers the requests not
func GatherRemainingRequests() {

}
