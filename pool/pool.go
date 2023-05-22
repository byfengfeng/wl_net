package pool

import (
	"sync"
)

type IPool interface {
	Get() any
	Put(x any)
}

func NewPool(f func() any) *sync.Pool {
	return &sync.Pool{New: f}
}
