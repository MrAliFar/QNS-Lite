package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"

	"example.com/request"

	"example.com/quantum"
)

type ModifiedGreedyProfile struct {
	Network       *graph.Grid
	isFinished    bool
	hasRecovery   bool
	RunTime       int
	pathAlgorithm string
}

func (mgp *ModifiedGreedyProfile) Build(topology string) {
	mgp.RunTime = 1
	mgp.hasRecovery = config.GetConfig().GetHasRecovery()
	mgp.pathAlgorithm = path.MODIFIED_GREEDY
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		mgp.Network = grid
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

////////////////// TODO: The requests should remain the same through iterations.
////////////////// Take them out.
func (mgp *ModifiedGreedyProfile) Run(reqs []*request.Request, maxItr int) {
	links := mgp.Network.GetLinks()
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	//var cntr int
	whichPath := make([]int, len(reqs))
	if !isOpportunistic {
		for !mgp.isFinished {
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			///////////////////////////////// Check the following commented isReady.
			//isReady := true
			mgp.RunTime++
			// EG() also handles lifetimes.
			quantum.EG(links)
			if !mgp.hasRecovery {
				numReached, whichPath = noRecoveryRun(mgp.Network, reqs, whichPath, numReached, mgp.RunTime, false)
			}
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mgp.isFinished = true
			}
		}
	} else {
		numReached = 0
		for !mgp.isFinished {
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			//numReached = 0
			//k := config.GetConfig().GetOpportunismDegree()
			//isReady := true
			mgp.RunTime++
			quantum.EG(links)
			if !mgp.hasRecovery {
				numReached, whichPath = noRecoveryRunOPP(mgp.Network, reqs, whichPath, numReached, mgp.RunTime, false)
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

func (mgp *ModifiedGreedyProfile) GetNetwork() graph.Topology {
	return mgp.Network
}

func (mgp *ModifiedGreedyProfile) GetRunTime() int {
	return mgp.RunTime
}

func (mgp *ModifiedGreedyProfile) GetHasRecovery() bool {
	return mgp.hasRecovery
}

func (mgp *ModifiedGreedyProfile) GetPathAlgorithm() string {
	return mgp.pathAlgorithm
}
