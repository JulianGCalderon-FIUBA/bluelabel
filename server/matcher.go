package main

import (
	"net"
)

type matcher struct {
	waiting chan net.Conn
}

func newMatcher() matcher {
	return matcher{
		make(chan net.Conn),
	}
}

func (m *matcher) match(connection net.Conn, game game) {
	select {
	case other := <-m.waiting:
		game(other, connection)
	case m.waiting <- connection:
	}
}
