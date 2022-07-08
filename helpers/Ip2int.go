package helpers

import (
	"encoding/binary"
	"net"
)

// Convert IP string to uint32, e.g.: 127.0.0.1 to 2130706433
func Ip2int(ip string) uint32 {
	var ipParsed net.IP = net.ParseIP(ip)

	if len(ipParsed) == 16 {
		return binary.BigEndian.Uint32(ipParsed[12:16])
	}

	return binary.BigEndian.Uint32(ipParsed)
}
