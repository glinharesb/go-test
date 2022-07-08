package helpers

import (
	"crypto/sha1"
	"fmt"
)

func TransformToSha1(password *string) {
	hash := sha1.Sum([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}
