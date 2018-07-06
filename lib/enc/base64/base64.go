package base64

import(
	"fmt"
	"errors"
)

const (
	encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	StdPadding = '='
)

func EncodeToString(src []byte) string {
	buf := make([]byte, EncodedLen(src))
	Encode(buf, src)
	return string(buf)
}

func EncodedLen(src []byte) int {
	// Fill up to closest multiple of 12
	// First make it divisible by 3 -> each 3 in byte = 4 in b64
	return ((len(src) + 2) / 3) * 4
}

func Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	i, j := 0, 0
	// Work for the simple case only
	simpleC := (len(src) / 3 * 3)

	for i < simpleC {
		val := uint(src[i+0])<<16 | uint(src[i+1])<<8 | uint(src[i+2])

		dst[j+0] = encodeStd[val>>18&0x3F]
		dst[j+1] = encodeStd[val>>12&0x3F]
		dst[j+2] = encodeStd[val>>6&0x3F]
		dst[j+3] = encodeStd[val&0x3F]

		i += 3
		j += 4
	}
	remain := len(src) - i

	// We are done
	if remain == 0 {
		return
	}

	// There's at least 1 remaining
	val := uint(src[i]) << 16
	if remain == 2 {
		val |= uint(src[i + 1]) << 8
	}

	dst[j] 		= encodeStd[val>>18&0x3F]
	dst[j+1]	= encodeStd[val>>12&0x3F]

	switch remain {
	case 1:
		dst[j+2]	= StdPadding
		dst[j+3]	= StdPadding
	case 2:
		dst[j+2]	= encodeStd[val>>6&0x3F]
		dst[j+3]	= StdPadding
	}
}


func DecodedLen(n int) int {
	return (n / 4) * 3
}

var InvalidStringSize = errors.New("Invalid size of string for base64 string")

type InvalidCharacter byte

func (i InvalidCharacter) Error() string{
	// Carefull with infinte loops if you print i without processing it will call this again - you dufus
	return fmt.Sprintf("Invalid character - not of base64: %v", byte(i))
}

// DecodeString returns the bytes represented by the base64 string s.
func DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, DecodedLen(len(s)))
	n, err := Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

func Decode(dst, src []byte) (n int, err error) {
	if len(src) == 0 {
		return 0, nil
	}

	// Each cycle gets 3 char
	var i int
	srcI := 0

	// The last 4 bytes might have padding - treat them differently
	for i = 0; i < len(src)/4; i++ {
		if ok := convertB64Block(dst[n:], src[srcI:]); ok {
			srcI += 4
			n += 3
		} else {
			var dstInc, srcInc int
			dstInc, srcInc, err = checkForQuantum(dst[n:], src[srcI:])
			if err != nil {
				fmt.Println(err)
				fmt.Println(dst)
				return n, InvalidCharacter(src[srcI])
			}

			srcI += srcInc
			n += dstInc
			i = srcI / 4
		}
	}

	return n, nil
}

// Character was not in alfabet - might be '\n' '\r' or padding
func checkForQuantum(dst, src []byte) (n, srcI int, err error){
	// Convert 4 char in 64 to 3 bytes
	var dstBytes [4]byte
	dlen := 3
	for i := 0; i < len(dstBytes); i++ {
		// needed more char to be able to decode
		if(srcI == len(src)){
			return 0, 0, errors.New("Need more info to be able to decode")
		}

		c := src[srcI]
		srcI++

		out, ok := convertB64ToBin(c)
		if ok {
			dstBytes[i] = out
			continue
		}

		// Ignore '\n' and '\r'
		if c == '\n' || c == '\r' {
			i--
			continue
		}

		if c != StdPadding {
			// Blow up - This value can't be parsed
			return 0, 0, InvalidCharacter(c)
		}

		// We've reached the end and there's padding
		switch i {
		case 0, 1:
			// incorrect padding
			// padding can only appear in the last or 2 last bytes
			return 0, 0, errors.New(fmt.Sprintf("Error - Padding malformed %#U", c))
		case 2:
			// Expecting '=='
			// First '=' was read in c
			for ; srcI < len(src); srcI++{
				// Ignore '\n' and '\r'
				if src[srcI] == '\n' || src[srcI] == '\r' {
					continue
				}
				break
			}

			// Did it stop because it can't read more?
			if srcI == len(src) {
				return 0, 0, InvalidCharacter(src[srcI])
			}

			// We're looking for a '=', did we find it?
			if src[srcI] != StdPadding {
				return 0, 0, InvalidCharacter(src[srcI])
			}

			// I'll just ignore everything else
		}

		// Stop looking if we found the first
		dlen = i
		break
	}

	val := uint(dstBytes[0])<<18 | uint(dstBytes[1])<<12 | uint(dstBytes[2])<<6 | uint(dstBytes[3])
	dstBytes[0] = byte(val>>16&0xFF)
	dstBytes[1] = byte(val>>8&0xFF)
	dstBytes[2] = byte(val&0xFF)

	switch dlen {
	case 4:
		dst[2] = dstBytes[2]
		fallthrough
	case 3:
		dst[1] = dstBytes[1]
		fallthrough
	case 2:
		dst[0] = dstBytes[0]
	}

	return dlen, srcI, nil
}

func convertB64Block(dst, src []byte) bool {
	var dval uint32
	var n byte
	var ok bool

	if n, ok = convertB64ToBin(src[0]); !ok {
		return false
	}
	dval |= uint32(n) << 18
	if n, ok = convertB64ToBin(src[1]); !ok {
		return false
	}
	dval |= uint32(n) << 12
	if n, ok = convertB64ToBin(src[2]); !ok {
		return false
	}
	dval |= uint32(n) << 6
	if n, ok = convertB64ToBin(src[3]); !ok {
		return false
	}
	dval |= uint32(n)

	dst[0] = byte(dval>>16&0xFF)
	dst[1] = byte(dval>>8&0xFF)
	dst[2] = byte(dval&0xFF)

	return true
}

func convertB64ToBin(c byte) (byte, bool) {
	switch {
	case 'A' <= c && c <= 'Z':
		return c - 'A', true
	case 'a' <= c && c <= 'z':
		return c - 'a' + 26, true
	case '0' <= c && c <= '9':
		return c - '0' + 52, true
	case c == '+':
		return 62, true
	case c == '/':
		return 63, true
	}

	return 0, false
}
