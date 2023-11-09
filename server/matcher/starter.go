package matcher

import "net"

type MatchStarter interface {
	StartMatch(...net.Conn)
}

type MatchStarterFunc func(...net.Conn)

func (m MatchStarterFunc) StartMatch(connections ...net.Conn) {
	m(connections...)
}
