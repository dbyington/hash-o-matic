package main

import (
    "github.com/dbyington/hash-o-matic/handlers"
    "log"
    "net/http"
    "context"
    "time"
    "os"
    "os/signal"
    "syscall"
)

const LISTEN_ADDR = ":8080"

func main() {

    handler := handlers.LogHandler(http.DefaultServeMux)
    server := &http.Server{Addr: LISTEN_ADDR, Handler: handler}

    // Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
    interruptChan := make(chan os.Signal)
    signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

    // Create a channel and associated handler for PUTs to /shutdown
    shutdownChan := make(chan bool)
    ShutdownHandler := handlers.BuildShutdownHandler(shutdownChan)

    // Create channel to signal all done
    doneChan := make(chan bool)

    // routes
    // If only '/hash/' was configured requests to '/hash' would be redirected
    // that should not happen, hence the separate routes
    http.HandleFunc("/hash", handlers.HashHandler)
    http.HandleFunc("/hash/", handlers.HashHandler)
    http.HandleFunc("/shutdown", ShutdownHandler)

    // wait for interrupt or shutdown call
    go SelectChannel(server, interruptChan, shutdownChan, doneChan)

    log.Printf("Server listening on: %s", server.Addr)
    err := server.ListenAndServe()
    if err != http.ErrServerClosed {
        log.Fatalf("listen: %s\n", err)
    }

    <-doneChan
    log.Println("Shutdown complete.")

}

func SelectChannel(
    server *http.Server,
    interruptChan chan os.Signal,
    shutdownChan chan bool,
    doneChan chan bool) {

    select {
    case n := <-interruptChan:
        log.Printf("Received signal %s; shutting down\n", n.String())
        StopServer(server)

    case _ = <-shutdownChan:
        log.Print("Received call to /shutdown, shutting down\n")
        StopServer(server)
    }

    close(doneChan)
}

func StopServer(server *http.Server) {

    log.Print("Stopping server")
    // create context with a timeout of 5 seconds to allow requests in-flight to finish
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err := server.Shutdown(ctx)
    if err != nil {
        log.Fatalf("Error while shutting down: %s", err)
    }
}