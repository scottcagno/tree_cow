// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree_test

import (
	"../tree"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 1000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		go func() {
			if ok := root.Add(r); !ok {
				t.Fatal("Add Failed", i)
			}
		}()
	}
}

func TestGet(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 1000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		go func() {
			if ok := root.Add(r); !ok {
				t.Fatal("Add Failed", i)
			}
		}()
	}
	time.Sleep(time.Second * 2)
	for i := 0; i < size; i++ {
		re := string(root.Get([]byte(strconv.Itoa(i))))
		x := strconv.Itoa(i)
		if re != x {
			fmt.Println("response...", re, "compare...", x)
			t.Fatal("Get Failed", i)
		}
	}
}

func TestSet(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 1000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		root.Add(r)
	}
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		if ok := root.Set(r); !ok {
			t.Fatal("Set Failed", i)
		}
	}
	for i := 0; i < size; i++ {
		if string(root.Get([]byte(strconv.Itoa(i)))) != strconv.Itoa(i+1) {
			t.Fatal("Get Failed", i)
		}
	}
}
func TestDel(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 1000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		root.Add(r)
	}
	for i := 0; i < size; i++ {
		if ok := root.Del([]byte(strconv.Itoa(i))); !ok {
			t.Fatal("Delete Failed", i)
		}
		if root.Get([]byte(strconv.Itoa(i))) != nil {
			t.Fatal("Get Failed", i)
		}
	}
}
