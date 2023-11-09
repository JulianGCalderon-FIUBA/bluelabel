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

	matcher := matcher.NewMatcher(2, matcher.MatchStarterFunc(play_game))

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(c, "Looking for lobby...\n")
		go matcher.Match(c)
	}
}
