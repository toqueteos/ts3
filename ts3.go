package ts3

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	DefaultPort    = "10011"
	VerificationID = "TS3"
)

var (
	// ts3.Dial max timeout
	DialTimeout    = 1 * time.Second
	CommandTimeout = 500 * time.Millisecond
)

type notification func(string, string)

type Conn struct {
	conn          net.Conn
	rbuf          *bufio.Reader
	wbuf          *bufio.Writer
	cResult       chan string
	cNotification chan string
	cError        chan ErrorMsg
	notifyCb      notification
	cmdSync       *sync.WaitGroup
	cmdLast       time.Time
	cmdTimeout    time.Duration
}

type ErrorMsg struct {
	Id  int
	Msg string
}

// Dial connects to a local/remote TS3 server. A default port is appended to
// `addr` if user doesn't provide one.
func Dial(addr string, whitelisted bool) (*Conn, error) {
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
	if err != nil {
		return nil, err
	}

	// Create buffers
	rbuf := bufio.NewReader(conn)
	wbuf := bufio.NewWriter(conn)

	timeout := time.Duration(350 * time.Millisecond)
	if whitelisted {
		timeout = time.Duration(0 * time.Millisecond)
	}

	// Allocate connection object
	ts3conn := &Conn{
		conn:          conn,
		rbuf:          rbuf,
		wbuf:          wbuf,
		cResult:       make(chan string),
		cNotification: make(chan string),
		cError:        make(chan ErrorMsg),
		notifyCb:      nil,
		cmdSync:       new(sync.WaitGroup),
		cmdLast:       time.Now(),
		cmdTimeout:    timeout,
	}

	// Read VerificationID
	line, err = rbuf.ReadString('\n')

	// Then check if it's a TS3 server
	if !strings.Contains(line, VerificationID) {
		return nil, errors.New("Invalid VerificationID")
	}

	// Read welcome message
	line, err = rbuf.ReadString('\n')

	// Sync socket incoming data, to ts3conn
	go ts3conn.sync()

	return ts3conn, nil
}

// Close closes underlying TCP Conn to local/remote server.
func (this *Conn) Close() error {
	return this.conn.Close()
}

// Cmd sends a request to a server and waits for its response.
func (this *Conn) Cmd(cmd string) (string, ErrorMsg) {
	//Make sure only 1 command runs at a time per connection
	this.cmdSync.Wait()
	this.cmdSync.Add(1)

	//Make sure we timeout nicely per command to avoid spamming
	diff := time.Since(this.cmdLast)
	if diff < this.cmdTimeout {
		<-time.After(this.cmdTimeout - diff)
	}

	var (
		result string   //Holds end result
		temp   string   //Holds temp result
		err    ErrorMsg //Holds error message
	)

	// Write the cmd to the socket, ending with a \n
	this.send(cmd + "\n")

	// Block on a channel that will recieve the socket READ that is NOT a notification or error
	// Block on a channel that will recieve the socket READ that is NOT a notification or result
	done := false
	for !done {
		select {
		case temp = <-this.cResult:
			result += temp
			continue
		case err = <-this.cError:
			done = true
		case <-time.After(CommandTimeout):
			//TODO Is there a better way to handle this?
			err.Id = 1
			err.Msg = "timeout"
			done = true
		}
	}

	this.cmdLast = time.Now()
	this.cmdSync.Done()
	return result, err
}

func (this *Conn) NotifyFunc(cb notification) {
	this.notifyCb = cb
}

func (this *Conn) send(p string) (int, error) {
	b := []byte(p)
	// Double IAC chars
	bytes.Replace(b, []byte{0xff}, []byte{0xff, 0xff}, -1)
	return this.conn.Write(b)
}

func (this *Conn) handleResponse(data string) {
	// Has to be done by line!
	// Sadly its possible, when spammed, to have an error and notification to stick together, and more
	// Should be a better way to solve this overall, but teamspeak is a bitch and doesnt have a standard way to end a reply

	lines := strings.Split(data, "\n")
	var result string

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "error"):
			this.cError <- parseError(line)
		case strings.HasPrefix(line, "notify"):
			if this.notifyCb != nil {
				split := strings.SplitN(line, " ", 2)
				go this.notifyCb(split[0], split[1])
			}
		default:
			result += line
		}
	}

	if result != "" {
		this.cResult <- result
	}
}

// Write writes the contents of p into the buffer. It returns the number of bytes written.
func (this *Conn) Write(p []byte) (int, error) {
	s := string(p)
	s = strings.Replace(s, "\r", "", -1)

	this.handleResponse(s)
	return len(s), nil
}

// cp copies from an io.Reader to an io.Writer
func (this *Conn) sync() {
	for {
		io.Copy(this, this.conn)
	}
}
