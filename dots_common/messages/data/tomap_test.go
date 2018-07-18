package data_messages

import(
  "encoding/json"
  "fmt"
  "reflect"
  "testing"
)

type Marshal_test struct{
            Test string
}

type Test_s2 struct {
	    E	*string	`json:"TestE"`
	    F	int	`json:"TestF"`
	    G	[]int	`json:"TestG"`
	    H	[]byte	`json:"TestH"`
}

type Test_s1 struct {
	    A	string	`json:"TestA"`
	    B	int	`json:"TestB"`
	    C	float64	`json:"TestC"`
	    D	Test_s2	`json:"TestD"`
}

func make_struct() Test_s1 {
	strc := Test_s1{}
	strc.A = "test"
	strc.B = 100
	strc.C = 2.3
	strc.D = Test_s2{}
	strc.D.E = nil
	strc.D.F = 50
	mt := Marshal_test{"TEST_CODE"}
	strc.D.H,_ = json.Marshal(mt)
	return strc
}

func Test_ToMap(t *testing.T){
	src := make_struct()
	t.Log("struct: ",src)
	m,_ := ToMap(src,Content_All)
	t.Log("map: ",m)
	t.Log("終了")
}

type P struct {
  C  C  `json:"c"`
  PC *C `json:"pc"`
}

type C struct {
  N int `json:"n"`
}

func Test_Recursive(t *testing.T) {
  p := P { C: C{12}, PC: &C{34} }
  r, err := ToMap(p, Content_All) // r should be { "c": { "n": 12 }, "pc": { "n": 34 } }
  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }

  if c, ok := r["c"]; !ok {
    t.Error("map does not contains key 'c'")
  } else {
    if m, ok := c.(map[string]interface{}); !ok {
      t.Errorf("value for key 'c' is not map[string]interface{}: %#+v", c)
    } else {
      if n, ok := m["n"]; !ok {
        t.Error("map for key 'c' does not contains key 'n'")
      } else {
        if n != 12 {
          t.Errorf("value for key 'n' is not 12: %#+v", n)
        }
      }
      if len(m) != 1 {
        t.Errorf("map size is not 1: %d", len(m))
      }
    }
  }

  if pc, ok := r["pc"]; !ok {
    t.Error("map does not contains key 'pc'")
  } else {
    if m, ok := pc.(map[string]interface{}); !ok {
      t.Errorf("value for key 'pc' is not map[string]interface{}: %#+v", pc)
    } else {
      if n, ok := m["n"]; !ok {
        t.Error("map for key 'pc' does not contains key 'n'")
      } else {
        if n != 34 {
          t.Errorf("value for key 'n' is not 34: %#+v", n)
        }
      }
      if len(m) != 1 {
        t.Errorf("map size is not 1: %d", len(m))
      }
    }
  }
}

type NT struct {
  I  interface{}    `json:"i"`
  N  interface{}    `json:"n"`
  F  interface{}    `json:"f"`
  B  []byte         `json:"b"`
  M  map[string]int `json:"m"`
  BN []byte         `json:"bn"`
  AN [0]byte        `json:"an"`
  MN map[string]int `json:"mn"`
  AE [1]interface{}  `json:"ae"`
}

type F struct { N int }

func (*F) String() string { return "hoge" /* Does not return "<nil>" for nil */ }

func Test_FilterNil(t *testing.T) {
  nt := NT{}
  nt.N = (*int)(nil)
  nt.F = (*F)(nil)
  nt.BN = make([]byte, 0)
  nt.AN = [0]byte{}
  nt.MN = make(map[string]int)
  nt.AE = [1]interface{}{ nil }

  r, err := ToMap(nt, Content_All) // r should be { "bn": []interface{}{}, "an": []interface{}{}, "mn": map[string]interface{}{} }
  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }

  if _, ok := r["i"]; ok {
    t.Errorf("map contains key 'i': %#+v", r["i"])
  }
  if _, ok := r["n"]; ok {
    t.Errorf("map contains key 'n': %#+v", r["n"])
  }
  if _, ok := r["f"]; ok {
    t.Errorf("map contains key 'f': %#+v", r["f"])
  }
  if _, ok := r["b"]; ok {
    t.Errorf("map contains key 'b': %#+v", r["b"])
  }
  if _, ok := r["m"]; ok {
    t.Errorf("map contains key 'm': %#+v", r["m"])
  }
  if v, ok := r["bn"]; !ok {
    t.Error("map does not contains key 'bn'")
  } else {
    value := reflect.ValueOf(v)
    if value.Kind() != reflect.Slice || value.Len() != 0 {
      t.Error("value for key 'bn' is not empty slice.")
    }
  }
  if v, ok := r["an"]; !ok {
    t.Error("map does not contains key 'an'")
  } else {
    value := reflect.ValueOf(v)
    if value.Kind() != reflect.Slice || value.Len() != 0 {
      t.Error("value for key 'bn' is not empty slice.")
    }
  }
  if v, ok := r["mn"]; !ok {
    t.Error("map does not contains key 'mn'")
  } else {
    value := reflect.ValueOf(v)
    if value.Kind() != reflect.Map || value.Len() != 0 {
      t.Error("value for key 'bn' is not empty slice.")
    }
  }
  if _, ok := r["ae"]; ok {
    t.Errorf("map contains key 'ae': %#+v", r["ae"])
  }

  if len(r) != 3 {
    t.Errorf("map size expected: 3, actual: %d", len(r))
  }
}

type M struct {
  X int `json:"x"`
}

type MP struct {
  Y int `json:"Y"`
}

func (m M) MarshalJSON() ([]byte, error) {
  return json.Marshal(fmt.Sprintf("M(%d)", m.X))
}

func (p *MP) MarshalJSON() ([]byte, error) {
  return json.Marshal(fmt.Sprintf("MP(%d)", p.Y))
}

type MT struct {
  M   M   `json:"m"`
  PM  *M  `json:"pm"`
  EM  *M  `json:"em"`
  MP  MP  `json:"mp"`
  PMP *MP `json:"pmp"`
  EMP *MP `json:"emp"`
}

func Test_Marshaler(t *testing.T) {
  mt := MT {
    M: M{X: 1},
    PM: &M{X:2},
    EM: nil,
    MP: MP{Y: 3},
    PMP: &MP{Y: 4},
    EMP: nil,
  }
  r, err := ToMap(mt, Content_All)
  // r should be { "m": M{X:1}, "pm": M{X:2}, "mp": MP{Y:3}, "pmp": MP{Y:4} }

  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }

  if _, ok := r["em"]; ok {
    t.Errorf("map contains key 'em': %#+v", r["em"])
  }
  if _, ok := r["emp"]; ok {
    t.Errorf("map contains key 'emp': %#+v", r["emp"])
  }

  if _, ok := r["m"]; !ok {
    t.Error("map does not contain key 'm'")
  } else {
    if _, ok := r["m"].(M); !ok {
      t.Errorf("value for key 'm' is not M: %#+v", r["m"])
    } else {
      if r["m"].(M).X != 1 {
        t.Errorf("X of value for key 'm' is not 1.")
      }
    }
  }

  if _, ok := r["pm"]; !ok {
    t.Error("map does not contain key 'pm'")
  } else {
    if _, ok := r["pm"].(M); !ok {
      t.Errorf("value for key 'pm' is not *M: %#+v", r["pm"])
    } else {
      if r["pm"].(M).X != 2 {
        t.Errorf("X of value for key 'pm' is not 2.")
      }
    }
  }

  if _, ok := r["mp"]; !ok {
    t.Error("map does not contain key 'mp'")
  } else {
    if _, ok := r["mp"].(MP); !ok {
      t.Errorf("value for key 'mp' is not MP: %#+v", r["mp"])
    } else {
      if r["mp"].(MP).Y != 3 {
        t.Errorf("X of value for key 'mp' is not 3.")
      }
    }
  }

  if _, ok := r["pmp"]; !ok {
    t.Error("map does not contain key 'pmp'")
  } else {
    if _, ok := r["pmp"].(MP); !ok {
      t.Errorf("value for key 'pmp' is not *MP: %#+v", r["pmp"])
    } else {
      if r["pmp"].(MP).Y != 4 {
        t.Errorf("X of value for key 'pmp is not 4.")
      }
    }
  }

  if len(r) != 4 {
    t.Errorf("map size expected: 4, actual: %d", len(r))
  }
}

type Y struct {
  U int
  C int `yang:"config"`
  N int `yang:"nonconfig"`
}

func Test_ContentFilter(t *testing.T) {
  y := Y{ U: 1, C: 2, N: 3 }
  r, err := ToMap(y, Content_Config)
  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }
  if _, ok := r["U"]; !ok {
    t.Errorf("map does not contains key 'U': %#+v", r)
  }
  if _, ok := r["C"]; !ok {
    t.Errorf("map does not contains key 'C': %#+v", r)
  }
  if _, ok := r["N"]; ok {
    t.Errorf("map contains key 'N': %#+v", r)
  }

  r, err = ToMap(y, Content_NonConfig)
  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }
  if _, ok := r["U"]; !ok {
    t.Errorf("map does not contains key 'U': %#+v", r)
  }
  if _, ok := r["C"]; ok {
    t.Errorf("map contains key 'C': %#+v", r)
  }
  if _, ok := r["N"]; !ok {
    t.Errorf("map does not contains key 'N': %#+v", r)
  }

  r, err = ToMap(y, Content_All)
  if err != nil {
    t.Errorf("ToMap error: %v", err)
  }
  if _, ok := r["U"]; !ok {
    t.Errorf("map does not contains key 'U': %#+v", r)
  }
  if _, ok := r["C"]; !ok {
    t.Errorf("map does not contains key 'C': %#+v", r)
  }
  if _, ok := r["N"]; !ok {
    t.Errorf("map does not contains key 'N': %#+v", r)
  }
}

type Bin []byte

func (e Bin) MarshalJSON() ([]byte, error) {
  return nil, nil
}

func Test_isJsonMarshaler(t *testing.T) {
  if true != isJsonMarshaler(make(Bin, 0)) {
    t.Errorf("Bin is not json.Marshaler")
  }
}
