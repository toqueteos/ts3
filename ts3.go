// TeamSpeak 3 Server Query library
//
// Reference: http://goo.gl/OpJXz
package ts3

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	DefaultPort    = "10011"
	VerificationID = "TS3"
)

var (
	WriteRetries = 3
	//
	ReadAfter    = 100 * time.Millisecond
	DialTimeout  = time.Second
	ReadTimeout  = 250 * time.Millisecond
	WriteTimeout = 500 * time.Millisecond
)

type Conn struct {
	addr string
	conn net.Conn
	// In holds messages written to server
	// Out holds messages read from server
	In, Out chan string
	quit    chan bool
}

// Dial connects to a local/remote TS3 server. A default port is appended to
// `addr` if user doesn't provide one.
func Dial(addr string) (*Conn, error) {
	var (
		err  error
		line string
	)

	// Append DefaultPort if user didn't specify one
	if !strings.Contains(addr, ":") {
		addr += ":" + DefaultPort
	}

	// Try to establish connection
	c, err := net.DialTimeout("tcp", addr, DialTimeout)
	fatal(err, fmt.Sprintf("Connection error: %v\n", err))

	// Allocate connection object
	t := &Conn{
		addr: addr,
		conn: c,
		In:   make(chan string),
		Out:  make(chan string),
		quit: make(chan bool),
	}

	// Buffer to read from TCP socket; Read first line
	buf := bufio.NewReader(c)
	line, err = buf.ReadString('\n')
	fatal(err, "Couldn't identify server.")
	fmt.Print(line)

	// Then check if it's a TS3 server
	if !strings.Contains(line, VerificationID) {
		log.Fatal("Not a TeamSpeak 3 server.")
	}

	// Show welcome message
	line, err = buf.ReadString('\n')
	fatal(err, "Couldn't recv welcome message.")
	fmt.Print(line)

	// Init workers
	go t.outWorker()
	go t.inWorker()

	t.In <- "version"
	fmt.Println("Checking server version...")
	fmt.Println(<-t.Out)

	return t, err
}

// Close closes underlying TCP Conn to local/remote server.
func (c *Conn) Close() error {
	// Two workers need two quit signals
	c.quit <- true
	c.quit <- true

	return c.conn.Close()
}

// Cmd sends an arbitrary command to the server and expects an output.
func (c *Conn) Cmd(s string) {
	fmt.Println("> " + s)
	c.In <- s
	fmt.Println(<-c.Out)
}

// outWorker reads server responses and sends them to `t.Out` channel.
func (c *Conn) outWorker() {
	for {
		select {
		// Check if there's something to receive
		case <-time.After(ReadAfter):
			log.Println("I'm reading something!")

			var (
				end   bool
				chunk string
			)

			c.conn.SetReadDeadline(time.Now().Add(ReadTimeout))

			// Read until end of response (error footer)
			buf := bufio.NewReader(c.conn)
			for !end {
				log.Println("Inside retry for loop!")
				line, err := buf.ReadString('\n')
				fatal(err, fmt.Sprint(err))

				// Trim network control chars and save line to chunk
				line = StringsTrimNet(line)
				chunk += line

				// Determine end of line
				if strings.HasPrefix(line, "error id=") ||
					strings.Contains(line, "notify") {
					c.Out <- chunk
					// EOL found, finish iteration
					end = true
				}
			}
		case <-c.quit:
			return
		}
	}
}

// outWorker writes server requests and sends them to `t.In` channel.
func (c *Conn) inWorker() {
	var (
		retry = true
		tries int
	)

	buf := bufio.NewWriter(c.conn)

	for {
		select {
		case line := <-c.In:
			// Ensure line is sent to remote telnetd
			for retry && (tries < WriteRetries) {
				c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))

				_, err := buf.WriteString(line)

				if err != nil {
					tries++
				} else {
					// Exit loop
					retry = false
				}
			}

			// Reset write ensurance
			retry, tries = true, 0
		case <-c.quit:
			return
		}
	}
}

// init registers log package prefix and removes any flags
func init() {
	log.SetPrefix("ts3> ")
	log.SetFlags(0)
}

// fatal exits application if encounters an error
func fatal(err error, s string) {
	if err != nil {
		log.Fatal(s)
	}
}

// info logs an error without exitting from application
func info(err error, s string) {
	if err != nil {
		log.Println(s)
	}
}

// info logs an error without exitting from application
func infoErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
