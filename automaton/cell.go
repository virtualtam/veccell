package automaton

// A Cell can be alive or dead.
type Cell struct {
	Alive bool // FIXME not thread-safe
}
