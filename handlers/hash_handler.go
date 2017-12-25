package handlers

import (
    "net/http"
    "time"
    "github.com/dbyington/hash-o-matic/util"
    "strings"
)

const FORM_KEY = "password"

func HashHandler(res http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodPost {
        HashPostMethod(res, req)
        return
    }

    // no more supported methods to match return 404, write header first
    res.WriteHeader(http.StatusNotFound)
    res.Write([]byte("404 page not found\n"))
}

func HashPostMethod(res http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    password, exists := req.PostForm[FORM_KEY]
    if exists {
        // requisite sleep, simulate slow connection
        time.Sleep(5 * time.Second)
        res.WriteHeader(http.StatusCreated)
        res.Write([]byte(util.HashString(strings.Join(password, ""))))
    } else {
        res.WriteHeader(http.StatusBadRequest)
        res.Write([]byte("Missing key " + FORM_KEY))
    }
}
