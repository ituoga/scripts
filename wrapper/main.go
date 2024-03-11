package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"text/template"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	flagListen = flag.String("lt", "", "")
	flagStdin  = flag.Bool("stdin", true, "")

	argCommand = []string{}
)

func main() {
	flag.Parse()
	if *flagListen == "" {
		*flagListen = os.Getenv("APP_LISTEN")
		if *flagListen == "" {
			panic("listen topic cant be empty. use -lt [topic]")
		}
	}
	found := 0
	for i, v := range os.Args {
		if v == "--" {
			found = i + 1
		}
	}
	argCommand = os.Args[found:]

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	_ = nc

	_, err = nc.Subscribe(*flagListen, func(msg *nats.Msg) {
		go func(data []byte) {
			rmsg := new(nats.Msg)
			if rmsg.Header == nil {
				rmsg.Header = make(nats.Header)
			}

			log.Printf("%s", data)
			var buffoutput, bufferror bytes.Buffer
			command := argCommand[0]
			argz := []string{}
			if len(argCommand) > 1 {
				argz = argCommand[1:]
			}

			largz := []string{}
			largz = append(largz, argz...)

			for i, v := range argz {
				var buff []byte
				bwr := bytes.NewBuffer(buff)
				t, err := template.New(fmt.Sprintf("v%d", i)).Parse(v)
				if err != nil {
					panic(err)
				}
				err = t.Execute(bwr, map[string]string{
					"Payload": string(data),
					"Topic":   msg.Subject,
				})
				if err != nil {
					panic(err)
				}
				largz[i] = bwr.String()
			}

			cmd := exec.Command(command, largz...)

			cmd.Stdout = &buffoutput
			cmd.Stderr = &bufferror
			stdin, err := cmd.StdinPipe()
			if err != nil {
				rmsg.Header.Set("error", err.Error())
				msg.RespondMsg(rmsg)
			}

			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			err = cmd.Start()
			if err != nil {
				rmsg.Header.Set("error", err.Error())
				msg.RespondMsg(rmsg)
			}

			if *flagStdin {
				_, err = io.Copy(stdin, bytes.NewReader(data))
				if err != nil {
					rmsg.Header.Set("error", err.Error())
					msg.RespondMsg(rmsg)
				}
				stdin.Close()
			}

			done := make(chan error, 1)

			go func() {
				done <- cmd.Wait()
			}()

			select {
			case err := <-done:
				if err != nil {
					rmsg.Header.Set("error", err.Error())
					msg.RespondMsg(rmsg)
					return
				}
			case <-time.After(60 * time.Second):

				pgid, err := syscall.Getpgid(cmd.Process.Pid)
				if err != nil {
					rmsg.Header.Set("error", err.Error())
					msg.RespondMsg(rmsg)
					return
				}
				if err := syscall.Kill(-pgid, 15); err != nil {
					rmsg.Header.Set("error", err.Error())
					msg.RespondMsg(rmsg)
					return
				}
				rmsg.Header.Set("error", "command timeout")
				msg.RespondMsg(rmsg)
				return
			}

			rmsg.Header.Add("error", bufferror.String())
			rmsg.Data = buffoutput.Bytes()

			msg.RespondMsg(rmsg)
		}(msg.Data)
	})
	if err != nil {
		panic(err)
	}
	log.Println("Service running... " + *flagListen)
	<-make(chan struct{})
}
