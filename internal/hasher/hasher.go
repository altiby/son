package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Hasher struct {
	sha hash.Hash
}

func New() *Hasher {
	return &Hasher{sha256.New()}
}

func (h Hasher) Hash(in string) (string, error) {
	_, err := h.sha.Write([]byte(in))
	defer h.sha.Reset()
	if err != nil {
		return "", err
	}

	value := h.sha.Sum(nil)
	return hex.EncodeToString(value[:]), nil
}
