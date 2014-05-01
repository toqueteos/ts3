// TeamSpeak 3 Server Query library
//
// Reference: http://goo.gl/OpJXz
package ts3

import (
	"bufio"
	"bytes"
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

var (
	// Custom logger
	log = stdlog.New(os.Stdout, "ts3> ", stdlog.LstdFlags)
	// ts3.Dial max timeout
	DialTimeout = 1 * time.Second
)

// Holds incoming, outgoing and notification requests (1) and responses (2 & 3)
type RWChan struct{ In, Out, Not chan string }

type Conn struct {
	conn net.Conn
	rw   RWChan
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
	conn, err := net.DialTimeout("tcp", addr, DialTimeout)
	fatal(err, fmt.Sprintf("Connection error: %v\n", err))

	// Allocate connection object
	ts3conn := &Conn{
		conn: conn,
		rw: RWChan{
			In:  make(chan string),
			Out: make(chan string),
			Not: make(chan string),
		},
	}

	rbuf := bufio.NewReader(conn)

	// Buffer to read from TCP socket; Read first line
	line, err = rbuf.ReadString('\n')
	fatal(err, "Couldn't identify server.")
	
	// Then check if it's a TS3 server
	if !strings.Contains(line, VerificationID) {
		log.Fatal("Not a TeamSpeak 3 server.")
	}

	// Show welcome message
	line, err = rbuf.ReadString('\n')
	fatal(err, "Couldn't recv welcome message.")

	// Copy flow: writer (request) -> conn -> reader (response)
	go cp(ts3conn, conn)
	go cp(conn, ts3conn)

	return ts3conn
}

// Read reads data from buffer into p doubling any IAC chars found (0xff), more
// info on RFC 854 (Telnet).  It returns the number of bytes read into p.
func (conn *Conn) Read(p []byte) (int, error) {
	b := []byte(<-conn.rw.In)
	// Double IAC chars
	bytes.Replace(b, []byte{0xff}, []byte{0xff, 0xff}, -1)
	copy(p, b)
	return len(b), nil
}

// Write writes the contents of p into the buffer. It returns the number of
// bytes written.
func (conn *Conn) Write(p []byte) (int, error) {
	s := string(p)
	conn.rw.Out <- s
	return len(p), nil
}

// Close closes underlying TCP Conn to local/remote server.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Cmd sends a request to a server and waits for its response.
func (c *Conn) Cmd(cmd string) string {
	var s string

	c.rw.In <- cmd + "\n"

	// Some commands output two lines
	var end bool
	for !end {
		select {
		case line := <-c.rw.Out:
			s += line
			if strings.HasPrefix(s, "notify") {
				c.rw.Not <- s
				end = true
			}
			if strings.HasPrefix(s, "error id=") {
				end = true
			}
		case <-time.After(500 * time.Millisecond):
			end = true
		}
	}

	return trimNet(s)
}

// Chans returns a low-level interaction
func (c *Conn) Chans() *RWChan {
	return &c.rw
}

// cp copies from an io.Reader to an io.Writer
func cp(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	fatal(err)
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
