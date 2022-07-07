package network

import (
	"fmt"
	"go-test/config"
	"go-test/crypto"
	"net"
	"time"
)

type Packet struct {
	Conn    net.Conn
	XteaKey [4]uint32
}

func (p *Packet) ParsePacket(packet []byte) {
	msg := Network{
		Buffer: packet,
		Pos:    0,
	}

	packetSize := msg.GetU16()
	fmt.Printf("[!] packetSize: %d\n", packetSize)

	// FIX
	checksum := msg.GetU32()
	fmt.Printf("[!] checksum: %d\n", checksum)

	// FIX
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

		p.XteaKey[0] = decryptedMsg.GetU32()
		p.XteaKey[1] = decryptedMsg.GetU32()
		p.XteaKey[2] = decryptedMsg.GetU32()
		p.XteaKey[3] = decryptedMsg.GetU32()
	}

	fmt.Printf("[!] Xtea: %v\n", p.XteaKey)

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

	VERSION_MIN := uint16(config.ConfigInstance.VersionMin)
	VERSION_MAX := uint16(config.ConfigInstance.VersionMax)
	if (VERSION_MIN == VERSION_MAX) && (version != VERSION_MIN) {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol %s allowed!", FormatVersion(VERSION_MIN)), version)
	} else if version < VERSION_MIN || version > VERSION_MAX {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol between %s and %s allowed!",
			FormatVersion(VERSION_MIN),
			FormatVersion(VERSION_MAX)),
			version)
	}

	p.GetCharacterList(accountName, password, "", version)
}

func (p *Packet) DisconnectClient(msg string, version uint16) {
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

	p.SendPacket(outputMsg)
}

func (p *Packet) SendPacket(outputMsg Network) {
	outputMsg.XteaEncrypt(p.XteaKey)

	outputMsg.AddChecksum()

	outputMsg.AddSize()

	outputMsg.Buffer = outputMsg.Buffer[outputMsg.Header:outputMsg.Pos]
	p.Conn.Write(outputMsg.Buffer)
}

func (p *Packet) GetCharacterList(accountName string, password string, token string, version uint16) {
	outputMsg := Network{
		Buffer: make([]byte, 8192),
		Pos:    0,
		Header: 10,
	}
	outputMsg.Pos = outputMsg.Header

	motd := "Hello!"
	characters := []string{"Account Manager"}

	// motd
	outputMsg.AddU8(0x14)
	outputMsg.AddString(fmt.Sprintf("1\n%s", motd))

	// session key
	if version >= 1074 {
		dateNow := time.Now().Unix()
		outputMsg.AddU8(0x28)
		outputMsg.AddString(fmt.Sprintf("%s\n%s\n%s\n%d", accountName, password, token, dateNow))
	}

	outputMsg.AddU8(0x64)

	if version >= 1010 {
		numberOfWorlds := 2
		outputMsg.AddU8(byte(numberOfWorlds))

		for i := 0; i < numberOfWorlds; i++ {
			outputMsg.AddU8(byte(i))
			outputMsg.AddString("Offline")
			outputMsg.AddString(config.ConfigInstance.GameIp)
			outputMsg.AddU16(uint16(config.ConfigInstance.GamePort))
			outputMsg.AddU8(0)
		}

		outputMsg.AddU8(byte(len(characters)))

		for _, character := range characters {
			outputMsg.AddU8(0)
			outputMsg.AddString(character)
		}
	} else {
		outputMsg.AddU8(byte(len(characters)))

		for _, character := range characters {
			outputMsg.AddString(character)

			outputMsg.AddString("Teste")
			outputMsg.AddU32(Ip2int(config.ConfigInstance.GameIp))
			outputMsg.AddU16(uint16(config.ConfigInstance.GamePort))

			if version >= 980 {
				outputMsg.AddU8(0)
			}
		}
	}

	// premium
	if version >= 1077 {
		outputMsg.AddU8(0)
		outputMsg.AddU8(0)
		outputMsg.AddU32(0)
	} else {
		outputMsg.AddU16(0)
	}

	p.SendPacket(outputMsg)
}
