package main

import "github.com/rostislaved/go-types-internals/maptype"

func main() {
	m := make(map[int]int)

	for i := 0; i < 1e3; i++ {
		m[i] = i
	}
	maptype.PrintMapInfo(m)
	ms := maptype.GetMapStruct(m)
	//
	ms.PrintBuckets()
}
