package nbmutex

import "sync/atomic"

// Mutex is a non-blocking mutual exclusion lock. The zero value for a Mutex is
// an unlocked mutex. Do not copy this after first use or you'll have a bad
// time.
//
// It can be used to replace sync.Mutex for the occasional case where the mutex
// should only be *successfully* acquired by one goroutine, but if the mutex
// can't be acquired, the caller should not be blocked.
//
// This is useful, for example, when the critical section performs work that
// potentially blocks for an extended period of time, such as making a request
// across the network. Code that periodically submits metrics to a remote
// system would benefit from this, for example.
type Mutex struct {
	// cnt tracks the number of goroutines in the critical section, which should
	// only ever be 1.
	cnt int32
}

var emptyFunc = func() {}

// TryLock attempts to acquire the mutex.
//
// If the mutex is acquired, ok will be true and unlock will be a function that
// *must* be called to release the mutex.
//
// If the mutex is not acquired, ok will be false.
//
// It is safe to call unlock regardless of the return value of ok.
func (m *Mutex) TryLock() (unlock func(), ok bool) {
	// note that the named returns are just for "documentation"
	if !atomic.CompareAndSwapInt32(&m.cnt, 0, 1) {
		return emptyFunc, false
	}

	unlocker := func() {
		if !atomic.CompareAndSwapInt32(&m.cnt, 1, 0) {
			// if this happens something is very broken with the implementation
			panic("unlock detected inconsistency")
		}
	}
	return unlocker, true
}
