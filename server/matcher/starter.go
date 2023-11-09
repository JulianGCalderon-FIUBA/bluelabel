package matcher

import "net"

// This interface is used for the `matcher` when a lobby is filled.
type MatchStarter interface {
	StartMatch(...net.Conn)
}

// If the `MatchStarter` is a function, rather than an interface, can be casted
// into a `MatchStarterFunc`, which already implements `MatchStarter`.
type MatchStarterFunc func(...net.Conn)

func (m MatchStarterFunc) StartMatch(connections ...net.Conn) {
	m(connections...)
}
