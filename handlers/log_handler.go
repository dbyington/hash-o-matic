package handlers

import (
    "net/http"
    "log"
)

func LogHandler(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        req.ParseForm()
        log.Printf("%s %s %s %s\n",
            req.RemoteAddr, req.Method, req.URL, req.UserAgent())
        handler.ServeHTTP(res, req)
    })
}
