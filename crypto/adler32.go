package crypto

import (
	"hash/adler32"
)

// CalculateAdler32Checksum
func Adler32(buffer []byte, offset int, size int) uint32 {
	// create a new Adler-32 hash
	h := adler32.New()

	// write the buffer to the hash from the specified offset and size
	h.Write(buffer[offset : offset+size])

	// return the hash value as a uint32
	return h.Sum32()
}
