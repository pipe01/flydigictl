package utils

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// From https://gist.github.com/zviadm/c234426882bfc8acba88f3503edaaa36

// Conditional variable implementation that uses channels for notifications.
// Only supports .Broadcast() method, however supports timeout based Wait() calls
// unlike regular sync.Cond.
type Cond struct {
	L sync.Locker
	n unsafe.Pointer
}

func NewCond(l sync.Locker) *Cond {
	c := &Cond{L: l}
	n := make(chan struct{})
	c.n = unsafe.Pointer(&n)
	return c
}

// Waits for Broadcast calls. Similar to regular sync.Cond, this unlocks the underlying
// locker first, waits on changes and re-locks it before returning.
func (c *Cond) Wait() {
	n := c.NotifyChan()
	c.L.Unlock()
	<-n
	c.L.Lock()
}

// Same as Wait() call, but will only wait up to a given timeout.
func (c *Cond) WaitWithTimeout(t time.Duration) {
	n := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-n:
	case <-time.After(t):
	}
	c.L.Lock()
}

// Returns a channel that can be used to wait for next Broadcast() call.
func (c *Cond) NotifyChan() <-chan struct{} {
	ptr := atomic.LoadPointer(&c.n)
	return *((*chan struct{})(ptr))
}

// Broadcast call notifies everyone that something has changed.
func (c *Cond) Broadcast() {
	n := make(chan struct{})
	ptrOld := atomic.SwapPointer(&c.n, unsafe.Pointer(&n))
	close(*(*chan struct{})(ptrOld))
}

type CondValue[T any] struct {
	*Cond
	Value *T
}

func NewCondValue[T any](l sync.Locker) *CondValue[T] {
	return &CondValue[T]{
		Cond: NewCond(l),
	}
}
