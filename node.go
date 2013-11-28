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

// b tree node
type Node struct {
	IndexData
	NodeRecordData
}

// locate node
func (self *Node) locate(k []byte) int {
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

// clone node
func (self *Node) clone(t *Tree) TreeNode {
	node := t.initNode()
	node.Keys = make([][]byte, len(self.Keys))
	copy(node.Keys, self.Keys)
	node.Children = make([]int32, len(self.Children))
	copy(node.Children, self.Children)
	return node
}

// split node
func (self *Node) split(t *Tree) (k []byte, l, r int32) {
	node := t.initNode()
	mid := t.GetNodeMax() / 2
	k = self.Keys[mid]
	node.Keys = make([][]byte, len(self.Keys[mid+1:]))
	copy(node.Keys, self.Keys[mid+1:])
	self.Keys = self.Keys[:mid]
	node.Children = make([]int32, len(self.Children[mid+1:]))
	copy(node.Children, self.Children[mid+1:])
	self.Children = self.Children[:mid+1]
	l, r = self.GetId(), node.GetId()
	return
}

// merge node
func (self *Node) mergeNode(lid, rid int32, idx int, t *Tree) int32 {
	l, r := t.getNode(lid), t.getNode(rid)
	if (len(l.Keys) + len(r.Keys)) > int(t.GetNodeMax()) {
		return -1
	}
	lclone := l.clone(t).(*Node)
	id := lclone.GetId()
	t.nodes[id] = lclone
	self.Children[idx] = id
	lclone.Keys = append(lclone.Keys, append([][]byte{self.Keys[idx]}, r.Keys...)...)
	lclone.Children = append(lclone.Children, r.Children...)
	self.Keys = append(self.Keys[:idx], self.Keys[idx+1:]...)
	self.Children = append(self.Children[:idx+1], self.Children[idx+2:]...)
	if len(lclone.Keys) > int(t.GetNodeMax()) {
		k, l, r := lclone.split(t)
		self.addOnce(k, l, r, t)
	}
	t.markDup(lid)
	t.markDup(rid)
	return lclone.GetId()
}

// merge leaf
func (self *Node) mergeLeaf(lid, rid int32, idx int, t *Tree) int32 {
	l, r := t.getLeaf(lid), t.getLeaf(rid)
	if (len(l.Vals) + len(r.Vals)) > int(t.GetLeafMax()) {
		return -1
	}
	lclone := l.clone(t).(*Leaf)
	id := lclone.GetId()
	t.nodes[id] = lclone
	self.Children[idx] = id
	if idx == len(self.Keys) {
		self.Children = self.Children[:idx]
		self.Keys = self.Keys[:idx-1]
	} else {
		self.Children = append(self.Children[:idx+1], self.Children[idx+2:]...)
		self.Keys = append(self.Keys[:idx], self.Keys[idx+1:]...)
	}
	lclone.Vals = append(lclone.Vals, r.Vals...)
	lclone.Keys = append(lclone.Keys, r.Keys...)
	t.markDup(lid)
	t.markDup(rid)
	return lclone.GetId()
}

// insert node
func (self *Node) add(r *Record, t *Tree) (bool, TreeNode) {
	idx := self.locate(r.Key)
	if ok, clonedTreeNode := t.nodes[self.Children[idx]].add(r, t); ok {
		clonedNode, _ := self.clone(t).(*Node)
		clonedNode.Children[idx] = clonedTreeNode.GetId()
		if len(clonedTreeNode.GetKeys()) > int(t.GetNodeMax()) {
			k, l, r := clonedTreeNode.split(t)
			clonedNode.addOnce(k, l, r, t)
		}
		t.markDup(self.GetId())
		return true, clonedNode
	}
	return false, nil
}

// insert key into tree node
func (self *Node) addOnce(k []byte, lid, rid int32, t *Tree) {
	idx := self.locate(k)
	if len(self.Keys) == 0 {
		self.Children = append([]int32{lid}, rid)
	} else {
		self.Children = append(self.Children[:idx+1], append([]int32{rid}, self.Children[idx+1:]...)...)
	}
	self.Keys = append(self.Keys[:idx], append([][]byte{k}, self.Keys[idx:]...)...)
}

// delete node
func (self *Node) del(k []byte, t *Tree) (bool, TreeNode, []byte) {
	idx := self.locate(k)
	if ok, clonedTreeNode, newk := t.nodes[self.Children[idx]].del(k, t); ok {
		clonedNode, _ := self.clone(t).(*Node)
		clonedNode.Children[idx] = clonedTreeNode.GetId()
		if self.GetId() == t.GetRoot() {
			t.cloneroot = clonedNode.GetId()
		}
		tmpk := newk
		if newk != nil {
			if clonedNode.replace(k, newk) {
				newk = nil
			}
		}
		if idx == 0 {
			idx = 1
		}
		if len(clonedNode.Keys) > 0 {
			var l int32
			if t.getLeaf(clonedNode.Children[idx-1]) != nil {
				l = clonedNode.mergeLeaf(clonedNode.Children[idx-1], clonedNode.Children[idx], idx-1, t)
				if idx == 1 && tmpk == nil {
					leaf := t.getLeaf(clonedNode.Children[0])
					if leaf != nil && len(leaf.Keys) > 0 {
						newk = leaf.Keys[0]
					}
				}
			} else {
				l = clonedNode.mergeNode(clonedNode.Children[idx-1], clonedNode.Children[idx], idx, t)
			}
			if l > 0 {
				clonedNode.Children[idx-1] = l
			}
		}
		t.markDup(self.GetId())
		return true, clonedNode, newk
	}
	return false, nil, nil
}

// replace node
func (self *Node) replace(oldk, newk []byte) bool {
	idx := self.locate(oldk) - 1
	if idx >= 0 {
		if bytes.Compare(self.Keys[idx], oldk) == 0 {
			self.Keys[idx] = newk
			return true
		}
	}
	return false
}

// get record
func (self *Node) get(k []byte, t *Tree) []byte {
	idx := self.locate(k)
	return t.nodes[self.Children[idx]].get(k, t)
}

// set/update node
func (self *Node) set(r *Record, t *Tree) (bool, TreeNode) {
	idx := self.locate(r.Key)
	if ok, clonedTreeNode := t.nodes[self.Children[idx]].set(r, t); ok {
		clonedNode, _ := self.clone(t).(*Node)
		id := clonedTreeNode.GetId()
		clonedNode.Children[idx] = id
		t.nodes[id] = clonedTreeNode
		t.markDup(self.GetId())
		return true, clonedNode
	}
	return false, nil
}
