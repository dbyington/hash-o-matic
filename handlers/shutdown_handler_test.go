package handlers

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "strings"
)

func TestBuildShutdownHandler(t *testing.T) {
    shutdown := make(chan bool)
    ShutdownHandler := BuildShutdownHandler(shutdown)
    called := make(chan bool)
    server := httptest.NewServer(http.HandlerFunc(ShutdownHandler))

    go func() {
        <-shutdown
        server.Close()
        called <- true
    }()

    t.Logf("Should return %d on GET", http.StatusNotFound)
    res, err := http.Get(server.URL)
    if err != nil {
        t.Logf("Error: %v", err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("%d not thrown", http.StatusNotFound)
    }

    t.Log("Should start graceful server shutdown on PUT")
    req, err := http.NewRequest(http.MethodPut, server.URL, strings.NewReader(""))
    if err != nil {
        t.Errorf("Failed PUT: %s", err.Error())
    }

    t.Log("Should respond to client request")
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Errorf("Failed: %s", err.Error())
    }
    t.Logf("Should return %d on success", http.StatusAccepted)
    if res.StatusCode != http.StatusAccepted {
        t.Errorf("PUT to /shutdown failed; expected %d got %d", http.StatusAccepted, res.StatusCode)
    }


    t.Log("shutdown channel should be set true, calling shutdown")
    shutdownCalled := <-called
    if !shutdownCalled {
        t.Error("Shutdown was not called")
    }

    t.Log("Server should no longer exist")
    if server.Config.Addr != "" {
        t.Error("Server not closed")
    }

}