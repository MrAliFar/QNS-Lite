package benchmark

import (
	"fmt"

	"example.com/config"
	"example.com/profile"
	"example.com/request"
)

type Benchmarker struct {
	keepReqs         bool
	regenerateReqs   bool
	refreshSources   bool
	Throughput       []float64
	TotalWaitingTime []int
	profile          profile.Profile
	reqs             []*request.Request
}

func (bm *Benchmarker) Set(itr int, prof string, topology string) {
	bm.Throughput = make([]float64, itr)
	bm.TotalWaitingTime = make([]int, itr)
	if prof == profile.MODIFIED_GREEDY {
		mgp := new(profile.ModifiedGreedyProfile)
		mgp.Build(topology)
		bm.profile = mgp
		bm.keepReqs = false
		bm.regenerateReqs = false
	} else {
		fmt.Println("Benchmark: Caution! The profile is not implemented.")
	}
}

func (bm *Benchmarker) Start(itr int, maxItr int) {
	///////////////////////// This might be unnecessary, since now we have the regenerateReqs()
	///////////////////////// function.
	if !bm.keepReqs {
		reqs := profile.GenRequests(config.GetConfig().GetNumRequests(), bm.profile.GetNetwork(), config.GetConfig().GetIsMultiPath(), bm.profile.GetPathAlgorithm())
		bm.reqs = reqs
	}
	sources := make([]*graph.Node, len(reqs))
	for m, req := range reqs {
		sources[m] = req.Src
	}
	for i := 0; i <= itr-1; i++ {
		//fmt.Println(*bm)
		//fmt.Println("Iteration", i)
		for m, req := range sources {
			reqs[m].Src = sources[m]
		}
		bm.profile.Run(bm.reqs, maxItr)
		//fmt.Println(*bm)
		bm.TotalWaitingTime[i] = bm.profile.GetRunTime()
		bm.profile.Clear()
		for _, req := range bm.reqs {
			request.ClearReq(req)
		}
		//fmt.Println("GOTCHA!", *bm)
	}
	bm.profile.Stop()
}

func (bm *Benchmarker) SetKeepReqs(keepReqs bool) {
	bm.keepReqs = keepReqs
}

func (bm *Benchmarker) RegenerateReqs() {
	bm.reqs = profile.GenRequests(config.GetConfig().GetNumRequests(), bm.profile.GetNetwork(), config.GetConfig().GetIsMultiPath(), bm.profile.GetPathAlgorithm())
}

func (bm *Benchmarker) AverageWaiting(maxItr int) float64 {
	sum := 0
	meanLength := len(bm.TotalWaitingTime)
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr {
			meanLength -= 1
			continue
		}
		sum += val
	}
	return float64(sum) / float64(meanLength)
}

func (bm *Benchmarker) VarianceWaiting(maxItr int) float64 {
	sum := float64(0)
	ave := bm.AverageWaiting(maxItr)
	varLength := len(bm.TotalWaitingTime)
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr {
			varLength--
			continue
		}
		sum += (float64(val) - ave) * (float64(val) - ave)
	}
	return float64(sum) / float64(varLength)
}
