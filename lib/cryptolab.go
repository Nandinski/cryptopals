// Cryptopals
package lib

import (
	"Cryptopals/lib/enc/hex"
	"Cryptopals/lib/enc/base64"
	"Cryptopals/lib/enc/xor"
	"fmt"
	"errors"
	"strings"
	"bufio"
	"os"
)

// *************************
//
// Set 1
//
// *************************

//
// Convert hex to base64
func ConvertHexTo64(hexString string) string {
	fmt.Printf("\nHex\n%s", hexString)
	fmt.Println()
	// Convert hex to string
	decodedHexString, err := hex.DecodeString(hexString)
	checkError(err)
	fmt.Printf("\nDecodedHex\n%s", decodedHexString)


	encodedbase64String := base64.EncodeToString(decodedHexString)
	fmt.Printf("\nBase64\n%s", encodedbase64String)

	return encodedbase64String
}

//
// Fixed XOR

func FixedXOR(s1, s2 string) string {
	fmt.Printf("\nHex1\n%s\n", s1)
	fmt.Print("XOR")
	fmt.Printf("\nHex2\n%s\n", s2)
	// Convert hex to string
	decodedHexString1, err := hex.DecodeString(s1)
	checkError(err)
	decodedHexString2, err := hex.DecodeString(s2)
	checkError(err)

	xoredValue, err := XOROverByteVector(decodedHexString1, decodedHexString2)
	checkError(err)

	hexEncodedXOR := hex.EncodeToString(xoredValue)

	return hexEncodedXOR
}

var TryingToXORDifSizedValues = errors.New("Can't perform XOR on different sized arrays")

func XOROverByteVector(b1, b2 []byte) ([]byte, error){
	if len(b1) != len(b2) {
		return b1, TryingToXORDifSizedValues
	}
	xored := make([]byte, len(b1))

	for i := range(b1) {
		xored[i] = b1[i]^b2[i]
	}

	return xored, nil
}

//
// Single-byte XOR cipher

func SingleKeyXORDecipher(hexString string) (rune, string) {
	// Convert hex to bytes
	decodedHexString, err := hex.DecodeString(hexString)
	checkError(err)
	// fmt.Printf("\nDecodedHex: %s\n", decodedHexString)

	key, plaintext := BreakSingleByteXOR(decodedHexString)
	return key, plaintext
}

func BreakSingleByteXOR(b []byte) (rune, string) {
	var key rune
	var plainTextMaxScore string
	maxScore := 0

	// Byte has 256 possibilities
	for i := 0; i < 256; i++ {
		plainText := xorByByte(byte(i), b)
		score := computeScore(plainText)

		// Equals prevents case where no option is valid (lowest score is 0)
		if score >= maxScore {
			maxScore = score
			plainTextMaxScore = plainText
			key = rune(i)
			// fmt.Println(rune(i), score)
		}


	}
	return key, plainTextMaxScore
}

func xorByByte(key byte, b []byte) string {
	decipheredText := make([]byte, len(b))

	for i := range b {
		decipheredText[i] = key ^ b[i]
	}

	return string(decipheredText)
}

var EnglishFrequencyTable = map[rune]float32 {
	' ' : 0.1918182,	// This was the key to break it!
	'a' : 0.0651738,
	'b' : 0.0124248,
	'c' : 0.0217339,
	'd' : 0.0349835,
	'e' : 0.1041442,
	'f' : 0.0197881,
	'g' : 0.0158610,
	'h' : 0.0492888,
	'i' : 0.0558094,
	'j' : 0.0009033,
	'k' : 0.0050529,
	'l' : 0.0331490,
	'm' : 0.0202124,
	'n' : 0.0564513,
	'o' : 0.0596302,
	'p' : 0.0137645,
	'q' : 0.0008606,
	'r' : 0.0497563,
	's' : 0.0515760,
	't' : 0.0729357,
	'u' : 0.0225134,
	'v' : 0.0082903,
	'w' : 0.0171272,
	'x' : 0.0013692,
	'y' : 0.0145984,
	'z' : 0.0007836,
}

func computeScore(text string) int {
	score := 0
	for _, c := range(strings.ToLower(text)) {
		if c < 32 || c > 127 {
			return score
		}

		if val, ok := EnglishFrequencyTable[c]; ok{
			score += int(val * 100)
		}
	}
	return score
}


//
// Detect single-character XOR

func DetectSingleCharacterXORInFile(filename string) (rune, string) {
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var key rune
	var plainTextMaxScore string
	maxScore := 0
	for scanner.Scan() {
		//scanner.Text() has line content
		keyLine, plaintext := SingleKeyXORDecipher(scanner.Text())
		score := computeScore(plaintext)

		if score >= maxScore {
			maxScore = score
			plainTextMaxScore = plaintext
			key = keyLine
		}
	}

	if err := scanner.Err(); err != nil {
		panic("Error reading file")
	}

	return key, plainTextMaxScore
}

//
// Implement repeating-key XOR
func EncodeWithRepeatingXOR(key string, plainText string) (string, error) {
	cypheredText, err := xor.EncodeWithRepeatingXOR(key, plainText)
	checkError(err)

	encoded := hex.EncodeToString(cypheredText)
	return encoded, nil
}

//
// Aux

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
