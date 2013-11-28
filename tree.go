// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

import (
	"sync"
	"sync/atomic"
	"time"
)

// b tree data struct
type Tree struct {
	TreeData
	sync.RWMutex
	nodes            []TreeNode
	state, cloneroot int32
	dupnodelist      []int32
}

// init new tree (with default values)
func InitTree() *Tree {
	self := new(Tree)
	self.nodes = make([]TreeNode, TreeSize)
	self.TreeData = TreeData{
		Size:    TreeSize,
		LeafMax: *LeafSize,
		NodeMax: *NodeSize,
		//	LeafCount: **int32(0),
		//	NodeCount: **int32(0),
		//	Index:     **int32(0),
	}
	//self.Version = uint32(0)
	self.Root = int32(self.initLeaf().GetId())
	self.state = StateNormal
	return self
}

// insert record into tree
func (self *Tree) Add(r *Record) bool {
	self.Lock()
	defer self.Unlock()
	*self.Version++
	ok, clonedNode := self.nodes[self.GetRoot()].add(r, self)
	if ok {
		if len(clonedNode.GetKeys()) > int(self.GetNodeMax()) {
			node := self.initNode()
			k, l, r := clonedNode.split(self)
			self.Root = int32(n.GetId())
			self.nodes[int(self.GetRoot())] = node
		} else {
			self.Root = int32(clonedNode.GetId())
		}
	} else {
		*self.Version--
	}
	return ok
}

// update a record in tree
func (self *Tree) Set(r *Record) bool {
	self.Lock()
	defer self.Unlock()
	*self.Version++
	ok, clonedNode := self.nodes[self.GetRoot()].set(r, self)
	if ok {
		self.markDup(self.GetRoot())
		self.Root = int32(clonedNode.GetId())
	} else {
		*self.Version--
	}
	return ok
}

// return value from record in tree
func (self *Tree) Get(k []byte) []byte {
	self.Lock()
	defer self.Unlock()
	return self.nodes[self.GetRoot()].get(k, self)
}

// delete record from tree
func (self *Tree) Del(k []byte) bool {
	self.Lock()
	defer self.Unlock()
	*self.Version++
	ok, clonedNode, _ := self.nodes[self.GetRoot()].del(k, self)
	if ok {
		if len(clonedNode.GetKeys()) == 0 {
			if node, ok := clonedNode.(*Node); ok {
				self.Root = int32(self.nodes[clonedNode.Clildren[0]].GetId())
				self.markDup(node.GetId())
			} else {
				self.Root = int32(clonedNode.GetId())
			}
		} else {
			self.Root = int32(clonedNode.GetId())
		}
	} else {
		*self.Version--
	}
	return ok
}

func (self *Tree) initId() int32 {
	var id int32
	if len(self.Free) > 0 {
		id = self.Free[len(self.Free)-1]
	} else {
		if self.GetIndex() >= self.GetSize() {
			self.nodes = append(self.nodes, make([]TreeNode, TreeSize)...)
			*self.Size += int32(TreeSize)
		}
		id = self.GetIndex()
		*self.Index++
	}
	return id
}

func (self *Tree) initNode() *Node {
	*self.NodeCount++
	id := self.initId()
	node := &Node{
		IndexData: IndexData{
			Id:      int32(id),
			Version: int32(self.GetVersion()),
		},
	}
	self.nodes[id] = node
	return node
}

func (self *Tree) getNode(id int32) *Node {
	if node, ok := self.nodes[id].(*Node); ok {
		return node
	}
	return nil
}

func (self *Tree) initLeaf() *Leaf {
	*self.LeafCount++
	id := self.initId()
	leaf := &Leaf{
		IndexData: IndexData{
			Id:      int32(id),
			Version: int32(self.GetVersion()),
		},
	}
	self.nodes[id] = leaf
	return leaf
}
func (self *Tree) getLeaf(id int32) *Leaf {
	if leaf, ok := self.nodes[id].(*Leaf); ok {
		return leaf
	}
	return nil
}

func (self *Tree) markDup(id int32) {
	self.dupnodelist = append(self.dupnodelist, id)
}

func (self *Tree) gc() {
	for {
		self.Lock()
		if atomic.CompareAndSwapInt32(&self.state, StateNormal, StateGC) {
			if len(self.dupnodelist) > 0 {
				id := self.dupnodelist[len(self.dupnodelist)-1]
				switch self.nodes[id].(type) {
				case *Node:
					*self.NodeCount--
				case *Leaf:
					*self.LeafCount--
				default:
					atomic.CompareAndSwapInt32(&self.state, StateGC, StateNormal)
					continue
				}
				self.Free = append(self.Free, id)
				self.dupnodelist = self.dupnodelist[:len(self.dupnodelist)-1]
				atomic.CompareAndSwapInt32(&self.state, StateGC, StateNormal)
			}
		} else {
			time.Sleep(time.Second)
		}
		self.Unlock()
	}
}
