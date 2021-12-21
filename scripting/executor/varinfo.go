package executor

import (
	"reflect"
	"unsafe"

	. "github.com/shellyln/takenoco/base"
)

// Implements interface Box
type VariableInfo struct {
	Flags AstOpCodeType
	Value interface{}
}

//
func (p *VariableInfo) GetAny() interface{} {
	return p.Value
}

//
func (p *VariableInfo) GetInt() int64 {
	// return p.Value.(int64)
	return *(*int64)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetUint() uint64 {
	// return p.Value.(uint64)
	return *(*uint64)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetFloat() float64 {
	// return p.Value.(float64)
	return *(*float64)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetBool() bool {
	// return p.Value.(bool)
	return *(*bool)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetString() string {
	// return p.Value.(string)
	return *(*string)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetBytes() []byte {
	// return p.Value.([]byte)
	return *(*[]byte)((*rawInterface2)(unsafe.Pointer(&p.Value)).Ptr)
}

//
func (p *VariableInfo) GetPointer() unsafe.Pointer {
	return p.Value.(unsafe.Pointer)
}

//
func (p *VariableInfo) SetAny(v interface{}) {
	p.Value = v
}

//
func (p *VariableInfo) SetInt(v int64) {
	p.Value = v
}

//
func (p *VariableInfo) SetUint(v uint64) {
	p.Value = v
}

//
func (p *VariableInfo) SetFloat(v float64) {
	p.Value = v
}

//
func (p *VariableInfo) SetBool(v bool) {
	p.Value = v
}

//
func (p *VariableInfo) SetString(v string) {
	p.Value = v
}

//
func (p *VariableInfo) SetBytes(v []byte) {
	p.Value = v
}

//
func (p *VariableInfo) SetPointer(v unsafe.Pointer) {
	p.Value = v
}

//
func (p *VariableInfo) Index(i int) Box {
	return ReflectionBox{Val: reflect.ValueOf(p.Value).Index(i)}
}

//
func (p *VariableInfo) MapIndex(k string) Box {
	return ReflectionBox{Val: reflect.ValueOf(p.Value).MapIndex(reflect.ValueOf(k))}
}

//
func (p *VariableInfo) ComplexMapIndex(k interface{}) Box {
	return ReflectionBox{Val: reflect.ValueOf(p.Value).MapIndex(reflect.ValueOf(k))}
}
