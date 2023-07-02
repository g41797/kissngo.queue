package kissngoqueue

import (
	"sync"
)

type Queue[T any] struct {
	lock   sync.Mutex
	items  chan []T
	empty  chan bool
	cancel chan struct{}
}

func NewQueue[T any]() *Queue[T] {
	qu := Queue[T]{
		items:  make(chan []T, 1),
		empty:  make(chan bool, 1),
		cancel: make(chan struct{}),
	}
	qu.empty <- true
	return &qu
}

func (q *Queue[T]) PutMT(item T) (status bool) {

	q.lock.Lock()
	cont := q.Put(item)
	q.lock.Unlock()

	return cont

}

func (q *Queue[T]) Put(item T) (status bool) {
	var items []T

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

func (q *Queue[T]) CancelMT() {
	q.lock.Lock()
	q.Cancel()
	q.lock.Unlock()
}

func (q *Queue[T]) Cancel() {
	close(q.cancel)
}

func (q *Queue[T]) Get() (item T, status bool) {
	var nilitem T
	var items []T
	select {
	case <-q.cancel:
		return nilitem, false
	case items = <-q.items:
	}
	item = items[0]
	if len(items) == 1 {
		q.empty <- true
	} else {
		q.items <- items[1:]
	}

	if checkCancel(q.cancel) {
		return nilitem, false
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
