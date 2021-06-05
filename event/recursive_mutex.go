package event

import (
	"sync"
	"time"

	tls "github.com/huandu/go-tls"
)

type RecursiveMutex struct {
	mutex            sync.Mutex
	internalMutex    sync.Mutex
	currentGoRoutine int64
	lockCount        uint64
}

func (rm *RecursiveMutex) Lock() {
	goRoutineID := tls.ID()

	for {
		rm.internalMutex.Lock()
		if rm.currentGoRoutine == 0 {
			rm.currentGoRoutine = goRoutineID
			break
		} else if rm.currentGoRoutine == goRoutineID {
			break
		} else {
			rm.internalMutex.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}
	}
	rm.lockCount++
	rm.internalMutex.Unlock()
}

func (rm *RecursiveMutex) Unlock() {
	rm.internalMutex.Lock()
	rm.lockCount--
	if rm.lockCount == 0 {
		rm.currentGoRoutine = 0
	}
	rm.internalMutex.Unlock()
}
