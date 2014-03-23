package idmap

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
)

type Map struct {
	buckets uint32
	counter uint64
	lookup  map[uint32]*Bucket
}

type Bucket struct {
	sync.RWMutex
	lookup map[string]uint64
}

func New(buckets int) *Map {
	b := uint32(buckets)
	m := &Map{
		buckets: b,
		lookup:  make(map[uint32]*Bucket, buckets),
	}
	for i := uint32(0); i < b; i++ {
		m.lookup[i] = &Bucket{
			lookup: make(map[string]uint64),
		}
	}
	return m
}

func (m *Map) Get(s string, create bool) uint64 {
	bucket := m.getBucket(s)
	bucket.RLock()
	id, exists := bucket.lookup[s]
	bucket.RUnlock()
	if exists {
		return id
	}
	if create == false {
		return 0
	}

	bucket.Lock()
	defer bucket.Unlock()

	id, exists = bucket.lookup[s]
	if exists {
		return id
	}

	id = atomic.AddUint64(&m.counter, 1)
	bucket.lookup[s] = id
	return id
}

func (m *Map) Remove(s string) {
	bucket := m.getBucket(s)
	bucket.Lock()
	delete(bucket.lookup, s)
	bucket.Unlock()
}

func (m *Map) getBucket(s string) *Bucket {
	h := fnv.New32a()
	h.Write([]byte(s))
	index := h.Sum32() % m.buckets
	return m.lookup[index]
}
