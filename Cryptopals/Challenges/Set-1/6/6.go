package main

import (
	// "Cryptopals/lib"
	"Cryptopals/lib/enc/base64"
	"fmt"
)

// File - It's been base64'd after being encrypted with repeating-key XOR (unkown size).
//
//  Decrypt it.
//
// Here's how:
//
// 1. Let KEYSIZE be the guessed length of the key; try values from 2 to (say) 40.
// 2. Write a function to compute the edit distance/Hamming distance between two strings. The Hamming distance is just the number of differing bits. The distance between:
//
// 		'this is a test'
//
// 		and
//
// 		'wokka wokka!!!'
//
// 		is 37. Make sure your code agrees before you proceed.


func main() {

	plainText := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"

	fmt.Println("Original text")
	fmt.Printf("%s\n", plainText)

	bytes := []byte(plainText)
	encodedbase64String := base64.EncodeToString(bytes)
	fmt.Println("Base64")
	fmt.Printf("%s\n", encodedbase64String)

	decodedbase64String, _ := base64.DecodeString(encodedbase64String)
	fmt.Println("String")
	fmt.Printf("%s\n", decodedbase64String)
}
