package pkg

import "time"

type Event string

// None event produced by a watcher timeout
func (e Event) None() bool {
	if string(e) == "" {
		return true
	}
	return false
}

type EventStream interface {
	Next(timeout time.Duration) chan Event
}
