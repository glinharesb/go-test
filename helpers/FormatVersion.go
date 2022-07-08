package helpers

import (
	"fmt"
	"strconv"
)

// Format client version, e.g.: 1098 to 10.98
func FormatVersion(version uint16) string {
	toString := strconv.FormatUint(uint64(version), 10)
	length := len(toString) - 2

	return fmt.Sprintf("%s.%s", toString[:length], toString[length:])
}
