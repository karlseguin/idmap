package idmap

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type Map32Tests struct{}

func Test_Map32(t *testing.T) {
	Expectify(new(Map32Tests), t)
}

func (_ Map32Tests) ReturnsANewId() {
	m := New(2)
	Expect(m.Get("over", true)).To.Equal(uint64(2))
	Expect(m.Get("9000", true)).To.Equal(uint64(1))
}

func (_ Map32Tests) ReturnsAnExistingId() {
	m := New(2)
	m.Get("over", true)
	m.Get("9000", true)
	Expect(m.Get("over", false)).To.Equal(uint64(2))
}

func (_ Map32Tests) GeneratesUniqueIdsAccrossBuckets() {
	m := New(4)
	Expect(m.Get("a", true)).To.Equal(uint64(1))
	Expect(m.Get("b", true)).To.Equal(uint64(2))
	Expect(m.Get("c", true)).To.Equal(uint64(3))
	Expect(m.Get("d", true)).To.Equal(uint64(4))
	Expect(m.Get("e", true)).To.Equal(uint64(5))
	Expect(m.Get("i", true)).To.Equal(uint64(9))
	Expect(m.Get("m", true)).To.Equal(uint64(13))
	Expect(m.Get("f", true)).To.Equal(uint64(6))
}

func (_ Map32Tests) DoesNotCreateANewId() {
	m := New(2)
	Expect(m.Get("over", false)).To.Equal(uint64(0))
}

func (_ Map32Tests) RemovesAnId() {
	m := New(2)
	m.Get("over", true)
	m.Remove("over")
	Expect(m.Get("over", false)).To.Equal(uint64(0))
}
