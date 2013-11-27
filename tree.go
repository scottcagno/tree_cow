// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

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
		Size:      int32(TreeSize),
		LeafMax:   int32(LeafSize),
		NodeMax:   int32(NodeSize),
		LeafCount: int32(0),
		NodeCount: int32(0),
		Index:     int32(0),
	}
	self.Version = uint32(0),
	self.Root = int32(self.newLeaf().GetId())
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
				node := self.newNode()
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
	return self.nodes[self.GetRoot()].fnd(k, self)
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