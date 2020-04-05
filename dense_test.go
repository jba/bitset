package bitset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDense(t *testing.T) {
	check := func(b bool) {
		t.Helper()
		if !b {
			t.Fatal("check failed")
		}
	}

	s := NewDense(100)
	check(s.Len() == 0)
	check(s.Empty())
	s.Add(17)
	check(s.Len() == 1)
	check(!s.Empty())
	check(s.Contains(17))
	check(!s.Contains(98))
	s.Add(98)
	check(s.Len() == 2)
	check(!s.Empty())
	check(s.Contains(17))
	check(s.Contains(98))
	s.Remove(13)
	check(s.Len() == 2)
	check(!s.Empty())
	check(s.Contains(17))
	check(s.Contains(98))
	s.Remove(17)
	check(s.Len() == 1)
	check(!s.Empty())
	check(!s.Contains(17))
	check(s.Contains(98))
	s.Clear()
	check(s.Len() == 0)
	check(s.Empty())

	s2 := NewDense(100)
	s.Add(17)
	s.Add(98)
	s2.Add(22)
	s2.Add(98)

	check(s.Equal(s))
	check(!s.Equal(s2))
	si := s.Copy()
	check(si.Equal(s))
	si.RemoveNotIn(s2)

	check(si.Len() == 1 && si.Contains(98))

	si.Complement()
	check(si.Len() == si.Cap()-1)
	check(!si.Contains(98))

	s = NewDense(30)
	s.Add(17)
	su := s.Copy()
	su.AddIn(s2)
	check(su.Len() == 3)
	check(su.Contains(17))
	check(su.Contains(22))
	check(su.Contains(98))

}

type s []uint

var tests = []struct {
	s1, s2       []uint
	union        []uint
	intersection []uint
	difference   []uint
}{
	{nil, nil, nil, nil, nil},
	{
		s{1}, s{2},
		s{1, 2}, s{}, s{1},
	},
	{
		s{5, 7, 8}, s{7, 9, 11},
		s{5, 7, 8, 9, 11},
		s{7},
		s{5, 8},
	},
	{
		s{2, 4, 6}, s{0, 2, 4, 6, 8},
		s{0, 2, 4, 6, 8},
		s{2, 4, 6},
		s{},
	},
	{
		s{0, 2, 4, 6, 8}, s{2, 4, 6},
		s{0, 2, 4, 6, 8},
		s{2, 4, 6},
		s{0, 8},
	},
	{
		s{10, 60, 90}, s{10, 90, 99},
		s{10, 60, 90, 99},
		s{10, 90},
		s{60},
	},
}

func denseFrom(us []uint) *Dense {
	d := NewDense(100)
	for _, u := range us {
		d.Add(u)
	}
	return d
}

func denseUints(d *Dense) []uint {
	var us []uint
	var u uint
	for u = 0; u < uint(d.Cap()); u++ {
		if d.Contains(u) {
			us = append(us, u)
		}
	}
	return us
}

func TestDenseBinaryFunctions(t *testing.T) {
	for _, test := range tests {
		d1 := denseFrom(test.s1)
		if got := d1.Elements(); !cmp.Equal(got, test.s1) {
			t.Errorf("got %v, want %v", got, test.s1)
		}
		d2 := denseFrom(test.s2)

		got := d1.Copy()
		got.AddIn(d2)
		if !got.Equal(denseFrom(test.union)) {
			t.Errorf("%v union %v: got %v, want %v", test.s1, test.s2, denseUints(got), test.union)
		}

		got = d1.Copy()
		got.RemoveNotIn(d2)
		if !got.Equal(denseFrom(test.intersection)) {
			t.Errorf("%v intersection %v: got %v, want %v", test.s1, test.s2, denseUints(got), test.intersection)
		}

		got = d1.Copy()
		got.RemoveIn(d2)
		if !got.Equal(denseFrom(test.difference)) {
			t.Errorf("%v difference %v: got %v, want %v", test.s1, test.s2, denseUints(got), test.difference)
		}
	}
}
