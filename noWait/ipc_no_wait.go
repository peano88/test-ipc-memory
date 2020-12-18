package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/siadat/ipc"
)

const (
	queueIn                = 234
	queueOut               = 235
	nrMsgs                 = 2000
	mtype                  = 456
	duration time.Duration = 20 * time.Millisecond
)

func main() {

	var memprofile = flag.String("memprofile", "", "write memory profile to this file")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	qid, err := ipc.Msgget(queueIn, ipc.IPC_CREAT|0600)
	if err != nil {
		log.Fatalf("Fatal error while getting id for queue %d", queueIn)
	}

	payload := make([]byte, 4096)
	buffer := bytes.NewBuffer(payload)
	buffer.WriteString("This is a completely trivial example")

	msg := &ipc.Msgbuf{
		Mtype: mtype,
		Mtext: buffer.Bytes(),
	}

	start := make(chan bool, 1)

	sendMsgs := func() {
		<-start
		for i := 0; i < nrMsgs; i++ {
			for {
				fmt.Printf("Sending message %d of length %d \n", i, len(msg.Mtext))
				if err := ipc.Msgsnd(qid, msg, ipc.IPC_NOWAIT); err != nil {
					if errors.Is(err, syscall.EAGAIN) {
						//time.Sleep(duration)
						continue
					}
					log.Fatalf("sending error: %v, index %d", err, i)
					return
				}

				break

			}

		}
	}

	go sendMsgs()
	msgChan := make(chan *ipc.Msgbuf, nrMsgs)
	receiveMsg := func() {
		for {
			respBuf := &ipc.Msgbuf{}
			time.Sleep(duration)
			if err = ipc.Msgrcv(qid, respBuf, ipc.IPC_NOWAIT); err != nil {
				if errors.Is(err, syscall.ENOMSG) {
					continue
				}
				log.Fatalf("receiving error: %v", err)
			}
			msgChan <- respBuf
		}
	}

	go receiveMsg()

	start <- true

	for i := 0; i < nrMsgs; i++ {
		<-msgChan
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
