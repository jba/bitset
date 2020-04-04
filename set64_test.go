package bitset

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sampleSet64() Set64 {
	var s Set64
	s.With(3, 63, 17)
	return s
}

func TestSet64String(t *testing.T) {
	for _, test := range []struct {
		set  Set64
		want string
	}{
		{Set64(0), "{}"},
		{Set64(8), "{3}"},
		{sampleSet64(), "{3, 17, 63}"},
		{Set64(math.MaxUint64), "{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}"},
	} {
		got := test.set.String()
		if got != test.want {
			t.Errorf("%d: got %q, want %q", test.set, got, test.want)
		}
	}
}

func TestSet64Add(t *testing.T) {
	for _, test := range []struct {
		in, want Set64
	}{
		{Set64(0), Set64(1 << 7)},
		{Set64(1 << 3), Set64From(3, 7)},
		{Set64From(3, 63, 17), Set64From(3, 7, 17, 63)},
		{Set64(math.MaxUint64), Set64(math.MaxUint64)},
	} {
		got := test.in
		got.Add(7)
		if got != test.want {
			t.Errorf("%s: got %q, want %q", test.in, got, test.want)
		}
	}
}

func TestSet64LenEmpty(t *testing.T) {
	for _, test := range []struct {
		in   Set64
		want int
	}{
		{Set64(0), 0},
		{Set64(1 << 3), 1},
		{Set64From(3, 63, 17), 3},
		{Set64(math.MaxUint64), 64},
	} {
		got := test.in.Len()
		if got != test.want {
			t.Errorf("%s: got %d, want %d", test.in, got, test.want)
		}
		if want := got == 0; test.in.Empty() != want {
			t.Errorf("%s.Empty: got %t, want %t", test.in, test.in.Empty(), want)
		}
	}
}

func TestElements(t *testing.T) {
	var a [10]uint8
	s := sampleSet64()
	for _, test := range []struct {
		n     int
		start uint8
		want  []uint8
	}{
		{0, 0, []uint8{}},
		{0, 10, []uint8{}},
		{1, 0, []uint8{3}},
		{1, 5, []uint8{17}},
		{1, 27, []uint8{63}},
		{2, 0, []uint8{3, 17}},
		{2, 5, []uint8{17, 63}},
		{2, 39, []uint8{63}},
		{2, 63, []uint8{63}},
		{2, 83, []uint8{}},
		{3, 0, []uint8{3, 17, 63}},
		{3, 10, []uint8{17, 63}},
		{3, 99, []uint8{}},
	} {
		n := s.Elements(a[:test.n], test.start)
		got := a[:n]
		if !cmp.Equal(got, test.want) {
			t.Errorf("%+v: got %v, want %v", test, got, test.want)
		}
	}
}

func TestElements64(t *testing.T) {
	var a [10]uint64
	s := sampleSet64()
	for _, test := range []struct {
		n     int
		start uint64
		high  uint64
		want  []uint64
	}{
		{0, 0, 0, []uint64{}},
		{0, 10, 0, []uint64{}},
		{1, 0, 0, []uint64{3}},
		{1, 5, 0, []uint64{17}},
		{1, 27, 0, []uint64{63}},
		{2, 0, 0, []uint64{3, 17}},
		{2, 5, 0, []uint64{17, 63}},
		{2, 39, 0, []uint64{63}},
		{2, 63, 0, []uint64{63}},
		{2, 83, 0, []uint64{}},
		{3, 0, 0, []uint64{3, 17, 63}},
		{3, 10, 0, []uint64{17, 63}},
		{3, 99, 0, []uint64{}},
		{3, 10, 64, []uint64{64 + 17, 64 + 63}},
		{3, 0, 256, []uint64{256 + 3, 256 + 17, 256 + 63}},
	} {
		n := s.elements64(a[:test.n], uint8(test.start), test.high)
		got := a[:n]
		if !cmp.Equal(got, test.want) {
			t.Errorf("%+v: got %v", test, got)
		}
	}
}

func TestPosition(t *testing.T) {
	s := sampleSet64()
	for _, test := range []struct {
		n   uint8
		pos int
		in  bool
	}{
		{0, 0, false},
		{1, 0, false},
		{2, 0, false},
		{3, 0, true},
		{4, 1, false},
		{10, 1, false},
		{16, 1, false},
		{17, 1, true},
		{20, 2, false},
		{62, 2, false},
		{63, 2, true},
	} {
		gotPos, gotIn := s.position(test.n)
		if gotPos != test.pos || gotIn != test.in {
			t.Errorf("Position(%d) = (%d, %t), want (%d, %t)", test.n, gotPos, gotIn, test.pos, test.in)
		}
	}
}

func naiveElementsUint8(s Set64) []uint8 {
	var els []uint8
	for i := 0; i < s.Cap(); i++ {
		u := uint8(i)
		if s.Contains(u) {
			els = append(els, u)
		}
	}
	return els
}
