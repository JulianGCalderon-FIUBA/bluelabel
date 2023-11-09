package matcher

import (
	"net"
)

type Matcher struct {
	waiting   chan net.Conn
	lobbySize int
	starter   MatchStarter
}

func NewMatcher(lobbySize int, starter MatchStarter) Matcher {
	return Matcher{
		make(chan net.Conn),
		lobbySize,
		starter,
	}
}

func (m *Matcher) Match(connection net.Conn) {
	select {
	case other := <-m.waiting:
		m.starter.StartMatch(connection, other)
	case m.waiting <- connection:
	}
}
