package main

import (
	"crypto/sha512"
	"encoding/base64"
)

func HashString(stringToHash string) string {

	// Sum512() argument must be type [size]byte
	byteArray := []byte(stringToHash)
	sum512byteArray := sha512.Sum512(byteArray)

	// EncodeToString() argument must be type []byte so first turn the byte array into a string
	sum512string := string(sum512byteArray[:])
	base64string := base64.StdEncoding.EncodeToString([]byte(sum512string))

	return base64string
}
