package bitset

//TODO: use sync.Pool?

import (
	"fmt"
	"reflect"
	"strings"
)

// Sparse is a sparse bitset. It can represent any uint64 and uses memory
// proportional to the number of elements it contains.
type Sparse struct {
	root *node // compact radix tree 7 levels deep.
}

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

// Equal reports whether two sparse bitsets have the same elements.
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

func (s *Sparse) elements(a []uint64, start uint64) int {
	if s.root == nil {
		return 0
	}
	return s.root.elements64high(a, start, 0)
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

// RemoveIn removes all the elements in s2 from s1.
// It sets s1 to the set difference of s1 and s2.
func (s1 *Sparse) RemoveIn(s2 *Sparse) {
	if s1.Empty() || s2.Empty() {
		return
	}
	s1.root.removeIn(s2.root)
}

func (s1 *Sparse) RemoveNotIn(s2 *Sparse) {
	if s1.Empty() {
		return
	}
	if s2.Empty() {
		s1.Clear()
		return
	}
	s1.root.removeNotIn(s2.root)
}

// TODO: rethink
// s becomes the intersection of the ss. It must not be
// one of the ss, and it is not part of the intersection.
// func (s *Sparse) Intersect(ss ...*Sparse) {
// 	s.Clear()
// 	var nodes []*node
// 	for _, t := range ss {
// 		if t.Empty() {
// 			return
// 		}
// 		nodes = append(nodes, t.root)
// 	}
// 	s.root = intersectNodes(nodes)
// }

func (s Sparse) String() string {
	if s.Empty() {
		return "{}"
	}
	els := make([]uint64, s.Len())
	s.elements(els, 0)
	var b strings.Builder
	fmt.Fprintf(&b, "{%d", els[0])
	for _, e := range els[1:] {
		fmt.Fprintf(&b, ", %d", e)
	}
	b.WriteByte('}')
	return b.String()
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
