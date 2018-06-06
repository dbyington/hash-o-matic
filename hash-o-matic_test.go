package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestStopServer(t *testing.T) {
	t.Log("StopServer")
	handler := http.HandlerFunc(MainHandler)

	// Cannot use an instance of httptest.NewServer() in place of *http.Server
	// So create both and use the addr from the test instance
	ts := httptest.NewServer(handler)
	server := &http.Server{Addr: ts.URL}
	defer server.Close()

	shutdownCalled := false
	wg.Add(1)
	server.RegisterOnShutdown(func() { shutdownCalled = true })
	server.ListenAndServe()

	time.Sleep(time.Second) // give the server time to start
	StopServer(server)
	time.Sleep(6 * time.Second) // ensure server is shutdown
	if !shutdownCalled {
		t.Error("Shutdown not called")
	}

}

func TestSelectShutdownChannel(t *testing.T) {
	t.Log("SelectShutdownChannel")
	handler := http.HandlerFunc(MainHandler)
	// Again create an instance of both for testing
	ts := httptest.NewServer(handler)
	server := &http.Server{Addr: ts.URL}

	shutdownCalled := false
	wg.Add(1)
	server.RegisterOnShutdown(func() { shutdownCalled = true })
	defer server.Close()
	server.ListenAndServe()

	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel and associated handler for PUTs to /shutdown
	shutdownChan := make(chan bool)
	go SelectChannel(server, interruptChan, shutdownChan)

	// Send to shutdownChan
	shutdownChan <- true

	// wait for done
	wg.Wait()

	// wait up to 6 seconds for server to exit
	for i := 0; i < 6; i++ {
		if shutdownCalled {
			continue
		} else {
			time.Sleep(time.Second)
		}
	}

	if !shutdownCalled {
		t.Error("Shutdown not called on shutdownChan")
	}
}

func TestSelectInterruptChannel(t *testing.T) {
	t.Log("SelectInterruptChannel")
	handler := http.HandlerFunc(MainHandler)
	// Again create an instance of both for testing
	ts := httptest.NewServer(handler)
	server := &http.Server{Addr: ts.Config.Addr}

	shutdownCalled := false
	wg.Add(1)

	server.RegisterOnShutdown(func() { shutdownCalled = true })
	defer server.Close()
	server.ListenAndServe()

	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel and associated handler for PUTs to /shutdown
	shutdownChan := make(chan bool)
	go SelectChannel(server, interruptChan, shutdownChan)

	interruptChan <- syscall.SIGINT

	// wait for done
	wg.Wait()

	// wait up to 6 seconds for server to exit
	for i := 0; i < 6; i++ {
		if shutdownCalled {
			continue
		} else {
			time.Sleep(time.Second)
		}
	}

	if !shutdownCalled {
		t.Error("Shutdown not called on interruptChan")
	}

}

func MainHandler(res http.ResponseWriter, req *http.Request) {
	return
}

func TestLogHandler(t *testing.T) {
	t.Log("LogHandler")
	var buf bytes.Buffer
	log.SetOutput(&buf)
	// cleanup when we exit
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	handler := LogHandler(http.DefaultServeMux)
	server := httptest.NewServer(http.Handler(handler))
	defer server.Close()
	wg.Add(1)

	expectedString := "on / using Go-http-client/1.1"
	_, err := http.Get(server.URL)
	if err != nil {
		t.Error("Error testing GET for LogHandler:", err)
	}
	out := buf.String()
	if !strings.Contains(out, expectedString) {
		t.Errorf("Expected '%s' to contain '%s'", out, expectedString)
	}

}
