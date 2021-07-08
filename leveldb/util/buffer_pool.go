// Copyright (c) 2014, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package util

import "sync"

// BufferPool is a 'buffer pool'.
type BufferPool struct {
	pools [6]sync.Pool

	baseline  [4]int
	baseline0 int
}

func (p *BufferPool) Get(n int) []byte {
	if p == nil {
		return make([]byte, n)
	}
	poolNum := p.poolNum(n)
	val := p.pools[poolNum].Get().(*[]byte)
	if cap(*val) > n {
		return (*val)[:n]
	}
	// slice too small, return to pool and allocate proper
	p.pools[poolNum].Put(val)
	return make([]byte, n)
}

func (p *BufferPool) Put(b []byte) {
	if p == nil {
		return
	}
	poolNum := p.poolNum(cap(b))
	p.pools[poolNum].Put(&b)
}

func (p *BufferPool) Close() {
	p = nil
}

func NewBufferPool(baseline int) *BufferPool {
	if baseline <= 0 {
		panic("baseline can't be <= 0")
	}
	p := &BufferPool{
		baseline0: baseline,
		baseline:  [...]int{baseline / 4, baseline / 2, baseline * 2, baseline * 4},
	}
	for i, cap := range []int{2, 2, 4, 4, 2, 1} {
		p.pools[i] = sync.Pool{
			New: func() interface{} {
				b := make([]byte, cap)
				return &b
			},
		}
	}
	return p
}

func (p *BufferPool) poolNum(n int) int {
	if n <= p.baseline0 && n > p.baseline0/2 {
		return 0
	}
	for i, x := range p.baseline {
		if n <= x {
			return i + 1
		}
	}
	return len(p.baseline) + 1
}
