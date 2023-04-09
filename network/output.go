package network

import (
	"encoding/binary"
	"go-test/crypto"
)

func (n *Network) AddU8(value byte) {
	n.Buffer[n.Pos] = value
	n.Pos++
}

func (n *Network) AddU16(value uint16) {
	binary.LittleEndian.PutUint16(n.Buffer[n.Pos:], value)
	n.Pos += 2
}

func (n *Network) AddU32(value uint32) {
	binary.LittleEndian.PutUint32(n.Buffer[n.Pos:], value)
	n.Pos += 4
}

func (n *Network) AddString(str string) {
	// string length
	n.AddU16(uint16(len(str)))

	for _, letter := range []byte(str) {
		n.AddU8(letter)
	}
}

func (n *Network) AddSize() {
	binary.LittleEndian.PutUint16(n.Buffer[n.Header-2:], uint16(n.Pos-n.Header))
	n.Header -= 2
}

func (n *Network) AddChecksum() {
	binary.LittleEndian.PutUint32(n.Buffer[n.Header-4:], crypto.CalculateAdler32Checksum(n.Buffer, n.Header, n.Pos-n.Header))
	n.Header -= 4
}
