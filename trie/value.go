package trie

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
)

// Value ...
type Value struct {
	Val  interface{}
	kind reflect.Value
}

// NewValue ...
func NewValue(val interface{}) *Value {
	t := val
	if v, ok := val.(*Value); ok {
		t = v.Val
	}

	return &Value{Val: t}
}

// NewValueFromBytes ...
func NewValueFromBytes(data []byte) *Value {
	if len(data) != 0 {
		data, _ := Decode(data, 0)
		return NewValue(data)
	}

	return NewValue(nil)
}

// NewSliceValue ...
func NewSliceValue(s interface{}) *Value {
	list := EmptyValue()

	if s != nil {
		if slice, ok := s.([]interface{}); ok {
			for _, val := range slice {
				list.Append(val)
			}
		} else if slice, ok := s.([]string); ok {
			for _, val := range slice {
				list.Append(val)
			}
		}
	}

	return list
}

// EmptyValue ...
func EmptyValue() *Value {
	return NewValue([]interface{}{})
}

// String ...
func (val *Value) String() string {
	return fmt.Sprintf("%x", val.Val)
}

// Type ...
func (val *Value) Type() reflect.Kind {
	return reflect.TypeOf(val.Val).Kind()
}

// IsNil ...
func (val *Value) IsNil() bool {
	return val.Val == nil
}

// Size ...
func (val *Value) Size() int {
	if data, ok := val.Val.([]interface{}); ok {
		return len(data)
	}

	return len(val.Bytes())
}

// Raw ...
func (val *Value) Raw() interface{} {
	return val.Val
}

// Interface ...
func (val *Value) Interface() interface{} {
	return val.Val
}

// Uint ....
func (val *Value) Uint() uint64 {
	if Val, ok := val.Val.(uint8); ok {
		return uint64(Val)
	} else if Val, ok := val.Val.(uint16); ok {
		return uint64(Val)
	} else if Val, ok := val.Val.(uint32); ok {
		return uint64(Val)
	} else if Val, ok := val.Val.(uint64); ok {
		return Val
	} else if Val, ok := val.Val.(int); ok {
		return uint64(Val)
	} else if Val, ok := val.Val.(uint); ok {
		return uint64(Val)
	} else if Val, ok := val.Val.([]byte); ok {
		return ReadVarint(bytes.NewReader(Val))
	} else if Val, ok := val.Val.(*big.Int); ok {
		return Val.Uint64()
	}

	return 0
}

// Byte ...
func (val *Value) Byte() byte {
	if Val, ok := val.Val.(byte); ok {
		return Val
	}

	return 0x0
}

// BigInt ...
func (val *Value) BigInt() *big.Int {
	if a, ok := val.Val.([]byte); ok {
		b := new(big.Int).SetBytes(a)

		return b
	} else if a, ok := val.Val.(*big.Int); ok {
		return a
	}
	return big.NewInt(int64(val.Uint()))
}

// Str ...
func (val *Value) Str() string {
	if a, ok := val.Val.([]byte); ok {
		return string(a)
	} else if a, ok := val.Val.(string); ok {
		return a
	} else if a, ok := val.Val.(byte); ok {
		return string(a)
	}

	return ""
}

// Bytes ...
func (val *Value) Bytes() []byte {
	if a, ok := val.Val.([]byte); ok {
		return a
	} else if s, ok := val.Val.(byte); ok {
		return []byte{s}
	} else if s, ok := val.Val.(string); ok {
		return []byte(s)
	} else if s, ok := val.Val.(*big.Int); ok {
		return s.Bytes()
	}

	return []byte{}
}

// Slice ...
func (val *Value) Slice() []interface{} {
	if d, ok := val.Val.([]interface{}); ok {
		return d
	}

	return []interface{}{}
}

// SliceFrom ...
func (val *Value) SliceFrom(from int) *Value {
	slice := val.Slice()

	return NewValue(slice[from:])
}

// SliceTo ...
func (val *Value) SliceTo(to int) *Value {
	slice := val.Slice()

	return NewValue(slice[:to])
}

// SliceFromTo ...
func (val *Value) SliceFromTo(from, to int) *Value {
	slice := val.Slice()

	return NewValue(slice[from:to])
}

// IsSlice ...
func (val *Value) IsSlice() bool {
	return val.Type() == reflect.Slice
}

// IsStr ...
func (val *Value) IsStr() bool {
	return val.Type() == reflect.String
}

// IsList ...
func (val *Value) IsList() bool {
	_, ok := val.Val.([]interface{})

	return ok
}

// IsEmpty ...
func (val *Value) IsEmpty() bool {
	return val.Val == nil || ((val.IsSlice() || val.IsStr()) && val.Size() == 0)
}

// Get ...
func (val *Value) Get(idx int) *Value {
	if d, ok := val.Val.([]interface{}); ok {
		if len(d) <= idx {
			return NewValue(nil)
		}

		if idx < 0 {
			return NewValue(nil)
		}

		return NewValue(d[idx])
	}

	return NewValue(nil)
}

// Copy ...
func (val *Value) Copy() *Value {
	switch v := val.Val.(type) {
	case *big.Int:
		return NewValue(new(big.Int).Set(v))
	case []byte:
		return NewValue(CopyBytes(v))
	default:
		return NewValue(val.Val)
	}
}

// Cmp ...
func (val *Value) Cmp(o *Value) bool {
	return reflect.DeepEqual(val.Val, o.Val)
}

// Encode ...
func (val *Value) Encode() []byte {
	return Encode(val.Val)
}

// Decode decodes value. It assume that the data passed is encoded.
func (val *Value) Decode() {
	v, _ := Decode(val.Bytes(), 0)
	val.Val = v
}

// AppendList ...
func (val *Value) AppendList() *Value {
	list := EmptyValue()
	val.Val = append(val.Slice(), list)

	return list
}

// Append ...
func (val *Value) Append(v interface{}) *Value {
	val.Val = append(val.Slice(), v)

	return val
}
