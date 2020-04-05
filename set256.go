package bitset

import (
	"strconv"
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

func (s *set256) add64(e uint64) { s.add(uint8(e)) }

func (s *set256) remove(n uint8) {
	s.sets[n/64].Remove(n % 64)
}

func (s *set256) remove64(e uint64) bool {
	s.remove(uint8(e))
	return s.empty()
}

func (s *set256) contains(n uint8) bool {
	return s.sets[n/64].Contains(n % 64)
}

func (s *set256) contains64(e uint64) bool { return s.contains(uint8(e)) }

func (s *set256) empty() bool {
	return s.sets[0].Empty() && s.sets[1].Empty() && s.sets[2].Empty() && s.sets[3].Empty()
}

func (s *set256) len() int {
	return s.sets[0].Len() + s.sets[1].Len() + s.sets[2].Len() + s.sets[3].Len()
}

func (s1 *set256) equal(b subber) bool {
	s2 := b.(*set256)
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

func (s1 *set256) removeIn(sub subber) (empty bool) {
	s2 := sub.(*set256)
	s1.sets[0].RemoveIn(s2.sets[0])
	s1.sets[1].RemoveIn(s2.sets[1])
	s1.sets[2].RemoveIn(s2.sets[2])
	s1.sets[3].RemoveIn(s2.sets[3])
	return s1.empty()
}

func (s1 *set256) removeNotIn(sub subber) (empty bool) {
	s2 := sub.(*set256)
	s1.sets[0].RemoveNotIn(s2.sets[0])
	s1.sets[1].RemoveNotIn(s2.sets[1])
	s1.sets[2].RemoveNotIn(s2.sets[2])
	s1.sets[3].RemoveNotIn(s2.sets[3])
	return s1.empty()
}

func (s *set256) elements(f func([]uint64) bool, offset uint64) bool {
	var buf [64]uint64
	for i, ss := range s.sets {
		n := ss.populate64(&buf)
		offset2 := offset + uint64(64*i)
		for j := range buf[:n] {
			buf[j] += offset2
		}
		if !f(buf[:n]) {
			return false
		}
	}
	return true
}

func (s set256) String() string {
	var b strings.Builder
	b.WriteByte('{')
	first := true
	s.elements(func(elts []uint64) bool {
		for _, e := range elts {
			if !first {
				b.WriteString(", ")
			}
			first = false
			b.WriteString(strconv.FormatUint(e, 10))
		}
		return true
	}, 0)
	b.WriteByte('}')
	return b.String()
}

func (s *set256) memSize() uint64 { return memSize(*s) }
