// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

// func aeshash32(p unsafe.Pointer, h uintptr) uintptr
TEXT ·aeshash32(SB),NOSPLIT|NOFRAME,$0-24
	MOVD	p+0(FP), R0
	MOVD	h+8(FP), R1
	MOVD	$ret+16(FP), R2
	MOVD	$·aeskeysched+0(SB), R3

	VEOR	V0.B16, V0.B16, V0.B16
	VLD1	(R3), [V2.B16]
	VLD1	(R0), V0.S[1]
	VMOV	R1, V0.S[0]

	AESE	V2.B16, V0.B16
	AESMC	V0.B16, V0.B16
	AESE	V2.B16, V0.B16
	AESMC	V0.B16, V0.B16
	AESE	V2.B16, V0.B16

	VST1	[V0.D1], (R2)
	RET

// func aeshash64(p unsafe.Pointer, h uintptr) uintptr
TEXT ·aeshash64(SB),NOSPLIT|NOFRAME,$0-24
	MOVD	p+0(FP), R0
	MOVD	h+8(FP), R1
	MOVD	$ret+16(FP), R2
	MOVD	$·aeskeysched+0(SB), R3

	VEOR	V0.B16, V0.B16, V0.B16
	VLD1	(R3), [V2.B16]
	VLD1	(R0), V0.D[1]
	VMOV	R1, V0.D[0]

	AESE	V2.B16, V0.B16
	AESMC	V0.B16, V0.B16
	AESE	V2.B16, V0.B16
	AESMC	V0.B16, V0.B16
	AESE	V2.B16, V0.B16

	VST1	[V0.D1], (R2)
	RET

// func aeshash(p unsafe.Pointer, h, size uintptr) uintptr
TEXT ·aeshash(SB),NOSPLIT|NOFRAME,$0-32
	MOVD	p+0(FP), R0
	MOVD	s+16(FP), R1
	MOVWU	h+8(FP), R3
	MOVD	$ret+24(FP), R2
	B	aeshashbody<>(SB)

// func aeshashstr(p unsafe.Pointer, h uintptr) uintptr
TEXT ·aeshashstr(SB),NOSPLIT|NOFRAME,$0-24
	MOVD	p+0(FP), R10 // string pointer
	LDP	(R10), (R0, R1) //string data/ length
	MOVWU	h+8(FP), R3
	MOVD	$ret+16(FP), R2 // return adddress
	B	aeshashbody<>(SB)
