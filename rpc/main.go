package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if len(os.Args) < 2 {
		panic("topic cant be empty ./app call.mq.topic")
	}
	var body string
	if len(os.Args) > 2 {
		body = os.Args[2]
	}
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	_ = nc
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Printf("%v", err)
	}

	_ = fi
	// Check for stdin pipe presence.
	// If in a pipeline, read from stdin.
	// Otherwise, skip to prevent hanging as the program awaits input.
	if !(fi.Mode()&os.ModeNamedPipe == 0) {
		b, _ := io.ReadAll(os.Stdin)
		body = string(b)
	}

	topic := os.Args[1]
	msg, err := nc.Request(topic, []byte(body), 60*time.Second)
	if err != nil {
		panic(err)
	}
	if msg.Header.Get("error") != "" {
		log.Printf("Error: %s", msg.Header.Get("error"))
		return
	}
	println(string(msg.Data))

}
