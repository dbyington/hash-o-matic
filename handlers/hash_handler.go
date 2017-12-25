package handlers

import (
    "net/http"
    "fmt"
    "time"
    "github.com/dbyington/hash-o-matic/util"
    "strings"
)

const FORM_KEY = "password"

func HashHandler(res http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    if req.Method == http.MethodPost {
        HashPostMethod(res, req)
        return
    }

    // no more supported methods to match return 404, write header first
    res.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(res, "404 page not found\n")
}

func HashPostMethod(res http.ResponseWriter, req *http.Request) {
    password, exists := req.PostForm[FORM_KEY]
    if exists {
        // requisite sleep, simulate slow connection
        time.Sleep(5 * time.Second)
        res.WriteHeader(http.StatusCreated)
        fmt.Fprintf(res, "%s", util.HashString(strings.Join(password, "")))
    } else {
        res.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(res, "Missing field '%s'\n", FORM_KEY)
    }
}
