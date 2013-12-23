// * 
// * Copyright, 2014 Scott Cagno. All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// * 
// * StarkDB - A distinct, clear, and complete DB.
// * 

package main

// data page
type d struct {
	c int
	d [2*kd + 1]de
	n, p *d
}

// data element
type de struct {
	k, v []byte
}

// move data page 'c' elements left
func (l *d) mvL(r *d, c int) {
	copy(l.d[l.c:], r.d[:c])
	copy(r.d[:], r.d[c:r.c])
	l.c += c
	r.c -= c
}

// move data page 'c' elements right
func (l *d) mvR(r *d, c int) {
	copy(r.d[c:], r.d[:r.c])
	copy(r.d[:c], l.d[l.c-c:])
	r.c += c
	l.c -= c
}