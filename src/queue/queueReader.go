package queue

import (
	"container/list"
)

type QueueReader struct {
	head *list.Element
	queue *Queue
}

func NewQueueReader(q *Queue) (QueueReader, bool) {
	qr := QueueReader{
		head:  q.Peek(),
		queue: q,
	}
	if qr.head == nil {
		return qr, false
	}
	return qr, true
}

func (qr *QueueReader) Next() *list.Element{
	qr.queue.lock.Lock()
	defer qr.queue.lock.Unlock()
	if qr.head == nil {
		//head is nil, peek the head to store as head
		qr.head = qr.queue.Peek()
	} else {
		//get the next element if present and store as head
		qr.head = qr.head.Next()
	}
	return qr.head
}