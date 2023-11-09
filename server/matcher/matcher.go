package matcher

import (
	"net"
)

type Matcher struct {
	waiting chan net.Conn
	starter MatchStarter
}

func NewMatcher(starter MatchStarter) Matcher {
	return Matcher{
		make(chan net.Conn),
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
