// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

import (
	"fmt"
)

func init() {}

type TreeData struct {
	Root               *int32
	NodeCount, NodeMax *int32
	LeafCount, LeafMax *int32
	Free               []int32
	Size, Index        *int32
	Version            *uint32
	Other_             []byte
}

func (self *TreeData) Reset() {
	*self = TreeData{}
}

func (self *TreeData) String() string {
	return fmt.Sprintf("%v", self)
}

func (self *TreeData) GetRoot() int32 {
	if self != nil && self.Root != nil {
		return *self.Root
	}
	return 0
}

func (self *TreeData) GetNodeCount() int32 {
	if self != nil && self.NodeCount != nil {
		return *self.NodeCount
	}
	return 0
}

func (self *TreeData) GetNodeMax() int32 {
	if self != nil && self.NodeMax != nil {
		return *self.NodeMax
	}
	return 0
}

func (self *TreeData) GetLeafCount() int32 {
	if self != nil && self.LeafCount != nil {
		return *self.LeafCount
	}
	return 0
}

func (self *TreeData) GetLeafMax() int32 {
	if self != nil && self.LeafMax != nil {
		return *self.LeafMax
	}
	return 0
}

func (self *TreeData) GetFree() []int32 {
	if self != nil {
		return self.Free
	}
	return nil
}

func (self *TreeData) GetSize() int32 {
	if self != nil && self.Size != nil {
		return *self.Size
	}
	return 0
}

func (self *TreeData) GetIndex() int32 {
	if self != nil && self.Index != nil {
		return *self.Index
	}
	return 0
}

func (self *TreeData) GetVersion() uint32 {
	if self != nil && self.Version != nil {
		return *self.Version
	}
	return 0
}

type IndexData struct {
	Id      *int32
	Version *uint32
	Other_  []byte
}

func (self *IndexData) Reset() {
	*self = IndexData{}
}

func (self *IndexData) String() string {
	return fmt.Sprintf("%v", self)
}

func (self *IndexData) GetId() int32 {
	if self != nil && self.Id != nil {
		return *self.Id
	}
	return 0
}

func (self *IndexData) GetVersion() uint32 {
	if self != nil && self.Version != nil {
		return *self.Version
	}
	return 0
}

type NodeRecordData struct {
	Children []int32
	Keys     [][]byte
	Other_   []byte
}

func (self *NodeRecordData) Reset() {
	*self = NodeRecordData{}
}

func (self *NodeRecordData) String() string {
	return fmt.Sprintf("%v", self)
}

func (self *NodeRecordData) GetChildren() []int32 {
	if self != nil {
		return self.Children
	}
	return nil
}

func (self *NodeRecordData) GetKeys() [][]byte {
	if self != nil {
		return self.Keys
	}
	return nil
}

type LeafRecordData struct {
	Keys, Vals [][]byte
	Other_     []byte
}

func (self *LeafRecordData) Reset() {
	*self = LeafRecordData{}
}

func (self *LeafRecordData) String() string {
	return fmt.Sprintf("%v", self)
}

func (self *LeafRecordData) GetKeys() [][]byte {
	if self != nil {
		return self.Keys
	}
	return nil
}

func (self *LeafRecordData) GetVils() [][]byte {
	if self != nil {
		return self.Vals
	}
	return nil
}
