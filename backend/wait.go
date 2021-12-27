package main

import (
	"context"
	"sync"
)

type Wait struct {
	C chan struct{}
	sync.Mutex
}

func (w *Wait) Wait(ctx context.Context) error {
	w.Lock()
	c := w.C
	if c == nil {
		c = make(chan struct{})
		w.C = c
	}
	w.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return nil
	}
}

func (w *Wait) Notify() {
	w.Lock()
	defer w.Unlock()

	if w.C == nil {
		return // no one is waiting
	}

	close(w.C)
	w.C = nil
}
