package client

import "sync"

type idGen struct {
	sync.Mutex
	value uint
}

func (a *idGen) inc() uint {
	a.Lock()
	defer a.Unlock()

	a.value = a.value + 1
	return a.value
}
