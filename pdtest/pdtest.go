package main

import (
	"fmt"
	"os"

	"github.com/grokify/pdfreader"
)

// Example program for pdfread.go

// The program takes a PDF file as argument and writes the MediaBoxes and
// defined fonts of the pages.

func main() {
	pd := pdfreader.Load(os.Args[1])
	if pd != nil {
		pg := pd.Pages()
		for k := range pg {
			fmt.Printf("Page %d - MediaBox: %s\n",
				k+1, pd.Att("/MediaBox", pg[k]))
			fonts := pd.PageFonts(pg[k])
			for l := range fonts {
				fontname := pd.Dic(fonts[l])["/BaseFont"]
				fmt.Printf("  %s = \"%s\"\n",
					l, fontname[1:])
			}
		}
	}
}
