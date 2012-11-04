package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/toqueteos/ts3"
)

func main() {
	in := make(chan string)
	q := make(chan bool)

	// Connect to Server Query daemon
	c, err := net.DialTimeout("tcp", "127.0.0.1:10011", time.Second)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	buf := bufio.NewReader(c)
	go stdin(in, q)

	go func(c net.Conn, in chan string, q chan bool) {
		for {
			select {
			case line := <-in:
				c.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
				c.Write([]byte(line))
			case <-q:
				return
			}
		}
	}(c, in, q)

	for {
		c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		s, err := buf.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if len(s) > 1 {
			fmt.Print(s)
		}
	}

}

func stdin(chIn chan string, q chan bool) {
	// Feed ts3.In
	in := bufio.NewReader(os.Stdin)

	for {
		// repl thingy
		line, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("ts3> stdin read error: %v\n", err)
		}
		line = ts3.StringsTrimNet(line)

		// Ignore empty lines
		if line != "\n" {
			chIn <- line

			if strings.HasPrefix(line, "quit") {
				q <- true
				break
			}
		}
	}
}
