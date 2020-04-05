package bitset

import (
	"fmt"
	"strings"
)

// A set256 represents a set of integers in the range [0, 256).
// It does so more efficiently than a Dense set of capacity 256.
// For efficiency, the methods of set256 perform no bounds checking on their
// arguments.
type set256 struct {
	sets [4]Set64
}

func (s *set256) copy() subber {
	c := *s
	return &c
}

func (s *set256) add(n uint8) {
	s.sets[n/64].Add(n % 64)
}

func (s *set256) remove(n uint8) {
	s.sets[n/64].Remove(n % 64)
}

func (s *set256) contains(n uint8) bool {
	return s.sets[n/64].Contains(n % 64)
}

func (s *set256) empty() bool {
	return s.sets[0].Empty() && s.sets[1].Empty() && s.sets[2].Empty() && s.sets[3].Empty()
}

func (s *set256) clear() {
	s.sets[0].Clear()
	s.sets[1].Clear()
	s.sets[2].Clear()
	s.sets[3].Clear()
}

func (s *set256) len() int {
	return s.sets[0].Len() + s.sets[1].Len() + s.sets[2].Len() + s.sets[3].Len()
}

func (set256) cap() int {
	return 256
}

func (s1 *set256) equal(s2 *set256) bool {
	return s1.sets[0] == s2.sets[0] &&
		s1.sets[1] == s2.sets[1] &&
		s1.sets[2] == s2.sets[2] &&
		s1.sets[3] == s2.sets[3]
}

// position returns the 0-based position of n in the set. If
// the set is {3, 8, 15}, then the position of 8 is 1.
// If n is not in the set, returns 0, false.
// If not a member, return where it would go.
// The second return value reports whether n is a member of b.
func (b *set256) position(n uint8) (int, bool) {
	var pos int
	i := n / 64
	switch i {
	case 1:
		pos = b.sets[0].Len()
	case 2:
		pos = b.sets[0].Len() + b.sets[1].Len()
	case 3:
		pos = b.sets[0].Len() + b.sets[1].Len() + b.sets[2].Len()
	}
	p, ok := b.sets[i].position(n % 64)
	return pos + p, ok
}

func (s1 *set256) addIn(sub subber) {
	s2 := sub.(*set256)
	s1.sets[0].AddIn(s2.sets[0])
	s1.sets[1].AddIn(s2.sets[1])
	s1.sets[2].AddIn(s2.sets[2])
	s1.sets[3].AddIn(s2.sets[3])
}

// c = a intersect b
// func (c *Set256) Intersect2(a, b *Set256) {
// 	c.sets[0] = a.sets[0] & b.sets[0]
// 	c.sets[1] = a.sets[1] & b.sets[1]
// 	c.sets[2] = a.sets[2] & b.sets[2]
// 	c.sets[3] = a.sets[3] & b.sets[3]
// }

// c cannot be one of sets
func (c *set256) intersectN(bs []*set256) {
	if len(bs) == 0 {
		c.clear()
		return
	}
	for i := 0; i < len(c.sets); i++ {
		c.sets[i] = bs[0].sets[i]
		for _, s := range bs[1:] {
			c.sets[i].RemoveNotIn(s.sets[i])
		}
	}
}

// Fill a with set elements, starting from start.
// Return the number added.
func (s *set256) elements8(a []uint8, start uint8) int {
	if len(a) == 0 {
		return 0
	}
	si := start / 64
	n := s.sets[si].elementsOr(a, start%64, si*64)
	for i := si + 1; i < 4; i++ {
		n += s.sets[i].elementsOr(a[n:], 0, i*64)
	}
	return n
}

func (s *set256) elements64high8(a []uint64, start uint8, high uint64) int {
	if len(a) == 0 {
		return 0
	}
	si := start / 64
	n := s.sets[si].elements64or(a, start%64, high|uint64(si*64))
	for i := si + 1; i < 4; i++ {
		n += s.sets[i].elements64or(a[n:], 0, high|uint64(i*64))
	}
	return n
}

func (s set256) String() string {
	var a [256]uint64
	n := s.elements64high(a[:], 0, 0)
	if n == 0 {
		return "{}"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "{%d", a[0])
	for _, e := range a[1:n] {
		fmt.Fprintf(&b, ", %d", e)
	}
	b.WriteByte('}')
	return b.String()
}

// For subber, used in node:

func (s *set256) add64(e uint64) { s.add(uint8(e)) }

func (s *set256) remove64(e uint64) bool {
	s.remove(uint8(e))
	return s.empty()
}

func (s *set256) contains64(e uint64) bool {
	return s.contains(uint8(e))
}

func (s *set256) memSize() uint64 { return memSize(*s) }

func (s *set256) elements64high(a []uint64, start, high uint64) int {
	return s.elements64high8(a, uint8(start), high)
}

func (s *set256) equalSub(b subber) bool {
	return s.equal(b.(*set256))
}
