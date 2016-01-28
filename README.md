# ts3 [![](https://godoc.org/github.com/toqueteos/ts3?status.svg)](http://godoc.org/github.com/toqueteos/ts3)

A TeamSpeak 3 Server Query library.

The API is stable but horrible. Hopefully it's about to change, this project was done quite some time ago and it needs a bit of love from my side.

# Installation

`go get github.com/toqueteos/ts3/...`

# Example

```go
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

func bot(conn *ts3.Conn) {
    defer conn.Cmd("quit")

    s := "version"
    r := conn.Cmd(s)
    fmt.Printf("> %s\n%s", s, r)
}
```

# Notifications support

Check `examples/bot_notifications.go`
