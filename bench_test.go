package hash_test

import (
	"math"
	"testing"
	"unsafe"

	"github.com/arr-ai/hash"
)

// hashableSeed implements hash.Hashable for benchmarking the Hashable path.
type hashableSeed struct{ v int }

func (h hashableSeed) Hash(seed uintptr) uintptr {
	return hash.Int(h.v, seed)
}

// Primitive types via Any

func BenchmarkAny_Bool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(true, 0)
	}
}

func BenchmarkAny_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(42, 0)
	}
}

func BenchmarkAny_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(int8(42), 0)
	}
}

func BenchmarkAny_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(int64(42), 0)
	}
}

func BenchmarkAny_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(uint64(42), 0)
	}
}

func BenchmarkAny_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(math.Pi, 0)
	}
}

func BenchmarkAny_String_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any("hello", 0)
	}
}

func BenchmarkAny_String_Long(b *testing.B) {
	s := "--------------------------------------------------------"
	for i := 0; i < b.N; i++ {
		hash.Any(s, 0)
	}
}

func BenchmarkAny_Uintptr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Any(uintptr(42), 0)
	}
}

func BenchmarkAny_UnsafePointer(b *testing.B) {
	x := 42
	p := unsafe.Pointer(&x)
	for i := 0; i < b.N; i++ {
		hash.Any(p, 0)
	}
}

func BenchmarkAny_Complex128(b *testing.B) {
	c := complex(math.Pi, math.E)
	for i := 0; i < b.N; i++ {
		hash.Any(c, 0)
	}
}

// Hashable interface path

func BenchmarkAny_Hashable(b *testing.B) {
	h := hashableSeed{42}
	for i := 0; i < b.N; i++ {
		hash.Any(h, 0)
	}
}

// Reflection paths (types not in type switch)

func BenchmarkAny_Struct(b *testing.B) {
	type S struct {
		X int
		Y string
	}
	s := S{42, "hello"}
	for i := 0; i < b.N; i++ {
		hash.Any(s, 0)
	}
}

func BenchmarkAny_Array(b *testing.B) {
	a := [5]int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		hash.Any(a, 0)
	}
}

func BenchmarkAny_SliceInterface(b *testing.B) {
	s := []interface{}{1, "hello", true}
	for i := 0; i < b.N; i++ {
		hash.Any(s, 0)
	}
}

// Direct typed calls for comparison

func BenchmarkDirect_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Int(42, 0)
	}
}

func BenchmarkDirect_String_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.String("hello", 0)
	}
}

func BenchmarkDirect_String_Long(b *testing.B) {
	s := "--------------------------------------------------------"
	for i := 0; i < b.N; i++ {
		hash.String(s, 0)
	}
}

func BenchmarkDirect_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Uint64(42, 0)
	}
}

func BenchmarkDirect_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Float64(math.Pi, 0)
	}
}

func BenchmarkDirect_Bool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.Bool(true, 0)
	}
}
