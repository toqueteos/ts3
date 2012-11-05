package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/toqueteos/ts3"
)

func main() {
	conn := ts3.Dial(":10011")
	defer conn.Close()

	// repl thingy
	in := bufio.NewReader(os.Stdin)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			defer conn.Cmd("quit")
			return
		}

		// Ignore empty lines
		if line != "\n" {
			fmt.Print(conn.Cmd(line))

			if strings.HasPrefix(line, "quit") {
				return
			}
		}
	}
}
