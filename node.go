package bitset

import "fmt"

// A node is a compact radix tree element.
// It behaves like a 256-element array of subnodes, indexed by one byte of the
// element. In fact, only the non-empty subnodes are represented; the bitset
// field stores this set and the subnodes field contains the non-empty subnodes
// in order.
type node struct {
	shift    uint // how many bits to shift elements right
	bitset   set256
	subnodes []subnode // if shift > 0
}

type subnode struct {
	index uint8 // the index in the full 256-element array
	sub   subber
}

// subber is the interface satisifed by nodes of the tree.
// It is implemented by node, for interior nodes, and set256, for leaves.
type subber interface {
	add64(uint64)
	remove64(uint64) bool // return true if empty
	contains64(uint64) bool
	elements(func([]uint64) bool, uint64) bool
	//elements64high(a []uint64, start, high uint64) int
	len() int
	memSize() uint64
	equal(subber) bool
	copy() subber
	addIn(subber)
	removeIn(subber) bool
	removeNotIn(subber) bool
}

func (n *node) newSubber() subber {
	if n.shift == 8 {
		return &set256{}
	} else {
		return &node{shift: n.shift - 8}
	}
}

func (n *node) copyNode() *node {
	if n == nil {
		return nil
	}
	n2 := *n
	for i, sn := range n2.subnodes {
		n2.subnodes[i].sub = sn.sub.copy()
	}
	return &n2
}

func (n *node) copy() subber { return n.copyNode() }

func (n *node) add64(e uint64) {
	index := uint8(e >> n.shift)
	pos, found := n.bitset.position(index)
	var sub subber
	if found {
		sub = n.subnodes[pos].sub
	} else {
		sub = n.newSubber()
		n.insertSubnode(pos, subnode{index: index, sub: sub})
	}
	sub.add64(e)
}

func (n *node) remove64(e uint64) (empty bool) {
	// assert node is not empty
	index := uint8(e >> n.shift)
	pos, found := n.bitset.position(index)
	if !found {
		return false // we weren't empty coming in
	}
	assert(n.subnodes[pos].index == index)
	sub := n.subnodes[pos].sub
	if sub.remove64(e) {
		if len(n.subnodes) == 1 {
			// No need to clean up, we're finished.
			return true
		}
		n.deleteSubnode(pos)
	}
	return false
}

func (n *node) insertSubnode(pos int, sn subnode) {
	newsubs := make([]subnode, len(n.subnodes)+1)
	copy(newsubs, n.subnodes[:pos])
	newsubs[pos] = sn
	copy(newsubs[pos+1:], n.subnodes[pos:])
	n.subnodes = newsubs
	n.bitset.add(sn.index)
}

func (n *node) deleteSubnode(pos int) {
	// TODO: really shrink memory
	index := n.subnodes[pos].index
	copy(n.subnodes[pos:], n.subnodes[pos+1:])
	n.subnodes = n.subnodes[:len(n.subnodes)-1]
	n.bitset.remove(index)
}

func (n *node) contains64(e uint64) bool {
	index := uint8(e >> n.shift)
	p, found := n.bitset.position(index)
	if !found {
		return false
	}
	return n.subnodes[p].sub.contains64(e)
}

func (n1 *node) equal(sub subber) bool {
	n2 := sub.(*node)
	if !n1.bitset.equal(&n2.bitset) {
		fmt.Printf("bitsets unequal: %s, %s\n", n1.bitset, n2.bitset)
		return false
	}
	for i, sn1 := range n1.subnodes {
		if !sn1.sub.equal(n2.subnodes[i].sub) {
			return false
		}
	}
	return true
}

func (n *node) len() int {
	t := 0
	for _, s := range n.subnodes {
		t += s.sub.len()
	}
	return t
}

func (n *node) memSize() uint64 {
	sz := memSize(*n)
	for _, s := range n.subnodes {
		sz += memSize(s)
		sz += s.sub.memSize()
	}
	return sz
}

func (n *node) elements(f func([]uint64) bool, offset uint64) bool {
	for _, sn := range n.subnodes {
		if !sn.sub.elements(f, offset+uint64(sn.index)<<n.shift) {
			return false
		}
	}
	return true
}

func (n1 *node) addIn(s subber) {
	n2 := s.(*node)
	assert(n1.shift == n2.shift)
	// Merge the lists of subnodes.
	i1 := 0
	i2 := 0
	for i1 < len(n1.subnodes) && i2 < len(n2.subnodes) {
		sn1 := n1.subnodes[i1]
		sn2 := n2.subnodes[i2]
		switch {
		case sn1.index < sn2.index:
			// n1 has a chunk of elements that n2 does not.
			// Skip over it.
			i1++

		case sn1.index > sn2.index:
			// n2 has elements that n1 does not. Add a subnode to n1
			// that is a copy of n2's subnode.
			n1.insertSubnode(i1, subnode{index: sn2.index, sub: sn2.sub.copy()})
			i1++
			i2++

		default:
			// sn1 and sn2 have the same index. Merge their contents.
			sn1.sub.addIn(sn2.sub)
			i1++
			i2++
		}
	}
	// If there are more n2 subnodes, copy them in.
	for i2 < len(n2.subnodes) {
		sn2 := n2.subnodes[i2]
		n1.insertSubnode(i1, subnode{index: sn2.index, sub: sn2.sub.copy()})
		i1++
		i2++
	}
}

func (n1 *node) removeIn(s subber) (empty bool) {
	n2 := s.(*node)
	assert(n1.shift == n2.shift)
	i1 := 0
	i2 := 0
	removed := false
	for i1 < len(n1.subnodes) && i2 < len(n2.subnodes) {
		sn1 := n1.subnodes[i1]
		sn2 := n2.subnodes[i2]
		switch {
		case sn1.index < sn2.index:
			// n1 has a chunk of elements that n2 does not.
			// Skip over it.
			i1++

		case sn1.index > sn2.index:
			// n2 has elements that n1 does not. Skip.
			i2++

		default:
			// sn1 and sn2 have the same index.
			if sn1.sub.removeIn(sn2.sub) {
				n1.bitset.remove(sn1.index)
				removed = true
			}
			i1++
			i2++
		}
	}
	if n1.bitset.empty() {
		return true
	}
	if !removed {
		return false
	}
	n1.adjustSubnodes()
	return false
}

func (n *node) adjustSubnodes() {
	// Change subnodes to match bitset.
	sns := n.subnodes
	n.subnodes = n.subnodes[:0]
	for _, sn := range sns {
		if n.bitset.contains(sn.index) {
			n.subnodes = append(n.subnodes, sn)
		}
	}
}

func (n1 *node) removeNotIn(s subber) (empty bool) {
	n2 := s.(*node)
	assert(n1.shift == n2.shift)
	i1 := 0
	i2 := 0
	removed := false
	for i1 < len(n1.subnodes) && i2 < len(n2.subnodes) {
		sn1 := n1.subnodes[i1]
		sn2 := n2.subnodes[i2]
		switch {
		case sn1.index < sn2.index:
			// n1 has elements that n2 does not. Remove them.
			n1.bitset.remove(sn1.index)
			removed = true
			i1++

		case sn1.index > sn2.index:
			// n2 has elements that n1 does not. Skip.
			i2++

		default:
			// sn1 and sn2 have the same index.
			if sn1.sub.removeNotIn(sn2.sub) {
				n1.bitset.remove(sn1.index)
				removed = true
			}
			i1++
			i2++
		}
	}
	// If there are more n1 subnodes, remove them.
	for i1 < len(n1.subnodes) {
		n1.bitset.remove(n1.subnodes[i1].index)
		removed = true
		i1++
	}
	if n1.bitset.empty() {
		return true
	}
	if !removed {
		return false
	}
	n1.adjustSubnodes()
	return false
}
