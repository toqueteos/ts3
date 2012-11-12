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
		`sendtextmessage targetmode=2 target=1 msg=Bot\smessage!`,
	}

	for _, s := range cmds {
		// Send command and wait for its response
		r := conn.Cmd(s)
		// Display as:
		//     > request
		//     response
		fmt.Printf("> %s\n%s", s, r)
		// Wait a bit after each command so we don't get banned. By default you
		// can issue 10 commands within 3 seconds.  More info on the
		// WHITELISTING AND BLACKLISTING section of TS3 ServerQuery Manual
		// (http://goo.gl/OpJXz).
		time.Sleep(350 * time.Millisecond)
	}
}
