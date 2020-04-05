package bitset

// Dense is a standard bitset, represented as a sequence of bits. See Sparse in
// this package for a more memory-efficient storage scheme for sparse bitsets.
type Dense struct {
	sets []Set64
}

// NewDense creates a set capable of representing values in the range
// [0, capacity), at least. The Cap method reports the exact capacity.
// NewDense panics if capacity is negative.
func NewDense(capacity int) *Dense {
	return &Dense{
		sets: setslice(capacity),
	}
}

func setslice(capacity int) []Set64 {
	if capacity == 0 {
		return nil
	}
	if capacity < 0 {
		panic("negative capacity")
	}
	return make([]Set64, (capacity-1)/64+1)
}

func (s *Dense) Cap() int {
	return len(s.sets) * 64
}

func (s *Dense) Len() int {
	sz := 0
	for _, t := range s.sets {
		sz += t.Len()
	}
	return sz
}

func (s *Dense) Empty() bool {
	for _, t := range s.sets {
		if !t.Empty() {
			return false
		}
	}
	return true
}

func (s *Dense) Copy() *Dense {
	newSets := make([]Set64, len(s.sets))
	copy(newSets, s.sets)
	return &Dense{sets: newSets}
}

func (s *Dense) Add(u uint) {
	s.sets[u/64].Add(uint8(u % 64))
}

func (s *Dense) Remove(u uint) {
	s.sets[u/64].Remove(uint8(u % 64))
}

func (s *Dense) Contains(u uint) bool {
	return s.sets[u/64].Contains(uint8(u % 64))
}

func (s *Dense) Clear() {
	for i := range s.sets { // can't use _, t because it copies
		s.sets[i].Clear()
	}
}

func (s *Dense) SetCap(newCapacity int) {
	newSets := setslice(newCapacity)
	copy(newSets, s.sets)
	s.sets = newSets
}

// Equal reports whether s2 has the same elements as s1. It may have a different capacity.
func (s1 *Dense) Equal(s2 *Dense) bool {
	if len(s1.sets) > len(s2.sets) {
		s1, s2 = s2, s1
	}
	// Here, len(s1.sets) <= len(s2.sets).
	for i, t1 := range s1.sets {
		if t1 != s2.sets[i] {
			return false
		}
	}
	for _, t2 := range s2.sets[len(s1.sets):] {
		if t2 != 0 {
			return false
		}
	}
	return true
}

func (s *Dense) Complement() {
	for i := 0; i < len(s.sets); i++ {
		s.sets[i].Complement()
	}
}

func (s1 *Dense) AddIn(s2 *Dense) {
	if s1.Cap() < s2.Cap() {
		// TODO: Grow s1 less if it's not necessary, or panic.
		s1.SetCap(s2.Cap())
	}
	for i, t2 := range s2.sets {
		s1.sets[i].AddIn(t2)
	}
}

func (s1 *Dense) RemoveIn(s2 *Dense) {
	min := len(s1.sets)
	if min > len(s2.sets) {
		min = len(s2.sets)
	}
	for i := 0; i < min; i++ {
		s1.sets[i].RemoveIn(s2.sets[i])
	}
}

func (s1 *Dense) RemoveNotIn(s2 *Dense) {
	min := len(s1.sets)
	if min > len(s2.sets) {
		min = len(s2.sets)
	}
	for i := 0; i < min; i++ {
		s1.sets[i].RemoveNotIn(s2.sets[i])
	}
	for i := min; i < len(s1.sets); i++ {
		s1.sets[i].Clear()
	}
}

// Elements calls f on successive slices of the set's elements, from lowest to
// highest. If f returns false, the iteration stops. The slice passed to f will
// be reused when f returns.
func (s *Dense) Elements(f func([]uint) bool) {
	var buf [64]uint
	for i, t := range s.sets {
		n := t.populate(&buf)
		offset := uint(64 * i)
		for j := range buf[:n] {
			buf[j] += offset
		}
		if !f(buf[:n]) {
			break
		}
	}
}
