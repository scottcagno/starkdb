// * 
// * Copyright, 2014 Scott Cagno. All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// * 
// * StarkDB - A distinct, clear, and complete DB.
// * 

package main

// node and leaf degrees
const kx = 128
const kd = 64

// compare a and b, return int
type Cmp func(a, b []byte) int

// "zero" values
var zd  d
var zde de
var zx  x
var zxe xe

// clear/reset
func clr(q interface{}) {
	switch x := q.(type) {
	case *x:
		for i := 0; i <= x.c; i++ {
			clr(x.x[i].ch)
		}
		*x = zx
	case *d:
		*x = zd
	}
}