package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"

	"example.com/request"

	"example.com/path"

	"example.com/quantum"
)

type modifiedGreedyProfile struct {
	network     graph.Grid
	isFinished  bool
	hasRecovery bool
	runTime     int
}

func (mgp modifiedGreedyProfile) Build(topology string) {
	mgp.runTime = 1
	mgp.hasRecovery = config.GetConfig().GetHasRecovery()
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (mgp modifiedGreedyProfile) Run(numRequests int) {
	var priority []int
	priority = make([]int, numRequests)
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		priority[i] = 1
	}
	ids := mgp.network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, mgp.network.GetType(), mgp.runTime)
	if err == nil {
		fmt.Println("Profile Run: Error in request generation!")
		return
	}
	//for i := 0; i < num; i++ {
	//	fmt.Println(i, reqs[i].Src)
	//	fmt.Println(reqs[i].Dest)
	//	for _, node := range reqs[i].Paths[0] {
	//		fmt.Println("PATHS FOR THIS REQUEST", node.ID)
	//	}
	//}
	//fmt.Println(reqs)
	//fmt.Println("Request generation error:", err)

	// TODO: hasRecovery can be forced here...........................................................
	if mgp.hasRecovery {
		path.PF(mgp.network, reqs, "modified greedy", false)
	} else {
		config := config.GetConfig()
		configPointer := &config
		configPointer.SetAggressiveness(1)
		path.PF(mgp.network, reqs, "modified greedy", true)
	}

	//check := 2
	//fmt.Println(reqs[check].Src, reqs[check].Dest)
	//for _, node := range paths[check] {
	//	fmt.Println("PATH:", node.ID)
	//}
	//for _, link := range links {
	//	fmt.Println("Before", link.IsActive)
	//}
	links := mgp.network.GetLinks()
	numReached := 0
	for !mgp.isFinished {
		isReady := true
		mgp.runTime++
		quantum.EG(links)
		if !mgp.hasRecovery {
			for _, req := range reqs {
				if req.HasReached == true {
					continue
				}
				for i := 1; i <= len(req.Paths[0])-1; i++ {
					isReady = isReady && mgp.network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsActive
				}
				if isReady == true {
					numReached += quantum.ES(req, mgp.network, mgp.runTime)
				}
				isReady = true
			}
		}
		if numReached == len(reqs) {
			mgp.isFinished = true
		}
	}
}

func (mgp modifiedGreedyProfile) Stop() {
	mgp.isFinished = true
}
