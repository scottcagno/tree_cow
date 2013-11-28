// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

const (
	TreeSize = 1 << 10 // 1024
	NodeSize = 1 << 6  // 64
	LeafSize = 1 << 5  // 32
)

const (
	isNode = iota
	isLeaf
)

const (
	StateNormal = iota
	StateDump
	StateGC
)
