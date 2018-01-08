package handlers

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "os"
    "bytes"
    "log"
    "strings"
)

func TestLogHandler(t *testing.T) {

    var buf bytes.Buffer
    log.SetOutput(&buf)
    // cleanup when we exit
    defer func() {
        log.SetOutput(os.Stderr)
    }()

    handler := LogHandler(http.DefaultServeMux)
    server := httptest.NewServer(http.Handler(handler))
    defer server.Close()

    expectedString := "GET / Go-http-client/1.1"
    _, err := http.Get(server.URL)
    if err != nil {
        t.Error("Error testing GET for LogHandler:", err)
    }
    out := buf.String()
    if ! strings.Contains(out, expectedString) {
        t.Errorf("Expected '%s' to contain '%s'", out, expectedString)
    }


}