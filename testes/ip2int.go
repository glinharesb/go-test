package main

import (
	"fmt"
	"strconv"
)

func main() {
	var version uint16 = 1098
	fmt.Println(FormatVersion(version))

}

func FormatVersion(version uint16) string {
	toString := strconv.FormatUint(uint64(version), 10)
	length := len(toString) - 2
	return fmt.Sprintf("%s.%s", toString[:length], toString[length:])
}
