package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "time"
    "os"
    "syscall"
    "reflect"
)

func TestBuildServer(t *testing.T) {
    handler := http.HandlerFunc(MainHandler)
    ts := httptest.NewServer(handler)
    server := BuildServer(ts.URL)
    defer server.Close()
    if server.Addr != ts.URL {
        t.Errorf("Expected new server to have URL %s got %s", ts.URL, server.Addr)
    }
}

func TestBuildChannels(t *testing.T) {
    shutChan, intChan, doneChan := BuildChannels()
    
    if reflect.TypeOf(shutChan).String() != "chan bool" {
        t.Errorf("Expected shutChan to be a 'chan bool' got %s", reflect.TypeOf(shutChan))
    }

    if reflect.TypeOf(intChan).String() != "chan os.Signal" {
        t.Errorf("Expected intChan to be a 'chan bool' got %s", reflect.TypeOf(intChan))
    }

    if reflect.TypeOf(doneChan).String() != "chan bool" {
        t.Errorf("Expected doneChan to be a 'chan bool' got %s", reflect.TypeOf(doneChan))
    }
}

func TestStopServer(t *testing.T) {
    handler := http.HandlerFunc(MainHandler)

    // Cannot use an instance of httptest.NewServer() in place of *http.Server
    // So create both and use the addr from the test instance
    ts := httptest.NewServer(handler)
    server := &http.Server{Addr: ts.URL}
    defer server.Close()

    shutdownCalled := false
    server.RegisterOnShutdown(func() {shutdownCalled = true})
    server.ListenAndServe()

    time.Sleep(time.Second) // give the server time to start
    StopServer(server)
    time.Sleep(6 * time.Second) // ensure server is shutdown
    if ! shutdownCalled {
        t.Error("Shutdown not called")
    }

}

func TestSelectShutdownChannel(t *testing.T) {

    handler := http.HandlerFunc(MainHandler)
    // Again create an instance of both for testing
    ts := httptest.NewServer(handler)
    server := &http.Server{Addr: ts.URL}

    shutdownCalled := false
    server.RegisterOnShutdown(func() { shutdownCalled = true })
    defer server.Close()
    server.ListenAndServe()

    intTestChan, downTestChan, doneTestChan := CreateTestChannels()
    go SelectChannel(server, intTestChan, downTestChan, doneTestChan)

    // Send to downTestChan
    downTestChan <- true

    // wait for done
    <-doneTestChan

    // wait up to 6 seconds for server to exit
    for i := 0; i < 6; i++ {
        if shutdownCalled {
            continue
        } else {
            time.Sleep(time.Second)
        }
    }

    if ! shutdownCalled {
        t.Error("Shutdown not called on downTestChan")
    }
}

func TestSelectInterruptChannel(t *testing.T) {
    handler := http.HandlerFunc(MainHandler)
    // Again create an instance of both for testing
    ts := httptest.NewServer(handler)
    server := &http.Server{Addr: ts.Config.Addr}

    shutdownCalled := false
    server.RegisterOnShutdown(func() { shutdownCalled = true })
    defer server.Close()
    server.ListenAndServe()

    intTestChan, downTestChan, doneTestChan := CreateTestChannels()
    go SelectChannel(server, intTestChan, downTestChan, doneTestChan)

    intTestChan <- syscall.SIGINT

    // wait for done
    <- doneTestChan

    // wait up to 6 seconds for server to exit
    for i := 0; i < 6; i++ {
        if shutdownCalled {
            continue
        } else {
            time.Sleep(time.Second)
        }
    }

    if ! shutdownCalled {
        t.Error("Shutdown not called on intTestChan")
    }

}

func CreateTestChannels() (
    intTestChan chan os.Signal,
    downTestChan chan bool,
    doneTestChan chan bool) {
    // create handlers
    intTestChan = make(chan os.Signal, 1)
    downTestChan = make(chan bool)
    doneTestChan = make(chan bool)

    return intTestChan, downTestChan, doneTestChan
}

func MainHandler(res http.ResponseWriter, req *http.Request) {
    return
}