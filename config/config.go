package config

// The global configuration variable
var config Config

// The Config struct gathers all of the configuration parameters.
type Config struct {
	size     int
	memory   int
	lifetime int
	// The number of paths each request tries to reserve for itself.
	aggressiveness int
	p_gen          float64
	p_swap         float64
	hasRecovery    bool
}

// init() initializes the config variable once the package is imported.
func init() {
	config.size = 5
	config.memory = 4
	config.lifetime = 100
	config.p_gen = 0.5
	config.p_swap = 0.5
	config.aggressiveness = 1
	config.hasRecovery = false
}

// GetConfig returns the configuration.
func GetConfig() Config {
	return config
}

// SetConfig allows to manually set the configuration.
func SetConfig(size, memory, lifetime int, p_gen, p_swap float64, hasRecovery bool) {
	config.size = size
	config.memory = memory
	config.lifetime = lifetime
	config.p_gen = p_gen
	config.p_swap = p_swap
	config.hasRecovery = hasRecovery
}

// GetSize returns the size of the network.
func (conf Config) GetSize() int {
	return conf.size
}

// GetMemory returns the memory of every node.
func (conf Config) GetMemory() int {
	return conf.memory
}

// GetLifetime returns the life time of every entangled pair.
func (conf Config) GetLifetime() int {
	return conf.lifetime
}

// GetPGen returns the probability of a successful entanglement generation.
func (conf Config) GetPGen() float64 {
	return conf.p_gen
}

// GetPSwap returns the probability of a successful entanglement swapping.
func (conf Config) GetPSwap() float64 {
	return conf.p_swap
}

// GetAggressiveness returns the aggressiveness of the requests.
func (conf Config) GetAggressiveness() int {
	return conf.aggressiveness
}

func (conf Config) GetHasRecovery() bool {
	return conf.hasRecovery
}

// TODO: CHECK THIS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func (conf *Config) SetAggressiveness(aggressiveness int) {
	conf.aggressiveness = aggressiveness
}
