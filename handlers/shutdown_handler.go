package handlers

import (
    "net/http"
)

func BuildShutdownHandler(shutdownChan chan bool) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        if req.Method == http.MethodPut {
            // shutdown issued here
            res.WriteHeader(http.StatusAccepted)
            res.Write([]byte("shutting down..."))
            shutdownChan <- true
        } else {
            // no more supported methods to match return 404, write header first
            res.WriteHeader(http.StatusNotFound)
            res.Write([]byte("404 page not found\n"))
        }
    }

    }
