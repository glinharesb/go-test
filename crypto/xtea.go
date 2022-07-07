package crypto

import "unsafe"

func XteaEncrypt(buffer *[]byte, key [4]uint32) {
	const delta = 0x9e3779b9
	const num_rounds = 32

	var u32 []uint32 = (*(*[]uint32)(unsafe.Pointer(buffer)))

	for i := 2; i < len(u32)/4; i += 2 {
		u32[0] = 0

		for j := 0; j < num_rounds; j++ {
			u32[i] += (((u32[i+1] << 4) ^ (u32[i+1] >> 5)) + u32[i+1]) ^ (u32[0] + key[u32[0]&3])
			u32[0] += delta
			u32[i+1] += (((u32[i] << 4) ^ (u32[i] >> 5)) + u32[i]) ^ (u32[0] + key[(u32[0]>>11)&3])
		}
	}
}
