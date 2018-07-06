package xor

import (
	"errors"
)

var InvalidEmptyKey = errors.New("enc/xor: key can't be empty")


func EncodeWithRepeatingXOR(key, plainText string) ([]byte, error) {
	out := make([]byte, len(plainText))
	src := []byte(plainText)

	if len(key) == 0 {
		return out, InvalidEmptyKey
	}

	ki := 0
	for i := range(src) {
		out[i] = src[i] ^ key[ki]
		ki = (ki + 1) % len(key)
	}

	return out, nil
}
