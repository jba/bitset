package bitset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sampleSet256() set256 {
	var s set256
	s.add(3)
	s.add(63)
	s.add(17)
	s.add(70)
	s.add(200)
	s.add(201)
	s.add(192)
	return s
}

func TestBasics256(t *testing.T) {
	s := sampleSet256()
	want := "{3, 17, 63, 70, 192, 200, 201}"
	got := s.String()
	if got != want {
		t.Errorf("s.String() = %q, want %q", got, want)
	}
	if !s.equal(&s) {
		t.Fatal("not equal")
	}
	if !cmp.Equal(naiveElementsUint64(&s), []uint64{3, 17, 63, 70, 192, 200, 201}) {
		t.Errorf("%s: wrong elements", s)
	}
	if s.len() != 7 {
		t.Error("wrong size")
	}
	if s.empty() {
		t.Error("shouldn't be empty")
	}
	var z set256
	if !z.empty() {
		t.Error("should be empty")
	}
}

func TestPosition256(t *testing.T) {
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

func naiveElementsUint64(s *set256) []uint64 {
	var els []uint64
	for i := 0; i < 256; i++ {
		if s.contains(uint8(i)) {
			els = append(els, uint64(i))
		}
	}
	return els
}
