package helpers

import (
	"crypto/sha512"
	"fmt"
)

func TransformToSha512(password *string) {
	hash := sha512.Sum512([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}
