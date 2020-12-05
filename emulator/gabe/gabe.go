package gabe

import (
	"fmt"

	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/lib/queue"
)

type Gabe struct {
	//out    chan byte       // for 'display'
	InBuf  queue.QueueByte // for 'keyboard'
	Data    byte
	command byte
}

func New() (*Gabe, error) {
	//console := Gabe{make(chan byte, 200), queue.NewQueueByte(200)}
	gabe := Gabe{queue.NewQueueByte(200), 0, 0}
	return &gabe, nil
}

func (console *Gabe) Dump(address uint32) []byte {
	return nil // XXX - todo
}

func (console *Gabe) String() string {
	return "Gabe area"
}

func (console *Gabe) Shutdown() {
}

func (console *Gabe) Clear() { // Maybe Reset?
}

func (console *Gabe) Size() uint32 {
	return 0x100 // XXX: something
}

func (console *Gabe) Read(address uint32) byte {
	//mylog.Logger.Log(fmt.Sprintf("."))
	switch {
	case address == 0xAF1060:
		return console.Data
		if console.InBuf.Len() > 0 {
			return *console.InBuf.Dequeue()
		} else {
			return 0
		}
	case address == 0xAF1064:			// we support only bit 0 
		return console.command
		/*
		if console.InBuf.Len() > 0 {
			mylog.Logger.Log(fmt.Sprintf("gabe: read from addr %6X 1 returned", address))
			return 2
		} else {
			return 0
		}
		*/
	default:
		mylog.Logger.Log(fmt.Sprintf("gabe: read from addr %6X is not implemented, 0 returned", address))
		return 0
	}
}
// taken from FoenixIDE
func (console *Gabe) Write(address uint32, val byte) {
	switch {
	case address == 0xAF1060:
		if val == 0x69 {				// 
			console.command = 1		// 
		}
		if val == 0xEE {				// echo
			console.command = 1		// 
		}
		if val == 0xF4 {				// kbd reset
			console.InBuf.Enqueue(0xFA)		// self-test result
			console.Data = 0xFA
			console.command = 1
		}
		if val == 0xF6 {				// 
			console.command = 1
		}
	case address == 0xAF1064:
		if val == 0x20 {				// 
			console.command = 1
		}
		if val == 0x60 {				// 
			console.command = 1
		}
		if val == 0xA8 {				// 
			console.command = 1
		}
		if val == 0xA9 {				// 
			console.InBuf.Enqueue(0x00)		// 
			console.Data = 0
			console.command = 0x01
		}
		if val == 0xAA {				// self-test
			console.InBuf.Enqueue(0x55)		// self-test result
			console.Data = 0x55
			console.command = 0x01
		}
		if val == 0xAB {				// self-test
			console.InBuf.Enqueue(0x00)		// 
			console.Data = 0
		}
		if val == 0xD4 {				// 
			console.command = 0x01
		}
		if val == 0x00 {
			console.command = 0x00
		}
	default:
		mylog.Logger.Log(fmt.Sprintf("gabe: write to addr %6X val %2X is not implemented, 0 returned", address, val))
	}
	return
}
