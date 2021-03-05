package profile

const (
	MODIFIED_GREEDY = "modified greedy"
)

type Profile interface {
	Build(topology string)
	Run(numRequests int)
	Stop()
	Clear()
	GetRunTime() int
}

// Each profile will have a unique profile id.
//func BuildProfile(profileID int) profile {

//}
