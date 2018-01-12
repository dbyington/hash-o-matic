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

const DEFAULT_LISTEN_ADDR = ":8080"

func main() {

    serverAddress := DEFAULT_LISTEN_ADDR
    if len(os.Getenv("PORT")) > 0 {
        serverAddress = ":" + os.Getenv("PORT")
    }

    // create the server, channels, and routes
    server := BuildServer(serverAddress)
    shutdownChan, interruptChan, doneChan := BuildChannels()
    BuildRouteHandlers(shutdownChan)

    // setup to wait for an interrupt or shutdown call
    go SelectChannel(server, interruptChan, shutdownChan, doneChan)

    log.Printf("Server listening on: %s", server.Addr)
    err := server.ListenAndServe()
    if err != http.ErrServerClosed {
        log.Fatalf("listen: %s\n", err)
    }

    <-doneChan
    log.Println("Shutdown complete.")

}

func BuildServer(address string) (server *http.Server) {
    if address == "" {
        address = DEFAULT_LISTEN_ADDR
    }
    log.Printf("Address %s", address)
    handler := handlers.LogHandler(http.DefaultServeMux)
    server = &http.Server{Addr: address, Handler: handler}
    return server
}

func BuildChannels() (
    shutdownChan chan bool,
    interruptChan chan os.Signal,
    doneChan chan bool) {

    // Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
    interruptChan = make(chan os.Signal)
    signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

    // Create a channel and associated handler for PUTs to /shutdown
    shutdownChan = make(chan bool)

    // Create channel to signal all done
    doneChan = make(chan bool)

    return shutdownChan, interruptChan, doneChan
}

func BuildRouteHandlers(shutdownChan chan bool) {
    ShutdownHandler := handlers.BuildShutdownHandler(shutdownChan)
    http.HandleFunc("/hash", handlers.HashHandler)
    http.HandleFunc("/hash/", handlers.HashHandler)
    http.HandleFunc("/shutdown", ShutdownHandler)
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