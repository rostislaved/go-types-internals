package slicetype

import (
	"unsafe"

	"github.com/rostislaved/go-types-internals/general"
	"github.com/rostislaved/go-types-internals/interfacetype"
)

type SliceStruct struct {
	InitialValue interface{}
	Type         SliceType
	Value        Slice
}

func GetSliceStruct(value any) SliceStruct {
	ts := general.GetTypeStruct[SliceType, Slice](value)

	return SliceStruct{
		InitialValue: ts.InitialValue,
		Type:         ts.Type,
		Value:        ts.Value,
	}
}

// internal/abi/type.go:478
type SliceType struct {
	Type interfacetype.T_type
	Elem *interfacetype.T_type // slice element type
}

// runtime/slice.go:15
type Slice struct {
	Array unsafe.Pointer
	Len   int
	Cap   int
}
