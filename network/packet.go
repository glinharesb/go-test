package network

import (
	"fmt"
	"go-test/config"
	"go-test/crypto"
	"go-test/database"
	"go-test/helpers"
	"log"
	"net"
	"time"
)

type Packet struct {
	Conn        net.Conn
	XteaKey     [4]uint32
	HasChecksum bool
}

func (p *Packet) ParsePacket(packet []byte) error {
	msg := Network{
		Buffer: packet,
		Pos:    0,
	}

	msg.GetU16() // packet size

	checksum := msg.PeekU32()
	if checksum == crypto.CalculateAdler32Checksum(msg.Buffer, msg.Pos+4, len(msg.Buffer)-msg.Pos-4) {
		msg.GetU32() // read checksum
		p.HasChecksum = true
	}

	packetType := msg.GetU8()
	// status check
	if packetType == 0xFF {
		return nil
	}

	if packetType != 0x01 {
		return fmt.Errorf("invalid packet type: %d, should be 1", packetType)
	}

	msg.GetU16() // os

	version := msg.GetU16()
	log.Printf(">> Version: %d\n", version)

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
	} else {
		decryptedMsg = msg
	}

	var accountName string
	if version >= 840 {
		accountName = decryptedMsg.GetString()
	} else {
		accountName = fmt.Sprint(decryptedMsg.GetU32())
	}

	password := decryptedMsg.GetString()
	log.Printf(">> Account name: %s\n", accountName)
	log.Printf(">> Password: %s\n", password)

	if version >= 1061 {
		msg.GetU8()     // ogl info 1
		msg.GetU8()     // ogl info 2
		msg.GetString() // GPU
		msg.GetString() // GPU version
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
		log.Printf(">> Token: %s\n", accountToken)

		if version >= 1074 {
			decryptAuthPacket.GetU8() // stay logged > 0
		}
	}

	p.ValidateVersion(version)

	if len(accountName) == 0 || accountName == "0" {
		p.DisconnectClient("Invalid account name.", version)
	}

	if len(password) == 0 {
		p.DisconnectClient("Invalid password.", version)
	}

	var account database.Account
	if !loginserverAuthentication(accountName, password, &account) {
		p.DisconnectClient("Account name or password is not correct.", version)
	}

	p.GetCharacterList(account, accountToken, version)

	return nil
}

func (p *Packet) ValidateVersion(version uint16) {
	VERSION_MIN := uint16(config.ConfigInstance.VersionMin)
	VERSION_MAX := uint16(config.ConfigInstance.VersionMax)

	if (VERSION_MIN == VERSION_MAX) && (version != VERSION_MIN) {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol %s allowed!",
			helpers.FormatVersion(VERSION_MIN)),
			version)
	} else if version < VERSION_MIN || version > VERSION_MAX {
		p.DisconnectClient(fmt.Sprintf("Only clients with protocol between %s and %s allowed!",
			helpers.FormatVersion(VERSION_MIN),
			helpers.FormatVersion(VERSION_MAX)),
			version)
	}
}

func loginserverAuthentication(accountName string, password string, account *database.Account) bool {
	*account, _ = database.DatabaseInstance.LoadAccountByName(accountName)
	if len(account.Name) == 0 {
		return false
	}

	switch config.ConfigInstance.EncryptionType {
	case "sha1":
		helpers.TransformToSha1(&password)
	case "sha256":
		helpers.TransformToSha256(&password)
	case "sha512":
		helpers.TransformToSha512(&password)
	}

	return account.Password == password
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

	characters, _ := database.DatabaseInstance.LoadCharactersById(account.Id)

	// motd
	if len(config.ConfigInstance.Motd) > 0 {
		outputMsg.AddU8(0x14)
		outputMsg.AddString(fmt.Sprintf("1\n%s", config.ConfigInstance.Motd))
	}

	// session key
	if version >= 1074 {
		outputMsg.AddU8(0x28)
		outputMsg.AddString(fmt.Sprintf("%s\n%s\n%s\n%d", account.Name, account.Password, token, time.Now().Unix()))
	}

	// character list
	outputMsg.AddU8(0x64)

	if version >= 1010 {
		worldsLength := 2
		outputMsg.AddU8(byte(worldsLength)) // worlds quantity

		for i := 0; i < worldsLength; i++ {
			outputMsg.AddU8(byte(i))                                 // world id
			outputMsg.AddString(config.ConfigInstance.ServerName)    // server name, online/offline status or world name
			outputMsg.AddString(config.ConfigInstance.GameIp)        // server ip
			outputMsg.AddU16(uint16(config.ConfigInstance.GamePort)) // server port
			outputMsg.AddU8(0)                                       // world preview
		}

		outputMsg.AddU8(byte(len(characters))) // characters quantity

		for _, character := range characters {
			outputMsg.AddU8(0)             // world id
			outputMsg.AddString(character) // character name
		}
	} else {
		outputMsg.AddU8(byte(len(characters))) // characters quantity

		for _, character := range characters {
			outputMsg.AddString(character)                                 // character name
			outputMsg.AddString(config.ConfigInstance.ServerName)          // server name, online/offline status or world name
			outputMsg.AddU32(helpers.Ip2int(config.ConfigInstance.GameIp)) // server ip
			outputMsg.AddU16(uint16(config.ConfigInstance.GamePort))       // server port

			if version >= 980 {
				outputMsg.AddU8(0) // world preview
			}
		}
	}

	// premium
	if version > 1077 {
		// account status: 0 - ok, 1 - frozen, 2 - suspended
		outputMsg.AddU8(0)

		// premium status
		if account.Premdays > 0 {
			outputMsg.AddU8(1) // premium
		} else {
			outputMsg.AddU8(0) // free
		}

		outputMsg.AddU32(uint32(account.Premdays)) // timestamp
	} else {
		outputMsg.AddU16(uint16(account.Premdays)) // days
	}

	p.SendPacket(outputMsg)
}
