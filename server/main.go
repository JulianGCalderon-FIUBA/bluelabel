package main

import (
	"bluelabel/server/matcher"
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

	matcher := matcher.NewMatcher(matcher.MatchStarterFunc(play_game))

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(c, "Esperando usuario...\n")
		go matcher.Match(c)
	}
}
