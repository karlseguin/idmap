package idmap

import (
	"hash/fnv"
	"sync"
)

type Map struct {
	bucketMask uint32
	jump       uint64
	lookup     map[uint32]*Bucket
}

type Bucket struct {
	sync.RWMutex
	counter uint64
	lookup  map[string]uint64
}

// Create a new id mapper. For greater throughput, values are internally
// sharded into X buckets, where buckets must be a power of 2
func New(buckets int) *Map {
	if buckets == 0 || ((buckets&(^buckets+1)) == buckets) == false {
		buckets = 16
	}

	b := uint32(buckets)
	m := &Map{
		bucketMask: b - 1,
		jump:       uint64(buckets),
		lookup:     make(map[uint32]*Bucket, buckets),
	}
	for i := uint32(0); i < b; i++ {
		m.lookup[i] = &Bucket{
			counter: uint64(i + 1),
			lookup:  make(map[string]uint64),
		}
	}
	return m
}

// Get the id for the given string,
// optionally creating one if it doesn't exist
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

	id = bucket.counter
	bucket.lookup[s] = id
	bucket.counter += m.jump
	return id
}

// Remove the value from the map
func (m *Map) Remove(s string) {
	bucket := m.getBucket(s)
	bucket.Lock()
	delete(bucket.lookup, s)
	bucket.Unlock()
}

func (m *Map) getBucket(s string) *Bucket {
	h := fnv.New32a()
	h.Write([]byte(s))
	return m.lookup[h.Sum32()&m.bucketMask]
}
