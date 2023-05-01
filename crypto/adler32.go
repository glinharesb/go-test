package crypto

import (
	"hash/adler32"
)

func CalculateAdler32Checksum(buffer []byte, offset int, size int) uint32 {
	h := adler32.New()
	h.Write(buffer[offset : offset+size])

	return h.Sum32()
}
