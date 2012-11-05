TeamSpeak 3 Server Query Library

# Simple usage

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

# Notifications support

Check `examples/bot_notifications.go`
