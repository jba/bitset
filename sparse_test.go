package bitset

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func sparseFrom(us ...uint64) *Sparse {
	s := NewSparse()
	for _, e := range us {
		s.Add64(e)
	}
	return s
}

func TestSparseBasics(t *testing.T) {
	check := func(b bool) {
		t.Helper()
		if !b {
			t.Fatal("check failed")
		}
	}
	var s Sparse

	check(s.Empty())
	s.Add(0)
	check(!s.Empty())
	check(s.Contains(0))
	check(!s.Contains(1))

	s.Add(492409)
	check(!s.Empty())
	check(s.Contains(0))
	check(!s.Contains(1))
	check(s.Contains(492409))

	s.Remove(0)
	check(!s.Empty())
	check(!s.Contains(0))
	check(!s.Contains(1))
	check(s.Contains(492409))

	s.Remove(492409)
	check(s.Empty())
	check(!s.Contains(0))
	check(!s.Contains(1))
	check(!s.Contains(492409))
}

// TODO: use cover to make sure we're hitting everything
func TestSparseAddIn(t *testing.T) {
	for _, test := range []struct {
		in1, in2 []uint64
	}{
		{nil, nil},
		{nil, []uint64{1}},
		{[]uint64{17, 99}, []uint64{3, 500, 1000}},
	} {
		s1 := sparseFrom(test.in1...)
		s2 := sparseFrom(test.in2...)
		s1.AddIn(s2)
		want := sparseFrom(uUnion(test.in1, test.in2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in1, test.in2, s1, want)
		}

		s1 = sparseFrom(test.in2...)
		s2 = sparseFrom(test.in1...)
		s1.AddIn(s2)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in2, test.in1, s1, want)
		}
	}

	const sz = 100
	const n = 100
	for i := 0; i < n; i++ {
		u1 := uRandSlice(sz)
		u2 := uRandSlice(sz)
		s1 := sparseFrom(u1...)
		s2 := sparseFrom(u2...)
		s1.AddIn(s2)
		want := sparseFrom(uUnion(u1, u2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", u1, u2, s1, want)
		}
	}
}

func TestSparseRemoveIn(t *testing.T) {
	for _, test := range []struct {
		in1, in2 []uint64
	}{
		{nil, nil},
		{nil, []uint64{1}},
		{[]uint64{17, 99}, []uint64{3, 500, 1000}},
		{[]uint64{5000, 7000, 9000, 11000}, []uint64{2000, 5000, 7000, 11000}},
	} {
		s1 := sparseFrom(test.in1...)
		s2 := sparseFrom(test.in2...)
		s1.RemoveIn(s2)
		want := sparseFrom(uDifference(test.in1, test.in2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in1, test.in2, s1, want)
		}

		s1 = sparseFrom(test.in2...)
		s2 = sparseFrom(test.in1...)
		s1.RemoveIn(s2)
		want = sparseFrom(uDifference(test.in2, test.in1)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in2, test.in1, s1, want)
		}
	}

	const sz = 100
	const n = 100
	for i := 0; i < n; i++ {
		u1 := uRandSlice(sz)
		u2 := uRandSlice(sz)
		s1 := sparseFrom(u1...)
		s2 := sparseFrom(u2...)
		s1.RemoveIn(s2)
		want := sparseFrom(uDifference(u1, u2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", u1, u2, s1, want)
		}
	}
}

func TestSparseRemoveNotIn(t *testing.T) {
	t.Skip()
	for _, test := range []struct {
		in1, in2 []uint64
	}{
		{nil, nil},
		{nil, []uint64{1}},
		{[]uint64{17, 99}, []uint64{3, 500, 1000}},
		{[]uint64{5000, 7000, 9000, 11000}, []uint64{2000, 5000, 7000, 11000}},
	} {
		s1 := sparseFrom(test.in1...)
		s2 := sparseFrom(test.in2...)
		s1.RemoveNotIn(s2)
		want := sparseFrom(uIntersection(test.in1, test.in2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in1, test.in2, s1, want)
		}

		s1 = sparseFrom(test.in2...)
		s2 = sparseFrom(test.in1...)
		s1.RemoveNotIn(s2)
		want = sparseFrom(uIntersection(test.in2, test.in1)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", test.in2, test.in1, s1, want)
		}
	}

	const sz = 100
	const n = 100
	for i := 0; i < n; i++ {
		u1 := uRandSlice(sz)
		u2 := uRandSlice(sz)
		s1 := sparseFrom(u1...)
		s2 := sparseFrom(u2...)
		s1.RemoveNotIn(s2)
		want := sparseFrom(uIntersection(u1, u2)...)
		if !s1.Equal(want) {
			t.Errorf("%v, %v: got %s, want %s", u1, u2, s1, want)
		}
	}
}

func TestLots(t *testing.T) {
	var s Sparse
	nums := uRandSlice(1e3)
	for i, n := range nums {
		if s.Len() != i {
			t.Fatalf("s.Size() = %d, want %d", s.Len(), i)
		}
		s.Add64(n)

	}
	for _, n := range nums {
		if !s.Contains64(n) {
			t.Errorf("does not contain %d", n)
		}
	}
	for i, n := range nums {
		got := s.Len()
		want := len(nums) - i
		if got != want {
			t.Fatalf("s.Size() = %d, want %d", got, want)
		}
		s.Remove64(n)
	}
	for _, n := range nums {
		if s.Contains64(n) {
			t.Errorf("does contain %d", n)
		}
	}
}

type uslice []uint64

func (u uslice) Len() int           { return len(u) }
func (u uslice) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u uslice) Less(i, j int) bool { return u[i] < u[j] }

func TestSparseElements1(t *testing.T) {
	var s Sparse
	els := []uint64{3, 17, 300, 12345, 1e8}
	for _, e := range els {
		s.Add64(e)
	}
	if !s.Contains(1e8) {
		t.Fatal("no 1e8")
	}
	a := make([]uint64, len(els), len(els))
	n := s.elements(a, 0)
	got := a[:n]
	if !reflect.DeepEqual(got, els) {
		t.Fatalf("got %v, want %v", got, els)
	}
}

func TestSparseElements2(t *testing.T) {
	var s Sparse
	nums := uRandSlice(1e3)
	for _, n := range nums {
		s.Add64(n)
	}
	sort.Sort(uslice(nums))
	a := make([]uint64, len(nums), len(nums))
	if s.Len() != len(nums) {
		t.Fatalf("size: got %d", s.Len())
	}

	n := s.elements(a, 0)
	if n != len(nums) {
		t.Fatalf("len: got %d, want %d", n, len(nums))
	}
	if !reflect.DeepEqual(a[:n], nums) {
		t.Fatal("not equal")
	}
}

func set(x ...uint64) []uint64 { return x }

func TestString(t *testing.T) {
	for _, test := range []struct {
		els  []uint64
		want string
	}{
		{nil, "{}"},
		{set(9), "{9}"},
		{set(9, 1e4, 99), "{9, 99, 10000}"},
	} {
		got := sparseFrom(test.els...).String()
		if got != test.want {
			t.Errorf("%v: got %q, want %q", test.els, got, test.want)
		}
	}
}

// func TestIntersect(t *testing.T) {
// 	for _, test := range []struct {
// 		els1, els2, want []uint64
// 	}{
// 		{nil, nil, nil},
// 		{set(9), nil, nil},
// 		{nil, set(9), nil},
// 		{set(9), set(9), set(9)},
// 		{set(9), set(9, 10), set(9)},
// 		{set(9, 99, 1e8), set(99, 1e8+1), set(99)},
// 	} {
// 		s1 := sparseFrom(test.els1...)
// 		s2 := sparseFrom(test.els2...)
// 		want := sparseFrom(test.want...)
// 		var got Sparse
// 		got.Intersect(s1, s2)
// 		if !got.Equal(want) {
// 			t.Errorf("%s & %s = %v, want %v", s1, s2, got, want)
// 		}
// 	}
// }

// func TestConsecutive(t *testing.T) {
// 	for _, start := range []uint64{0, 100, 1e8} {
// 		for _, sz := range []int{0, 1, 2, 3, 4, 5, 64, 256, 512, 1000, 10000, 100000} {
// 			var s SparseSet
// 			for i := 0; i < sz; i++ {
// 				s.Add(uint64(i) + start)
// 			}
// 			fmt.Printf("consec: size=%d, start=%d, bytes=%d\n", s.Size(), start, s.MemSize())
// 		}
// 	}

// 	fmt.Printf("memsize of set256: %d\n", memSize(Set256{}))
// 	fmt.Printf("memsize of node: %d\n", memSize(node{}))
// 	fmt.Printf("memsize of *node: %d\n", memSize(&node{}))
// }
// func TestMemSize(t *testing.T) {
// 	var s SparseSet
// 	for i := 0; i < 3; i++ {
// 		fmt.Printf("Add %d\n", i)
// 		s.Add(uint64(i))
// 	}
// 	fmt.Printf("consec: size=%d, bytes=%d\n", s.Size(), s.MemSize())
// }

func uUnion(u1, u2 []uint64) []uint64 {
	m1 := uMap(u1)
	for _, u := range u2 {
		m1[u] = true
	}
	return uSlice(m1)
}

func uIntersection(u1, u2 []uint64) []uint64 {
	m1 := uMap(u1)
	m2 := uMap(u2)
	for u := range m1 {
		if !m2[u] {
			delete(m1, u)
		}
	}
	return uSlice(m1)
}

func uDifference(u1, u2 []uint64) []uint64 {
	m1 := uMap(u1)
	for _, u := range u2 {
		delete(m1, u)
	}
	return uSlice(m1)
}

func uMap(us []uint64) map[uint64]bool {
	m := map[uint64]bool{}
	for _, u := range us {
		m[u] = true
	}
	return m
}

func uSlice(m map[uint64]bool) []uint64 {
	if len(m) == 0 {
		return nil
	}
	s := make([]uint64, 0, len(m))
	for u := range m {
		s = append(s, u)
	}
	return s
}

func uRandSlice(n int) []uint64 {
	s := make([]uint64, n)
	for i := 0; i < len(s); i++ {
		s[i] = uRand()
	}
	return s
}

func uRand() uint64 {
	lo := uint64(rand.Uint32())
	hi := uint64(rand.Uint32())
	return (hi << 32) | lo
}
