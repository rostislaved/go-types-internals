package interfacetype

import "unsafe"

func GetEmptyInterface(mm any) *Iface {
	up := unsafe.Pointer(&mm)

	ei := (*Iface)(up)

	return ei
}

// Еще надо посмотреть сюда: /reflect/value.go:206
type Iface struct { // runtime2.go
	Itab unsafe.Pointer // Как я понял, тут могут быть разные типы. В случае мапы Maptype
	Data unsafe.Pointer
}

type itab struct {
	Inter *Interfacetype
	_type *T_type
	Hash  uint32 // copy of T_type.hash. Used for type switches.
	_     [4]byte
	Fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type Interfacetype struct {
	Typ     T_type
	Pkgpath Name
	Mhdr    []Imethod
}

type T_type struct { // структура называется _type, но переименовал в T_type, чтобы была экспортируемой
	Size       uintptr
	Ptrdata    uintptr // size of memory prefix holding all pointers
	Hash       uint32
	Tflag      Tflag
	Align      uint8
	FieldAlign uint8
	Kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	Gcdata    *byte
	Str       NameOff
	PtrToThis TypeOff
}

type (
	Tflag   uint8
	NameOff int32
	TypeOff int32
	Name    struct {
		bytes *byte
	}
)

type Imethod struct {
	Name NameOff
	Ityp TypeOff
}
