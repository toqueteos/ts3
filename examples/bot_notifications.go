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

// bot is a simple bot that registers itself for channel#1 text messages (using
// notifications).
func bot(conn *ts3.Conn) {
	defer conn.Cmd("quit")

	var cmdList = []string{
		// Login
		"login serveradmin 123456",
		// Choose virtual server
		"use 1",
		// Update nickname
		`clientupdate client_nickname=My\sBot`,
		// Register to channel id=1 text messages
		"servernotifyregister event=textchannel id=1",
	}

	// Chans returns a struct with three `chan string`.  We want `ch.Not` the
	// one that contains notifications.
	ch := conn.Chans()

	for _, cmdReq := range cmdList {
		// Send request to server (feed in a command)
		ch.In <- cmdReq + "\n"
		// Wait for its response
		cmdResp := <-ch.Out
		// Display as:
		//     > request
		//     response
		fmt.Printf("> %s\n%s", cmdReq, r)
		// Wait a bit after each command so we don't get banned. By default you
		// can issue 10 commands within 3 seconds.  More info on the
		// WHITELISTING AND BLACKLISTING section of TS3 ServerQuery Manual
		// (http://goo.gl/OpJXz).
		time.Sleep(350 * time.Millisecond)
	}

	// Keep
	for m := range nch {
		fmt.Printf("Notification: %s", m)
	}
}
