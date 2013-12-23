// * 
// * Copyright, 2014 Scott Cagno. All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// * 
// * StarkDB - A distinct, clear, and complete DB.
// * 

package main

import (
	"io"
)

// b+tree
type Tree struct {
	c     int
	cmp   Cmp
	first, last *d
	r     interface{}
	ver   int64
}

// init and return a new b+tree
func InitTree(cmp Cmp) *Tree {
	return &Tree{
		cmp: cmp,
	}
}

// remove all key/val pairs from the tree
func (t *Tree) Clear() {
	if t.r == nil {
		return
	}
	clr(t.r)
	t.c, t.first, t.last, t.r = 0, nil, nil, nil
	t.ver++
}

// ?
func (t *Tree) cat(p *x, q, r *d, pi int) {
	t.ver++
	q.mvL(r, r.c)
	if r.n != nil {
		r.n.p = q
	} else {
		t.last = q
	}
	q.n = r.n
	if p.c > 1 {
		p.extract(pi)
		p.x[pi].ch = q
	} else {
		t.r = q
	}
}

// ?
func (t *Tree) catX(p, q, r *x, pi int) {
	t.ver++
	q.x[q.c].sep = p.x[pi].sep
	copy(q.x[q.c+1:], r.x[:r.c])
	q.c += r.c + 1
	q.x[q.c].ch = r.x[r.c].ch
	if p.c > 1 {
		p.c--
		pc := p.c
		if pi < pc {
			p.x[pi].sep = p.x[pi+1].sep
			copy(p.x[pi+1:], p.x[pi+2:pc+1])
			p.x[pc].ch = p.x[pc+1].ch
			p.x[pc].sep = nil
			p.x[pc+1].ch = nil
		}
		return
	}

	t.r = q
}

// remove k's key/val pair if it exists.
func (t *Tree) Delete(k []byte) (ok bool) {
	pi := -1
	var p *x
	q := t.r
	if q == nil {
		return
	}

	for {
		var i int
		i, ok = t.find(q, k)
		if ok {
			switch x := q.(type) {
			case *x:
				dp := x.x[i].sep
				switch {
				case dp.c > kd:
					t.extract(dp, 0)
				default:
					if x.c < kx && q != t.r {
						t.underflowX(p, &x, pi, &i)
					}
					pi = i + 1
					p = x
					q = x.x[pi].ch
					ok = false
					continue
				}
			case *d:
				t.extract(x, i)
				if x.c >= kd {
					return
				}

				if q != t.r {
					t.underflow(p, x, pi)
				} else if t.c == 0 {
					t.Clear()
				}
			}
			return
		}

		switch x := q.(type) {
		case *x:
			if x.c < kx && q != t.r {
				t.underflowX(p, &x, pi, &i)
			}
			pi = i
			p = x
			q = x.x[i].ch
		case *d:
			return
		}
	}
}

// remove data page at 'i'
func (t *Tree) extract(q *d, i int) {
	t.ver++
	q.c--
	if i < q.c {
		copy(q.d[i:], q.d[i+1:q.c+1])
	}
	q.d[q.c] = zde
	t.c--
	return
}

// locate index of data page
func (t *Tree) find(q interface{}, k []byte) (i int, ok bool) {
	var mk []byte
	l := 0
	switch x := q.(type) {
	case *x:
		h := x.c - 1
		for l <= h {
			m := (l + h) >> 1
			mk = x.x[m].sep.d[0].k
			switch cmp := t.cmp(k, mk); {
			case cmp > 0:
				l = m + 1
			case cmp == 0:
				return m, true
			default:
				h = m - 1
			}
		}
	case *d:
		h := x.c - 1
		for l <= h {
			m := (l + h) >> 1
			mk = x.d[m].k
			switch cmp := t.cmp(k, mk); {
			case cmp > 0:
				l = m + 1
			case cmp == 0:
				return m, true
			default:
				h = m - 1
			}
		}
	}
	return l, false
}

// return first item in key collating order, nil if empty
func (t *Tree) First() (k []byte, v []byte) {
	if q := t.first; q != nil {
		q := &q.d[0]
		k, v = q.k, q.v
	}
	return
}

// return val and true if it exists, else nil and false
func (t *Tree) Get(k []byte) (v []byte, ok bool) {
	q := t.r
	if q == nil {
		return
	}

	for {
		var i int
		if i, ok = t.find(q, k); ok {
			switch x := q.(type) {
			case *x:
				return x.x[i].sep.d[0].v, true
			case *d:
				return x.d[i].v, true
			}
		}
		switch x := q.(type) {
		case *x:
			q = x.x[i].ch
		default:
			return
		}
	}
}

// insert key and val
func (t *Tree) insert(q *d, i int, k []byte, v []byte) *d {
	t.ver++
	c := q.c
	if i < c {
		copy(q.d[i+1:], q.d[i:c])
	}
	c++
	q.c = c
	q.d[i].k, q.d[i].v = k, v
	t.c++
	return q
}

// return last item of the tree in key collating order
// else return nil, nil
func (t *Tree) Last() (k []byte, v []byte) {
	if q := t.last; q != nil {
		q := &q.d[q.c-1]
		k, v = q.k, q.v
	}
	return
}

// Len returns the number of items in the tree.
func (t *Tree) Len() int {
	return t.c
}

// handle overflow
func (t *Tree) overflow(p *x, q *d, pi, i int, k []byte, v []byte) {
	t.ver++
	l, r := p.siblings(pi)

	if l != nil && l.c < 2*kd {
		l.mvL(q, 1)
		t.insert(q, i-1, k, v)
		return
	}

	if r != nil && r.c < 2*kd {
		if i < 2*kd {
			q.mvR(r, 1)
			t.insert(q, i, k, v)
		} else {
			t.insert(r, 0, k, v)
		}
		return
	}

	t.split(p, q, pi, i, k, v)
}

// Seek returns an Enumerator positioned on a an item such that k >= item's
// key. ok reports if k == item.key The Enumerator's position is possibly
// after the last item in the tree.
func (t *Tree) Seek(k []byte) (e *Enumerator, ok bool) {
	q := t.r
	if q == nil {
		e = &Enumerator{nil, false, 0, k, nil, t, t.ver}
		return
	}

	for {
		var i int
		if i, ok = t.find(q, k); ok {
			switch x := q.(type) {
			case *x:
				e = &Enumerator{nil, ok, 0, k, x.x[i].sep, t, t.ver}
				return
			case *d:
				e = &Enumerator{nil, ok, i, k, x, t, t.ver}
				return
			}
		}
		switch x := q.(type) {
		case *x:
			q = x.x[i].ch
		case *d:
			e = &Enumerator{nil, ok, i, k, x, t, t.ver}
			return
		}
	}
}

// SeekFirst returns an enumerator positioned on the first KV pair in the tree,
// if any. For an empty tree, err == io.EOF is returned and e will be nil.
func (t *Tree) SeekFirst() (e *Enumerator, err error) {
	q := t.first
	if q == nil {
		return nil, io.EOF
	}

	return &Enumerator{nil, true, 0, q.d[0].k, q, t, t.ver}, nil
}

// SeekLast returns an enumerator positioned on the last KV pair in the tree,
// if any. For an empty tree, err == io.EOF is returned and e will be nil.
func (t *Tree) SeekLast() (e *Enumerator, err error) {
	q := t.last
	if q == nil {
		return nil, io.EOF
	}

	return &Enumerator{nil, true, q.c - 1, q.d[q.c-1].k, q, t, t.ver}, nil
}

// Set sets the value associated with k.
func (t *Tree) Set(k []byte, v []byte) {
	pi := -1
	var p *x
	q := t.r
	if q != nil {
		for {
			i, ok := t.find(q, k)
			if ok {
				switch x := q.(type) {
				case *x:
					x.x[i].sep.d[0].v = v
				case *d:
					x.d[i].v = v
				}
				return
			}

			switch x := q.(type) {
			case *x:
				if x.c > 2*kx {
					t.splitX(p, &x, pi, &i)
				}
				pi = i
				p = x
				q = x.x[i].ch
			case *d:
				switch {
				case x.c < 2*kd:
					t.insert(x, i, k, v)
				default:
					t.overflow(p, x, pi, i, k, v)
				}
				return
			}
		}
	}

	z := t.insert(&d{}, 0, k, v)
	t.r, t.first, t.last = z, z, z
	return
}

func (t *Tree) split(p *x, q *d, pi, i int, k []byte, v []byte) {
	t.ver++
	r := &d{}
	if q.n != nil {
		r.n = q.n
		r.n.p = r
	} else {
		t.last = r
	}
	q.n = r
	r.p = q

	copy(r.d[:], q.d[kd:2*kd])
	for i := range q.d[kd:] {
		q.d[kd+i] = zde
	}
	q.c = kd
	r.c = kd
	if pi >= 0 {
		p.insert(pi, r, r)
	} else {
		t.r = newX(q).insert(0, r, r)
	}
	if i > kd {
		t.insert(r, i-kd, k, v)
		return
	}

	t.insert(q, i, k, v)
}

func (t *Tree) splitX(p *x, pp **x, pi int, i *int) {
	t.ver++
	q := *pp
	r := &x{}
	copy(r.x[:], q.x[kx+1:])
	q.c = kx
	r.c = kx
	if pi >= 0 {
		p.insert(pi, q.x[kx].sep, r)
	} else {
		t.r = newX(q).insert(0, q.x[kx].sep, r)
	}
	q.x[kx].sep = nil
	for i := range q.x[kx+1:] {
		q.x[kx+i+1] = zxe
	}
	if *i > kx {
		*pp = r
		*i -= kx + 1
	}
}

func (t *Tree) underflow(p *x, q *d, pi int) {
	t.ver++
	l, r := p.siblings(pi)

	if l != nil && l.c+q.c >= 2*kd {
		l.mvR(q, 1)
	} else if r != nil && q.c+r.c >= 2*kd {
		q.mvL(r, 1)
		r.d[r.c] = zde // GC
	} else if l != nil {
		t.cat(p, l, q, pi-1)
	} else {
		t.cat(p, q, r, pi)
	}
}

func (t *Tree) underflowX(p *x, pp **x, pi int, i *int) {
	t.ver++
	var l, r *x
	q := *pp

	if pi >= 0 {
		if pi > 0 {
			l = p.x[pi-1].ch.(*x)
		}
		if pi < p.c {
			r = p.x[pi+1].ch.(*x)
		}
	}

	if l != nil && l.c > kx {
		q.x[q.c+1].ch = q.x[q.c].ch
		copy(q.x[1:], q.x[:q.c])
		q.x[0].ch = l.x[l.c].ch
		q.x[0].sep = p.x[pi-1].sep
		q.c++
		*i++
		l.c--
		p.x[pi-1].sep = l.x[l.c].sep
		return
	}

	if r != nil && r.c > kx {
		q.x[q.c].sep = p.x[pi].sep
		q.c++
		q.x[q.c].ch = r.x[0].ch
		p.x[pi].sep = r.x[0].sep
		copy(r.x[:], r.x[1:r.c])
		r.c--
		rc := r.c
		r.x[rc].ch = r.x[rc+1].ch
		r.x[rc].sep = nil
		r.x[rc+1].ch = nil
		return
	}

	if l != nil {
		*i += l.c + 1
		t.catX(p, l, q, pi-1)
		*pp = l
		return
	}

	t.catX(p, q, r, pi)
}