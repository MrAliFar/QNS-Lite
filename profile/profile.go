package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/quantum"
	"example.com/request"
)

const (
	MODIFIED_GREEDY    = "modified greedy"
	NONOBLIVIOUS_LOCAL = "nonoblivious local"
)

type Profile interface {
	Build(topology string)
	Run(reqs []*request.Request, maxItr int)
	Stop()
	Clear()
	GetRunTime() int
	GetHasRecovery() bool
	GetNetwork() graph.Topology
	GetPathAlgorithm() string
}

func GenRequests(numRequests int, network graph.Topology, isMultiPath bool, algorithm string) []*request.Request {
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
	path.PF(network, reqs, algorithm)

	/*for _, req := range reqs {
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
	}*/
	return reqs
}

func noRecoveryRun(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int) (int, []int) {
	//numReached := 0
	isReady := true
	//whichPath := make([]int, len(reqs))
	var cntr int
	for reqNum, req := range reqs {
		//fmt.Println("Run - The req is: ", reqNum, " The path is: ", whichPath[reqNum])
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				//fmt.Println("Req ", reqNum, "Freed the resources.")
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
					//fmt.Println("Freeing link", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
				}
			}
			req.CanMove = false
			continue
		}
		cntr = 0
		if !req.CanMove {
			for which, _ := range req.Paths {
				//fmt.Println("Which is ", which, "req is ", reqNum)
				cntr = 0
				isReady = true
				for i := 1; i <= len(req.Paths[which])-1; i++ {
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					if link.IsReserved == false {
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							isReady = isReady && link.IsActive
							cntr++
						} else {
							//fmt.Println("1--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
							//fmt.Println("link is", link.ID)
							isReady = false
						}
					}
					if cntr == 0 {
						isReady = false
					}
					if !isReady {
						///////////////// IMPORTANT!!! CHECK THIS!!!!!
						break
					}
					// Solve the isReady issue.
				}
				if isReady {
					//fmt.Println("Run - The req is ", reqNum, "which is ", which)
					whichPath[reqNum] = which
					break
				}
			}
		} else {
			for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
				if link.IsReserved == false {
					isReady = isReady && link.IsActive
					cntr++
				} else {
					if link.Reservation == reqNum || link.Reservation == -1 {
						isReady = isReady && link.IsActive
						cntr++
					} else {
						//fmt.Println("2--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
						//fmt.Println("link is", link.ID)
						isReady = false
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
		}
		//fmt.Println("Request", reqNum, isReady)
		if isReady {
			// req.CanMove shows the fact that the request has previously reserved the
			// path, and is only trying to swap its way to the end.
			if !req.CanMove {
				//fmt.Println("Req ", reqNum, "is reserving for path ", whichPath[reqNum])
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
					//fmt.Println("Reserving link: ", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
				}
			}
			req.CanMove = true
			//fmt.Println("------------------------LENGTH IS: ", len(req.Paths[whichPath[reqNum]]), "WHICH IS: ", whichPath[reqNum], "LENGTH OF PATHS IS: ", len(req.Paths))
			numReached += quantum.ES(req, network, runTime, whichPath[reqNum])
			//if req.HasReached {
			//fmt.Println("Req ", reqNum, " Has reached!")
			//}
		}
		isReady = true
	}
	return numReached, whichPath
}

func noRecoveryRunOPP(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int) (int, []int) {
	isReady := true
	oppCntr := 0
	k := config.GetConfig().GetOpportunismDegree()
	//whichPath := make([]int, len(reqs))
	var cntr int
	for reqNum, req := range reqs {
		oppCntr = 0
		if req.CanMove {
			for i := 1; i <= req.Position-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
			}
		}
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
				}
			}
			req.CanMove = false
			continue
		}
		cntr = 0

		// req.Position starts from 1. Check this!!!!!!!!!!!!!!!!!!!!!!!!!
		if !req.CanMove {
			for which, _ := range req.Paths {
				for i := req.Position; i <= len(req.Paths[which])-1; i++ {
					//fmt.Println("Request num", reqNum, "position is", req.Position)
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					//fmt.Println("link is reserved", link.IsReserved)
					if link.IsReserved == false {
						//fmt.Println("link not reserved. Link activation is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
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
				if oppCntr >= k {
					whichPath[reqNum] = which
					break
				}
			}
		} else {
			for i := req.Position; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				//fmt.Println("Request num", reqNum, "position is", req.Position)
				link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
				//fmt.Println("link is reserved", link.IsReserved)
				if link.IsReserved == false {
					//fmt.Println("link not reserved. Link activation is", link.IsActive)
					isReady = isReady && link.IsActive
					cntr++
				} else {
					if link.Reservation == reqNum || link.Reservation == -1 {
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
		}
		//fmt.Println("Request", reqNum, oppCntr >= k)
		if oppCntr >= k {
			//if !req.CanMove {
			for i := req.Position; i <= req.Position+oppCntr-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true
			numReached += quantum.ES(req, network, runTime, whichPath[reqNum])
		} else if (len(req.Paths[whichPath[reqNum]]) - req.Position) <= oppCntr {
			//if !req.CanMove {
			for i := req.Position; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true
			numReached += quantum.ES(req, network, runTime, whichPath[reqNum])
			fmt.Println("Fill in here. Maybe the remaining links are less than k, but are ready nonetheless.")
		}
		isReady = true
	}
	return numReached, whichPath
}

// Each profile will have a unique profile id.
//func BuildProfile(profileID int) profile {

//}
