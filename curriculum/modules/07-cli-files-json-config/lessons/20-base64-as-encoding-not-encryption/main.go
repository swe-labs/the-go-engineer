package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func EncodeDecode(original []byte) (encoded string, decoded []byte, err error) {
	encoded = base64.URLEncoding.EncodeToString(original)
	decoded, err = base64.URLEncoding.DecodeString(encoded)
	return
}

func main() {
	// Text data
	text := []byte("Base64 is encoding, not encryption!")
	enc, dec, err := EncodeDecode(text)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Original: %s\n", string(text))
	fmt.Printf("Encoded:  %s\n", enc)
	fmt.Printf("Decoded:  %s\n", string(dec))

	// Binary data
	binary := []byte{0x00, 0xFF, 0x80, 0x40, 0x7F, 0x01}
	enc2, dec2, err := EncodeDecode(binary)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\nBinary original: %v\n", binary)
	fmt.Printf("Binary encoded:  %s\n", enc2)
	fmt.Printf("Binary decoded:  %v\n", dec2)

	// Verify URL safety
	if !strings.ContainsAny(enc, "+/") {
		fmt.Println("\nURLEncoding is URL-safe (no + or / characters)")
	}
	if strings.HasSuffix(enc, "=") {
		fmt.Printf("Has padding: yes (%d pad chars)\n", strings.Count(enc, "="))
	}
}
