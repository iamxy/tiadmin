package pkg

type Event string

type EventStream interface {
	Next(stopc <-chan struct{}) chan Event
}
