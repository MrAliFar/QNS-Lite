package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"

	"example.com/request"

	"example.com/quantum"
)

type ModifiedGreedyProfile struct {
	Network     *graph.Grid
	isFinished  bool
	hasRecovery bool
	RunTime     int
}

func (mgp *ModifiedGreedyProfile) Build(topology string) {
	mgp.RunTime = 1
	mgp.hasRecovery = config.GetConfig().GetHasRecovery()
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		mgp.Network = grid
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (mgp *ModifiedGreedyProfile) GetNetwork() graph.Topology {
	return mgp.Network
}

////////////////// TODO: The requests should remain the same through iterations.
////////////////// Take them out.
func (mgp *ModifiedGreedyProfile) Run(reqs []*request.Request) {
	//var priority []int
	//priority = make([]int, numRequests)
	// Priority for the requests
	//for i := 0; i < numRequests; i++ {
	//	priority[i] = 1
	//}
	//ids := mgp.network.GetNodeIDs()
	//reqs, err := request.RG(numRequests, ids, priority, mgp.network.GetType(), mgp.RunTime)
	//if err != nil {
	//	fmt.Println("Profile Run: Error in request generation!", err)
	//	return
	//}

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

	//if mgp.hasRecovery {
	//	path.PF(mgp.network, reqs, "modified greedy")
	//} else {
	//	config := config.GetConfig()
	//	configPointer := &config
	//	configPointer.SetAggressiveness(1)
	//	path.PF(mgp.network, reqs, "modified greedy")
	//}

	//for _, req := range reqs {
	//	n1 := req.Src
	//	n2 := req.Dest
	//	fmt.Println(*n1)
	//	fmt.Println(*n2)
	//	fmt.Println(len(req.Paths[0]))
	//	lenn := len(req.Paths)
	//	for i := 0; i <= lenn-1; i++ {
	//		for _, nodede := range req.Paths[i] {
	//			fmt.Println("The next node for path", i+1)
	//			fmt.Println(*nodede)
	//		}
	//	}
	//}

	//check := 2
	//fmt.Println(reqs[check].Src, reqs[check].Dest)
	//for _, node := range paths[check] {
	//	fmt.Println("PATH:", node.ID)
	//}
	//for _, link := range links {
	//	fmt.Println("Before", link.IsActive)
	//}
	links := mgp.Network.GetLinks()
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	maxItr := 5000
	var cntr int
	if !isOpportunistic {
		for !mgp.isFinished {
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			isReady := true
			mgp.RunTime++
			// EG() also handles lifetimes.
			quantum.EG(links)
			if !mgp.hasRecovery {
				for reqNum, req := range reqs {
					if req.HasReached == true {
						// Release the reserved links
						// Here, req.CanMove is used to release the links, and is set to false
						// to prevent extra work every time the request enters this if statement.
						if req.CanMove {
							for i := 1; i <= len(req.Paths[0])-1; i++ {
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsReserved = false
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).Reservation = -1
							}
						}
						req.CanMove = false
						continue
					}
					cntr = 0
					for i := 1; i <= len(req.Paths[0])-1; i++ {
						link := mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1])
						if link.IsReserved == false {
							isReady = isReady && link.IsActive
							cntr++
						} else {
							if link.Reservation == reqNum {
								isReady = isReady && link.IsActive
								cntr++
							}
						}
						if cntr == 0 {
							isReady = false
						}
						if !isReady {
							break
						}
						// Solve the isReady issue.
					}
					//fmt.Println("Request", reqNum, isReady)
					if isReady == true {
						// req.CanMove shows the fact that the request has previously reserved the
						// path, and is only trying to swap its way to the end.
						if !req.CanMove {
							for i := 1; i <= len(req.Paths[0])-1; i++ {
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsReserved = true
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).Reservation = reqNum
							}
						}
						req.CanMove = true
						numReached += quantum.ES(req, mgp.Network, mgp.RunTime)
					}
					isReady = true
				}
			}
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mgp.isFinished = true
			}
		}
	} else {
		numReached = 0
		oppCntr := 0
		for !mgp.isFinished {
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			//numReached = 0
			k := config.GetConfig().GetOpportunismDegree()
			isReady := true
			mgp.RunTime++
			quantum.EG(links)
			if !mgp.hasRecovery {
				for reqNum, req := range reqs {
					oppCntr = 0
					if req.HasReached == true {
						// Release the reserved links
						// Here, req.CanMove is used to release the links, and is set to false
						// to prevent extra work every time the request enters this if statement.
						if req.CanMove {
							for i := 1; i <= len(req.Paths[0])-1; i++ {
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsReserved = false
								mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).Reservation = -1
							}
						}
						req.CanMove = false
						continue
					}
					cntr = 0
					// req.Position starts from 1. Check this!!!!!!!!!!!!!!!!!!!!!!!!!
					for i := req.Position; i <= len(req.Paths[0])-1; i++ {
						//fmt.Println("Request num", reqNum, "position is", req.Position)
						link := mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1])
						//fmt.Println("link is reserved", link.IsReserved)
						if link.IsReserved == false {
							//fmt.Println("link not reserved. Link activation is", link.IsActive)
							isReady = isReady && link.IsActive
							cntr++
						} else {
							if link.Reservation == reqNum {
								//fmt.Println("corresponding reservation.")
								isReady = isReady && link.IsActive
								cntr++
							}
						}
						if isReady == true {
							//fmt.Println("oppCntr increment. oppCntr is:", oppCntr)
							oppCntr++
						} else {
							break
						}
						if cntr == 0 {
							isReady = false
							break
						}
					}
					//fmt.Println("Request", reqNum, oppCntr >= k)
					if oppCntr >= k {
						//if !req.CanMove {
						for i := req.Position; i <= req.Position+oppCntr-1; i++ {
							mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsReserved = true
							mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).Reservation = reqNum
						}
						//}
						req.CanMove = true
						numReached += quantum.ES(req, mgp.Network, mgp.RunTime)
					} else if (len(req.Paths[0]) - req.Position) <= oppCntr {
						//if !req.CanMove {
						for i := req.Position; i <= len(req.Paths[0])-1; i++ {
							mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).IsReserved = true
							mgp.Network.GetLinkBetween(req.Paths[0][i], req.Paths[0][i-1]).Reservation = reqNum
						}
						//}
						req.CanMove = true
						numReached += quantum.ES(req, mgp.Network, mgp.RunTime)
						fmt.Println("Fill in here. Maybe the remaining links are less than k, but are ready nonetheless.")
					}
					isReady = true
				}
			}
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mgp.isFinished = true
			}
		}
	}
}

func (mgp *ModifiedGreedyProfile) Stop() {
	mgp.isFinished = true
}

func (mgp *ModifiedGreedyProfile) Clear() {
	mgp.isFinished = false
	mgp.RunTime = 1
	mgp.Network.Clear()
}

func (mgp *ModifiedGreedyProfile) GetRunTime() int {
	return mgp.RunTime
}

func (mgp *ModifiedGreedyProfile) GetHasRecovery() bool {
	return mgp.hasRecovery
}
