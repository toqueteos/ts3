package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/toqueteos/ts3"
)

func main() {
	conn, err := net.Dial("tcp", ":10011")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	errc := make(chan error)
	go cp(os.Stdout, conn, errc)
	// go cp(conn, os.Stdin, errc)

	// repl thingy
	in := bufio.NewReader(os.Stdin)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			return
		}
		line = ts3.StringsTrimNet(line)

		// Ignore empty lines
		if line != "\n" {
			conn.Write([]byte(line))

			if strings.HasPrefix(line, "quit") {
				return
			}
		}
	}

	log.Fatal(<-errc)
}

func cp(dst io.Writer, src io.Reader, errc chan<- error) {
	_, err := io.Copy(dst, src)
	errc <- err
}
