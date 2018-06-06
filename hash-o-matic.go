package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_LISTEN_ADDR = ":8080"

var (
	wg             sync.WaitGroup
	shutdownChan   chan bool
	interruptChan  chan os.Signal
	HashCount      int64
	HashCountMutex sync.Mutex
	Hashes         map[int64]string
	HashesMutex    sync.Mutex
)

type redirect struct {
	target string
	code   int
}

const ShutdownTimeout = 30 * time.Second

func main() {

	serverAddress := DEFAULT_LISTEN_ADDR
	if len(os.Getenv("PORT")) > 0 {
		serverAddress = ":" + os.Getenv("PORT")
	}

	HashesMutex.Lock()
	Hashes = make(map[int64]string)
	HashesMutex.Unlock()

	// handlers for configured endpoints
	http.HandleFunc("/", RedirectHandler)
	http.HandleFunc("/hash", HashPostHandler)
	http.HandleFunc("/hash/", HashGetHandler)
	http.HandleFunc("/shutdown", ShutdownHandler)

	// basic logging of connections
	handler := LogHandler(http.DefaultServeMux)

	// create the server, used in SelectChannel for graceful shutdown
	server := &http.Server{Addr: serverAddress, Handler: handler}

	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
	interruptChan = make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel and associated handler for PUTs to /shutdown
	shutdownChan = make(chan bool)

	// setup to wait for an interrupt or shutdown call
	go SelectChannel(server, interruptChan, shutdownChan)

	log.Printf("Server listening on: %s", server.Addr)
	// add server to waitgroup
	wg.Add(1)
	go server.ListenAndServe()

	// wait for everything to finish
	wg.Wait()

	log.Println("Shutdown complete.")

}

func SelectChannel(
	server *http.Server,
	interruptChan chan os.Signal,
	shutdownChan chan bool) {

	select {
	case n := <-interruptChan:
		log.Printf("Received signal %s; shutting down\n", n.String())
		StopServer(server)

	case _ = <-shutdownChan:
		log.Printf("Received call to /shutdown, shutting down\n")
		StopServer(server)
	}

}

func StopServer(server *http.Server) {
	defer wg.Done()
	log.Println("Stopping server...")

	// create context with a max timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error while shutting down: %s", err)
	}
}

func ShutdownHandler(res http.ResponseWriter, req *http.Request) {
	// PUT because shutdown is a modification
	if req.Method == http.MethodPut {
		// shutdown issued here
		res.WriteHeader(http.StatusAccepted)
		res.Write([]byte("shutting down..."))
		shutdownChan <- true
	} else {
		// no more supported methods to match return 404, write header first
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte("Method not allowed\n"))
	}
}

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		log.Printf("Incoming %s request from %s on %s using %s\n",
			req.Method, req.RemoteAddr, req.URL, req.UserAgent())
		handler.ServeHTTP(res, req)
		log.Printf("%s Request from %s Complete\n", req.Method, req.RemoteAddr)
		log.Printf("%d hashes saved\n", HashCount)
	})
}

func getRedirectTargets() (redirectMap map[string]redirect) {
	redirectMap = make(map[string]redirect)
	redirectMap["/"] = redirect{
		"https://github.com/dbyington/hash-o-matic#readme",
		http.StatusFound,
	}

	return redirectMap
}

func mapRedirect(requestUrl string) (redirectUrl redirect) {

	var redir redirect
	var err error
	var ok bool
	redirectMap := getRedirectTargets()

	if redir, ok = redirectMap[requestUrl]; ok {
		if err != nil {
			log.Printf("Error parsing redirect url %s, got %s", redir.target, err.Error())
		}
	}
	return redir
}

func RedirectHandler(res http.ResponseWriter, req *http.Request) {

	redirectTo := mapRedirect(req.RequestURI)
	if redirectTo.code > 0 {
		http.Redirect(res, req, redirectTo.target, redirectTo.code)
	}
}
