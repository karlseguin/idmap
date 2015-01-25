package idmap

import (
	"hash/fnv"
	"sync"
)

type Map32 struct {
	bucketMask uint32
	jump       uint32
	lookup     map[uint32]*Bucket32
}

type Bucket32 struct {
	sync.RWMutex
	counter uint32
	lookup  map[string]uint32
}

// Create a new id mapper. For greater throughput, values are internally
// sharded into X buckets, where buckets must be a power of 2
func New32(buckets int) *Map32 {
	if buckets == 0 || ((buckets&(^buckets+1)) == buckets) == false {
		buckets = 16
	}

	b := uint32(buckets)
	m := &Map32{
		bucketMask: b - 1,
		jump:       uint32(buckets),
		lookup:     make(map[uint32]*Bucket32, buckets),
	}
	for i := uint32(0); i < b; i++ {
		m.lookup[i] = &Bucket32{
			counter: uint32(i + 1),
			lookup:  make(map[string]uint32),
		}
	}
	return m
}

// Get the id for the given string,
// optionally creating one if it doesn't exist
func (m *Map32) Get(s string, create bool) uint32 {
	bucket := m.getBucket(s)
	bucket.RLock()
	id, exists := bucket.lookup[s]
	bucket.RUnlock()
	if exists || create == false {
		return id
	}
	bucket.Lock()
	defer bucket.Unlock()

	id, exists = bucket.lookup[s]
	if exists {
		return id
	}

	id = bucket.counter
	bucket.lookup[s] = id
	bucket.counter += m.jump
	return id
}

// Remove the value from the map
func (m *Map32) Remove(s string) {
	bucket := m.getBucket(s)
	bucket.Lock()
	delete(bucket.lookup, s)
	bucket.Unlock()
}

func (m *Map32) getBucket(s string) *Bucket32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return m.lookup[h.Sum32()&m.bucketMask]
}
