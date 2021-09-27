
package queue

import (
	_ "fmt"
	"sync"
)

type QueueByte struct {
	lock    *sync.Mutex
	Values  []byte
	index   int
}

func NewQueueByte(length int) QueueByte {
	return QueueByte{&sync.Mutex{}, make([]byte, length), 0}
}

func (q *QueueByte) Len() int {
	return len(q.Values)
}

func (q *QueueByte) Enqueue(x byte) {
	if q.index == len(q.Values) {
		q.index = 0
	}

	q.lock.Lock()
	q.Values[q.index] = x
	q.lock.Unlock()
	q.index = q.index + 1
	//fmt.Printf("queue byte: %d %v\n", q.index, q.Values)
	return
}

func (q *QueueByte) Dequeue() byte {
	if q.index == 0 {
		q.index = len(q.Values)
	}

	q.index = q.index - 1
	q.lock.Lock()
	x := q.Values[q.index]
	q.lock.Unlock()
	//fmt.Printf("dequeue byte: %d %v\n", q.index, q.Values)
	return x
}
