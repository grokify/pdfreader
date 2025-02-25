// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

package graf

import (
	"github.com/grokify/pdfreader/fancy"
	"github.com/grokify/pdfreader/ps"
	"github.com/grokify/pdfreader/stacks"
	"github.com/grokify/pdfreader/strm"
	"github.com/grokify/pdfreader/util"
)

type DrawerColor interface {
	RGB(rgb [][]byte) string
	CMYK(cmyk [][]byte) string
	Gray(g []byte) string
}

type Drawer interface {
	CloseDrawing()
	ClosePath()
	Concat(s [][]byte)
	CurveTo(s [][]byte)
	DropPath()
	EOFill()
	EOFillAndStroke()
	Fill()
	FillAndStroke()
	LineTo(s [][]byte)
	MoveTo(s [][]byte)
	Rectangle(s [][]byte)
	SetIdentity()
	Stroke()
}

type DrawerConfig interface {
	SetCMYKFill(s [][]byte)
	SetCMYKStroke(s [][]byte)
	SetColors(DrawerColor)
	SetFlat(a []byte)
	SetGrayFill(a []byte)
	SetGrayStroke(a []byte)
	SetLineCap(a []byte)
	SetLineJoin(a []byte)
	SetLineWidth(a []byte)
	SetMiterLimit(a []byte)
	SetRGBFill(s [][]byte)
	SetRGBStroke(s [][]byte)
}

type DrawerConfigT struct {
	FillColor   string
	StrokeColor string
	LineWidth   string
	LineCap     string
	LineJoin    string
	MiterLimit  string
	Flat        string
	color       DrawerColor
}

func (t *DrawerConfigT) SetLineWidth(a []byte) {
	t.LineWidth = string(a)
}
func (t *DrawerConfigT) SetLineCap(a []byte) {
	t.LineCap = string(a)
}
func (t *DrawerConfigT) SetLineJoin(a []byte) {
	t.LineJoin = string(a)
}
func (t *DrawerConfigT) SetMiterLimit(a []byte) {
	t.MiterLimit = string(a)
}
func (t *DrawerConfigT) SetFlat(a []byte) {
	t.Flat = string(a)
}

type TextConfig interface {
	SetCharSpace(a []byte)
	SetFontAndSize(s [][]byte)
	SetLeading(a []byte)
	SetRender(a []byte)
	SetRise(a []byte)
	SetScale(a []byte)
	SetWordSpace(a []byte)
}

type TextConfigT struct {
	CharSpace string
	WordSpace string
	Scale     string
	Leading   string
	Render    string
	Rise      string
	Font      string
	FontSize  string
}

func (t *TextConfigT) SetCharSpace(a []byte) {
	t.CharSpace = string(a)
}
func (t *TextConfigT) SetWordSpace(a []byte) {
	t.WordSpace = string(a)
}
func (t *TextConfigT) SetScale(a []byte) {
	t.Scale = string(a)
}
func (t *TextConfigT) SetLeading(a []byte) {
	t.Leading = string(a)
}
func (t *TextConfigT) SetRender(a []byte) {
	t.Render = string(a)
}
func (t *TextConfigT) SetRise(a []byte) {
	t.Rise = string(a)
}

type DrawerText interface {
	TMoveTo(s [][]byte)
	TNextLine()
	TSetMatrix(s [][]byte)
	TShow(a []byte)
}

type DocumentMarker interface {
}

type PdfDrawerT struct {
	Stack        stacks.Stack
	Ops          map[string]func(pd *PdfDrawerT)
	CurrentPoint [][]byte
	ConfigD      *DrawerConfigT
	TConfD       *TextConfigT
	Write        *util.OutT
	Draw         Drawer
	Config       DrawerConfig
	TConf        TextConfig
	Text         DrawerText
	Marker       DocumentMarker
}

var PdfOps = map[string]func(pd *PdfDrawerT){
	"B": func(pd *PdfDrawerT) {
		pd.Draw.FillAndStroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"B*": func(pd *PdfDrawerT) {
		pd.Draw.EOFillAndStroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"F": func(pd *PdfDrawerT) {
		pd.Draw.Fill()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"S": func(pd *PdfDrawerT) {
		pd.Draw.Stroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"b": func(pd *PdfDrawerT) {
		pd.Draw.ClosePath()
		pd.CurrentPoint = nil
		pd.Draw.FillAndStroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"b*": func(pd *PdfDrawerT) {
		pd.Draw.ClosePath()
		pd.CurrentPoint = nil
		pd.Draw.EOFillAndStroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"c": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(6)
		pd.Draw.CurveTo(a)
		pd.CurrentPoint = a[4:6]
	},
	"cm": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(6)
		pd.Draw.Concat(a)
		pd.CurrentPoint = a[4:6]
	},
	"f": func(pd *PdfDrawerT) {
		pd.Draw.Fill()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"f*": func(pd *PdfDrawerT) {
		pd.Draw.EOFill()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"h": func(pd *PdfDrawerT) {
		pd.Draw.ClosePath()
		pd.CurrentPoint = nil
	},
	"l": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(2)
		pd.Draw.LineTo(a)
		pd.CurrentPoint = a
	},
	"m": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(2)
		pd.Draw.MoveTo(a)
		pd.CurrentPoint = a
	},
	"n": func(pd *PdfDrawerT) {
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"re": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(4)
		pd.Draw.Rectangle(a)
		pd.CurrentPoint = nil
	},
	"s": func(pd *PdfDrawerT) {
		pd.Draw.ClosePath()
		pd.Draw.Stroke()
		pd.Draw.DropPath()
		pd.CurrentPoint = nil
	},
	"v": func(pd *PdfDrawerT) {
		c := pd.CurrentPoint
		a := pd.Stack.Drop(4)
		pd.Draw.CurveTo([][]byte{c[0], c[1], a[0], a[1], a[2], a[3]})
		pd.CurrentPoint = a[2:4]
	},
	"y": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(4)
		pd.Draw.CurveTo([][]byte{a[0], a[1], a[2], a[3], a[2], a[3]})
		pd.CurrentPoint = a[2:4]
	},
	"G": func(pd *PdfDrawerT) {
		pd.Config.SetGrayStroke(pd.Stack.Pop())
		pd.Ops["SC"] = pd.Ops["G"]
	},
	"J": func(pd *PdfDrawerT) {
		pd.Config.SetLineCap(pd.Stack.Pop())
	},
	"K": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(4)
		pd.Config.SetCMYKStroke(a)
		pd.Ops["SC"] = pd.Ops["K"]
	},
	"M": func(pd *PdfDrawerT) {
		pd.Config.SetMiterLimit(pd.Stack.Pop())
	},
	"RG": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(3)
		pd.Config.SetRGBStroke(a)
		pd.Ops["SC"] = pd.Ops["RG"]
	},
	"g": func(pd *PdfDrawerT) {
		pd.Config.SetGrayFill(pd.Stack.Pop())
		pd.Ops["sc"] = pd.Ops["g"]
	},
	"gs": func(pd *PdfDrawerT) {
		// FIXME!
		pd.Draw.SetIdentity()
		pd.Stack.Pop()
	},
	"i": func(pd *PdfDrawerT) {
		pd.Config.SetFlat(pd.Stack.Pop())
	},
	"j": func(pd *PdfDrawerT) {
		pd.Config.SetLineJoin(pd.Stack.Pop())
	},
	"k": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(4)
		pd.Config.SetCMYKFill(a)
		pd.Ops["sc"] = pd.Ops["k"]
	},
	"rg": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(3)
		pd.Config.SetRGBFill(a)
		pd.Ops["sc"] = pd.Ops["rg"]
	},
	"w": func(pd *PdfDrawerT) {
		pd.Config.SetLineWidth(pd.Stack.Pop())
	},
	"TL": func(pd *PdfDrawerT) {
		pd.TConf.SetLeading(pd.Stack.Pop())
	},
	"Tc": func(pd *PdfDrawerT) {
		pd.TConf.SetCharSpace(pd.Stack.Pop())
	},
	"Tf": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(2)
		pd.TConf.SetFontAndSize(a)
	},
	"Tr": func(pd *PdfDrawerT) {
		pd.TConf.SetRender(pd.Stack.Pop())
	},
	"Ts": func(pd *PdfDrawerT) {
		pd.TConf.SetRise(pd.Stack.Pop())
	},
	"Tw": func(pd *PdfDrawerT) {
		pd.TConf.SetWordSpace(pd.Stack.Pop())
	},
	"Tz": func(pd *PdfDrawerT) {
		pd.TConf.SetScale(pd.Stack.Pop())
	},
	"'": func(pd *PdfDrawerT) {
		pd.Text.TNextLine()
		pd.Text.TShow(pd.Stack.Pop())
	},
	"BT": func(pd *PdfDrawerT) {
		pd.Text.TSetMatrix(nil)
	},
	"ET": func(pd *PdfDrawerT) {
	},
	"T*": func(pd *PdfDrawerT) {
		pd.Text.TNextLine()
	},
	"TD": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(2)
		pd.TConf.SetLeading(util.Bytes(strm.Neg(string(a[1]))))
		pd.Text.TMoveTo(a)
	},
	"TJ": func(pd *PdfDrawerT) {
		pd.Text.TShow(pd.Stack.Pop())
	},
	"Td": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(2)
		pd.Text.TMoveTo(a)
	},
	"Tj": func(pd *PdfDrawerT) {
		pd.Text.TShow(pd.Stack.Pop())
	},
	"Tm": func(pd *PdfDrawerT) {
		a := pd.Stack.Drop(6)
		pd.Text.TSetMatrix(a)
	},
	"\"": func(pd *PdfDrawerT) {
		t := pd.Stack.Drop(3)
		pd.TConf.SetWordSpace(t[0])
		pd.TConf.SetCharSpace(t[1])
		pd.Text.TNextLine()
		pd.Text.TShow(t[3])
	},
	"BDC": func(pd *PdfDrawerT) {
		pd.Stack.Drop(2)
	},
	"BMC": func(pd *PdfDrawerT) {
		pd.Stack.Pop()
	},
	"DP": func(pd *PdfDrawerT) {
		pd.Stack.Drop(2)
	},
	"EMC": func(pd *PdfDrawerT) {
	},
	"MP": func(pd *PdfDrawerT) {
		pd.Stack.Pop()
	},
}

func (pd *PdfDrawerT) Interpret(rdr fancy.Reader) {
	for {
		t, _ := ps.Token(rdr)
		if len(t) == 0 {
			break
		}
		if f, ok := pd.Ops[string(t)]; ok {
			f(pd)
		} else {
			pd.Stack.Push(t)
		}
	}
}

// "constructor"

func NewPdfDrawer() *PdfDrawerT {
	r := new(PdfDrawerT)
	r.Stack = stacks.NewStack(1024)
	r.Ops = make(map[string]func(pd *PdfDrawerT))
	for k := range PdfOps {
		r.Ops[k] = PdfOps[k]
	}
	r.ConfigD = new(DrawerConfigT)
	r.Config = r.ConfigD
	r.TConfD = new(TextConfigT)
	r.TConf = r.TConfD
	r.Text = r.TConfD
	r.Write = new(util.OutT)
	return r
}

// few glue code to get interfaces working.

func (t *DrawerConfigT) SetCMYKFill(s [][]byte) {
	t.FillColor = t.color.CMYK(s)
}
func (t *DrawerConfigT) SetCMYKStroke(s [][]byte) {
	t.StrokeColor = t.color.CMYK(s)
}
func (t *DrawerConfigT) SetGrayFill(a []byte) {
	t.FillColor = t.color.Gray(a)
}
func (t *DrawerConfigT) SetGrayStroke(a []byte) {
	t.StrokeColor = t.color.Gray(a)
}
func (t *DrawerConfigT) SetRGBFill(s [][]byte) {
	t.FillColor = t.color.RGB(s)
}
func (t *DrawerConfigT) SetRGBStroke(s [][]byte) {
	t.StrokeColor = t.color.RGB(s)
}
func (t *DrawerConfigT) SetColors(hook DrawerColor) {
	t.color = hook
}

func (t *TextConfigT) SetFontAndSize(a [][]byte) {
	t.Font = string(a[0])
	t.FontSize = string(a[1])
}
func (t *TextConfigT) TMoveTo(s [][]byte)    {}
func (t *TextConfigT) TNextLine()            {}
func (t *TextConfigT) TSetMatrix(s [][]byte) {}
func (t *TextConfigT) TShow(a []byte)        {}
