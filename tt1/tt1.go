// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

// Type1 font tester.
package main

import (
	"fmt"
	"os"

	"github.com/grokify/pdfreader/fancy"
	"github.com/grokify/pdfreader/pfb"
	"github.com/grokify/pdfreader/type1"
	"github.com/grokify/pdfreader/util"
)

// use this program with a pfa-font - it is only here for testing

func dumpT1(i *type1.TypeOneI) {
	for k := range i.Fonts {
		fmt.Printf("Font: %s %s\n", k, i.Fonts[k])
		df := i.Dic(i.Fonts[k])
		for l := range df {
			fmt.Printf("  %s %s\n", l, df[l])
		}
		fmt.Printf("\nFontInfo:\n")
		d := i.Dic(string(df["/FontInfo"]))
		for l := range d {
			fmt.Printf("  %s %s\n", l, d[l])
		}
		/*
		   fmt.Printf("\n\nCharStrings:");
		   d = i.Dic(string(df["/CharStrings"]));
		   for l := range d {
		     fmt.Printf("  %s %v\n", l, d[l])
		   }
		*/
	}
}

func main() {
	a, _ := os.ReadFile(os.Args[1])
	if a[0] == 128 {
		a = pfb.Decode(a)
	}
	g := type1.Read(fancy.SliceReader(a))
	fmt.Printf("%v\n", util.StringArray(g.St.Dump()))
	dumpT1(g)
}
