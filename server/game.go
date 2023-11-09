package main

import (
	"io"
	"net"
)

func play_game(ps ...net.Conn) {
	for i, c1 := range ps {
		io.WriteString(c1, "Found lobby\n")

		for j, c2 := range ps {
			if i == j {
				continue
			}

			go io.Copy(c1, c2)
		}
	}
}
