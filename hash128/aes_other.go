//go:build !amd64 && !arm64

package hash128

import "unsafe"

// No single-pass assembly on this architecture; the software fallback is
// always used and these are never called.
func useAES() bool { return false }

func aeshash32H128(unsafe.Pointer) H128        { panic("unreachable") }
func aeshash64H128(unsafe.Pointer) H128        { panic("unreachable") }
func aeshashH128(unsafe.Pointer, uintptr) H128 { panic("unreachable") }
func aeshashstrH128(unsafe.Pointer) H128       { panic("unreachable") }
