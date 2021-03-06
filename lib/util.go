package lib

import (
	"Cryptopals/lib/enc/base64"
	"io/ioutil"
)

// ReadBase64FromFile reads bytes from file in base 64
func ReadBase64FromFile(filename string) ([]byte, error) {
	dataB64, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	bytes, err := base64.DecodeString(string(dataB64))
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
