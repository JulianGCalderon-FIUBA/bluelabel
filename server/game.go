package main

import (
	"io"
	"net"
)

type game func(p1, p2 net.Conn)

func play_game(p1, p2 net.Conn) {
	io.WriteString(p1, "Usuario encontrado\n")
	io.WriteString(p2, "Usuario encontrado\n")

	go io.Copy(p1, p2)
	io.Copy(p2, p1)
}
