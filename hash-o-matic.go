package main

import (
    "github.com/dbyington/hash-o-matic/handlers"
    "log"
    "net/http"
)

const LISTEN_ADDR = ":8080"
const HASH_URL = "/hash"

func main() {
    http.HandleFunc(HASH_URL, handlers.HashHandler)
    log.Fatal(http.ListenAndServe(LISTEN_ADDR,
        handlers.LogHandler(http.DefaultServeMux)))
}


