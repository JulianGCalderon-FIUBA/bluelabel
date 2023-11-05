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

func (m *matcher) match(c net.Conn, g game) {
	select {
	case other := <-m.waiting:
		g(other, c)
	case m.waiting <- c:
	}
}
