
package lib

import (
	"io/ioutil"
	"Cryptopals/lib/enc/base64"
)

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
