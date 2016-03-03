package ts3

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDial(t *testing.T) {
	conn, err := Dial("127.0.0.1", true)
	if err != nil {
		t.Fatalf("Dial failed with error %q", err)
	}

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
