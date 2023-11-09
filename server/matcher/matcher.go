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
	m := Matcher{
		make(chan net.Conn),
		lobbySize,
		starter,
	}

	go m.loop()

	return m
}

func (m *Matcher) loop() {
	for {
		connections := make([]net.Conn, m.lobbySize)
		for i := 0; i < m.lobbySize; i++ {
			connections[i] = <-m.waiting
		}
		go m.starter.StartMatch(connections...)
	}
}

func (m *Matcher) Match(connection net.Conn) {
	m.waiting <- connection
}
