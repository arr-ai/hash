# hash

Hash is a fast non-cryptographic hash function library, primarily intended for
use in hash tables. Significant portions of this library are lifted from Go
internals.

## Documentation

https://godoc.org/github.com/arr-ai/hash

## 128-bit hashing

The `hash128` subpackage computes 128-bit hashes in a single pass, using AES
instructions on amd64 and arm64 and falling back to two seeded 64-bit hashes
elsewhere. It is seedless: a type implements `hash128.Hashable` by returning
its `H128` once, and composite values combine part hashes with `Mix`
(order-sensitive) or `Xor` (order-independent) instead of re-hashing per seed.
The seeded `Hashable` interface in the root package is unchanged; `hash128.Any`
accepts implementations of either.


## License

Portions of this software are taken from the Go source code at
https://github.com/golang/go. These portions are licensed as per the LICENSE-GO
file.

All other sources are licensed under the Apache 2.0 license as per the LICENSE
file.
