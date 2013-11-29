// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree_test

import (
	"../tree"
	"strconv"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 100000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		go func() {
			if ok := root.Add(r); !ok {
				t.Fatal("Insert Failed", i)
			}
		}()
	}
}

func TestGet(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 100000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		go root.Add(r)
	}
	time.Sleep(time.Second * 10)
	for i := 0; i < size; i++ {
		if string(root.Get([]byte(strconv.Itoa(i)))) != strconv.Itoa(i) {
			t.Fatal("Find Failed", i)
		}
	}
}

func TestSet(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 100000
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
			t.Fatal("Update Failed", i)
		}
	}
	for i := 0; i < size; i++ {
		if string(root.Get([]byte(strconv.Itoa(i)))) != strconv.Itoa(i+1) {
			t.Fatal("Find Failed", i)
		}
	}
}
func TestDel(t *testing.T) {
	root := tree.InitTreeSize(2, 3)
	size := 100000
	for i := 0; i < size; i++ {
		r := &tree.Record{
			[]byte(strconv.Itoa(i)),
			[]byte(strconv.Itoa(i)),
		}
		root.Add(r)
	}
	for i := 0; i < size; i++ {
		if ok := root.Del([]byte(strconv.Itoa(i))); !ok {
			t.Fatal("delete Failed", i)
		}
		if root.Get([]byte(strconv.Itoa(i))) != nil {
			t.Fatal("Find Failed", i)
		}
	}
}
