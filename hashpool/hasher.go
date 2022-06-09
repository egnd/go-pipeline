package hashpool

import (
	"github.com/cespare/xxhash"
)

// Hasher is a hash function type.
type Hasher func(string, uint64) uint64

// DefaultHasher is a default hashing method.
func DefaultHasher(val string, cnt uint64) uint64 {
	return xxhash.Sum64String(val) % cnt
}
