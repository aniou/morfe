
package mylog

import (
	"fmt"
	"github.com/aniou/go65c816/lib/queue"
	"time"
)

type MyLog struct {
	logMsg           chan string
	logBuf   	 queue.QueueString   // for 'keyboard'

}

var Logger MyLog;

func init() {
	fmt.Println("logger is initialized")
	Logger = MyLog{logMsg: make(chan string), logBuf: queue.NewQueueString(200)}
}


//func New() *MyLog {
//	return &MyLog{logMsg: make(chan string), logBuf: queue.NewQueueString(200)}
//}

func (l *MyLog) GetChannel() (chan string) {
	return l.logMsg
}

func (l *MyLog) Log(msg string) {
	//go func(message string) {
	//	l.logMsg<-message
	//}(msg)
	//l.logMsg<-msg
        fmt.Printf("| %s\n", msg)
	l.logBuf.Enqueue(msg)
}

func (l *MyLog) Len() int {
	return l.logBuf.Len()
}

func (l *MyLog) Dequeue() *string {
	return l.logBuf.Dequeue()
}

func (l *MyLog) ConsolePrinter() {
        go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if l.logBuf.Len() > 0 {
				fmt.Printf("> %s\n", *l.logBuf.Dequeue())
			}
		}
        }()
}

