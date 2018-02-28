package dots_common

import (
    "math"
    "math/big"
    "reflect"
    "github.com/ugorji/go/codec"
    "github.com/shopspring/decimal"
)

type decimalExt struct {}

func (_ *decimalExt) ConvertExt(v interface{}) interface{} {
    d := v.(*decimal.Decimal)
    c := d.Neg().Neg().Coefficient() // ensure initialized
    c64 := c.Int64()
    if c.Cmp(big.NewInt(c64)) == 0 { // c.IsInt64() is go1.9+ API
        return []interface{}{ d.Exponent(), c64 }
    } else {
        return nil
    }
}

func (_ *decimalExt) UpdateExt(dst interface{}, src interface{}) {
    s, ok := src.([]interface{})
    if !ok || len(s) != 2 {
        return
    }
    d := dst.(*decimal.Decimal)

    var e int32
    var c int64
    switch x := s[0].(type) {
    case int64:
        if x < math.MinInt32 || math.MaxInt32 < x {
            return
        }
        e = int32(x)
    case uint64:
        if math.MaxInt32 < x {
            return
        }
        e = int32(x)
    default:
        return
    }
    switch x := s[1].(type) {
    case int64:
        c = x
    case uint64:
        if math.MaxInt64 < x {
            return
        }
        c = int64(x)
    default:
        return
    }
    *d = decimal.New(c, e)
}

func NewCborHandle() *codec.CborHandle {
    h := new(codec.CborHandle)
    h.SetInterfaceExt(reflect.TypeOf(decimal.Decimal{}), 4, &decimalExt{})
    return h
}
