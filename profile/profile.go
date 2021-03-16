package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/request"
)

const (
	MODIFIED_GREEDY = "modified greedy"
)

type Profile interface {
	Build(topology string)
	Run(reqs []*request.Request)
	Stop()
	Clear()
	GetRunTime() int
	GetHasRecovery() bool
	GetNetwork() graph.Topology
}

func GenRequests(numRequests int, network graph.Topology, hasRecovery bool) []*request.Request {
	var priority []int
	priority = make([]int, numRequests)
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		priority[i] = 1
	}
	ids := network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, network.GetType(), 1)
	if err != nil {
		fmt.Println("Profile genRequests: Error in request generation!", err)
		return nil
	}
	if hasRecovery {
		path.PF(network, reqs, "modified greedy")
	} else {
		config := config.GetConfig()
		configPointer := &config
		configPointer.SetAggressiveness(1)
		path.PF(network, reqs, "modified greedy")
	}

	for _, req := range reqs {
		n1 := req.Src
		n2 := req.Dest
		fmt.Println(*n1)
		fmt.Println(*n2)
		fmt.Println(len(req.Paths[0]))
		lenn := len(req.Paths)
		for i := 0; i <= lenn-1; i++ {
			for _, nodede := range req.Paths[i] {
				fmt.Println("The next node for path", i+1)
				fmt.Println(*nodede)
			}
		}
	}
	return reqs
}

// Each profile will have a unique profile id.
//func BuildProfile(profileID int) profile {

//}
