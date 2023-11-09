package matcher

import (
	"net"
)

// A `matcher` is used to create fixed-sized lobbies of `net.Conn`. For every `n` connections
// it receives, it constructs a match with those connections.
type Matcher struct {
	waiting   chan net.Conn
	lobbySize int
	starter   MatchStarter
}

// Creates a new `matcher` with a given `lobbySize` and a `starter` function.
func NewMatcher(lobbySize int, starter MatchStarter) Matcher {
	m := Matcher{
		make(chan net.Conn),
		lobbySize,
		starter,
	}

	go m.loop()

	return m
}

// Loops on the waiting channel, starting a match when enough connections are sent.
func (m *Matcher) loop() {
	for {
		connections := make([]net.Conn, m.lobbySize)
		for i := 0; i < m.lobbySize; i++ {
			connections[i] = <-m.waiting
		}
		go m.starter.StartMatch(connections...)
	}
}

// Tries to find a match for the given connection. When enough connections are sent, a match is started. This
// functions is blocking until the loop goroutine retreives the value (not until a match is found). This
// delay should be really small.
func (m *Matcher) Match(connection net.Conn) {
	m.waiting <- connection
}
