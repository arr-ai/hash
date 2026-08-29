package hash128

import (
	"math"
	"testing"
	"unsafe"
)

type selfHasher struct{ h H128 }

func (s selfHasher) Hash128() H128 { return s.h }

type seededHasher uint64

func (s seededHasher) Hash(seed uintptr) uintptr { return uintptr(s) + seed }

func assertEqual(t *testing.T, a, b H128, msg string) {
	t.Helper()
	if a != b {
		t.Errorf("%s: %v != %v", msg, a, b)
	}
}

func assertDistinct(t *testing.T, a, b H128, msg string) {
	t.Helper()
	if a == b {
		t.Errorf("%s: %v == %v", msg, a, b)
	}
}

func TestDeterministicWithinProcess(t *testing.T) {
	assertEqual(t, String("hello"), String("hello"), "String")
	assertEqual(t, Uint64(42), Uint64(42), "Uint64")
	assertEqual(t, Float64(1.5), Float64(1.5), "Float64")
	assertEqual(t, Bytes([]byte{1, 2, 3}), Bytes([]byte{1, 2, 3}), "Bytes")
}

func TestStringEqualsBytes(t *testing.T) {
	for _, s := range []string{"", "a", "0123456789abcdef", "0123456789abcdefg",
		"a string somewhat longer than thirty-two bytes", string(make([]byte, 200))} {
		assertEqual(t, String(s), Bytes([]byte(s)), "String vs Bytes: "+s)
		hdr := (*stringHeader)(unsafe.Pointer(&s))
		assertEqual(t, String(s), Mem(hdr.data, uintptr(hdr.len)), "String vs Mem: "+s)
	}
}

func TestDistinctInputsDistinctHashes(t *testing.T) {
	assertDistinct(t, String("a"), String("b"), "strings")
	assertDistinct(t, String(""), String("\x00"), "empty vs NUL")
	assertDistinct(t, Uint64(1), Uint64(2), "uint64")
	assertDistinct(t, Uint32(1), Uint64(1), "uint32 vs uint64 width")
	assertDistinct(t, Int(7), String("7"), "int vs string")
	// Every length class of the variable-length body.
	seen := map[H128]int{}
	for n := 0; n <= 300; n++ {
		h := Bytes(make([]byte, n))
		if m, dup := seen[h]; dup {
			t.Errorf("zero buffers of length %d and %d collide", m, n)
		}
		seen[h] = n
	}
}

func TestRunesHashesContent(t *testing.T) {
	assertEqual(t, Runes([]rune("héllo")), Runes([]rune("héllo")), "same runes")
	assertDistinct(t, Runes([]rune("héllo")), Runes([]rune("hello")), "different runes")
	assertEqual(t, Runes(nil), Runes([]rune{}), "nil vs empty")
}

func TestFloats(t *testing.T) {
	assertEqual(t, Float64(0), Float64(math.Copysign(0, -1)), "+0 == -0")
	assertEqual(t, Float32(0), Float32(float32(math.Copysign(0, -1))), "+0 == -0 (32)")
	assertDistinct(t, Float64(0), Uint64(0), "zero float vs zero int")
	assertDistinct(t, Float64(math.NaN()), Float64(math.NaN()), "NaN != NaN")
	assertEqual(t, Float64(2.5), Float64(2.5), "ordinary float")
}

func TestCombinators(t *testing.T) {
	a, b, c := String("a"), String("b"), String("c")
	assertEqual(t, a.Xor(b), b.Xor(a), "Xor commutes")
	assertEqual(t, a.Xor(b).Xor(c), c.Xor(b).Xor(a), "Xor associates")
	assertEqual(t, a.Xor(a), H128{}, "Xor self-inverse")
	assertDistinct(t, a.Mix(b), b.Mix(a), "Mix is order-sensitive")
	assertDistinct(t, a.Mix(b).Mix(c), a.Mix(c).Mix(b), "Mix chains are order-sensitive")
	assertEqual(t, a.Mix(b), a.Mix(b), "Mix deterministic")
	assertDistinct(t, a.Mix(b), a.Xor(b), "Mix differs from Xor")
	if a.Seeded(0) == a.Seeded(1) {
		t.Error("Seeded should vary with seed")
	}
	seeded := a.Seeded(3)
	if a.Seeded(3) != seeded {
		t.Error("Seeded should be deterministic")
	}
	if !(H128{}).IsZero() || a.IsZero() {
		t.Error("IsZero")
	}
}

func TestAnyDispatch(t *testing.T) {
	h := H128{1, 2}
	assertEqual(t, Any(selfHasher{h}), h, "Hashable")
	assertEqual(t, Any(seededHasher(5)), H128{5, 6}, "legacy seeded Hashable")
	assertEqual(t, Any("s"), String("s"), "string")
	assertEqual(t, Any(3), Int(3), "int")
	assertEqual(t, Any(uint64(3)), Uint64(3), "uint64")
	assertEqual(t, Any(1.5), Float64(1.5), "float64")
	assertEqual(t, Any(true), Bool(true), "bool")
	assertEqual(t, Any([]byte("b")), Bytes([]byte("b")), "[]byte")
	assertEqual(t, Any(struct{ A, B int }{1, 2}), Any(struct{ A, B int }{1, 2}), "reflected struct")
	assertDistinct(t, Any(struct{ A, B int }{1, 2}), Any(struct{ A, B int }{2, 1}), "reflected struct order")
}

func TestSoftwareFallbackAgreesWithItself(t *testing.T) {
	saved := useAsm
	defer func() { useAsm = saved }()
	useAsm = false
	assertEqual(t, String("x"), Bytes([]byte("x")), "software String vs Bytes")
	assertEqual(t, Uint64(9), Uint64(9), "software deterministic")
	assertDistinct(t, Uint64(9), Uint64(10), "software distinct")
	for n := 0; n <= 40; n++ {
		assertEqual(t, Bytes(make([]byte, n)), Mem(unsafe.Pointer(&make([]byte, n+1)[0]), uintptr(n)), "software Mem")
	}
}

func BenchmarkString32(b *testing.B) {
	s := "0123456789abcdef0123456789abcdef"
	for i := 0; i < b.N; i++ {
		String(s)
	}
}

func BenchmarkUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64(uint64(i))
	}
}

func BenchmarkMix(b *testing.B) {
	h, o := String("a"), String("b")
	for i := 0; i < b.N; i++ {
		h = h.Mix(o)
	}
}
