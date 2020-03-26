package netconsole

import (
	"fmt"
	"net"

	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/lib/queue"
)

type Console struct {
	out    chan byte       // for 'display'
	InBuf  queue.QueueByte // for 'keyboard'
	logger *mylog.MyLog
}

func NewNetConsole(logger *mylog.MyLog) (*Console, error) {
	console := Console{make(chan byte, 200), queue.NewQueueByte(200), logger}
	go console.handle()
	return &console, nil
}

func (console *Console) Dump(address uint32) []byte {
	return nil // XXX - todo
}

func (console *Console) String() string {
	return "netconsole area"
}

func (console *Console) Shutdown() {
}

func (console *Console) Clear() { // Maybe Reset?
}

func (console *Console) Size() uint32 {
	return 0x100 // XXX: something
}

func (console *Console) Read(address uint32) byte {
	switch {
	case address == 0x000F8B: //  act like KEY_BUFFER_RPOS
		if console.InBuf.Len() > 0 {
			return 1
		} else {
			return 0
		}
	case address == 0x000F8D: //  act like KEY_BUFFER_WPOS
		return 0
	case address == 0x000F00:
		if console.InBuf.Len() > 0 {
			return *console.InBuf.Dequeue()
		} else {
			return 0
		}
	default:
		return 0
	}
}

func (console *Console) Write(address uint32, val byte) {
	switch {
	case address == 0x00EFE: // random, less conflicting with c256 addr
		//console.logger.Log(fmt.Sprintf("netconsole: out %s", val))
		console.out <- val
	default:
	}
}

// in goroutine
func (console *Console) handle() {
	var canWrite bool
	deadConn := make(chan net.Conn)
	for {
		console.logger.Log(fmt.Sprintf("netconsole: listening started"))
		l, err := net.Listen("tcp", ":12321")
		if err != nil {
			console.logger.Log(fmt.Sprintf("netconsole listen error: %s", err))
		}
		conn, err := l.Accept()
		if err != nil {
			console.logger.Log(fmt.Sprintf("netconsole accept error: %s", err))
		}
		canWrite = true

		go func() {
			p := make([]byte, 1)
			//console.logger.Log(fmt.Sprintf("netconsole: goroutine started"))
			for {
				size, err := conn.Read(p)
				if err != nil {
					deadConn <- conn
					break
				} else {
					if size > 0 {
						val := p[0]
						console.InBuf.Enqueue(val)
					}
				}
			}
			//console.logger.Log(fmt.Sprintf("netconsole: goroutine closed"))
		}()

		for canWrite {
			select {
			case deadOne := <-deadConn:
				_ = deadOne.Close()
				console.logger.Log(fmt.Sprintf("netconsole: client closed"))
				canWrite = false
			case val := <-console.out:
				conn.Write([]byte{val})
			}
		}
	}
}
