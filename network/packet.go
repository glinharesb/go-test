package network

import (
	"fmt"
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"net"
	"time"
)

type Packet struct {
	Conn        net.Conn
	XteaKey     [4]uint32
	HasChecksum bool
}

func (p *Packet) ParsePacket(packet []byte) {
	msg := Network{
		Buffer: packet,
		Pos:    0,
	}

	msg.GetU16() // packet size

	checksum := msg.PeekU32()
	if checksum == crypto.Adler32(msg.Buffer, msg.Pos+4, len(msg.Buffer)-msg.Pos-4) {
		msg.GetU32() // read checksum
		p.HasChecksum = true
	}

	packetType := msg.GetU8()
	// status check
	if packetType == 0xFF {
		return
	}

	if packetType != 0x01 {
		fmt.Printf("Invalid packet type: %d, should be 1\n", packetType)
		return
	}

	msg.GetU16() // os

	version := msg.GetU16()
	fmt.Printf("[!] version: %d\n", version)

	if version >= 980 {
		msg.GetU32() // read version
	}

	if version >= 1071 {
		msg.GetU16() // content revision
		msg.GetU16() // unkown
	} else {
		msg.GetU32() // data signature
	}

	msg.GetU32() // spr signature
	msg.GetU32() // pic signature

	if version >= 980 {
		msg.GetU8() // preview state
	}

	var decryptedMsg Network
	if version >= 770 {
		decryptedMsg = Network{
			Buffer: msg.GetBytes(128),
			Pos:    0,
		}

		// RSA decrypt
		decryptedMsg.Buffer = crypto.RsaInstance.Decrypt(decryptedMsg.Buffer)

		// XTEA keys
		p.XteaKey[0] = decryptedMsg.GetU32()
		p.XteaKey[1] = decryptedMsg.GetU32()
		p.XteaKey[2] = decryptedMsg.GetU32()
		p.XteaKey[3] = decryptedMsg.GetU32()
	}

	var accountName string
	if version >= 840 {
		accountName = decryptedMsg.GetString()
	} else {
		accountName = fmt.Sprint(decryptedMsg.GetU32())
	}

	password := decryptedMsg.GetString()
	fmt.Printf("[!] accountName: %s\n", accountName)
	fmt.Printf("[!] password: %s\n", password)

	if version >= 1061 {
		msg.GetU8()
		msg.GetU8()

		gpu := msg.GetString()
		fmt.Printf("[!] gpu: %s\n", gpu)

		gpuVersion := msg.GetString()
		fmt.Printf("[!] gpuVersion: %s\n", gpuVersion)
	}

	accountToken := ""
	if version >= 1072 {
		decryptAuthPacket := Network{
			Buffer: msg.GetBytes(128),
			Pos:    0,
		}

		// RSA decrypt
		decryptAuthPacket.Buffer = crypto.RsaInstance.Decrypt(decryptAuthPacket.Buffer)

		accountToken = decryptAuthPacket.GetString()
		fmt.Printf("[!] accountToken: %s\n", accountToken)

		if version >= 1074 {
			decryptAuthPacket.GetU8()
		}
	}

	p.ValidateVersion(version)

	if accountName == "" || accountName == "0" {
		p.DisconnectClient("Invalid account name.", version)
	}

	if password == "" {
		p.DisconnectClient("Invalid password.", version)
	}

	var account database.Account
	if !loginserverAuthentication(accountName, password, &account) {
		p.DisconnectClient("Account name or password is not correct.", version)
	}

	p.GetCharacterList(account, accountToken, version)
}

func (p *Packet) ValidateVersion(version uint16) {
	VERSION_MIN := uint16(config.ConfigInstance.VersionMin)
	VERSION_MAX := uint16(config.ConfigInstance.VersionMax)

	if (VERSION_MIN == VERSION_MAX) && (version != VERSION_MIN) {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol %s allowed!",
			FormatVersion(VERSION_MIN)),
			version)
	} else if version < VERSION_MIN || version > VERSION_MAX {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol between %s and %s allowed!",
			FormatVersion(VERSION_MIN),
			FormatVersion(VERSION_MAX)),
			version)
	}
}

func loginserverAuthentication(accountName string, password string, account *database.Account) bool {
	*account = database.DatabaseInstance.LoadAccountByName(accountName)
	if account.Name == "" {
		return false
	}

	if account.Password != password {
		return false
	}

	return true
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
	if p.XteaKey[0] != 0 {
		outputMsg.XteaEncrypt(p.XteaKey)
	}

	if p.HasChecksum {
		outputMsg.AddChecksum()
	}

	outputMsg.AddSize()

	outputMsg.Buffer = outputMsg.Buffer[outputMsg.Header:outputMsg.Pos]
	p.Conn.Write(outputMsg.Buffer)
}

func (p *Packet) GetCharacterList(account database.Account, token string, version uint16) {
	outputMsg := Network{
		Buffer: make([]byte, 8192),
		Pos:    0,
		Header: 10,
	}
	outputMsg.Pos = outputMsg.Header

	characters := database.DatabaseInstance.LoadCharactersById(account.Id)

	// motd
	outputMsg.AddU8(0x14)
	outputMsg.AddString(fmt.Sprintf("1\n%s", config.ConfigInstance.Motd))

	// session key
	if version >= 1074 {
		outputMsg.AddU8(0x28)
		outputMsg.AddString(fmt.Sprintf("%s\n%s\n%s\n%d", account.Name, account.Password, token, time.Now().Unix()))
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
	if version > 1077 {
		outputMsg.AddU8(0)

		if account.Premdays > 0 {
			outputMsg.AddU8(1)
		} else {
			outputMsg.AddU8(0)
		}

		outputMsg.AddU32(uint32(account.Premdays))
	} else {
		outputMsg.AddU16(uint16(account.Premdays))
	}

	p.SendPacket(outputMsg)
}
