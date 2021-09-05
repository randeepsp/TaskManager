package queue

import (
	"container/list"
	"errors"
	"sync"
)

//struct to encompass a list with mutex to make concurrent safe
type Queue struct {
	lock sync.Mutex
	list *list.List
}

//returns a new queue
func NewQueue() *Queue {
	list := list.New()
	return &Queue{list: list}
}

//returns the size of the queue
func(q *Queue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.list.Len()
}

//adds element add the end of the queue
func (q *Queue) Enqueue (val interface{}) *list.Element{
	q.lock.Lock()
	defer q.lock.Unlock()
	e := q.list.PushBack(val)
	return e
}

//removes and returns the element, waits if there is no element in the queue
func (q *Queue) Dequeue() *list.Element {
	//acquire lock before removing from list
	q.lock.Lock()
	defer q.lock.Unlock()
	e := q.list.Front()
	q.list.Remove(e)
	return e
}

//removes the element, returns err if no element found
func (q *Queue) RemoveElement(element *list.Element) error {
	//acquire lock before removing from list
	q.lock.Lock()
	defer q.lock.Unlock()
	e := q.list.Remove(element)
	if e != nil {
		return  nil
	}
	return errors.New("no such element")
}

//Return the element in the queue without removing it
func (q *Queue) Peek() *list.Element {
	e := q.list.Front()
	return e
}

