package hash

// Hashable represents a type that can evaluate its own hash.
type Hashable interface {
	Hash(seed uintptr) uintptr
}
