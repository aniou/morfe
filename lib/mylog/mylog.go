
package mylog

import (
	"github.com/aniou/go65c816/lib/queue"
)

type MyLog struct {
	logMsg           chan string
	logBuf   	 queue.QueueString   // for 'keyboard'

}

func New() *MyLog {
	return &MyLog{logMsg: make(chan string), logBuf: queue.NewQueueString(200)}
}

func (l *MyLog) GetChannel() (chan string) {
	return l.logMsg
}

func (l *MyLog) Log(msg string) {
	//go func(message string) {
	//	l.logMsg<-message
	//}(msg)
	//l.logMsg<-msg
	l.logBuf.Enqueue(msg)
}

func (l *MyLog) Len() int {
	return l.logBuf.Len()
}

func (l *MyLog) Dequeue() *string {
	return l.logBuf.Dequeue()
}

