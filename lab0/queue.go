package lab0

import (
	"sync"
)

// Queue is a simple FIFO queue that is unbounded in size.
// Push may be called any number of times and is not
// expected to fail or overwrite existing entries.
type Queue[T any] struct {
	contents []T
}

// NewQueue returns a new queue which is empty.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		contents: make([]T, 0),
	}
}

// Push adds an item to the end of the queue.
func (q *Queue[T]) Push(t T) {
	q.contents = append(q.contents, t)
}

// Pop removes an item from the beginning of the queue
// and returns it unless the queue is empty.
//
// If the queue is empty, returns the zero value for T and false.
//
// If you are unfamiliar with "zero values", consider revisiting
// this section of the Tour of Go: https://go.dev/tour/basics/12
func (q *Queue[T]) Pop() (T, bool) {
	if len(q.contents) == 0 {
		var dflt T
		return dflt, false
	}
	first_element := q.contents[0]
	q.contents = q.contents[1:]
	return first_element, true
}

// ConcurrentQueue provides the same semantics as Queue but
// is safe to access from many goroutines at once.
//
// You can use your implementation of Queue[T] here.
//
// If you are stuck, consider revisiting this section of
// the Tour of Go: https://go.dev/tour/concurrency/9
type ConcurrentQueue[T any] struct {
	contents []T
	lock     sync.Mutex
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		contents: make([]T, 0),
	}

}

// Push adds an item to the end of the queue
func (q *ConcurrentQueue[T]) Push(t T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.contents = append(q.contents, t)

}

// Pop removes an item from the beginning of the queue.
// Returns a zero value and false if empty.
func (q *ConcurrentQueue[T]) Pop() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.contents) == 0 {
		var dflt T
		return dflt, false
	}


	first_element := q.contents[0]
	q.contents = q.contents[1:]
	
	return first_element, true
}
