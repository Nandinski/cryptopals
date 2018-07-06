package main

import (
	"Cryptopals/lib"
	// "Cryptopals/lib/enc/base64"
	"fmt"
)

/*
File - It's been base64'd after being encrypted with repeating-key XOR (unkown size).

 Decrypt it.

Here's how:

1. Let KEYSIZE be the guessed length of the key; try values from 2 to (say) 40.
2. Write a function to compute the edit distance/Hamming distance between two strings. The Hamming distance is just the number of differing bits. The distance between:
		'this is a test' and 'wokka wokka!!!'
		is 37. Make sure your code agrees before you proceed.
3. For each KEYSIZE, take the first KEYSIZE worth of bytes, and the second KEYSIZE worth of bytes, and find the edit distance between them.
	  Normalize this result by dividing by KEYSIZE.
*/
func main() {
	cypheredText, err := lib.ReadBase64FromFile("6.txt")
	fmt.Println()
	fmt.Println(err)

	fmt.Println(len(cypheredText))

	// Try to find the key length
	// for i := 2; i < 40; i++ {
	//
	// }

}
