package main

import (
	"fmt"
	"time"

	"github.com/toqueteos/ts3"
)

func main() {
	conn := ts3.Dial(":10011")
	defer conn.Close()

	bot(conn)
}


// bot is a simple bot that checks version, signs in and sends a text message to
// channel#1 then exits.
func bot(conn *ts3.Conn) {
	defer conn.Cmd("quit")

	var cmds = []string{
		// Show version
		"version",
		// Login
		"login serveradmin 123456",
		// Choose virtual server
		"use 1",
		// Update nickname
		`clientupdate client_nickname=My\sBot`,
		// "clientlist",
		// Send message to channel with id=1
		"sendtextmessage targetmode=2 target=1 msg=Bot\smessage!",
	}

	for _, s := range cmds {
		// Send command and wait for its response
		r := conn.Cmd(s)
		// Display as:
		//     > request
		//     response
		fmt.Printf("> %s\n%s", s, r)
		// Wait a bit after each command so we don't get banned.  serveradmin
		// doesn't need this if bot's ip is on whitelist.
		time.Sleep(200 * time.Millisecond)
	}
}
