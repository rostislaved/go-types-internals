package main

import (
	"github.com/rostislaved/go-types-internals/maptype"
)

func main() {
	m := make(map[int]int)

	for i := 1; i < 10; i++ {
		m[i] = i
	}

	mapStruct := maptype.GetMapStruct(m) // retrieve map internals

	// Print general map info: [# elements], [# buckets], [LoadFactor]
	mapStruct.PrintMapInfo()

	// Prints:
	// Elements | Buckets | LoadFactor
	//       9 |       2 | 4.5

	// Print Buckets(and OldBuckets): key:value
	mapStruct.PrintBuckets()

	// Prints
	// Buckets:
	//	Bucket: 0
	//	   Element: [0] = 1
	//	   Element: [1] = 2
	//	   Element: [2] = 6
	//	   Element: [3] = 8
	//	   Element: [4] = 0
	//	   Element: [5] = 0
	//	   Element: [6] = 0
	//	   Element: [7] = 0
	//	Bucket: 1
	//	   Element: [0] = 3
	//	   Element: [1] = 4
	//	   Element: [2] = 5
	//	   Element: [3] = 7
	//	   Element: [4] = 9
	//	   Element: [5] = 0
	//	   Element: [6] = 0
	//	   Element: [7] = 0
	//	OldBuckets:
	//	[]

	// Alternative to previous func in case the value type is not supported
	maptype.PrintBucketsGeneric[int64](mapStruct)

	// Get buckets [][][]unsafe.Pointer
	buckets := mapStruct.GetBuckets()
	_ = buckets
}
