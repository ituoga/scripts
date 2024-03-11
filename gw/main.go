package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	_ = nc
	http.HandleFunc("/mq", func(w http.ResponseWriter, r *http.Request) {
		rmsg := new(nats.Msg)
		rmsg.Subject = r.URL.Query().Get("t")

		if r.Method == "GET" {

		} else if r.Method == "POST" {
			b, _ := io.ReadAll(r.Body)
			rmsg.Data = b
		}
		msg, err := nc.RequestMsg(rmsg, 60*time.Second)
		if err != nil {
			panic(err)
		}
		if msg.Header.Get("error") != "" {
			w.WriteHeader(500)
			w.Write([]byte(msg.Header.Get("error")))
			return
		}
		w.Write(msg.Data)
	})

	port := "80"
	if p, ok := os.LookupEnv("APP_PORT"); ok {
		port = p
	}

	http.ListenAndServe(":"+port, nil)
}
