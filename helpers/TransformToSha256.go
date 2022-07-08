package helpers

import (
	"crypto/sha256"
	"fmt"
)

func TransformToSha256(password *string) {
	hash := sha256.Sum256([]byte(*password))
	*password = fmt.Sprintf("%x", hash)
}
