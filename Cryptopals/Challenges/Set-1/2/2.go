package main

import (
	"Cryptopals/lib"
	"fmt"
)


// Fixed XOR
//
// Write a function that takes two equal-length buffers and produces their XOR combination.
//
// If your function works properly, then when you feed it the string:
//
// 1c0111001f010100061a024b53535009181c
//
// ... after hex decoding, and when XOR'd against:
//
// 686974207468652062756c6c277320657965
//
// ... should produce:
//
// 746865206b696420646f6e277420706c6179


func main() {
	s1 := "1c0111001f010100061a024b53535009181c"

	s2 := "686974207468652062756c6c277320657965"

	xored := lib.FixedXOR(s1, s2)
	fmt.Print(xored)

	return
}
