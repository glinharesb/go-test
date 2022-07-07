package network

import (
	"fmt"
	"go-test/crypto"
	"net"
)

const VERSION_MIN = 860
const VERSION_MAX = 1000

type Packet struct {
	Conn net.Conn
	Xtea [4]uint32
}

func (p *Packet) ParsePacket(packet []byte) {
	msg := Network{
		Buffer: packet,
		Pos:    0,
	}

	packetSize := msg.GetU16()
	fmt.Printf("[!] packetSize: %d\n", packetSize)

	checksum := msg.GetU32()
	fmt.Printf("[!] checksum: %d\n", checksum)

	adler := crypto.Adler32(msg.Buffer, msg.Pos, len(msg.Buffer)-msg.Pos)
	fmt.Printf("[!] adler: %d\n", adler)

	packetType := msg.GetU8()
	fmt.Printf("[!] packetType: %d\n", packetType)

	msg.GetU16()

	version := msg.GetU16()
	fmt.Printf("[!] version: %d\n", version)

	if version >= 980 {
		msg.GetU32()
	}

	if version >= 1071 {
		msg.GetU16()
		msg.GetU16()
	} else {
		msg.GetU32()
	}

	sprSignature := msg.GetU32()
	fmt.Printf("[!] sprSignature: %d\n", sprSignature)

	picSignature := msg.GetU32()
	fmt.Printf("[!] picSignature: %d\n", picSignature)

	if version >= 980 {
		previewState := msg.GetU8()
		fmt.Printf("[!] previewState: %d\n", previewState)
	}

	var decryptedMsg Network
	if version >= 770 {
		decryptedMsg = Network{
			Buffer: msg.GetBytes(128),
			Pos:    0,
		}
		decryptedMsg.Buffer = crypto.RsaInstance.Decrypt(decryptedMsg.Buffer)

		p.Xtea[0] = decryptedMsg.GetU32()
		p.Xtea[1] = decryptedMsg.GetU32()
		p.Xtea[2] = decryptedMsg.GetU32()
		p.Xtea[3] = decryptedMsg.GetU32()
	}

	fmt.Printf("[!] Xtea: %v\n", p.Xtea)

	accountName := decryptedMsg.GetString()
	fmt.Printf("[!] accountName: %s\n", accountName)

	password := decryptedMsg.GetString()
	fmt.Printf("[!] password: %s\n", password)

	if version >= 1061 {
		msg.GetU8()
		msg.GetU8()

		gpu := msg.GetString()
		fmt.Printf("[!] gpu: %s\n", gpu)

		gpuVersion := msg.GetString()
		fmt.Printf("[!] gpuVersion: %s\n", gpuVersion)
	}

	var accountToken string
	if version >= 1072 {
		decryptAuthPacket := Network{
			Buffer: msg.GetBytes(128),
			Pos:    0,
		}
		decryptAuthPacket.Buffer = crypto.RsaInstance.Decrypt(decryptAuthPacket.Buffer)

		accountToken = decryptAuthPacket.GetString()
		fmt.Printf("[!] accountToken: %s\n", accountToken)

		if version >= 1074 {
			decryptAuthPacket.GetU8()
		}
	}

	// if (VERSION_MIN == VERSION_MAX) && (version != VERSION_MIN) {
	// 	p.disconnectClient(fmt.Sprintf("Only clients with protocol %d allowed!", VERSION_MIN), version)
	// } else if version < VERSION_MIN || version > VERSION_MAX {
	// 	p.disconnectClient(fmt.Sprintf("Only clients with protocol between %d and %d allowed!", VERSION_MIN, VERSION_MAX), version)
	// }

	p.disconnectClient("Teste123u89u89u89u89u98 u89u89u98u89", version)
}

func (p *Packet) disconnectClient(msg string, version uint16) {
	outputMsg := Network{
		Buffer: make([]byte, 8192),
		Pos:    0,
		Header: 10,
	}
	outputMsg.Pos = outputMsg.Header

	if version >= 1076 {
		outputMsg.AddU8(0x0b)
	} else {
		outputMsg.AddU8(0x0a)
	}

	outputMsg.AddString(msg)

	p.sendPacket(outputMsg)
}

func (p *Packet) sendPacket(outputMsg Network) {
	// xtea encrypt
	outputMsg.XteaEncrypt(p.Xtea)

	// add checksum
	outputMsg.AddChecksum()

	// add size
	outputMsg.AddSize()

	// write
	outputMsg.Buffer = outputMsg.Buffer[outputMsg.Header:outputMsg.Pos]

	p.Conn.Write(outputMsg.Buffer)
}
