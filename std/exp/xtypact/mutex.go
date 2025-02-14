package xtypact

import "sync"

// NewMutex returns [Mutex] with the given value.
func NewMutex[T any](value T) Mutex[T] {
	return Mutex[T]{
		value: value,
	}
}

// Mutex is a value holding mutex using [sync.Mutex].
type Mutex[T any] struct {
	lock  sync.Mutex
	value T
}

// Lock acquires a lock on the mutex and returns the value.
//
// See [sync.Mutex] for more details.
func (m *Mutex[T]) Lock() T {
	m.lock.Lock()

	return m.value
}

// Unlock releases the lock on the mutex.
//
// See [sync.Mutex] for more details.
func (m *Mutex[T]) Unlock() {
	m.lock.Unlock()
}

// TryLock tries to lock the mutex and reports the whether it succeeded.
//
// See [sync.Mutex] for more details.
func (m *Mutex[T]) TryLock() (T, bool) {
	if m.lock.TryLock() {
		return m.value, true
	}

	var zero T

	return zero, false
}

// UnsafeSet updates the internal value to val.
//
// WARN: It is the responsibility of the caller to acquire a lock first!
func (m *Mutex[T]) UnsafeSet(val T) {
	m.value = val
}

// Set acquires a lock and updates the internal value to val.
func (m *Mutex[T]) Set(val T) {
	m.lock.Lock()
	m.value = val
	m.lock.Unlock()
}

// WithLock acquires a lock, calls fn with the value and updates
// the internal value with the one returned by fn.
// The lock is released once this method returns.
func (m *Mutex[T]) WithLock(fn func(T) T) {
	m.lock.Lock()
	m.value = fn(m.value)
	m.lock.Unlock()
}
