# go-types-internals
Get access to runtime implementation of go types

For any supported type you can get its internal representation with help of functions:

```go
func GetMapStruct(value any) MapStruct  
func GetSliceStruct(value any) SliceStruct
```

Each type's struct has a form of this struct:
```go
type TypeStruct struct {
	InitialValue interface{}
	Type         TypeT
	Value        ValueT
}
```

and has its own handy methods (for example for printing map info)

**Content:**  
**[1. Map](#1)**  
**[2. Slice](#2)**  
**[_. Empty Interface](#_)**  



# <a name="2">2. Map:</a>
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
	InitialValue  interface{} // map itself
	Type          Maptype     // type info
	Value         Hmap        // map data
}
```

There are some handy functions:

### 1. Print general map info: [# elements], [# buckets], [LoadFactor]:
```go
mapStruct.PrintMapInfo()
```
	// Elements | Buckets | LoadFactor
	//        9 |       2 | 4.5


### 2. Print Buckets(and OldBuckets): key:value
```go
mapStruct.PrintBuckets()

// Alternative to previous func in case the value type is not supported
maptype.PrintBucketsGeneric[int64](mapStruct)

```
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

# <a name="2">2. Slice:</a>
```go
package main

import (
	"github.com/rostislaved/go-types-internals/slicetype"
    "fmt"
)

func main() {
	s := []int{1, 2, 3}

	ss := slicetype.GetSliceStruct(s)

	fmt.Println(ss.Value.Len)
	fmt.Println(ss.Value.Cap)
	fmt.Println(ss.Value.Array)

	fmt.Println(ss.Type.Elem.Align)
}
```
no handy methods yet :) only access to fields

# <a name="_">_. Empty Interface:</a>
```go
iFace := GetEmptyInterface(anyValue)
```

it is used mainly for others function, but I left it exported
```go
type Iface struct { 
	Itab unsafe.Pointer
	Data unsafe.Pointer
}
```

