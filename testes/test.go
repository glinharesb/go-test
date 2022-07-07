package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"unsafe"
)

func main() {
	buffer, err := ioutil.ReadFile("buffer_antes")
	if err != nil {
		log.Fatal(err)
	}

	buffer = buffer[:24]

	fmt.Println("[!] Decrypted Original: ")
	fmt.Println(hex.Dump(buffer[:50]))

	key := [4]uint32{3727812935, 1206674243, 4196718677, 2568959323}
	xteaEncrypt(&buffer, key)
	fmt.Println("[!] Encrypted by func: ")
	fmt.Println(hex.Dump(buffer[:50]))

	bufferOriginal, err2 := ioutil.ReadFile("buffer_depois")
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("[!] Encrypted Original: ")
	fmt.Println(hex.Dump(bufferOriginal[:50]))
}

func xteaEncrypt(buffer *[]byte, key [4]uint32) {
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