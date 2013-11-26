// Copyright (c) 2009 Helmar Wodtke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// The MIT License is an OSI approved license and can
// be found at
//   http://www.opensource.org/licenses/mit-license.php

package type1

import (
  "github.com/nathankerr/pdfreader/fancy"
  "fmt"
  "github.com/nathankerr/pdfreader/hex"
  "github.com/nathankerr/pdfreader/ps"
  "github.com/nathankerr/pdfreader/stacks"
  "github.com/nathankerr/pdfreader/strm"
  "github.com/nathankerr/pdfreader/util"
)

// Type1 Font "interpreter" - this PS btw.

type DicT struct {
  Defs map[string][]byte
  Name []byte
}

type TypeOneI struct {
  Rdr    fancy.Reader
  St     stacks.Stack
  DicSt  [64]*DicT
  DicSp  int
  Dicts  [1024]*DicT
  DicNo  int
  Arrays [1024][][]byte
  ArraNo int
  Marker int
  Args   [][]byte
  Done   bool
  Fonts  map[string]string
}

func (t *TypeOneI) NewDic() (r []byte) {
  rs := fmt.Sprintf("D%d", t.DicNo)
  r = util.Bytes(rs)
  t.Dicts[t.DicNo] = new(DicT)
  t.Dicts[t.DicNo].Defs = make(map[string][]byte)
  t.Dicts[t.DicNo].Name = r
  t.DicNo++
  return
}

func (t *TypeOneI) NewArray(size int) (r []byte) {
  r = util.Bytes(fmt.Sprintf("A%d", t.ArraNo))
  t.Arrays[t.ArraNo] = make([][]byte, size)
  t.ArraNo++
  return
}

func NewInterpreter() *TypeOneI {
  r := new(TypeOneI)
  r.St = stacks.NewStack(1024)
  r.NewDic()
  r.DicSt[0] = r.Dicts[0]
  r.Done = false
  r.Fonts = make(map[string]string)
  return r
}

const EEXEC_KEY = 55665
const CHARSTRING_KEY = 4330

func T1Decrypt(r int, s []byte) []byte {
  p := make([]byte, len(s))
  for k := range s {
    p[k] = s[k] ^ byte(r>>8)
    r = ((r+int(s[k]))*52845 + 22719) & 65535
  }
  return p
}

func eexec(rdr fancy.Reader) []byte {
  fpos, _ := rdr.Seek(0, 1)
  b := fancy.ReadAll(rdr)
  cnt := 0
  pos := 0
  k := 0
  for ; cnt < 256 && k < len(b); k++ {
    switch b[k] {
    case 32, 10, 13, 9:
    case '0':
      cnt++
    default:
      cnt = 0
      pos = k + 1
    }
  }
  b = b[0:pos]
  rdr.Seek(fpos+int64(k), 0)
  if hex.IsHex(b[0]) {
    b = hex.Decode(string(b))
  }
  return T1Decrypt(EEXEC_KEY, b)[4:]
}

func (t *TypeOneI) op_ifelse(a [][]byte) {
  p := a[2]
  if string(a[0]) == "true" {
    p = a[1]
  }
  if len(p) > 2 && p[0] == '{' {
    proceed(t, fancy.SliceReader(p[1:len(p)-1]))
  }
}

var Ops = map[string]func(t *TypeOneI){
  "+": func(t *TypeOneI) {
  },
  "array": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(t.NewArray(strm.Int(string(a), 1)))
  },
  "begin": func(t *TypeOneI) {
    a := t.St.Pop()
    if a[0] != 'D' {
      panic("Wrong dictionary!\n")
    }
    t.DicSp++
    t.DicSt[t.DicSp] = t.Dicts[strm.Int(string(a[1:]), 1)]
  },
  "bind": func(t *TypeOneI) {
  },
  "cleartomark": func(t *TypeOneI) {
    a := t.St.Pop()
    for string(a) != "mark" {
      a = t.St.Pop()
    }
  },
  "closefile": func(t *TypeOneI) {
    a := t.St.Pop()
    t.Done = true
    _ = a
  },
  "currentdict": func(t *TypeOneI) {
    t.St.Push(t.DicSt[t.DicSp].Name)
  },
  "currentfile": func(t *TypeOneI) {
    t.St.Push([]byte{'?'})
  },
  "def": func(t *TypeOneI) {
    a := t.St.Drop(2)
    t.DicSt[t.DicSp].Defs[string(a[0])] = a[1]
  },
  "definefont": func(t *TypeOneI) {
    a := t.St.Drop(2)
    t.Fonts[string(a[0])] = string(a[1])
    t.St.Push(util.Bytes("<FONT>")) // FIXME, we need this.
    _ = a
  },
  "defineresource": func(t *TypeOneI) {
    a := t.St.Drop(3)
    t.St.Push([]byte{'?'})
    _ = a
  },
  "dict": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(t.NewDic())
    _ = a
  },
  "dup": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(a)
    t.St.Push(a)
  },
  "end": func(t *TypeOneI) {
    t.DicSp--
  },
  "exch": func(t *TypeOneI) {
    a := t.St.Drop(2)
    a0 := a[0]
    t.St.Push(a[1])
    t.St.Push(a0)
  },
  "executeonly": func(t *TypeOneI) {
  },
  "findresource": func(t *TypeOneI) {
    a := t.St.Drop(2)
    t.St.Push([]byte{'?'})
    _ = a
  },
  "for": func(t *TypeOneI) {
    a := t.St.Drop(4)
    // FIXME
    _ = a
  },
  "get": func(t *TypeOneI) {
    a := t.St.Drop(2)
    i := strm.Int(string(a[0][1:]), 1)
    if a[0][0] == 'D' {
      t.St.Push(t.Dicts[i].Defs[string(a[1])])
    } else if a[0][0] == 'A' {
      t.St.Push(t.Arrays[i][strm.Int(string(a[1]), 1)])
    } else {
      panic("Can not 'get' from!\n")
    }
  },
  "index": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(t.St.Index(strm.Int(string(a), 1) + 1))
  },
  "known": func(t *TypeOneI) {
    a := t.St.Drop(2)
    t.St.Push(util.Bytes("false")) // FIX ME knows nothing ;)
    _ = a
  },
  "noaccess": func(t *TypeOneI) {
  },
  "pop": func(t *TypeOneI) {
    a := t.St.Pop()
    _ = a
  },
  "put": func(t *TypeOneI) {
    a := t.St.Drop(3)
    if a[0][0] == 'D' {
      t.Dicts[strm.Int(string(a[0][1:]), 1)].Defs[string(a[1])] = a[2]
    } else if a[0][0] == 'A' {
      t.Arrays[strm.Int(string(a[0][1:]), 1)][strm.Int(string(a[1]), 1)] = a[2]
    } else {
      panic("Wrong dictionary or array!\n")
    }
  },
  "readonly": func(t *TypeOneI) {
  },
  "readstring": func(t *TypeOneI) {
    a := t.St.Drop(2)
    c, _ := t.Rdr.Read(a[1])
    t.St.Push(a[1][0:c])
    t.St.Push(util.Bytes("true"))
  },
  "string": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(make([]byte, strm.Int(string(a), 1)))
  },
  "userdict": func(t *TypeOneI) {
    t.St.Push(util.Bytes("D0"))
  },
  "where": func(t *TypeOneI) {
    a := t.St.Pop()
    t.St.Push(util.Bytes("false"))
    _ = a
  },
}

// used to solve an inititalization loop
func init() {
  Ops["eexec"] = func(t *TypeOneI) {
    a := t.St.Pop()
    b := eexec(t.Rdr)
    old := t.Rdr
    t.Rdr = fancy.SliceReader(b)
    proceed(t, t.Rdr)
    t.Rdr = old
    t.Done = false
    _ = a
  }

  Ops["if"] = func(t *TypeOneI) {
    a := t.St.Drop(2)
    t.op_ifelse([][]byte{a[0], a[1], []byte{}})
  }

  Ops["ifelse"] = func(t *TypeOneI) {
    a := t.St.Drop(3)
    t.op_ifelse(a)
  }
}

func find(i *TypeOneI, s string) (r []byte, ok bool) {
  for k := i.DicSp; k >= 0 && !ok; k-- {
    r, ok = i.DicSt[k].Defs[s]
  }
  return
}

func proceed(i *TypeOneI, rdr fancy.Reader) {
  for !i.Done {
    t, _ := ps.Token(rdr)
    //    fmt.Printf("Stack: %v\n", util.StringArray(i.St.Dump()));
    //    fmt.Printf("--- %s\n", t);
    if len(t) < 1 {
      break
    }
    b, _ := rdr.ReadByte()
    if b > 32 {
      rdr.UnreadByte()
    }
    if len(t) == 0 {
      break
    }
    if d, ok := find(i, "/"+string(t)); ok {
      if d[0] == '{' {
        proceed(i, fancy.SliceReader(d[1:len(d)-1]))
      } else {
        i.St.Push(d)
      }
    } else if f, ok := Ops[string(t)]; ok {
      f(i)
    } else {
      i.St.Push(t)
    }
  }
  return
}

func Read(rdr fancy.Reader) (r *TypeOneI) {
  r = NewInterpreter()
  r.Rdr = rdr
  proceed(r, rdr)
  return
}

func (i *TypeOneI) Dic(id string) map[string][]byte {
  if id[0] != 'D' {
    panic("Wrong dictionary!\n")
  }
  idn := strm.Int(id[1:], 1)
  return i.Dicts[idn].Defs
}
