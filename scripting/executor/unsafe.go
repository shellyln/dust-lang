package executor

import (
	"unsafe"
)

//
type rawInterface struct {
	Typ uintptr
	Ptr uintptr
}

//
type rawInterface2 struct {
	Typ uintptr
	Ptr unsafe.Pointer
}

//
type rawStringHeader struct {
	Data uintptr
	Len  int
}

//
type rawStringHeader2 struct {
	Data unsafe.Pointer
	Len  int
}

//
type rawSliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

//
type rawSliceHeader2 struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
