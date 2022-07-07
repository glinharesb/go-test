package network

import "encoding/binary"

func (n *Network) GetBytes(count int) []byte {
	result := n.Buffer[n.Pos : n.Pos+count]
	n.Pos += count
	return result
}

func (n *Network) GetU8() byte {
	result := n.Buffer[0]
	n.Pos += 1
	return result
}

func (n *Network) GetU16() uint16 {
	result := binary.LittleEndian.Uint16(n.Buffer[n.Pos:])
	n.Pos += 2
	return result
}

func (n *Network) GetU32() uint32 {
	result := binary.LittleEndian.Uint32(n.Buffer[n.Pos:])
	n.Pos += 4
	return result
}

func (n *Network) GetString() string {
	length := int(n.GetU16())

	if n.check(length) {
		return ""
	}

	result := n.Buffer[n.Pos : n.Pos+length]
	n.Pos += length
	return string(result)
}
