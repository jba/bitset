package bitset

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

// elementsOf returns all elements of s in ascending order.
func elementsOf(s *Sparse) []uint64 {
	var elts []uint64
	for e := range s.Elements() {
		elts = append(elts, e)
	}
	return elts
}

func TestSparseCopy(t *testing.T) {
	orig := sparseFrom(1, 1000, 500000, 1<<40)
	cp := orig.Copy()
	if !cp.Equal(orig) {
		t.Fatalf("copy not equal: got %s, want %s", cp, orig)
	}
	// Mutating the copy must not affect the original.
	cp.Add64(7)
	cp.Remove64(1000)
	if orig.Contains64(7) {
		t.Error("adding to copy changed original")
	}
	if !orig.Contains64(1000) {
		t.Error("removing from copy changed original")
	}
	// Mutating the original must not affect the copy.
	orig2 := sparseFrom(2, 3, 4)
	cp2 := orig2.Copy()
	orig2.Add64(9)
	orig2.Remove64(2)
	if cp2.Contains64(9) {
		t.Error("adding to original changed copy")
	}
	if !cp2.Contains64(2) {
		t.Error("removing from original changed copy")
	}
	// Copy of empty set.
	var empty Sparse
	ec := empty.Copy()
	if !ec.Empty() {
		t.Error("copy of empty set is not empty")
	}
}

func TestSparseAddInNoAlias(t *testing.T) {
	// Values share a deep prefix so they land in the same set256 leaf.
	s1 := sparseFrom(1)
	s2 := sparseFrom(256, 257)
	s1.AddIn(s2)
	// Mutating s1 must not corrupt s2 (which was spliced in via copy()).
	s1.Add64(258)
	if s2.Contains64(258) {
		t.Error("mutating s1 after AddIn corrupted s2")
	}
	// And mutating s2 must not corrupt s1.
	s2.Add64(259)
	if s1.Contains64(259) {
		t.Error("mutating s2 after AddIn corrupted s1")
	}
	// Also exercise the deeper-node sharing case.
	a := sparseFrom(1 << 40)
	b := sparseFrom((1<<40)+1, (1<<40)+2)
	a.AddIn(b)
	a.Add64((1 << 40) + 3)
	if b.Contains64((1 << 40) + 3) {
		t.Error("deep AddIn aliased subtree")
	}
}

func TestSparseEqualNegative(t *testing.T) {
	for _, test := range []struct {
		a, b []uint64
		want bool
	}{
		{nil, nil, true},
		{[]uint64{1, 2, 3}, []uint64{1, 2, 3}, true},
		{[]uint64{1, 2, 3}, []uint64{1, 2, 4}, false},  // differing element, same bitset shape
		{[]uint64{1, 2, 3}, []uint64{1, 2}, false},     // superset
		{[]uint64{1, 2}, []uint64{1, 2, 3}, false},     // subset
		{[]uint64{1}, nil, false},                      // one empty
		{nil, []uint64{1}, false},                      // other empty
		{[]uint64{1 << 40}, []uint64{1 << 41}, false},  // different high byte
		{[]uint64{256, 257}, []uint64{256, 258}, false}, // same leaf, diff element
	} {
		a := sparseFrom(test.a...)
		b := sparseFrom(test.b...)
		if got := a.Equal(b); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.a, test.b, got, test.want)
		}
	}
}

func TestSparseElementsEarlyStop(t *testing.T) {
	s := sparseFrom(5, 9, 70, 300, 12345, 1e8)
	var got []uint64
	for e := range s.Elements() {
		got = append(got, e)
		break // stop after the first element
	}
	if len(got) == 0 {
		t.Fatal("iterator never yielded")
	}
	// Whatever we got must be a prefix of the sorted elements.
	want := elementsOf(sparseFrom(5, 9, 70, 300, 12345, 1e8))
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("early-stop elements not a prefix: got %v, want prefix of %v", got, want)
		}
	}
	if len(got) == len(want) {
		t.Errorf("early stop did not halt iteration: got all %d elements", len(got))
	}
}

func TestSparseRemoveAll(t *testing.T) {
	nums := uDedupSort(uRandSlice(1000))
	s := sparseFrom(nums...)
	if s.Len() != len(nums) {
		t.Fatalf("Len = %d, want %d", s.Len(), len(nums))
	}
	// Remove in shuffled order.
	order := rand.Perm(len(nums))
	for i, idx := range order {
		s.Remove64(nums[idx])
		if got, want := s.Len(), len(nums)-i-1; got != want {
			t.Fatalf("after removing %d elements: Len = %d, want %d", i+1, got, want)
		}
	}
	if !s.Empty() {
		t.Error("set not empty after removing all elements")
	}
	if s.root != nil {
		t.Error("root not nil after removing all elements")
	}
}

func TestSparseBoundaryValues(t *testing.T) {
	vals := []uint64{
		0,
		math.MaxUint64,
		1 << 56, 2 << 56, 255 << 56, // top byte
		1, 2, 255, // bottom byte
		1<<56 | 1, 1<<56 | 2, // differ in top and bottom
	}
	s := NewSparse()
	for _, v := range vals {
		s.Add64(v)
	}
	for _, v := range vals {
		if !s.Contains64(v) {
			t.Errorf("does not contain %d", v)
		}
	}
	if s.Len() != len(vals) {
		t.Errorf("Len = %d, want %d", s.Len(), len(vals))
	}
	// Elements come back sorted and complete.
	got := elementsOf(s)
	want := append([]uint64(nil), vals...)
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if len(got) != len(want) {
		t.Fatalf("Elements len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Elements = %v, want %v", got, want)
		}
	}
	// Removing the extremes works.
	s.Remove64(0)
	s.Remove64(math.MaxUint64)
	if s.Contains64(0) || s.Contains64(math.MaxUint64) {
		t.Error("extreme values still present after removal")
	}
}

func TestSparseLenMatchesElements(t *testing.T) {
	for i := 0; i < 50; i++ {
		s := sparseFrom(uRandSlice(200)...)
		if got, want := len(elementsOf(s)), s.Len(); got != want {
			t.Fatalf("Elements count %d != Len %d", got, want)
		}
	}
}

// TestSparseModel cross-checks Sparse against a map oracle while interleaving
// every mutating operation, including Copy.
func TestSparseModel(t *testing.T) {
	const ops = 3000
	const valRange = 5000 // keep collisions likely

	s := NewSparse()
	model := map[uint64]bool{}

	check := func(tag string) {
		t.Helper()
		if got, want := s.Len(), len(model); got != want {
			t.Fatalf("%s: Len = %d, want %d", tag, got, want)
		}
		got := elementsOf(s)
		want := make([]uint64, 0, len(model))
		for v := range model {
			want = append(want, v)
		}
		sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
		if len(got) != len(want) {
			t.Fatalf("%s: Elements len %d != %d", tag, len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("%s: Elements mismatch at %d: %v vs %v", tag, i, got, want)
			}
		}
	}

	// helper to build a random small set + its model
	randSet := func() (*Sparse, map[uint64]bool) {
		m := map[uint64]bool{}
		ns := NewSparse()
		for i := 0; i < rand.Intn(20); i++ {
			v := uint64(rand.Intn(valRange))
			ns.Add64(v)
			m[v] = true
		}
		return ns, m
	}

	for i := 0; i < ops; i++ {
		switch rand.Intn(8) {
		case 0, 1: // Add
			v := uint64(rand.Intn(valRange))
			s.Add64(v)
			model[v] = true
		case 2: // Remove
			v := uint64(rand.Intn(valRange))
			s.Remove64(v)
			delete(model, v)
		case 3: // Contains consistency
			v := uint64(rand.Intn(valRange))
			if s.Contains64(v) != model[v] {
				t.Fatalf("Contains(%d) = %t, want %t", v, s.Contains64(v), model[v])
			}
		case 4: // AddIn
			o, om := randSet()
			s.AddIn(o)
			for v := range om {
				model[v] = true
			}
			// o must be unchanged by later s mutations: spot check now.
			if !o.Equal(sparseFromMap(om)) {
				t.Fatal("AddIn mutated argument")
			}
		case 5: // RemoveIn
			o, om := randSet()
			s.RemoveIn(o)
			for v := range om {
				delete(model, v)
			}
		case 6: // RemoveNotIn
			o, om := randSet()
			s.RemoveNotIn(o)
			for v := range model {
				if !om[v] {
					delete(model, v)
				}
			}
		case 7: // Copy then mutate copy; original (s) must be unaffected
			cp := s.Copy()
			if !cp.Equal(s) {
				t.Fatal("Copy not equal to original")
			}
			cp.Add64(uint64(rand.Intn(valRange)))
			cp.Remove64(uint64(rand.Intn(valRange)))
			// s and model unchanged; verified by check below.
		}
		check("op")
	}
}

func sparseFromMap(m map[uint64]bool) *Sparse {
	s := NewSparse()
	for v := range m {
		s.Add64(v)
	}
	return s
}
