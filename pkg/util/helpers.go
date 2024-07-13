package util

import "sync"

// generates a thread safe map to update the controls on the operation view

type MutexMap struct {
	sync.Mutex
	mmap map[string]bool
}

func NewEmptyMutexMap() *MutexMap {
	MutexMap := MutexMap{
		mmap: make(map[string]bool),
	}

	return &MutexMap
}

func NewMutexMap(initialValues map[string]bool) *MutexMap {
	return &MutexMap{
		mmap: initialValues,
	}
}

func (m *MutexMap) Write(key string, value bool) {
	m.Lock()
	defer m.Unlock()
	m.mmap[key] = value
}

func (m *MutexMap) Get(key string) bool {
	m.Lock()
	defer m.Unlock()
	return m.mmap[key]
}
