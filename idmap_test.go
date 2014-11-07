package idmap

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type MapTests struct {}

func Test_Map(t *testing.T) {
	Expectify(new(MapTests), t)
}

func (_ *MapTests) TestIdMapReturnsANewId() {
	m := New(2)
	Expect(m.Get("over", true)).To.Equal(uint64(1))
	Expect(m.Get("9000", true)).To.Equal(uint64(2))
}

func (_ *MapTests) TestIdMapReturnsAnExistingId() {
	m := New(2)
	m.Get("over", true)
	m.Get("9000", true)
	Expect(m.Get("over", false)).To.Equal(uint64(1))
}

func (_ *MapTests) TestIdMapDoesNotCreateANewId() {
	m := New(2)
	Expect(m.Get("over", false)).To.Equal(uint64(0))
}

func (_ *MapTests) TestIdMapRemovesAnId() {
	m := New(2)
	m.Get("over", true)
	m.Remove("over")
	Expect(m.Get("over", false)).To.Equal(uint64(0))
}
