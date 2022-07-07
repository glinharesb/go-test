package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func Ip2int(ip string) uint32 {
	var ipParsed net.IP = net.ParseIP(ip)

	if len(ipParsed) == 16 {
		return binary.BigEndian.Uint32(ipParsed[12:16])
	}

	return binary.BigEndian.Uint32(ipParsed)
}

func FormatVersion(version uint16) string {
	toString := strconv.FormatUint(uint64(version), 10)
	length := len(toString) - 2
	return fmt.Sprintf("%s.%s", toString[:length], toString[length:])
}
