package bitset

import "fmt"

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func (s *Sparse) dump() {
	s.root.dump(0)
}

func (n *node) dump(level int) {
	indent(level)
	if n == nil {
		fmt.Println("nil")
	} else {
		fmt.Printf("shift %d, bitset %s\n", n.shift, n.bitset)
		for i, s := range n.subnodes {
			indent(level)
			fmt.Printf("%d: index %d\n", i, s.index)
			switch b := s.sub.(type) {
			case *node:
				b.dump(level + 1)
			case *set256:
				indent(level + 1)
				fmt.Printf("%s\n", b)
			}
		}
	}
}

func indent(level int) {
	for i := 0; i < level; i++ {
		fmt.Print("  ")
	}
}
