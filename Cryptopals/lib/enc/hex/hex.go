package hex

import (
	"errors"
	"fmt"
)


const hextable = "0123456789abcdef"

var InvalidStringSize = errors.New("enc/hex: odd length hex string")

// InvalidByteError values describe errors resulting from an invalid byte in a hex string.
type InvalidCharacter byte

func (e InvalidCharacter) Error() string {
	return fmt.Sprintf("enc/hex: invalid byte value: %#U", rune(e))
}

func Decode(dst, src []byte) (int, error) {
	if len(src) % 2 != 0 {
		return 0, InvalidStringSize
	}

	// Each cycle gets a char
	var i int
	for i = 0; i < len(src)/2; i++ {
		f, ok := convertHexToBin(src[2*i])
		if !ok {
			return i, InvalidCharacter(src[2*i])
		}
		s, ok := convertHexToBin(src[2*i + 1])
		if !ok {
			return i, InvalidCharacter(src[2*i + 1])
		}

		//Append first to second 4 bit
		dst[i] = (f << 4) | s
	}

	return i, nil
}

func convertHexToBin(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	}

	return 0, false
}

// Input needs to be even (2 elements = 1 byte)
// s should contain only hexadecimal characters
func DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, DecodedLen(len(s)))

	n, err := Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

func DecodedLen(n int) int {
	return n / 2
}

func EncodedLen(n int) int {
	return n * 2
}

func Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	for i := range src {
		val := src[i]

		dst[2*i] = hextable[val>>4&0xF]
		dst[2*i + 1] = hextable[val&0xF]
	}
}

func EncodeToString(src []byte) string{
	buf := make([]byte, EncodedLen(len(src)))
	Encode(buf, src)
	return string(buf)
}
