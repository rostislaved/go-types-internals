package maptype

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/rostislaved/go-types-internals/interfacetype"
)

//type MapStruct[T any] struct {
//	Hmap    Hmap
//	Maptype Maptype
//}

type MapStruct struct {
	Map     interface{}
	Hmap    Hmap
	Maptype Maptype
}

func GetMapStruct(mapValue any) MapStruct {
	ei := interfacetype.GetEmptyInterface(mapValue)

	hmapValue := *(*Hmap)(ei.Data)
	maptype := *(*Maptype)(ei.Itab)

	return MapStruct{
		Map:     mapValue,
		Hmap:    hmapValue, // value
		Maptype: maptype,   // type
	}
}

//func PrintMapInfo(m any) {
//	mapStruct := GetMapStruct(m)
//
//	hmap := mapStruct.Hmap
//
//	fmt.Printf("%8s | %7s | %10s\n", "Elements", "Buckets", "LoadFactor")
//	fmt.Printf("%8d | %7d | %.1f\n", hmap.Count, 1<<hmap.B, float32(hmap.Count)/float32(int(1<<hmap.B)))
//}

func (m MapStruct) PrintMapInfo() {
	hmap := m.Hmap

	fmt.Printf("%8s | %7s | %10s\n", "Elements", "Buckets", "LoadFactor")
	fmt.Printf("%8d | %7d | %.1f\n", hmap.Count, 1<<hmap.B, float32(hmap.Count)/float32(int(1<<hmap.B)))
}

type Buckets struct {
	// [Bucket][Bucket, overflow-Bucket...][Elemetn]
	Buckets    [][][]unsafe.Pointer
	OldBuckets [][][]unsafe.Pointer // can be nil
}

func (m MapStruct) GetBuckets() Buckets {
	bucketsUP := m.Hmap.Buckets

	b := m.getBuckets(bucketsUP)

	oldbucketsUP := m.Hmap.Oldbuckets

	ob := m.getBuckets(oldbucketsUP)

	return Buckets{
		Buckets:    b,
		OldBuckets: ob,
	}
}

func (m MapStruct) getBuckets(bucketsUP unsafe.Pointer) [][][]unsafe.Pointer {
	if bucketsUP == nil {
		return nil
	}

	numberOfBuckets := uintptr(1 << m.Hmap.B)

	m1 := make([][][]unsafe.Pointer, numberOfBuckets)

	// Проход по бакетам
	for bucketNumber := uintptr(0); bucketNumber < numberOfBuckets; bucketNumber++ {
		bucketAddress := (*Bmap)(add(bucketsUP, bucketNumber*uintptr(m.Maptype.Bucketsize)))

		overflowNumber := -1
		// проход по цепочке: бакет, оверфлоу_бакет
		for ; bucketAddress != nil; bucketAddress = bucketAddress.overflow(&m.Maptype) {
			overflowNumber++

			m2 := make([]unsafe.Pointer, bucketCnt)
			m1[int(bucketNumber)] = append(m1[int(bucketNumber)], m2)
			// Проход по элементам бакета
			for bucketElem := uintptr(0); bucketElem < bucketCnt; bucketElem++ {
				upValue := add(unsafe.Pointer(bucketAddress), DataOffset+bucketCnt*uintptr(m.Maptype.Keysize)+bucketElem*uintptr(m.Maptype.Elemsize))

				m1[bucketNumber][overflowNumber][bucketElem] = upValue

			}
		}
	}
	return m1
}

func PrintBucketsGeneric[T any](m MapStruct) {
	b := m.GetBuckets()

	fmt.Println("Buckets:")
	printBucketsGeneric[T](b.Buckets)

	fmt.Println("OldBuckets:")
	if b.OldBuckets != nil {
		printBucketsGeneric[T](b.OldBuckets)
	} else {
		fmt.Println("[]")
	}
}

func printBucketsGeneric[T any](b [][][]unsafe.Pointer) {
	// Бакеты
	for i := 0; i < len(b); i++ {
		fmt.Printf("Bucket: %d\n", i)

		// Элементы
		for j := 0; j < len(b[i]); j++ {
			if j != 0 {
				fmt.Printf("Bucket: %d (overflow)\n", i)
			}

			for k := 0; k < len(b[i][j]); k++ {
				valueUp := b[i][j][k]

				str := *(*T)(valueUp) // Вот тут

				fmt.Printf("   Element: [%d] = %v\n", k, str)
			}
		}
	}
}

func (m MapStruct) PrintBuckets() {
	b := m.GetBuckets()

	fmt.Println("Buckets:")
	m.printBuckets(b.Buckets)

	fmt.Println("OldBuckets:")
	if b.OldBuckets != nil {
		m.printBuckets(b.OldBuckets)
	} else {
		fmt.Println("[]")
	}
}

func (m MapStruct) printBuckets(b [][][]unsafe.Pointer) {
	kind := m.Maptype.Elem.Kind
	// Бакеты
	for i := 0; i < len(b); i++ {
		fmt.Printf("Bucket: %d\n", i)

		// Элементы
		for j := 0; j < len(b[i]); j++ {
			if j != 0 {
				fmt.Printf("Bucket: %d (overflow №%d)\n", i, j)
			}

			for k := 0; k < len(b[i][j]); k++ {
				valueUp := b[i][j][k]

				val := getValueByUnsafePointerAndKind(valueUp, kind) // TODO скорее всего не все типы учтены

				fmt.Printf("   Element: [%d] = %v\n", k, val)
			}
		}
	}
}

func getValueByUnsafePointerAndKind(valueUp unsafe.Pointer, kind uint8) any {
	reflectKind := reflect.Kind(kind)

	var val interface{}
	switch reflectKind {
	case reflect.Bool:
		val = *(*bool)(valueUp)
	case reflect.Int:
		val = *(*int)(valueUp)
	case reflect.Int8:
		val = *(*int8)(valueUp)
	case reflect.Int16:
		val = *(*int16)(valueUp)
	case reflect.Int32:
		val = *(*int32)(valueUp)
	case reflect.Int64:
		val = *(*int64)(valueUp)
	case reflect.Uint:
		val = *(*uint)(valueUp)
	case reflect.Uint8:
		val = *(*uint8)(valueUp)
	case reflect.Uint16:
		val = *(*uint16)(valueUp)
	case reflect.Uint32:
		val = *(*uint32)(valueUp)
	case reflect.Uint64:
		val = *(*uint64)(valueUp)
	case reflect.Uintptr:
		val = *(*uintptr)(valueUp)
	case reflect.Float32:
		val = *(*float32)(valueUp)
	case reflect.Float64:
		val = *(*float64)(valueUp)
	case reflect.Complex64:
		val = *(*complex64)(valueUp)
	case reflect.Complex128:
		val = *(*complex128)(valueUp)
	case reflect.String:
		val = *(*string)(valueUp)
	default:
		panic("Я не придумал, как сделать для этого типа")
		// val = *(*slice)(valueUp)

	}

	return val
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func (b *Bmap) overflow(t *Maptype) *Bmap {
	const PtrSize = 4 << (^uintptr(0) >> 63)
	return *(**Bmap)(add(unsafe.Pointer(b), uintptr(t.Bucketsize)-PtrSize))
}

type Hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.// Make sure this stays in sync with the compiler's definition.
	Count     int // # live cells == size of map. Must be first (used by len() builtin)
	Flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	Noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	Hash0     uint32 // hash seed

	Buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	Oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	Nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	Extra *Mapextra // optional fields
}

// mapextra holds fields that are not present on all maps.
type Mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, Bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	Overflow    *[]*Bmap
	Oldoverflow *[]*Bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	NextOverflow *Bmap
}

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
)

// A bucket for a Go map.
type Bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	Tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

const DataOffset = unsafe.Offsetof(struct {
	b Bmap
	v int64
}{}.v)

type Maptype struct {
	Typ    interfacetype.T_type
	Key    *interfacetype.T_type
	Elem   *interfacetype.T_type
	Bucket *interfacetype.T_type // internal type representing a hash bucket
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	Keysize    uint8  // size of key slot
	Elemsize   uint8  // size of elem slot
	Bucketsize uint16 // size of bucket
	Flags      uint32
}
