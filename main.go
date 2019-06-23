package main

import (
    "crypto/sha512"
    "encoding/base64"
)

func main () {

}

func hashString (s string) string {
    sum512byteArray := sha512.Sum512([]byte(s))
    return base64.StdEncoding.EncodeToString([]byte(string(sum512byteArray[:])))
}
