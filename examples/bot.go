package main

import (
	"fmt"
	"github.com/toqueteos/ts3"
)

func main() {
	t, err := ts3.Dial("127.0.0.1")
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Cmd("help")
	t.Cmd("login serveradmin 123456")
	t.Cmd("use 1")
	t.Cmd(`clientupdate client_nickname=My\sBot`)
	t.Cmd("quit")
}
