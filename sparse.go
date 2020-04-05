package bitset

//TODO: use sync.Pool?

import (
	"reflect"
	"strconv"
	"strings"
)

// Sparse is a sparse bitset. It can represent any uint64 and uses memory
// proportional to the number of elements it contains.
type Sparse struct {
	root *node // compact radix tree 7 levels deep.
}

// NewSparse creates a new Sparse bitset.
func NewSparse() *Sparse {
	return &Sparse{}
}

func (s *Sparse) init() {
	s.root = &node{shift: 64 - 8}
}

// Add adds n to s.
func (s *Sparse) Add(n uint) { s.Add64(uint64(n)) }

// Remove removes n from s.
func (s *Sparse) Remove(n uint) { s.Remove64(uint64(n)) }

// Contains reports whether s contains s.
func (s *Sparse) Contains(n uint) bool { return s.Contains64(uint64(n)) }

// Add64 adds n to s.
func (s *Sparse) Add64(n uint64) {
	if s.root == nil {
		s.init()
	}
	s.root.add64(n)
}

// Remove64 removes n from s.
func (s *Sparse) Remove64(n uint64) {
	if s.root == nil {
		return
	}
	if s.root.remove64(uint64(n)) {
		s.root = nil
	}
}

// Contains64 reports whether s contains s.
func (s *Sparse) Contains64(n uint64) bool {
	if s.root == nil {
		return false
	}
	return s.root.contains64(n)
}

// Empty reports whether s has no elements.
func (s *Sparse) Empty() bool {
	return s.root == nil
}

// Clear removes all elements from s.
func (s *Sparse) Clear() {
	s.root = nil
}

// Equal reports whether two bitsets have the same elements.
func (s1 *Sparse) Equal(s2 *Sparse) bool {
	if s1.root == nil || s2.root == nil {
		return s1.root == s2.root
	}
	return s1.root.equal(s2.root)
}

// Copy returns a copy of s.
func (s *Sparse) Copy() *Sparse {
	c := NewSparse()
	s.root = s.root.copyNode()
	return c
}

// Len returns the number of elements in s.
func (s *Sparse) Len() int {
	if s.root == nil {
		return 0
	}
	return s.root.len()
}

// AddIn adds all the elements in s2 to s1.
// It sets s1 to the union of s1 and s2.
func (s1 *Sparse) AddIn(s2 *Sparse) {
	if s2.Empty() {
		return
	}
	if s1.Empty() {
		s1.init()
	}
	s1.root.addIn(s2.root)
}

// RemoveIn removes from s1 all the elements that are in s2.
// It sets s1 to the set difference of s1 and s2.
func (s1 *Sparse) RemoveIn(s2 *Sparse) {
	if s1.Empty() || s2.Empty() {
		return
	}
	if s1.root.removeIn(s2.root) {
		s1.root = nil
	}
}

// RemoveNotIn removes from s1 all the elements that are not in s2.
// It sets s1 to the intersection of s1 and s2.
func (s1 *Sparse) RemoveNotIn(s2 *Sparse) {
	if s1.Empty() {
		return
	}
	if s2.Empty() {
		s1.Clear()
		return
	}
	if s1.root.removeNotIn(s2.root) {
		s1.root = nil
	}
}

// String returns a representation of s in standard set notation.
func (s *Sparse) String() string {
	var b strings.Builder
	b.WriteByte('{')
	first := true
	s.Elements(func(elts []uint64) bool {
		for _, e := range elts {
			if !first {
				b.WriteString(", ")
			}
			first = false
			b.WriteString(strconv.FormatUint(e, 10))
		}
		return true
	})
	b.WriteByte('}')
	return b.String()
}

// Elements calls f on successive slices of the set's elements, from lowest to
// highest. If f returns false, the iteration stops. The slice passed to f will
// be reused when f returns.
func (s *Sparse) Elements(f func([]uint64) bool) {
	if s.root == nil {
		return
	}
	s.root.elements(f, 0)
}

func (s *Sparse) memSize() uint64 {
	sz := memSize(*s)
	if s.root != nil {
		sz += s.root.memSize()
	}
	return sz
}

func memSize(x interface{}) uint64 {
	return uint64(reflect.TypeOf(x).Size())
}
