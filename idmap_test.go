package idmap

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestIdMapReturnsANewId(t *testing.T) {
	spec := gspec.New(t)
	m := New(2)
	spec.Expect(m.Get("over", true)).ToEqual(uint64(1))
	spec.Expect(m.Get("9000", true)).ToEqual(uint64(2))
}

func TestIdMapReturnsAnExistingId(t *testing.T) {
	spec := gspec.New(t)
	m := New(2)
	m.Get("over", true)
	m.Get("9000", true)
	spec.Expect(m.Get("over", false)).ToEqual(uint64(1))
}

func TestIdMapDoesNotCreateANewId(t *testing.T) {
	spec := gspec.New(t)
	m := New(2)
	spec.Expect(m.Get("over", false)).ToEqual(uint64(0))
}

func TestIdMapRemovesAnId(t *testing.T) {
	spec := gspec.New(t)
	m := New(2)
	m.Get("over", true)
	m.Remove("over")
	spec.Expect(m.Get("over", false)).ToEqual(uint64(0))
}
