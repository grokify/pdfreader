// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

package cmapi

import (
  "github.com/nathankerr/pdfreader/cmapt"
  "github.com/nathankerr/pdfreader/fancy"
  "github.com/nathankerr/pdfreader/ps"
  "github.com/nathankerr/pdfreader/stacks"
  "github.com/nathankerr/pdfreader/xchar"
)

// CMap "interpreter" - this PS btw.

type CharMapperT struct {
  Ranges, Uni *cmapt.CMapT
}

func New() *CharMapperT {
  r := new(CharMapperT)
  r.Ranges = cmapt.New()
  r.Uni = cmapt.New()
  return r
}

type CharMapperI struct {
  Target *CharMapperT
  St     stacks.Stack
  Dic    map[string][]byte
  Marker int
  Args   [][]byte
}

func NewInterpreter(t *CharMapperT) *CharMapperI {
  r := new(CharMapperI)
  r.Target = t
  r.St = stacks.NewStack(1024)
  r.Dic = make(map[string][]byte)
  return r
}

var Ops = map[string]func(t *CharMapperI){
  "begin": func(t *CharMapperI) {
    a := t.St.Pop()
    _ = a
  },
  "beginbfchar": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "beginbfrange": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "begincidchar": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "begincidrange": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "begincmap": func(t *CharMapperI) {
  },
  "begincodespacerange": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "beginnotdefchar": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "beginnotdefrange": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "beginrearrangedfont": func(t *CharMapperI) {
    t.Args = t.St.Drop(2)
    t.Marker = t.St.Depth()
  },
  "beginusematrix": func(t *CharMapperI) {
    t.Args = t.St.Drop(1)
    t.Marker = t.St.Depth()
  },
  "currentdict": func(t *CharMapperI) {
    t.St.Push([]byte{'?'})
  },
  "def": func(t *CharMapperI) {
    a := t.St.Drop(2)
    t.Dic[string(a[0])] = a[1]
  },
  "defineresource": func(t *CharMapperI) {
    a := t.St.Drop(3)
    t.St.Push([]byte{'?'})
    _ = a
  },
  "dict": func(t *CharMapperI) {
    a := t.St.Pop()
    t.St.Push([]byte{'?'})
    _ = a
  },
  "dup": func(t *CharMapperI) {
    a := t.St.Pop()
    t.St.Push(a)
    t.St.Push(a)
  },
  "end": func(t *CharMapperI) {
  },
  "endbfchar": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    for k := 0; k < len(a); k += 2 {
      t.Target.Uni.Add(ps.StrInt(ps.String(a[k])), ps.StrInt(ps.String(a[k+1])))
    }
  },
  "endbfrange": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    for k := 0; k < len(a); k += 3 {
      // leaving the array expression as it is: invalidate - we do not have char names to unicode now
      t.Target.Uni.AddRange(ps.StrInt(ps.String(a[k])),
        ps.StrInt(ps.String(a[k+1])), ps.StrInt(ps.String(a[k+2])))
    }
  },
  "endcidchar": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "endcidrange": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "endcmap": func(t *CharMapperI) {
  },
  "endcodespacerange": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    for k := 0; k < len(a); k += 2 {
      to, l := ps.StrIntL(ps.String(a[k+1]))
      t.Target.Ranges.AddDef(int(a[k][0]), int(a[k+1][0])+1, l)
      t.Target.Ranges.AddDef(ps.StrInt(ps.String(a[k])), to+1, l) // just not used.
    }
  },
  "endnotdefchar": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "endnotdefrange": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "endrearrangedfont": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "endusematrix": func(t *CharMapperI) {
    a := t.St.Drop(t.St.Depth() - t.Marker)
    _ = a
  },
  "exch": func(t *CharMapperI) {
    a := t.St.Drop(2)
    t.St.Push(a[1])
    t.St.Push(a[0])
  },
  "findresource": func(t *CharMapperI) {
    a := t.St.Drop(2)
    t.St.Push([]byte{'?'})
    _ = a
  },
  "pop": func(t *CharMapperI) {
    a := t.St.Pop()
    _ = a
  },
  "usecmap": func(t *CharMapperI) {
    a := t.St.Pop()
    _ = a
  },
  "usefont": func(t *CharMapperI) {
    a := t.St.Pop()
    _ = a
  },
}

func Read(rdr fancy.Reader) (r *CharMapperT) {
  r = New()
  if rdr == nil { // make identity setup
    r.Uni.AddRange(0, 256, 0)
    r.Ranges.AddDef(0, 256, 1)
    return
  }
  cm := NewInterpreter(r)
  for {
    t, _ := ps.Token(rdr)
    if len(t) == 0 {
      break
    }
    if f, ok := Ops[string(t)]; ok {
      f(cm)
    } else {
      cm.St.Push(t)
    }
  }
  return
}

func Decode(s []byte, to *CharMapperT) (r []byte) {
  r = make([]byte, len(s)*6)
  p := 0
  for k := 0; k < len(s); {
    l := to.Ranges.Code(int(s[k]))
    a := ps.StrInt(s[k : k+l])
    k += l
    p += xchar.EncodeRune(to.Uni.Code(a), r[p:len(r)])
  }
  return r[0:p]
}
