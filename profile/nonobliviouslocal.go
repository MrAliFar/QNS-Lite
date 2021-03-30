package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/quantum"
	"example.com/request"
)

type NonObliviousLocal struct {
	Network       *graph.Grid
	isFinished    bool
	hasRecovery   bool
	RunTime       int
	pathAlgorithm string
}

func (nol *NonObliviousLocal) Build(topology string) {
	nol.RunTime = 1
	nol.hasRecovery = config.GetConfig().GetHasRecovery()
	nol.pathAlgorithm = path.MODIFIED_GREEDY
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		nol.Network = grid
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (nol *NonObliviousLocal) Run(reqs []*request.Request, maxItr int) {
	links := nol.Network.GetLinks()
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	//var cntr int
	whichPath := make([]int, len(reqs))

	if !isOpportunistic {
		for !nol.isFinished {
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			nol.RunTime++
			path.ClearReqPaths(reqs)
			// EG() also handles lifetimes.
			quantum.EG(links)
			path.PF(nol.Network, reqs, "nonoblivious local")
			if !nol.hasRecovery {
				numReached, whichPath = noRecoveryRun(nol.Network, reqs, whichPath, numReached, nol.RunTime)
			}
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				nol.isFinished = true
			}
		}
	} else {
		numReached = 0
		for !nol.isFinished {
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			//numReached = 0
			//k := config.GetConfig().GetOpportunismDegree()
			//isReady := true
			nol.RunTime++
			quantum.EG(links)
			if !nol.hasRecovery {
				numReached, whichPath = noRecoveryRunOPP(nol.Network, reqs, whichPath, numReached, nol.RunTime)
			}
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				nol.isFinished = true
			}
		}
	}
}

func (nol *NonObliviousLocal) Stop() {
	nol.isFinished = true
}

func (nol *NonObliviousLocal) Clear() {
	nol.isFinished = false
	nol.RunTime = 1
	nol.Network.Clear()
}

func (nol *NonObliviousLocal) GetNetwork() graph.Topology {
	return nol.Network
}

func (nol *NonObliviousLocal) GetRunTime() int {
	return nol.RunTime
}

func (nol *NonObliviousLocal) GetHasRecovery() bool {
	return nol.hasRecovery
}

func (nol *NonObliviousLocal) GetPathAlgorithm() string {
	return nol.pathAlgorithm
}
