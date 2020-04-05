package bitset

import (
	"math/bits"
	"strconv"
	"strings"
)

// Set64 is an efficient representation of a bitset that can represent integers
// in the range [0, 64).
// For efficiency, the methods of Set64 perform no bounds checking on their arguments.
type Set64 uint64

// Set64From constructs a Set64 from a list of elements.
func Set64From(els ...uint8) Set64 {
	var s Set64
	s.with(els...)
	return s
}

func (s *Set64) with(els ...uint8) {
	for _, e := range els {
		s.Add(e)
	}
}

// Add adds u to s.
func (s *Set64) Add(u uint8) {
	*s |= 1 << u
}

// Remove removes u from s.
func (s *Set64) Remove(u uint8) {
	*s &^= (1 << u)
}

// Contains reports whether s contains u.
func (s *Set64) Contains(u uint8) bool {
	return *s&(1<<u) != 0
}

// Empty reports whether s has no elements.
func (s Set64) Empty() bool {
	return s == 0
}

// Len returns the number of elements in s.
func (s Set64) Len() int {
	return bits.OnesCount64(uint64(s))
}

// func (Set64) Cap() int {
// 	return 64
// }

// Clear removes all elements from s.
func (s *Set64) Clear() {
	*s = 0
}

// Equal reports whether two bitsets have the same elements.
func (s1 Set64) Equal(s2 Set64) bool {
	return s1 == s2
}

// Complement replaces s with its complement.
func (s *Set64) Complement() {
	*s = ^*s
}

// AddIn adds all the elements in s2 to s1.
// It sets s1 to the union of s1 and s2.
func (s1 *Set64) AddIn(s2 Set64) {
	*s1 |= s2
}

// RemoveIn removes from s1 all the elements that are in s2.
// It sets s1 to the set difference of s1 and s2.
func (s1 *Set64) RemoveIn(s2 Set64) {
	s2.Complement()
	s1.RemoveNotIn(s2)
}

// RemoveNotIn removes from s1 all the elements that are not in s2.
// It sets s1 to the intersection of s1 and s2.
func (s1 *Set64) RemoveNotIn(s2 Set64) {
	*s1 &= s2
}

// append appends the elements of s to elts, in ascending order.
func (s Set64) append(elts []uint8) []uint8 {
	low, high := s.elementRange()
	for e := low; e < high; e++ {
		u := uint8(e)
		if s.Contains(u) {
			elts = append(elts, u)
		}
	}
	return elts
}

func (s Set64) populate(b *[64]uint) int {
	low, high := s.elementRange()
	i := 0
	for e := low; e < high; e++ {
		if s.Contains(uint8(e)) {
			(*b)[i] = uint(e)
			i++
		}
	}
	return i
}

func (s Set64) populate64(b *[64]uint64) int {
	low, high := s.elementRange()
	i := 0
	for e := low; e < high; e++ {
		if s.Contains(uint8(e)) {
			(*b)[i] = uint64(e)
			i++
		}
	}
	return i
}

func (s Set64) elementRange() (int, int) {
	return bits.TrailingZeros64(uint64(s)), 64 - bits.LeadingZeros64(uint64(s))
}

// String returns a representation of s in standard set notation.
func (s Set64) String() string {
	var buf [64]uint
	n := s.populate(&buf)
	var b strings.Builder
	b.WriteByte('{')
	first := true
	for _, e := range buf[:n] {
		if !first {
			b.WriteString(", ")
		}
		first = false
		b.WriteString(strconv.FormatUint(uint64(e), 10))
	}
	b.WriteByte('}')
	return b.String()
}

// position returns the 0-based position of n in the set. If the set
// is {3, 8, 15}, then the position of 8 is 1.  If n is not in the
// set, position returns the position n would be at if it were a
// member. The second return value reports whether n is a member of
// s.
func (s Set64) position(n uint8) (int, bool) {
	mask := uint64(1 << n)
	in := (uint64(s)&mask != 0)
	pos := bits.OnesCount64(uint64(s) & (mask - 1))
	return pos, in
}
