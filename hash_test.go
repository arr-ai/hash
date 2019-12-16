package hash

import (
	"math"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestHash64(t *testing.T) {
	if Interface(uint64(0), 0) == 0 {
		t.Error()
	}
}

func TestHash64String(t *testing.T) {
	if Interface("hello", 0) == 0 {
		t.Error()
	}
}

func TestHashMatchesEquality(t *testing.T) {
	t.Logf("%d unique elements", len(cornucopia))
	total := 0
	falsePositives := 0
	for _, seeds := range [][2]uintptr{{0, 0}, {0, 1}, {0, 2}, {1, 1}, {1, 2}} {
		aSeed := seeds[0]
		bSeed := seeds[1]
		for _, a := range cornucopia {
			for _, b := range cornucopia {
				if aSeed == bSeed && a == b {
					assert.Equal(t, Interface(a, aSeed), Interface(b, bSeed),
						"a=%v b=%v hash(a)=%v hash(b)=%v",
						a, b, Interface(a, aSeed), Interface(b, aSeed))
				} else if Interface(a, aSeed) == Interface(b, bSeed) {
					h := Interface(a, aSeed)
					_ = Interface(b, bSeed)
					t.Logf("\nhash(%#v %[1]T, %v) ==\nhash(%#v %[3]T, %v) == %d",
						a, aSeed, b, bSeed, h)
					falsePositives++
				}
				total++
			}
		}
	}
	assert.LessOrEqual(t, falsePositives, total/100, total)
}

func BenchmarkHash(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		Interface(cornucopia[r.Int()%len(cornucopia)], 0)
	}
}

var cornucopia = func() []interface{} {
	x := 42
	result := []interface{}{
		false,
		true,
		&x,
		&[]int{43}[0],
		&[]string{"hello"}[0],
		uintptr(unsafe.Pointer(&x)),
		unsafe.Pointer(nil),
		unsafe.Pointer(&x),
		unsafe.Pointer(uintptr(unsafe.Pointer(&x))),
		[...]int{},
		[...]int{1, 2, 3, 4, 5},
		[...]int{5, 4, 3, 2, 1},
	}

	// The following number lists are massive overkill, but it can't hurt.

	for _, i := range []int64{
		-43, -42, -10, -1, 0, 1, 10, 42,
		math.MaxInt64, math.MaxInt64 - 1,
		math.MinInt64, math.MinInt64 + 1,
	} {
		result = append(result, int(i), int8(i), int16(i), int32(i), i)
	}

	for _, i := range []uint64{0, 42} {
		result = append(result, uint(i), uint8(i), uint16(i), uint32(i), i)
	}

	floats := []float64{
		0, 42, math.Pi,
		math.MaxFloat32, math.SmallestNonzeroFloat32,
		math.MaxFloat64, math.SmallestNonzeroFloat64,
	}

	for _, f := range floats {
		result = append(result, float32(f), f)
	}

	for _, re := range floats {
		for _, im := range floats {
			result = append(result, complex(float32(re), float32(im)))
			result = append(result, complex(re, im))
		}
	}

	for _, s := range []string{
		"",
		"a",
		"b",
		"hello",
		"-------------------------------------------------------",
		"--------------------------------------------------------",
		"--------------------------------------------------------\000",
	} {
		result = append(result, s)
	}

	// Dedupe
	m := map[interface{}]struct{}{}
	for _, i := range result {
		m[i] = struct{}{}
	}

	for i := range m {
		result = append(result, i)
	}

	return result
}()
