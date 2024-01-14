// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// LZW decoder for PDF.
package lzw

import (
	"github.com/grokify/pdfreader/crush"
)

const (
	lzwEOD       = 257
	lzwReset     = 256
	lzwDicSize   = 4096
	lzwStartBits = 9
	lzwStartUTOK = 258
)

type lzwDecoder struct {
	bits   *crush.BitT
	bc, cp int
	early  bool
}

func (lzw *lzwDecoder) reset() {
	lzw.bc = lzwStartBits
	lzw.cp = lzwStartUTOK - 1
}

func newLzwDecoder(s []byte, early bool) (lzw *lzwDecoder) {
	lzw = new(lzwDecoder)
	lzw.bits = crush.NewBits(s)
	lzw.early = early
	lzw.reset()
	return
}

func (lzw *lzwDecoder) update() bool {
	if lzw.cp < lzwDicSize-1 {
		lzw.cp++
		cmp := lzw.cp
		if lzw.early {
			cmp++
		}
		switch cmp {
		case 512:
			lzw.bc = 10
		case 1024:
			lzw.bc = 11
		case 2048:
			lzw.bc = 12
		}
		return true
	}
	return false
}

func (lzw *lzwDecoder) token() (r int) {
	for {
		r = lzw.bits.Get(lzw.bc)
		if r != lzwReset {
			break
		}
		lzw.reset()
	}
	return
}

func DecodeToSlice(s []byte, out []byte, early bool) (r int) {
	lzw := newLzwDecoder(s, early)
	dict := make([][]byte, lzwDicSize)
	for i := 0; i <= 255; i++ {
		dict[i] = []byte{byte(i)}
	}
	for c := lzw.token(); c != lzwEOD; c = lzw.token() {
		k := r
		for i := 0; i < len(dict[c]); i++ {
			out[r] = dict[c][i]
			r++
		}
		if lzw.update() {
			dict[lzw.cp] = out[k : r+1]
		}
	}
	return
}

func CalculateLength(s []byte, early bool) (r int) {
	lzw := newLzwDecoder(s, early)
	dict := make([]int, lzwDicSize)
	for i := 0; i <= 255; i++ {
		dict[i] = 1
	}
	for c := lzw.token(); c != lzwEOD; c = lzw.token() {
		r += dict[c]
		if lzw.update() {
			dict[lzw.cp] = dict[c] + 1
		}
	}
	return
}

func Decode(s []byte, early bool) []byte {
	r := make([]byte, CalculateLength(s, early)+1)
	return r[0:DecodeToSlice(s, r, early)]
}
