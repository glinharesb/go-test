package helpers

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// Format client version, e.g.: 1098 to 10.98
func FormatVersion(version uint16) string {
	toString := strconv.FormatUint(uint64(version), 10)
	length := len(toString) - 2

	return fmt.Sprintf("%s.%s", toString[:length], toString[length:])
}

// Convert IP string to uint32, e.g.: 127.0.0.1 to 2130706433
func Ip2int(ip string) uint32 {
	var ipParsed net.IP = net.ParseIP(ip)

	if len(ipParsed) == 16 {
		return binary.BigEndian.Uint32(ipParsed[12:16])
	}

	return binary.BigEndian.Uint32(ipParsed)
}

func TransformToSha1(password *string) {
	hash := sha1.Sum([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}

func TransformToSha256(password *string) {
	hash := sha256.Sum256([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}

func TransformToSha512(password *string) {
	hash := sha512.Sum512([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}