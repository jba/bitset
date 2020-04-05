package bitset

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	if !cmp.Equal(got, els) {
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
	if !cmp.Equal(a[:n], nums) {
		t.Fatal("not equal")
	}
}

func TestString(t *testing.T) {
	for _, test := range []struct {
		els  []uint64
		want string
	}{
		{nil, "{}"},
		{[]uint64{9}, "{9}"},
		{[]uint64{3000, 2000, 1000, 3000}, "{1000, 2000, 3000}"},
		{[]uint64{9, 1e4, 99}, "{9, 99, 10000}"},
	} {
		got := sparseFrom(test.els...).String()
		if got != test.want {
			t.Errorf("%v: got %q, want %q", test.els, got, test.want)
		}
	}
}

func TestMemSize(t *testing.T) {
	s := NewSparse()
	for i := 1000; i <= 1e6; i += 1000 {
		s.Add(uint(i))
	}
	fmt.Printf("size=%d, bytes=%d\n", s.Len(), s.memSize())
}

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
