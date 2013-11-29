package main

import (
	"fmt"
	"tree"
)

func main() {
	t := tree.InitTreeSize(2, 3)
	for i := 1; i < 5; i++ {
		r := &tree.Record{
			Key: []byte(fmt.Sprintf("k_%d", i)),
			Val: []byte(fmt.Sprintf("v_%d", i)),
		}
		t.Add(r)
		t.TreeStats()
	}
	t.TreeInfo()
}
