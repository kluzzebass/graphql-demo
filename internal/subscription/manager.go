package subscription

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Manager[T any] struct {

	// mutex for map access
	mutex sync.Mutex

	// slice of subscribers
	Subscribers map[uuid.UUID]*Subscription[T]
	Queue       []T
}

func NewManager[T any]() *Manager[T] {
	m := &Manager[T]{
		Subscribers: make(map[uuid.UUID]*Subscription[T]),
	}

	// start goroutines
	go m.distributor()
	go m.sendor()
	go m.evictor()

	return m
}

// add a new subscriber
func (m *Manager[T]) Subscribe(ctx context.Context, ff func(msg T) bool) *Subscription[T] {
	sub := NewSubscription(ctx, ff)

	m.mutex.Lock()
	m.Subscribers[uuid.New()] = sub
	m.mutex.Unlock()

	return sub
}

// send messages to subscribers
func (m *Manager[T]) Publish(msgs []T) {

	m.mutex.Lock()
	m.Queue = append(m.Queue, msgs...)
	m.mutex.Unlock()

}

// periodically distribute messages to subscribers
func (m *Manager[T]) distributor() {
	for {
		if len(m.Queue) > 0 {
			// pop a message from the queue
			m.mutex.Lock()
			msg := m.Queue[0]
			m.Queue = m.Queue[1:]
			m.mutex.Unlock()

			// push message onto subscriber queue
			for _, sub := range m.Subscribers {

				// filter messages based on subscription filter function
				if sub.Filter(msg) {
					sub.Push(&msg)
				}
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// periodically check for subscribers that are done
func (m *Manager[T]) evictor() {
	for {
		var deleteTheseKeys []uuid.UUID
		m.mutex.Lock()
		for key, sub := range m.Subscribers {
			select {
			case <-sub.ctx.Done():
				deleteTheseKeys = append(deleteTheseKeys, key)
			default:
			}
		}

		for _, key := range deleteTheseKeys {
			m.Subscribers[key].CloseChan()
			delete(m.Subscribers, key)
		}
		m.mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (m *Manager[T]) sendor() {
	for {
		sleep := true
		m.mutex.Lock()
		for _, sub := range m.Subscribers {
			// pop message from subscriber queue
			msg, ok := sub.Pop()
			if !ok {
				continue
			}
			// send message to subscriber
			select {
			case sub.Chan() <- *msg:
				sleep = false
			default:
				// if channel is full, push message back onto queue
				sub.Unpop(msg)
			}
		}
		m.mutex.Unlock()
		if sleep {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
