package hasher

import (
	"crypto/sha256"
	"encoding/base64"
)

type Hasher interface {
	Hash(str string) string
	Compare(str, hash string) bool
}

type SHA256Hasher struct {
	Salt string
}

func (h *SHA256Hasher) Hash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str + h.Salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (h *SHA256Hasher) Compare(str, hash string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(str + h.Salt))

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)) == hash
}

func NewHasher(salt string) Hasher {
	return &SHA256Hasher{salt}
}
