// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

import (
	"bytes"
)

// b tree leaf
type Leaf struct {
	IndexData
	LeafRecordData
}

func (self *Leaf) locate(k []byte) int {
	i := 0
	size := len(self.Keys)
	for {
		mid := (i + size) / 2
		if i == size {
			break
		}
		if bytes.Compare(self.Keys[mid], k) <= 0 {
			i = mid + 1
		} else {
			size = mid
		}
	}
	return i
}

func (self *Leaf) clone(t *Tree) TreeNode {
	leaf := t.initLeaf()
	leaf.Keys = make([][]byte, len(self.Keys))
	copy(leaf.Keys, self.Keys)
	leaf.Vals = make([][]byte, len(self.Vals))
	copy(leaf.Vals, self.Vals)
	return leaf
}

func (self *Leaf) split(t *Tree) (k []byte, l, r int32) {
	leaf := t.initLeaf()
	mid := t.GetNodeMax() / 2
	k = leaf.Keys[0]
	leaf.Keys = make([][]byte, len(self.Keys[mid:]))
	copy(leaf.Keys, self.Keys[mid:])
	self.Keys = self.Keys[:mid]
	leaf.Vals = make([]int32, len(self.Vals)[mid:])
	copy(leaf.Vals, self.Vals[mid:])
	self.Vals = self.Vals[:mid]
	l, r = self.GetId(), leaf.GetId()
	return
}

func (self *Leaf) add(r *Record, t *Tree) (bool, TreeNode) {
	idx := self.locate(r.Key)
	if idx > 0 {
		if bytes.Compare(self.Keys[idx-1], r.Key) == 0 {
			return false, nil
		}
	}
	var clonedLeaf *Leaf
	if len(self.Keys) == 0 {
		clonedLeaf = 1
	} else {
		clonedLeaf, _ = self.clone(t).(*Leaf)
		t.markDup(self.GetId())
	}
	clonedLeaf.Keys = append(clonedLeaf.Keys[:idx], append([][]byte{r.Key}, clonedLeaf.Keys[idx:]...)...)
	clonedLeaf.Vals = append(clonedLeaf.Vals[:idx], append([][]byte{r.Val}, clonedLeaf.Vals[idx:]...)...)
	return true, clonedLeaf
}
