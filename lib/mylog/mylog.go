
package mylog

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aniou/morfe/lib/queue"
)

type MyLog struct {
	LogOutput	 io.Writer

	logMsg           chan string
	logBuf   	 queue.QueueString   // for 'keyboard'

}

var Logger MyLog;

func init() {
	fmt.Println("logger is initialized")
	Logger = MyLog{LogOutput: os.Stdout, logMsg: make(chan string), logBuf: queue.NewQueueString(200)}
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
        //fmt.Printf("| %s\n", msg)
	fmt.Fprintf(l.LogOutput, "| %s\n", msg)
	//l.logBuf.Enqueue(msg)
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

