package benchmark

import (
	"fmt"

	"example.com/config"
	"example.com/profile"
)

type Benchmarker struct {
	Throughput       []float64
	TotalWaitingTime []int
	profile          profile.Profile
}

func (bm *Benchmarker) Set(itr int, prof string, topology string) {
	bm.Throughput = make([]float64, itr)
	bm.TotalWaitingTime = make([]int, itr)
	if prof == profile.MODIFIED_GREEDY {
		mgp := new(profile.ModifiedGreedyProfile)
		mgp.Build(topology)
		bm.profile = mgp
	} else {
		fmt.Println("Benchmark: Caution! The profile is not implemented.")
	}
}

func (bm *Benchmarker) Start(itr int) {
	for i := 0; i <= itr-1; i++ {
		//fmt.Println(*bm)
		fmt.Println("Iteration", i)
		bm.profile.Run(config.GetConfig().GetNumRequests())
		//fmt.Println(*bm)
		bm.TotalWaitingTime[i] = bm.profile.GetRunTime()
		bm.profile.Clear()
		//fmt.Println("GOTCHA!", *bm)
	}
	bm.profile.Stop()
}

func (bm *Benchmarker) AverageWaiting(maxItr int) float64 {
	sum := 0
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr {
			continue
		}
		sum += val
	}
	return float64(sum) / float64(len(bm.TotalWaitingTime))
}

func (bm *Benchmarker) VarianceWaiting(maxItr int) float64 {
	sum := float64(0)
	ave := bm.AverageWaiting(maxItr)
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr {
			continue
		}
		sum += (float64(val) - ave) * (float64(val) - ave)
	}
	return float64(sum) / float64(len(bm.TotalWaitingTime))
}
