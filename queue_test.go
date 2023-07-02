package kissngoqueue

import (
	"testing"
	"time"
)

func TestPutGet(t *testing.T) {

	q := NewQueue[int]()

	ln := 2

	ni := 10
	st := stm()
	resch := make(chan int, ln)

	for l := 0; l < ln; l++ {
		go putItems(q, ni, st, resch, q.Put)
	}

	var got int
	for l := 0; l < ln; l++ {
		got += <-resch
	}

	want := ni * ln

	if want != got {
		t.Errorf("Expected '%d', but got '%d'", want, got)
		return
	}

	for i := 1; i <= got; i++ {
		_, ok := q.Get()
		if !ok {
			t.Errorf("Expected OK, but was cancelled")
			return
		}
	}

}

func TestPutGetCancel(t *testing.T) {

	q := NewQueue[int]()

	ln := 2
	ni := 10
	st := stm()
	resch := make(chan int, ln)

	for l := 0; l < ln; l++ {
		go putItems(q, ni, st, resch, q.Put)
	}

	time.Sleep(stm())

	q.Cancel()

	var got int
	for l := 0; l < ln; l++ {
		got += <-resch
	}

	want := ni * ln

	if got == want {
		t.Errorf("want: Cancel doesn't work")
	}

	_, ok := q.Get()
	if ok {
		t.Errorf("Get: Cancel doesn't work")
	}

}

func putItems(q *Queue[int], ni int, st time.Duration, pn chan<- int, put func(int) bool) {
	var i int
	for i = 1; i <= ni; i++ {
		mp := i
		if !put(mp) {
			pn <- i - 1
			return
		}
		time.Sleep(st)
	}
	pn <- i - 1
}

func TestPutGetCancelMT(t *testing.T) {

	q := NewQueue[int]()

	ln := 1
	ni := 10
	st := stm()
	resch := make(chan int, ln)

	for l := 0; l < ln; l++ {
		go putItems(q, ni, st, resch, q.PutMT)
	}

	time.Sleep(stm())

	q.CancelMT()

	var got int
	for l := 0; l < ln; l++ {
		got += <-resch
	}

	want := ni * ln

	if got == want {
		t.Errorf("want: Cancel doesn't work")
	}

	_, ok := q.Get()
	if ok {
		t.Errorf("Get: Cancel doesn't work")
	}

}

func stm() time.Duration {
	return time.Second / 10
}
