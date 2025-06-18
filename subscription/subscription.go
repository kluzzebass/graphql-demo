package subscription

import (
	"context"
	"sync"
)

type Subscription[T any] struct {
	out   chan T
	queue []*T
	mutex sync.Mutex

	ctx        context.Context
	filterFunc func(msg T) bool
}

func NewSubscription[T any](ctx context.Context, ff func(msg T) bool) *Subscription[T] {
	sub := &Subscription[T]{
		out:        make(chan T),
		ctx:        ctx,
		filterFunc: ff,
	}

	return sub
}

// close channel
func (s *Subscription[T]) CloseChan() {
	close(s.out)
}

// return channel for subscription
func (s *Subscription[T]) Chan() chan T {
	return s.out
}

// push message onto end of queue
func (s *Subscription[T]) Push(msg *T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.queue = append(s.queue, msg)
}

// pop message from front of queue
func (s *Subscription[T]) Pop() (*T, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.queue) == 0 {
		return nil, false
	}
	msg := s.queue[0]
	s.queue = s.queue[1:]
	return msg, true
}

// unpop message onto front of queue
func (s *Subscription[T]) Unpop(msg *T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.queue = append([]*T{msg}, s.queue...)
}

// filter message using filter function
func (s *Subscription[T]) Filter(msg T) bool {
	// if filter function is defined, use it
	if s.filterFunc != nil {
		return (s.filterFunc)(msg)
	}
	// otherwise, let all messages through
	return true
}
