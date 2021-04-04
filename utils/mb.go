package utils

import (
	"errors"
	"sync"
	"time"
)

// ErrClosed is returned when you add message to closed queue
var ErrClosed = errors.New("MB closed")

// ErrTooManyMessages means that adding more messages (at one call) than the limit
var ErrTooManyMessages = errors.New("Too many messages")

// NewMB returns a new MB with given queue size.
// size <= 0 means unlimited
func NewMB(size int) *MB {
	return &MB{
		cond: sync.NewCond(&sync.Mutex{}),
		size: size,
		read: make(chan struct{}),
	}
}

// MsgWithTimestamp with timestamp
type MsgWithTimestamp struct {
	msg     interface{}
	addedAt time.Time
}

// MB - message batching object
// Implements queue.
// Based on condition variables
type MB struct {
	msgs []MsgWithTimestamp

	cond *sync.Cond
	size int
	wait int
	read chan struct{}

	addCount, getCount         int64
	addMsgsCount, getMsgsCount int64
}

// WaitTimeoutOrMax it's Wait with limit of maximum returning array size or time out
// this method not supporting paused or closed
func (mb *MB) WaitTimeoutOrMax(timeout time.Duration, max int) (msgs []interface{}) {
	mb.cond.L.Lock()
	var timer *time.Timer = nil
	// check timeout
	isTimeout := false
	for {
		if len(mb.msgs) > 0 {
			leftDuration := mb.msgs[0].addedAt.Add(timeout).Sub(time.Now())
			if leftDuration > 0 {
				if timer == nil {
					timer = time.AfterFunc(leftDuration, func() {
						mb.cond.L.Lock()
						isTimeout = true
						mb.cond.Signal()
						mb.cond.L.Unlock()
					})
				}
			} else {
				// timeouted
				isTimeout = true
				break
			}
		}
		mb.cond.Wait()
		if isTimeout {
			// timeout
			break
		}
		if len(mb.msgs) >= max {
			break
		}
	}
	if timer != nil {
		timer.Stop()
		timer = nil
	}

	if len(mb.msgs) > max {
		msgs = dropTimestamp(mb.msgs[:max])
		mb.msgs = mb.msgs[max:]
	} else {
		msgs = dropTimestamp(mb.msgs)
		mb.msgs = make([]MsgWithTimestamp, 0)
	}
	mb.getCount++
	mb.getMsgsCount += int64(len(msgs))
	mb.unlockAdd()
	mb.cond.L.Unlock()
	return
}

// GetAll return all messages and flush queue
// Works on closed queue
func (mb *MB) GetAll() (msgs []interface{}) {
	mb.cond.L.Lock()
	msgs = dropTimestamp(mb.msgs)
	mb.msgs = make([]MsgWithTimestamp, 0)
	mb.getCount++
	mb.getMsgsCount += int64(len(msgs))
	mb.unlockAdd()
	mb.cond.L.Unlock()
	return
}

// Add - adds new messages to queue.
// When queue is closed - returning ErrClosed
// When count messages bigger then queue size - returning ErrTooManyMessages
// When the queue is full - wait until will free place
func (mb *MB) Add(msgs ...interface{}) (err error) {
add:
	mb.cond.L.Lock()
	if mb.size > 0 && len(mb.msgs)+len(msgs) > mb.size {
		if len(msgs) > mb.size {
			mb.cond.L.Unlock()
			return ErrTooManyMessages
		}
		// limit reached
		mb.wait++
		mb.cond.L.Unlock()
		<-mb.read
		goto add
	}
	for _, msg := range msgs {
		mb.msgs = append(mb.msgs, MsgWithTimestamp{
			msg:     msg,
			addedAt: time.Now(),
		})
	}
	mb.addCount++
	mb.addMsgsCount += int64(len(msgs))
	mb.cond.L.Unlock()
	mb.cond.Signal()
	return
}

func (mb *MB) unlockAdd() {
	if mb.wait > 0 {
		for i := 0; i < mb.wait; i++ {
			mb.read <- struct{}{}
		}
		mb.wait = 0
	}
}

// Len returning current size of queue
func (mb *MB) Len() (l int) {
	mb.cond.L.Lock()
	l = len(mb.msgs)
	mb.cond.L.Unlock()
	return
}

func dropTimestamp(msgs []MsgWithTimestamp) []interface{} {
	result := make([]interface{}, 0)
	for _, m := range msgs {
		result = append(result, m.msg)
	}
	return result
}
