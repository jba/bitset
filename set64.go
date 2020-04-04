package bitset

import (
	"fmt"
	"math/bits"
	"strings"
)

// Set64 is an efficient representation of a bitset that can represent integers
// in the range [0, 64).
// For efficiency, the methods of Set64 perform no bounds checking on their arguments.
type Set64 uint64

func Set64From(els ...uint8) Set64 {
	var s Set64
	s.With(els...)
	return s
}

func (s *Set64) With(els ...uint8) {
	for _, e := range els {
		s.Add(e)
	}
}

func (s *Set64) Add(u uint8) {
	*s |= 1 << u
}

func (s *Set64) Remove(u uint8) {
	*s &^= (1 << u)
}

func (s *Set64) Contains(u uint8) bool {
	return *s&(1<<u) != 0
}

func (s Set64) Empty() bool {
	return s == 0
}

func (s Set64) Len() int {
	return bits.OnesCount64(uint64(s))
}

func (Set64) Cap() int {
	return 64
}

func (s *Set64) Clear() {
	*s = 0
}

func (s1 Set64) Equal(s2 Set64) bool {
	return s1 == s2
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

func (s1 *Set64) Intersect(s2 Set64) {
	*s1 &= s2
}

func (s1 *Set64) Union(s2 Set64) {
	*s1 |= s2
}

// TODO: complement, difference

// Elements populates els with at most len(els) elements of s, starting with
// start. That is, els[0] will be the smallest element of s that is greater than
// or equal to start. The return value is the number of elements added to els.
func (s Set64) Elements(els []uint8, start uint8) int {
	if len(els) == 0 {
		return 0
	}
	i := 0
	for b := start; b < 64 && i < len(els); b++ {
		if s.Contains(b) {
			els[i] = b
			i++
		}
	}
	return i
}

func (s Set64) elements64(a []uint64, start uint8, high uint64) int {
	if len(a) == 0 {
		return 0
	}
	i := 0
	for b := start; b < 64 && i < len(a); b++ {
		if s.Contains(b) {
			a[i] = high | uint64(b)
			i++
		}
	}
	return i
}

func (s Set64) String() string {
	if s.Empty() {
		return "{}"
	}
	var b strings.Builder
	b.WriteByte('{')
	first := true
	var i uint8
	for i = 0; i < 64; i++ {
		if s.Contains(i) {
			if !first {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%d", i)
			first = false
		}
	}
	b.WriteByte('}')
	return b.String()
}
