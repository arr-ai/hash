package hash128

import (
	"crypto/rand"
	"unsafe"

	"github.com/arr-ai/hash"
)

// useAsm is true when the single-pass AES assembly is in use. It is bound
// once at init; tests may clear it to exercise the software path.
var useAsm bool

// aeskeysched is the per-process key schedule read by the assembly. It is
// filled with random bytes at init so collisions are hard to engineer.
var aeskeysched [128]byte

// zeroFloat is the hash shared by +0 and -0.
var zeroFloat H128

func init() {
	if useAES() {
		if _, err := rand.Read(aeskeysched[:]); err != nil {
			panic(err)
		}
		useAsm = true
	}
	var f float64
	zeroFloat = mem64(unsafe.Pointer(&f)).Mix(H128{golden, mixSalt})
}

// nanFloat returns a fresh hash for a NaN, so NaN != NaN holds for hashing
// as it does for comparison.
func nanFloat() H128 {
	var r [8]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	return memN(unsafe.Pointer(&r[0]), uintptr(len(r)))
}

func mem32(p unsafe.Pointer) H128 {
	if useAsm {
		return aeshash32H128(p)
	}
	return softMem(p, 4)
}

func mem64(p unsafe.Pointer) H128 {
	if useAsm {
		return aeshash64H128(p)
	}
	return softMem(p, 8)
}

func memN(p unsafe.Pointer, n uintptr) H128 {
	if useAsm {
		return aeshashH128(p, n)
	}
	return softMem(p, n)
}

func str(p unsafe.Pointer) H128 {
	if useAsm {
		return aeshashstrH128(p)
	}
	return softString(*(*string)(p))
}

type stringHeader struct {
	data unsafe.Pointer
	len  int
}

// The software fallback routes through two seeded calls of the parent
// package, which has its own per-arch fast paths.
func softMem(p unsafe.Pointer, n uintptr) H128 {
	return softString(*(*string)(unsafe.Pointer(&stringHeader{p, int(n)})))
}

func softString(s string) H128 {
	return H128{uint64(hash.String(s, 0)), uint64(hash.String(s, 1))}
}
