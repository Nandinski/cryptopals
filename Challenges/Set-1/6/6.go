package main

import (
	"Cryptopals/lib"
	// "Cryptopals/lib/enc/base64"
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
// 		'this is a test' and 'wokka wokka!!!'
// 		is 37. Make sure your code agrees before you proceed.


func main() {

	s1 := "this is a test"
	s2 := "wokka wokka!!!"

	distance, _ := lib.HammingDistance([]byte(s1), []byte(s2))

	fmt.Println("Distance: ", distance)
}
