package main

import (
	"bluelabel/shared"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func read_from_connection(connection net.Conn) {
	fmt.Printf("%s Client\n", shared.Greet())
	for {
		tmp_buffer := make([]byte, 256)
		tmp_buffer_msg_len, err := connection.Read(tmp_buffer)

		if err != nil {
			break
		}

		fmt.Println(string(tmp_buffer[:tmp_buffer_msg_len]))
	}
}

func main() {
	connection, err := net.Dial("tcp", "localhost:4000")
	input_scanner := bufio.NewScanner(os.Stdin)

	if err == nil {
		go read_from_connection(connection)
		for {
			if input_scanner.Scan() {
				line := input_scanner.Text()
				io.WriteString(connection, line)
			}
		}
	}
}
