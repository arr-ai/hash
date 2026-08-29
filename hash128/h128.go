// Package hash128 provides 128-bit hashing computed in a single pass.
//
// On AES-capable amd64 and arm64 hardware the primitives run in assembly
// (derived from the Go runtime's aeshash); elsewhere they fall back to two
// seeded 64-bit hashes from the parent package. Hashes are randomised per
// process and must not be persisted.
//
// Types that can hash themselves implement Hashable. Composite values should
// hash each part once and combine the results with Mix (order-sensitive) or
// Xor (order-independent) rather than re-walking their parts per seed, which
// is the point of a seedless 128-bit interface: a value's hash can be
// computed once and cached.
package hash128

import (
	"fmt"
	"math/bits"
	"reflect"
	"unsafe"

	"github.com/arr-ai/hash"
)

// H128 is a 128-bit hash.
type H128 struct{ Lo, Hi uint64 }

// Hashable represents a type that can evaluate its own 128-bit hash.
type Hashable interface {
	Hash128() H128
}

func (h H128) String() string {
	return fmt.Sprintf("%016x:%016x", h.Lo, h.Hi)
}

// IsZero reports whether h is the zero hash.
func (h H128) IsZero() bool {
	return h.Lo == 0 && h.Hi == 0
}

// Xor combines two hashes symmetrically. Use it to hash unordered
// collections: the result is independent of the order the parts are
// combined in.
func (h H128) Xor(o H128) H128 {
	return H128{h.Lo ^ o.Lo, h.Hi ^ o.Hi}
}

// Mix folds o into h asymmetrically. Use it to hash ordered or keyed
// structures: h.Mix(o) != o.Mix(h) in general, and a.Mix(b).Mix(c) differs
// from a.Mix(c).Mix(b).
func (h H128) Mix(o H128) H128 {
	lo := fmix64(h.Lo ^ (o.Lo*golden + mixSalt))
	hi := fmix64(h.Hi ^ bits.RotateLeft64(o.Hi, 32) ^ lo)
	return H128{lo, hi}
}

// Seeded derives a 64-bit seeded hash from h, for callers that still speak
// the parent package's Hash(seed) interface.
func (h H128) Seeded(seed uintptr) uintptr {
	return uintptr(fmix64(h.Lo^(uint64(seed)*golden+mixSalt)) ^ h.Hi)
}

const (
	golden  = 0x9E3779B97F4A7C15
	mixSalt = 0x2545F4914F6CDD1D
)

// fmix64 is MurmurHash3's 64-bit finaliser.
func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

// Mem hashes n bytes at p.
func Mem(p unsafe.Pointer, n uintptr) H128 {
	return memN(p, n)
}

// Bytes hashes b.
func Bytes(b []byte) H128 {
	if len(b) == 0 {
		return memN(nil, 0)
	}
	return memN(unsafe.Pointer(&b[0]), uintptr(len(b)))
}

// String hashes s. String(s) == Bytes([]byte(s)).
func String(s string) H128 {
	return str(unsafe.Pointer(&s))
}

// Runes hashes r by its UTF-32 contents, without converting to a string.
func Runes(r []rune) H128 {
	if len(r) == 0 {
		return memN(nil, 0)
	}
	return memN(unsafe.Pointer(&r[0]), uintptr(len(r))*unsafe.Sizeof(r[0]))
}

// Bool hashes b.
func Bool(b bool) H128 {
	return memN(unsafe.Pointer(&b), 1)
}

// Int hashes x.
func Int(x int) H128 {
	return memN(unsafe.Pointer(&x), unsafe.Sizeof(x))
}

// Int64 hashes x.
func Int64(x int64) H128 {
	return mem64(unsafe.Pointer(&x))
}

// Uint hashes x.
func Uint(x uint) H128 {
	return memN(unsafe.Pointer(&x), unsafe.Sizeof(x))
}

// Uint8 hashes x.
func Uint8(x uint8) H128 {
	return memN(unsafe.Pointer(&x), 1)
}

// Uint16 hashes x.
func Uint16(x uint16) H128 {
	return memN(unsafe.Pointer(&x), 2)
}

// Uint32 hashes x.
func Uint32(x uint32) H128 {
	return mem32(unsafe.Pointer(&x))
}

// Uint64 hashes x.
func Uint64(x uint64) H128 {
	return mem64(unsafe.Pointer(&x))
}

// Uintptr hashes x.
func Uintptr(x uintptr) H128 {
	return memN(unsafe.Pointer(&x), unsafe.Sizeof(x))
}

// Float32 hashes f. +0 and -0 hash equal; every NaN hashes differently.
func Float32(f float32) H128 {
	switch {
	case f == 0:
		return zeroFloat
	case f != f:
		return nanFloat()
	}
	return mem32(unsafe.Pointer(&f))
}

// Float64 hashes f. +0 and -0 hash equal; every NaN hashes differently.
func Float64(f float64) H128 {
	switch {
	case f == 0:
		return zeroFloat
	case f != f:
		return nanFloat()
	}
	return mem64(unsafe.Pointer(&f))
}

// Any hashes i by dynamic type. Hashable values hash themselves; values
// implementing only the parent package's seeded hash.Hashable are hashed
// with seeds 0 and 1; other types are hashed by reflection via the parent
// package.
func Any(i interface{}) H128 { //nolint:cyclop
	switch k := i.(type) {
	case Hashable:
		return k.Hash128()
	case hash.Hashable:
		return H128{uint64(k.Hash(0)), uint64(k.Hash(1))}
	case string:
		return String(k)
	case int:
		return Int(k)
	case int64:
		return Int64(k)
	case uint64:
		return Uint64(k)
	case uint:
		return Uint(k)
	case uintptr:
		return Uintptr(k)
	case float64:
		return Float64(k)
	case bool:
		return Bool(k)
	case []byte:
		return Bytes(k)
	case reflect.Value:
		return H128{uint64(hash.Value(k, 0)), uint64(hash.Value(k, 1))}
	default:
		return H128{uint64(hash.Any(i, 0)), uint64(hash.Any(i, 1))}
	}
}
