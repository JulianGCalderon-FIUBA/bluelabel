package main

import (
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

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

		io.WriteString(c, "Esperando usuario...\n")
		go matcher.match(c, play_game)
	}
}
