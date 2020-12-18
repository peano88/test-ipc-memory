package main

import (
	"log"

	"github.com/siadat/ipc"
)

const (
	queueIn = 234
	queueOut = 235
	nrMsgs = 20000
	mtype = 456
)

func main() {
		qid, err := ipc.Msgget(queueIn, ipc.IPC_CREAT|0600)
		if err != nil {
			log.Fatalf("Fatal error while getting id for queue %d", queueIn)
		}

		msg := &ipc.Msgbuf{
			Mtype: mtype,
			Mtext: []byte("This is a completely trivial example"),
		}
	sendMsgs := func() {
		for i := 0; i < nrMsgs; i++ {
			if err := ipc.Msgsnd(qid, msg, 0); err != nil {
				log.Fatalf("sending error: %v", err)
				return
			}

		}
	}

	go sendMsgs()

	for i := 0; i < nrMsgs; i++ {
		respBuf := &ipc.Msgbuf{}
		if err = ipc.Msgrcv(qid, respBuf, 0); err != nil {
			log.Fatalf("receiving error: %v", err)
		}
		
	}
}
