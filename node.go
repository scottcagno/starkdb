// * 
// * Copyright, 2014 Scott Cagno. All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// * 
// * StarkDB - A distinct, clear, and complete DB.
// * 

package main

// index page
type x struct {
	c int
	x [2*kx + 2]xe	
}

// index element
type xe struct {
	ch  interface{}
	sep *d
}
	
// init and return new index page
func newX(ch0 interface{}) *x {
	r := &x{}
	r.x[0].ch = ch0
	return r
}

// remove intex page at 'i'
func (q *x) extract(i int) {
	q.c--
	if i < q.c {
		copy(q.x[i:], q.x[i+1:q.c+1])
		q.x[q.c].ch = q.x[q.c+1].ch
		q.x[q.c].sep = nil
		q.x[q.c+1] = zxe
	}
}

// add index page at 'i', return newly added index page
func (q *x) insert(i int, d *d, ch interface{}) *x {
	c := q.c
	if i < c {
		q.x[c+1].ch = q.x[c].ch
		copy(q.x[i+2:], q.x[i+1:c])
		q.x[i+1].sep = q.x[i].sep
	}
	c++
	q.c = c
	q.x[i].sep = d
	q.x[i+1].ch = ch
	return q
}

// return index page 'i's siblings
func (q *x) siblings(i int) (l, r *d) {
	if i >= 0 {
		if i > 0 {
			l = q.x[i-1].ch.(*d)
		}
		if i < q.c {
			r = q.x[i+1].ch.(*d)
		}
	}
	return
}