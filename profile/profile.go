package profile

type profile interface {
	Build(topology string)
	Run()
	Stop()
	Benchmark()
}

// Each profile will have a unique profile id.
func BuildProfile(profileID int) profile {

}
