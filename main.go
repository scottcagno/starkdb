// * 
// * Copyright, 2014 Scott Cagno. All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// * 
// * StarkDB - A distinct, clear, and complete DB.
// * 

package main

import (
	"bytes"
	"fmt"
)

// literal compare function
func cmp(a, b []byte) int {
	return bytes.Compare(a, b)
}

// main function
func main() {
	t := InitTree(cmp)
	t.Set([]byte{1}, []byte("hello, world!"))
	v, ok := t.Get([]byte{1})
	fmt.Printf("%s, (%v)[%d]\n", v, ok, t.ver)
}