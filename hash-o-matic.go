package main

import (
    "os"
    "fmt"
    "github.com/dbyington/hash-o-matic/util"
)

func main() {
    if (len(os.Args) > 1) {
        fmt.Printf("Your string as a base64 encoded sha512 hash is:\n'%s'\n", util.HashString(os.Args[1]))
    }
}
