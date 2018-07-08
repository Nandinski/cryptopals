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
	"sort"
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

func SingleKeyXORDecipher(hexString string) (byte, string) {
	// Convert hex to bytes
	decodedHexString, err := hex.DecodeString(hexString)
	checkError(err)
	// fmt.Printf("\nDecodedHex: %s\n", decodedHexString)

	plaintext, key, _ := BreakSingleByteXOR(decodedHexString)
	return key, plaintext
}

func BreakSingleByteXOR(b []byte) (string, byte, float64) {
	var key byte
	var plainTextMaxScore string
	var maxScore float64

	// Byte has 256 possibilities
	for i := 0; i < 256; i++ {
		plainText := xorByByte(byte(i), b)
		score := computeScore(plainText)

		// Equals prevents case where no option is valid (lowest score is 0)
		if score >= maxScore {
			maxScore = score
			plainTextMaxScore = plainText
			key = byte(i)
			// fmt.Println(rune(i), score)
		}
	}
	return plainTextMaxScore, key, maxScore
}

func xorByByte(key byte, b []byte) string {
	decipheredText := make([]byte, len(b))

	for i := range b {
		decipheredText[i] = key ^ b[i]
	}

	return string(decipheredText)
}

var EnglishFrequencyTable = map[rune]float64 {
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


func computeScore(text string) float64 {
	var score float64
	for _, char := range(strings.ToLower(text)) {
		if char > 126 {
			return 0
		}
		if val, ok := EnglishFrequencyTable[char]; ok{
			score += val
		}
	}
	return score
}


//
// Detect single-character XOR

func DetectSingleCharacterXORInFile(filename string) (byte, string) {
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var key byte
	var plainTextMaxScore string
	var maxScore float64
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
// Break repeating-key XOR

// Calculate the ammount of diferent bits
func HammingDistance(a, b []byte) int {
	if len(a) != len(b) {
		panic("HammingDistance: different Lengths")
	}

	diff := 0

	for i := 0; i < len(a); i++ {
		ba := a[i]
		bb := b[i]

		// Go through every bit of a byte
		var mask byte
		for j := 0; j < 8; j++ {
			mask = byte(1 << uint(j))
			if (ba & mask) != (bb & mask) {
				diff++
			}
		}
	}

	return diff
}

type scoreKey struct {
	score float64
	key int
}

func findRepeatingXORSize(cypheredText []byte) []int {
	var res []scoreKey
	for keySize := 2; keySize < 40; keySize++ {
		// for each 'keySize' size of text do
		averageHD :=  0.

		blockCount := len(cypheredText) / keySize
		for lBlockIndex := 0; lBlockIndex < blockCount - 1; lBlockIndex++ {
			normalizedHD := 0.
			leftBlock := cypheredText[lBlockIndex * keySize : (lBlockIndex + 1) * keySize]

			for rBlockIndex := lBlockIndex + 1; rBlockIndex < blockCount; rBlockIndex++ {
				rightBlock := cypheredText[rBlockIndex * keySize : (rBlockIndex + 1) * keySize]

				normalizedHD += float64(HammingDistance(leftBlock, rightBlock)) / float64(keySize)
			}
			averageHD = normalizedHD / float64(blockCount - lBlockIndex)
		}
		// a, b := cypheredText[:keySize*4], cypheredText[keySize*4 : keySize*4 * 2]
		// averageHD = float64(HammingDistance(a, b)) / float64(keySize)

		res = append(res, scoreKey{averageHD, keySize})
	}

	// Sort scores by lowest
	sort.Slice(res, func(i, j int) bool { return res[i].score < res[j].score })
	bestKeys := getKeyFromBestScores(res)
	return bestKeys[:3]
}

func BreakRepeatingXORCypher(cypheredText []byte) ([]byte, []byte){
	// Try to find the key length
	likelyKeySize := findRepeatingXORSize(cypheredText)

	var bestScore float64
	var bestKey []byte
	for _, keySize  := range likelyKeySize {
		column := make([]byte, len(cypheredText)/keySize)
		cypherKey := make([]byte, keySize)

		keyScore := 0.
		for col := 0; col < keySize; col++ {
			for row := 0; row < len(column) - 1; row++ {
				if row * keySize + col >= len(cypheredText) { continue } // Ignore it if you can't read it
				column[row] = cypheredText[row*keySize+col]
			}
			_, key, score := BreakSingleByteXOR(column)

			cypherKey[col] = key
			keyScore += score
		}

		if keyScore > bestScore {
			bestScore = keyScore
			bestKey = cypherKey
		}
	}

	decipheredText, _ := xor.EncodeWithRepeatingXOR(string(bestKey), string(cypheredText))

	return decipheredText, bestKey
}

//
// Aux

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getKeyFromBestScores(b []scoreKey) []int {
	var keys []int
	for _, scoreKeyEl := range b {
		keys = append(keys, scoreKeyEl.key)
	}
	return keys
}
