package general

import "github.com/rostislaved/go-types-internals/interfacetype"

type TypeStruct[TypeT, ValueT any] struct {
	InitialValue interface{}
	Type         TypeT
	Value        ValueT
}

func GetTypeStruct[TypeT, ValueT any](typeValue any) TypeStruct[TypeT, ValueT] {
	ei := interfacetype.GetEmptyInterface(typeValue)

	typ := *(*TypeT)(ei.Itab)
	value := *(*ValueT)(ei.Data)

	return TypeStruct[TypeT, ValueT]{
		InitialValue: typeValue,
		Type:         typ,
		Value:        value,
	}
}
