// TeamSpeak 3 Server Query library
//
// Reference: http://goo.gl/OpJXz
package ts3

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	DefaultPort    = "10011"
	VerificationID = "TS3"
)

// Custom logger
var log = stdlog.New(os.Stdout, "ts3> ", stdlog.LstdFlags)

var (
	ReadAfter    = 100 * time.Millisecond
	DialTimeout  = 1 * time.Second
	ReadTimeout  = 2 * time.Second
	WriteTimeout = 1500 * time.Millisecond
)

type Conn struct {
	addr string
	conn net.Conn
	// In holds messages written to server
	// Out holds messages read from server
	In, Out chan string
	quit    chan bool
	r       *bufio.Reader
	w       *bufio.Writer
}

// Dial connects to a local/remote TS3 server. A default port is appended to
// `addr` if user doesn't provide one.
func Dial(addr string) *Conn {
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
	ts3conn := &Conn{
		addr: addr,
		conn: c,
		In:   make(chan string),
		Out:  make(chan string),
		quit: make(chan bool),
		r:    bufio.NewReader(c),
		w:    bufio.NewWriter(c),
	}

	// Buffer to read from TCP socket; Read first line
	line, err = ts3conn.r.ReadString('\n')
	fatal(err, "Couldn't identify server.")
	fmt.Print(line)

	// Then check if it's a TS3 server
	if !strings.Contains(line, VerificationID) {
		log.Fatal("Not a TeamSpeak 3 server.")
	}

	// Show welcome message
	line, err = ts3conn.r.ReadString('\n')
	fatal(err, "Couldn't recv welcome message.")
	fmt.Print(line)

	// Copy flow: writer (request) -> conn -> reader (response)
	go cp(ts3conn.w, c)
	go cp(c, ts3conn.r)

	// Workers parse
	go ts3conn.inWorker()
	go ts3conn.outWorker()

	return ts3conn
}

// cp copies from an io.Reader to an io.Writer
func cp(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	fatal(err)
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

// inWorker writes (copies) incoming requests to its underlying io.Writer.
func (c *Conn) inWorker() {
	var err error

	for {
		select {
		case line := <-c.In:
			c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))

			_, err = c.w.WriteString(line)
			fatal(err)

			err = c.w.Flush()
			fatal(err)
		case <-c.quit:
			return
		}
	}
}

// inWorker reads (copies) incoming responses to its underlying io.Reader.
func (c *Conn) outWorker() {
	for {
		select {
		// Check if there's something to receive
		case <-time.After(ReadAfter):
			var (
				end   bool
				chunk string
			)

			// Read until end of response (error footer)
			for !end {
				line, err := c.r.ReadString('\n')
				fatal(err)

				// Trim network control chars
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

// fatal exits application if encounters an error
func fatal(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			log.Fatal(err)
		} else {
			log.Fatal(s)
		}
	}
}
