package framework

import (
	"gioui.org/io/event"
	"gioui.org/io/input"
)

type Listener interface {
	Interest() event.Filter
	OnEvent(event.Event)
}

type Watcher struct {
	list []Listener
}

func (w *Watcher) Add(l Listener) {
	w.list = append(w.list, l)
}

func (w *Watcher) Remove(l Listener) {
	for idx, ol := range w.list {
		if ol == l {
			w.list = append(w.list[:idx], w.list[idx+1:]...)
			break
		}
	}
}

func (w *Watcher) Trigger(q input.Source) {
	for _, l := range w.list {
		for {
			ev, ok := q.Event(l.Interest())
			if ok {
				l.OnEvent(ev)
			} else {
				break
			}
		}

	}
}
