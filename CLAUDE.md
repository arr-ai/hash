# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`github.com/arr-ai/hash` is a Go library providing fast non-cryptographic hash functions for use in hash tables. It exposes type-specific hash functions for all Go primitives, reflection-based hashing for structs/arrays, and a `Hashable` interface for custom types.

## Build Commands

```bash
GOWORK=off go build ./...               # Build (GOWORK=off if parent go.work has older Go version)
GOWORK=off go test ./...                # Run all tests
GOWORK=off go test -run TestHash        # Run a single test
GOWORK=off go test -bench .             # Run benchmarks
golangci-lint run                       # Lint (27 linters configured in .golangci.yml)
```

## Architecture

Single-package library (`package hash`) with no subpackages. All public API is in the root.

**Key files:**
- `hash.go` — Public API: type-dispatching `Any()` function plus typed functions (`Bool`, `Int`, `String`, `Float64`, etc.)
- `hashable.go` — `Hashable` interface (`Hash(seed Seed) Seed`)
- `alg.go` — Algorithm selection and initialization; runtime detection of AES-NI support on x86/ARM64
- `hash32.go` / `hash64.go` — Platform-width-specific hash implementations (build-tagged for 32-bit vs 64-bit architectures)
- `stubs.go` — Low-level utilities (`add`, `fastrand`, `noescape`)
- `arch_*.go` — Per-architecture constants (endianness)
- `asm_*.s` — Assembly implementations of AES-NI accelerated hashing for each architecture

**Hash algorithm selection** (two-tier, chosen at init):
1. AES-NI hardware path — used when CPU supports AES+SSSE3+SSE4.1 (x86) or AES (ARM64); implemented in assembly
2. Fallback pure-Go path — xxhash/cityhash-based; used on all other platforms

## Dependencies

- `golang.org/x/sys` — CPU feature detection for AES-NI capability

## Licensing

Dual-licensed: Apache 2.0 (`LICENSE`) for original code, BSD (`GO-LICENSE`) for Go-derived portions.
