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

func (self *Node) clone(t *Tree) TreeNode {
	node := t.initNode()
	node.Keys = make([][]byte, len(self.Keys))
	copy(node.Keys, self.Keys)
	node.Vals = make([][]byte, len(self.Vals))
	copy(node.Vals, self.Vals)
	return node
}

func (self *Node) split(t *Tree) (k []byte, l, r int32) {
	node := t.initNode()
	mid := t.GetNodeMax() / 2
	k = self.Keys[mid]
	node.Keys = make([][]byte, len(self.Keys[mid+1:]))
	copy(node.Keys, self.Keys[mid+1:])
	self.Keys = self.Keys[:mid]
	node.Children = make([]int32, len(self.Children)[mid+1:])
	copy(node.Children, self.Children[mid+1:])
	self.Children = self.Children[:mid+1]
	l, r = self.GetId(), node.GetId()
	return
}

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

func (self *Node) addOnce(k []byte, lid, rid int32, t *Tree) {
	idx := self.locate(k)
	if len(self.Keys) == 0 {
		self.Children = append([]int32{lid}, rid)
	} else {
		self.Children = append(self.Children[:idx+1], append([]int32{rid}, self.Children[idx+1:]...)...)
	}
	self.Keys = append(self.Keys[:idx], append([][]byte{k}, self.Keys[idx:]...)...)
}
