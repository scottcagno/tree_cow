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

// locate leaf
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

// clone leaf
func (self *Leaf) clone(t *Tree) TreeNode {
	leaf := t.initLeaf()
	leaf.Keys = make([][]byte, len(self.Keys))
	copy(leaf.Keys, self.Keys)
	leaf.Vals = make([][]byte, len(self.Vals))
	copy(leaf.Vals, self.Vals)
	return leaf
}

// split leaf
func (self *Leaf) split(t *Tree) (k []byte, l, r int32) {
	leaf := t.initLeaf()
	mid := t.GetLeafMax() / 2
	leaf.Vals = make([][]byte, len(self.Vals[mid:]))
	copy(leaf.Vals, self.Vals[mid:])
	leaf.Keys = make([][]byte, len(self.Keys[mid:]))
	self.Vals = self.Vals[mid:]
	copy(self.Keys, self.Keys[mid:])
	self.Keys = self.Keys[:mid]
	l, r, k = self.GetId(), leaf.GetId(), leaf.Keys[0]
	return
}

// add leaf
func (self *Leaf) add(r *Record, t *Tree) (bool, TreeNode) {
	idx := self.locate(r.Key)
	if idx > 0 {
		if bytes.Compare(self.Keys[idx-1], r.Key) == 0 {
			return false, nil
		}
	}
	var clonedLeaf *Leaf
	if len(self.Keys) == 0 {
		clonedLeaf = self
	} else {
		clonedLeaf, _ = self.clone(t).(*Leaf)
		t.markDup(self.GetId())
	}
	clonedLeaf.Keys = append(clonedLeaf.Keys[:idx], append([][]byte{r.Key}, clonedLeaf.Keys[idx:]...)...)
	clonedLeaf.Vals = append(clonedLeaf.Vals[:idx], append([][]byte{r.Val}, clonedLeaf.Vals[idx:]...)...)
	return true, clonedLeaf
}

// delete leaf
func (self *Leaf) del(k []byte, t *Tree) (bool, TreeNode, []byte) {
	deleted := false
	idx := self.locate(k) - 1
	if idx >= 0 {
		if bytes.Compare(self.Keys[idx], k) == 0 {
			deleted = true
		}
	}
	if deleted {
		clonedLeaf, _ := self.clone(t).(*Leaf)
		clonedLeaf.Keys = append(clonedLeaf.Keys[:idx], clonedLeaf.Keys[idx+1:]...)
		clonedLeaf.Vals = append(clonedLeaf.Vals[:idx], clonedLeaf.Vals[idx+1:]...)
		if self.GetId() == t.GetRoot() {
			t.cloneroot = clonedLeaf.GetId()
		}
		t.markDup(self.GetId())
		if idx == 0 && len(clonedLeaf.Keys) > 0 {
			return true, clonedLeaf, clonedLeaf.Keys[0]
		}
		return true, clonedLeaf, nil
	}
	return false, nil, nil
}

// get leaf
func (self *Leaf) get(k []byte, t *Tree) []byte {
	idx := self.locate(k) - 1
	if idx >= 0 {
		if bytes.Compare(self.Keys[idx], k) == 0 {
			return self.Vals[idx]
		}
	}
	return nil
}

// set/update leaf
func (self *Leaf) set(r *Record, t *Tree) (bool, TreeNode) {
	idx := self.locate(r.Key) - 1
	if idx >= 0 {
		if bytes.Compare(self.Keys[idx], r.Key) == 0 {
			clonedLeaf, _ := self.clone(t).(*Leaf)
			clonedLeaf.Vals[idx] = r.Val
			t.markDup(self.GetId())
			return true, clonedLeaf
		}
	}
	return false, nil
}
