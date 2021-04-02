package quantum

import (
	"math/rand"
	"time"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/request"
)

type swapper interface {
}

type generator interface {
	genEntanglement(links []*graph.Link)
}

var p_gen, p_swap float64

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	p_gen = config.GetConfig().GetPGen()
	p_swap = config.GetConfig().GetPSwap()
}

func genEntanglement(links []*graph.Link) {
	lifetime := config.GetConfig().GetLifetime()
	//var r float64
	for _, link := range links {
		if link.IsActive == true {
			if link.Age == lifetime {
				link.IsActive = false
			} else {
				link.Age++
			}
			continue
		}
		if link.IsActive == false {
			//r = rand.Float64()
			if rand.Float64() <= p_gen {
				//fmt.Println("Link successfully generated!")
				link.IsActive = true
				link.Age = 0
			}
		}
	}
}

func genEntangledLink(link *graph.Link) {
	if rand.Float64() <= p_gen {
		link.IsActive = true
	}
}

// TODO: Write a less dumb swap!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func swap(req *request.Request, path path.Path, network graph.Topology, roundNum int, changeSrc bool) bool {
	//fmt.Println(req.Position)
	if req.Position == len(path)-1 {
		//fmt.Println("SWAP: HAS REACHED 1")
		req.HasReached = true
		req.ServingTime = roundNum
		return true
	}
	///////////////////// If swapping is unsuccessful, do the links get destroyed?
	///////////////////// If swapping is unsuccessful, shall we unreserve the links, as
	///////////////////// we have done here?
	if req.Position == 1 {
		if len(path) == 1 {
			//fmt.Println("Inside quantum.swap: Request destination is:", req.Dest.ID, "request position is:", req.PositionID, "Last node of path is:", path[len(path)-1].ID)
		}
		network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsActive = false
	}
	//fmt.Println("Inside quantum swap. Position is:", req.Position, "PositionID is", req.PositionID, "path length is:", len(path))
	network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsActive = false
	if rand.Float64() <= p_swap {
		//fmt.Println("SWAP: Swap successful.")
		if req.Position == 1 {
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsReserved = false
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).Reservation = -1
		}
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsReserved = false
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).Reservation = -1
		//if req.Position == 1 {
		//	network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsActive = false
		//}
		//network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsActive = false
		//req.Position++
		req.Position++
		////////////////////////////////// VERY IMPORTANT!!!!! THE -1 IS NECESSARY HERE!!!
		req.PositionID = path[req.Position-1].ID
		if changeSrc {
			////////////////////////////////////// VERY IMPORTANT!!!!! RELATED TO THE
			////////////////////////////////////// VERY IMPORTANT COMMENT ABOVE!
			req.Src = path[req.Position-1]
		}
		if req.Position == len(path)-1 {
			//fmt.Println("SWAP: HAS REACHED 2")
			req.HasReached = true
			req.ServingTime = roundNum
			return true
		}
		//fmt.Println("Quantum: Dest is:", req.Dest.ID, "PositionID is:", req.PositionID)
		if graph.IsEqual(req.Dest.ID, req.PositionID) {
			//fmt.Println("Inside quantum.swap, has reached: Request destination is:", req.Dest.ID, "request position is:", req.PositionID, "Last node of path is:", path[len(path)-1].ID)
			req.HasReached = true
			req.ServingTime = roundNum
			return true
		}
	} else {
		/////////// Very important! Check this! dramatically decreases performance!
		req.Position = 1
		req.Src = req.InitialSrc
		//fmt.Println("swap failed!")
	}
	return false
}

func EG(links []*graph.Link) {
	genEntanglement(links)
}

// TODO: This function should also take the round number as an input..........................
func ES(req *request.Request, network graph.Topology, roundNum int, whichPath int, changeSrc bool) int {
	var reached bool
	numReached := 0
	reached = swap(req, req.Paths[whichPath], network, roundNum, changeSrc)
	if reached == true {
		numReached++
	}
	return numReached
}
