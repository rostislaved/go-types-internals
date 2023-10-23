package main

import (
	"github.com/rostislaved/go-types-internals/maptype"
)

func main() {
	m := make(map[int]int)

	for i := 1; i < 20; i++ {
		m[i] = i
	}
	// maptype.PrintMapInfo(m)
	b := maptype.GetMapStruct(m)
	b.PrintBuckets()
	maptype.PrintBucketsGeneric[int8](b)
}
