package kissngoqueue

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
	cont := q.PutNMT(item)
	q.lock.Unlock()

	return cont

}

func (q *Queue) PutNMT(item any) (status bool) {
	var items []any

	select {

	case <-q.cancel:
		return false

	case items = <-q.items:
	case <-q.empty:

	}

	if checkCancel(q.cancel) {
		return false
	}

	items = append(items, item)
	q.items <- items

	return true
}

func (q *Queue) Cancel() {
	q.lock.Lock()
	q.CancelNMT()
	q.lock.Unlock()
}

func (q *Queue) CancelNMT() {
	close(q.cancel)
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

	if checkCancel(q.cancel) {
		return nil, false
	}

	return item, true
}

func checkCancel(c chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}
