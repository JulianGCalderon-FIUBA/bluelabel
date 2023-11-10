package game

import (
	"io"
	"net"
)

func PlayGame(ps ...net.Conn) {
	for i, c1 := range ps {
		io.WriteString(c1, "Found lobby\n")

		others := make([]io.Writer, 0, len(ps)-1)
		for j, c2 := range ps {
			if i == j {
				continue
			}

			others = append(others, c2)
		}

		broadcast := io.MultiWriter(others...)

		go io.Copy(broadcast, c1)
	}
}
