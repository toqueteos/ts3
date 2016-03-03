package ts3

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const DefaultPort = 10011

type notification struct {
	Contents string
}

// Conn represents an established low level connection to a TS3 server.
type Conn struct {
	conn          net.Conn
	rw            *bufio.ReadWriter
	timeout       time.Duration
	stop          bool
	messages      chan message
	notify        bool
	notifications chan notification
}

// Dial establishes a raw connection to a TS3 server.
func Dial(addr string) (conn *Conn, err error) {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%s:%d", addr, DefaultPort)
	}

	conn = new(Conn)
	conn.conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Printf("net.Dial(tcp, %q) failed with error %q\n", addr, err)
		return nil, err
	}

	rb := bufio.NewReader(conn.conn)
	wb := bufio.NewWriter(conn.conn)
	conn.rw = bufio.NewReadWriter(rb, wb)
	conn.timeout = 5 * time.Second
	conn.messages = make(chan message)
	conn.notifications = make(chan notification)

	var output string
	for i := 0; i < 2; i++ {
		output, err = conn.readNext()
		if err != nil {
			return nil, err
		}
		log.Printf("< %q\n", output)
	}

	go conn.recv()

	return conn, nil
}

func (c *Conn) UseNotifications(value bool) { c.notify = value }

func (c *Conn) SetTimeout(timeout time.Duration) { c.timeout = timeout }

func (c *Conn) Close() error {
	c.stop = true
	return c.conn.Close()
}

func (c *Conn) recv() error {
	var (
		message message
		output  string
		err     error
	)
	for !c.stop {
		output, err = c.readNext()
		if err != nil {
			if nerr, ok := err.(net.Error); ok {
				if !nerr.Timeout() {
					log.Println("server response error", nerr.Error())
				}
			}

			continue
		}

		if strings.HasPrefix(output, "notifytextmessage ") {
			log.Println("N", output)
			if c.notify {
				// c.notifications <- output
			}

			continue
		}

		if strings.HasPrefix(output, "error ") {
			message.Error = NewErrorString(output)
			c.messages <- message

			message.Contents = ""
			message.Error = nil

			continue
		}

		message.Contents = output
	}

	return nil
}

func (c *Conn) readNext() (output string, err error) {
	c.conn.SetDeadline(time.Now().Add(c.timeout))
	output, err = c.rw.ReadString('\r')
	if err != nil {
		return
	}

	output = strings.TrimSpace(output)

	return
}

const EOL = "\n\r"

func (c *Conn) Send(input string) (response string, err error) {
	if !strings.HasSuffix(input, EOL) {
		input += EOL
	}

	log.Printf("> %q\n", input)
	if err = c.writeString(input); err != nil {
		return "", err
	}

	message := <-c.messages
	if len(message.Contents) > 0 {
		log.Printf("< Contents: %q\n", message.Contents)
	}
	if message.Error != nil {
		log.Printf("< Error: %q\n", message.Error)
	}

	return message.Contents, message.Error
}

func (c *Conn) writeString(input string) (err error) {
	if _, err = c.rw.WriteString(input); err != nil {
		return
	}

	return c.rw.Flush()
}
