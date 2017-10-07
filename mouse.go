package goglwindow

import (
	"fmt"
)

type Watcher struct {
	Channel chan Event
}

func NewWatcher() *Watcher {
	return &Watcher{
		Channel: make(chan Event, 256),
	}
}

type Mouse struct {
	mouseEvents []Event
	events      chan Event
	watchers    []*Watcher
}

func NewMouse() *Mouse {
	return &Mouse{
		events: make(chan Event, 256),
	}
}

func (m *Mouse) AddWatcher(w *Watcher) {
	m.watchers = append(m.watchers, w)
}

func (m *Mouse) Receive() chan Event {
	return m.events
}

func (m *Mouse) Events() []Event {
	return m.mouseEvents
}

func (m *Mouse) Consume() {
	// TODO remove in the future
	m.mouseEvents = nil
	done := false
	for !done {
		select {
		case event, ok := <-m.events:
			if ok {
				m.mouseEvents = append(m.mouseEvents, event)
				for _, w := range m.watchers {
					w.Channel <- event
				}
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			done = true
		}
	}
}
