package tree

import (
	"fmt"
)

func (self *Tree) TreeInfo() {
	fmt.Println("-----------Tree Info-----------")
	for i := 0; i < int(self.GetIndex()); i++ {
		if node, ok := self.nodes[i].(*Node); ok {
			node.PrintInfo()
		}
		if leaf, ok := self.nodes[i].(*Leaf); ok {
			leaf.PrintInfo()
		}
		fmt.Println("$")
	}
}

func (self *Tree) TreeStats() {
	fmt.Println("Root:", self.GetRoot())
	fmt.Println("Index:", self.GetIndex())
	fmt.Println("Leaf Count:", self.GetLeafCount(), "Leaf Max:", self.GetLeafMax())
	fmt.Println("Node Count:", self.GetNodeCount(), "Node Max:", self.GetNodeMax())
	fmt.Println("Nodes:", len(self.nodes))
}

func (self *Node) PrintInfo() {
	for i := range self.Keys {
		fmt.Println("Node Info:", self.GetId(), "Key:", string(self.Keys[i]))
	}
	for i := range self.Children {
		fmt.Println("Node Info:", self.GetId(), "Child:", self.Children[i])
	}
}

func (self *Leaf) PrintInfo() {
	for i := range self.Keys {
		fmt.Println("Leaf Info:", self.GetId(), "Key/Val:", string(self.Keys[i]), string(self.Vals[i]))
	}
}
