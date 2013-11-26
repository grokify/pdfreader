// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// string math
package strm

import (
	"math/big"
)

func operand(s string) (r int64, f int) {
	if len(s) < 1 {
		return 0, 1
	}
	sig := s[0] == '-'
	p := 0
	if sig {
		p++
	}
	for p < len(s) {
		if s[p] == '.' {
			f = 1
		} else {
			f *= 10
			r *= 10
			r += int64(s[p] - '0')
		}
		p++
	}
	if sig {
		r = -r
	}
	if f == 0 {
		f = 1
	}
	return
}

func Int64(s string, f int) int64 {
	ra, fa := operand(s)
	for fa < f {
		fa *= 10
		ra *= 10
	}
	return ra / int64(fa/f)
}

func Int(s string, f int) int { return int(Int64(s, f)) }

func twop(a, b string) (ra, rb int64, f int) {
	ra, f = operand(a)
	rb, fb := operand(b)
	for fb < f {
		fb *= 10
		rb *= 10
	}
	for f < fb {
		f *= 10
		ra *= 10
	}
	return
}

func String(a int64, f int) string {
	buf := make([]byte, 128)
	p := 0
	if a < 0 {
		buf[p] = '-'
		p++
		a = -a
	}
	var fu func(c int64)
	step := 1
	fu = func(c int64) {
		s := step
		step *= 10
		if c > 9 || step <= f {
			fu(c / 10)
		}
		buf[p] = '0' + byte(c%10)
		p++
		if f == s && f != 1 {
			buf[p] = '.'
			p++
		}
	}
	fu(a)
	return string(buf[0:p])
}

func Mul(a, b string) string {
	ra, _, f := twop(a, b)
	ar := big.NewRat(ra, int64(f))
	// br := big.NewRat(rb, int64(f))
	i := ar.Num()
	n := ar.Denom()
	nv := n.Int64()
	d := int64(1)
	for d%nv != 0 {
		d *= 10
	}
	i = i.Mul(i, big.NewInt(int64(d/nv)))
	if int64(f) < d {
		i = i.Div(i, big.NewInt(int64(d/int64(f))))
		d = int64(f)
	}
	return String(i.Int64(), int(d))
}

func Add(a, b string) string {
	ra, rb, f := twop(a, b)
	return String(ra+rb, f)
}

func Sub(a, b string) string {
	ra, rb, f := twop(a, b)
	return String(ra-rb, f)
}

func Neg(a string) string {
	if a[0] == '-' {
		return a[1:len(a)]
	}
	ra, f := operand(a)
	return String(-ra, f)
}

func Percent(c []byte) []byte { // convert 0..1 color lossless to percent
	r := make([]byte, len(c)+2)
	p := 0
	d := -111
	q := 0
	for p < len(c) {
		if d == p-3 {
			r[q] = '.'
			q++
		}
		if c[p] == '.' {
			d = p
		} else {
			r[q] = c[p]
			q++
		}
		p++
	}
	if d == -111 || d == p-1 {
		r[q] = '0'
		q++
		r[q] = '0'
		q++
	}
	if d == p-2 {
		r[q] = '0'
		q++
	}
	for p = 0; p < q-1 && r[p] == '0'; p++ {
	}
	return r[p:q]
}
