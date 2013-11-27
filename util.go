// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

type TreeNode interface {
	add(r *Record, t *Tree) (bool, TreeNode)
	set(r *Record, t *Tree) (bool, TreeNode)
	get(k []byte, t *Tree) []byte
	del(k []byte, t *Tree) (bool, TreeNode, []byte)
	split(t *Tree) (k []byte, l, r int32)
	locate(k []byte) int
	GetId() int32
	GetKeys() [][]byte
}
