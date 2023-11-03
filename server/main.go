package main

import (
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

type matcher struct {
	waiting chan net.Conn
}

func newMatcher() matcher {
	return matcher{
		make(chan net.Conn),
	}
}

func (m *matcher) match(c net.Conn) {
	io.WriteString(c, "Esperando a otro usuario\n")
	select {
	case other := <-m.waiting:
		game := newGame(other, c)
		game.start()
	case m.waiting <- c:
	}
}

type game struct {
	p1, p2 net.Conn
}

func newGame(p1, p2 net.Conn) game {
	return game{
		p1, p2,
	}
}

func (g *game) start() {
	io.WriteString(g.p1, "Usuario encontrado\n")
	io.WriteString(g.p2, "Usuario encontrado\n")

	go io.Copy(g.p1, g.p2)
	io.Copy(g.p2, g.p1)
}

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	matcher := newMatcher()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go matcher.match(c)
	}
}
