// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD License - sites.google.com/site/bsdc3license
// *
// * B-Tree :: balanced binary b-tree
// *

package tree

// b tree record
type Record struct {
	Key, Val []byte
}

// init new record
func InitRecord(k, v []byte) *Record {
	return &Record{k, v}
}
