package main

import "sync"

// PlayQueue is a FIFO queue for holding the URIs of songs to play
type PlayQueue struct {
	items []string
	lock  sync.Mutex
}

// New creates a new playqueue
func (q *PlayQueue) New() *PlayQueue {
	q.items = []string{}
	return q
}

// Enqueue adds a new item to the back of a queue.
func (q *PlayQueue) Enqueue(s string) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, s)
	return len(q.items)
}

// Pop retrieves the next (first) item of the queue
func (q *PlayQueue) Pop() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	ret := q.items[0]
	q.items = q.items[1:]

	return ret
}

// Peek returns the next item, but does not dequeue it.
func (q *PlayQueue) Peek() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.items[0]
}

// Size returns the current size (depth) of the queue
func (q *PlayQueue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items)
}
