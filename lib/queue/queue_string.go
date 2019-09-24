
package queue

import (
	"sync"
)

type QueueString struct {
	lock   *sync.Mutex
	Values []string
	maxSize int
}

func NewQueueString(length int) QueueString {
	return QueueString{&sync.Mutex{}, make([]string, 0), length}
}

func (q *QueueString) Len() int {
	return len(q.Values)
}

func (q *QueueString) Enqueue(x string) {
	for {
		if len(q.Values) < q.maxSize {
			q.lock.Lock()
			q.Values = append(q.Values, x)
			q.lock.Unlock()
		}
		return
	}
}

func (q *QueueString) Dequeue() *string {
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
