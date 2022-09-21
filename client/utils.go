package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

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

// Utility method to especially print structs in human readable JSON format
func pretty(obj interface{}) {
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}
