package main

import (
	"bluelabel/server/game"
	"bluelabel/server/matcher"
	"io"
	"log"
	"net"
)

const (
	listenAddr = "localhost:4000"
	lobbySize  = 3
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	matcher := matcher.NewMatcher(lobbySize, matcher.MatchStarterFunc(playGame))

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(c, "Looking for lobby...\n")
		go matcher.Match(c)
	}
}

func playGame(clientConnections ...net.Conn) {
	game := game.MakeGame(clientConnections...)
	game.PlayGame()
}
