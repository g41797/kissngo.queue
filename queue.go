package queue

import (
	"sync"
)

type Queue struct {
	lock   sync.Mutex
	items  chan []any
	empty  chan bool
	cancel chan struct{}
}

func NewQueue() *Queue {
	qu := Queue{
		items:  make(chan []any, 1),
		empty:  make(chan bool, 1),
		cancel: make(chan struct{}),
	}
	qu.empty <- true
	return &qu
}

func (q *Queue) Put(item any) (status bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	var items []any

	select {

	case <-q.cancel:
		return false

	case items = <-q.items:
	case <-q.empty:

	}

	items = append(items, item)
	q.items <- items

	return true
}

func (q *Queue) Cancel() {
	q.lock.Lock()
	close(q.cancel)
	q.lock.Unlock()
}

func (q *Queue) Get() (item any, status bool) {

	var items []any
	select {
	case <-q.cancel:
		return nil, false
	case items = <-q.items:
	}
	item = items[0]
	if len(items) == 1 {
		q.empty <- true
	} else {
		q.items <- items[1:]
	}
	return item, true
}
