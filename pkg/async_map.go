package pkg

import "sync"

type AsyncMap[T, U comparable] struct {
	mx *sync.Mutex
	m  map[T]U
}

func NewAsyncMap[T, U comparable]() *AsyncMap[T, U] {
	return &AsyncMap[T, U]{
		mx: new(sync.Mutex),
		m:  make(map[T]U),
	}
}

func (m *AsyncMap[T, U]) Get(k T) (U, bool) {
	m.mx.Lock()
	v, ok := m.m[k]
	m.mx.Unlock()
	return v, ok
}

func (m *AsyncMap[T, U]) Set(k T, v U) {
	m.mx.Lock()
	m.m[k] = v
	m.mx.Unlock()
	return
}
