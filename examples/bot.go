package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"

	"github.com/toqueteos/ts3"
)

func main() {
	log.SetPrefix("DEBUG) ")
	// log.SetFlags(0)

	conn, err := ts3.Dial("127.0.0.1:10011")
	if err != nil {
		log.Fatalf("Dial failed with error %q\n", err)
	}
	defer conn.Close()
	// conn.SetTimeout(5 * time.Second)

	bot(conn)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	<-ctrlc
	conn.Send("quit")
}

// bot is a simple bot that checks version, signs in and sends a text message to
// channel#1 then exits.
func bot(conn *ts3.Conn) {
	var commands = []string{
		// Show version
		"version",
		// Login
		"login serveradmin +2rR7ebO",
		// Choose virtual server
		"use 1",
		// Update nickname
		`clientupdate client_nickname=My\sBot`,
		// "clientlist",
		// Send message to channel with id=1
		`sendtextmessage targetmode=2 target=1 msg=Bot\smessage!`,

		`clientlist`,

		`channellist`,

		`foo`,

		`servernotifyregister event=server`,

		// `servernotifyregister event=channel`,

		`servernotifyregister event=textserver`,

		`servernotifyregister event=textchannel`,

		// `servernotifyregister event=textprivate`,

		// `quit`,
	}

	for _, input := range commands {
		// fmt.Println("> ", input)

		_, err := conn.Send(input)
		if err != nil {
			log.Println(err)
		}
	}

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		input := scan.Text()
		_, err := conn.Send(input)
		if err != nil {
			log.Println(err)
		}
	}
}
