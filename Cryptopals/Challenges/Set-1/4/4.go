package main

import (
	"Cryptopals/lib"
	"fmt"
)


// Detect single-character XOR
//
// One of the 60-character strings in the file 4.txt has been encrypted by single-character XOR.
//
// Find it.
//
// (Your code from #3 should help.)


func main() {
	fileName := "4.txt"

	key, decodedString := lib.DetectSingleCharacterXORInFile(fileName)

	fmt.Println("Key: ", string(key))
	fmt.Println("Decoded String: ", decodedString)
}
