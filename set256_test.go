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

func TestElements256(t *testing.T) {
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

func TestElements64_256(t *testing.T) {
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
	} {
		n := s.elements64or(a[:test.n], uint8(test.start), test.high)
		got := a[:n]
		if !cmp.Equal(got, test.want) {
			t.Errorf("%+v: got %v, want %v", test, got, test.want)
		}
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

func naiveElementsUint64(s interface {
	cap() int
	contains(uint8) bool
}) []uint64 {
	var els []uint64
	for i := 0; i < s.cap(); i++ {
		u := uint8(i)
		if s.contains(u) {
			els = append(els, uint64(u))
		}
	}
	return els
}

func TestIntersectN(t *testing.T) {
	var c set256
	b1 := sampleSet256()
	b2 := sampleSet256()
	c.intersectN([]*set256{&b1, &b2})
	if !c.equal(&b1) {
		t.Fatal("not equal")
	}
	c.intersectN([]*set256{&b1, &b2, &set256{}})
	if !c.empty() {
		t.Fatal("not empty")
	}
	var b3 set256
	b3.add(201)
	b3.add(188)
	b3.add(254)
	c.intersectN([]*set256{&b1, &b3})
	if c.len() != 1 || !c.contains(201) {
		t.Fatal("bad c")
	}
}
