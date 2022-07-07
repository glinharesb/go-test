package network

import (
	"go-test/crypto"
	"log"
)

type Network struct {
	Buffer []byte
	Pos    int
	Header int
}

func (n *Network) check(count int) bool {
	return n.Pos+count >= len(n.Buffer)
}

func (n *Network) XteaEncrypt(key [4]uint32) {
	n.AddSize()

	if (n.Pos-n.Header)%8 != 0 {
		toAdd := 8 - ((n.Pos - n.Header) % 8)
		for i := 0; i < toAdd; i++ {
			n.AddU8(0x33)
		}
	}

	// must have 8 reserved bytes
	if n.Header != 8 {
		log.Fatalf("invalid header size: %d", n.Header)
	}

	crypto.XteaEncrypt(&n.Buffer, key)
}
