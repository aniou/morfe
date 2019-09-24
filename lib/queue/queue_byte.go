
package queue

import (
	"sync"
)

type QueueByte struct {
	lock   *sync.Mutex
	Values []byte
	maxSize int
}

func NewQueueByte(length int) QueueByte {
	return QueueByte{&sync.Mutex{}, make([]byte, 0), length}
}

func (q *QueueByte) Len() int {
	return len(q.Values)
}

func (q *QueueByte) Enqueue(x byte) {
	for {
		if len(q.Values) < q.maxSize {
			q.lock.Lock()
			q.Values = append(q.Values, x)
			q.lock.Unlock()
		}
		return
	}
}

func (q *QueueByte) Dequeue() *byte {
	for {
		if len(q.Values) > 0 {
			q.lock.Lock()
			x := q.Values[0]
			q.Values = q.Values[1:]
			q.lock.Unlock()
			return &x
		}
		return nil
	}
}
