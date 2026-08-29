//go:build amd64

package hash128

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

// Implemented in asm_amd64.s. The routines only read through p during the
// call, so callers may pass pointers to locals without them escaping.
//
//go:noescape
func aeshash32H128(p unsafe.Pointer) H128

//go:noescape
func aeshash64H128(p unsafe.Pointer) H128

//go:noescape
func aeshashH128(p unsafe.Pointer, s uintptr) H128

//go:noescape
func aeshashstrH128(p unsafe.Pointer) H128

func useAES() bool {
	return cpu.X86.HasAES && // AESENC
		cpu.X86.HasSSSE3 && // PSHUFB
		cpu.X86.HasSSE41 // PINSR{D,Q}
}
