package main

import (
	"fmt"

	"example.com/benchmark"
)

func main() {
	//var grid *graph.Grid = new(graph.Grid)
	//grid.Build()
	//ids := grid.GetNodeIDs()
	//fmt.Println("IDs", ids)
	//links := grid.GetLinks()
	//fmt.Println("Links", links)
	//fmt.Println("Links", links[2][0][1].ID)
	//num := 5
	//var priority []int
	//priority = make([]int, num)
	//for i := 0; i < num; i++ {
	//	priority[i] = 1
	//}
	//reqs, err := request.RG(num, ids, priority, "grid")
	//for i := 0; i < num; i++ {
	//	fmt.Println(i, reqs[i].Src)
	//	fmt.Println(reqs[i].Dest)
	//	for _, node := range reqs[i].Paths[0] {
	//		fmt.Println("PATHS FOR THIS REQUEST", node.ID)
	//	}
	//}
	//fmt.Println(reqs)
	//fmt.Println("Request generation error:", err)
	//path.PF(grid, reqs, "modified greedy", false)
	//check := 2
	//fmt.Println(reqs[check].Src, reqs[check].Dest)
	//for _, link := range links {
	//	fmt.Println("Before", link.IsActive)
	//}
	//quantum.EG(links)
	//for _, link := range links {
	//	fmt.Println("After", link.IsActive)
	//}
	/////////////////////////////////////////////////////// Implement lifetime!!!!!!!!
	itr := 5000
	bm := new(benchmark.Benchmarker)
	bm.Set(itr, "modified greedy", "grid")
	bm.Start(itr)
	fmt.Println(*bm)
	fmt.Println("The average waiting time is:", bm.AverageWaiting(5000))
	fmt.Println("The variance of the waiting time is:", bm.VarianceWaiting(5000))
	//fmt.Println(*bm)
}
