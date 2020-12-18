package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"
	"syscall"
	"unsafe"

	"github.com/siadat/ipc"
)

var uintSize = bits.UintSize / 8


func prepareMsg(msg *ipc.Msgbuf) (*bytes.Buffer, error) {
	if len(msg.Mtext) > ipc.Msgmax() {
		return nil, fmt.Errorf("mtext is too large, %d > %d", len(msg.Mtext), ipc.Msgmax())
	}

	buf := make([]byte, uintSize+ ipc.Msgmax())
	buffer := bytes.NewBuffer(buf)
	buffer.Reset()
	var err error
	switch uintSize {
	case 4:
		err = binary.Write(buffer, binary.LittleEndian, uint32(msg.Mtype))
	case 8:
		err = binary.Write(buffer, binary.LittleEndian, uint64(msg.Mtype))
	}
	if err != nil {
		return nil, fmt.Errorf("Can't write binary: %v", err)
	}
	buffer.Write(msg.Mtext)

	return buffer, nil

}

func msgsndPrepared(qid uint, len int, buffer *bytes.Buffer, flags uint) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_MSGSND,
		uintptr(qid),
		uintptr(unsafe.Pointer(&buffer.Bytes()[0])),
		uintptr(len),
		uintptr(flags),
		0,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil

}