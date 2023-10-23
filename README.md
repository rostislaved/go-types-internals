# go-types-internals
Get access to runtime implementation of go types

**Content:**  
**[1. Empty Interface](#1)**  
**[2. Map](#2)**

# <a name="1">1. Empty Interface:</a>
```go
iFace := GetEmptyInterface(anyValue)
```
```go
type Iface struct { 
	Itab unsafe.Pointer
	Data unsafe.Pointer
}
```

# <a name="2">2. Map:</a>
# Map
Initialize some map:
```go
package main

import (
	"github.com/rostislaved/go-types-internals/maptype"
)

func main() {
        m := make(map[int]int)

	for i := 1; i < 10; i++ {
		m[i] = i
	}

	mapStruct := maptype.GetMapStruct(m)
}
```

```go
type MapStruct struct {
	Map     interface{} // map itself
	Hmap    Hmap        // map data
	Maptype Maptype     // type info
}
```

There are some handy functions:

### 1. Print general map info: [# elements], [# buckets], [LoadFactor]:
```go
mapStruct.PrintMapInfo()
```
	// Prints:
	// Elements | Buckets | LoadFactor
	//        9 |       2 | 4.5


### 2. Print Buckets(and OldBuckets): key:value
```go
mapStruct.PrintBuckets()

// Alternative to previous func in case the value type is not supported
maptype.PrintBucketsGeneric[int64](mapStruct)

```
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



### 3. Get buckets [][][]unsafe.Pointer
```
buckets := mapStruct.GetBuckets()
```


