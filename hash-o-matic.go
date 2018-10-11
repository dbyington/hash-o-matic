package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_LISTEN_ADDR = ":8080"

var (
	shutdownChan  chan bool
	interruptChan chan os.Signal
)

type Hashes struct {
	Count int64
	list  map[int64]string
	mutex sync.Mutex
	wg    sync.WaitGroup
}

type hashServer struct {
	http.Server
	wg    sync.WaitGroup
	mutex sync.Mutex
}

func main() {

	serverAddress := DEFAULT_LISTEN_ADDR
	if len(os.Getenv("PORT")) > 0 {
		serverAddress = ":" + os.Getenv("PORT")
	}

	hashes := &Hashes{list: make(map[int64]string)}
	// handlers for configured endpoints
	http.HandleFunc("/hash", hashes.PostHandler)
	http.HandleFunc("/hash/", hashes.GetHandler)
	http.HandleFunc("/shutdown", ShutdownHandler)

	// basic logging of connections
	handler := LogHandler(http.DefaultServeMux)

	// create the server, used in SelectChannel for graceful shutdown
	server := &hashServer{}
	server.Addr = serverAddress
	server.Handler = handler

	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
	interruptChan = make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel and associated handler for PUTs to /shutdown
	shutdownChan = make(chan bool)

	// setup to wait for an interrupt or shutdown call
	go server.SelectChannel()

	log.Printf("Server listening on: %s", server.Addr)
	// add server to waitgroup
	server.wg.Add(1)
	go server.ListenAndServe()

	// wait for server to finish
	server.wg.Wait()

	// wait for inflight hashes to be saved
	hashes.wg.Wait()

	log.Println("Shutdown complete.")

}

func (s *hashServer) SelectChannel() {
	select {
	case n := <-interruptChan:
		log.Printf("Received signal %s; shutting down\n", n.String())
		s.Stop()

	case _ = <-shutdownChan:
		log.Printf("Received call to /shutdown, shutting down\n")
		s.Stop()
	}

}

func (s *hashServer) Stop() {
	defer s.wg.Done()
	log.Println("Stopping server...")

	// create context with a max timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
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
	})
}

func (h *Hashes) PostHandler(res http.ResponseWriter, req *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	if req.Method != http.MethodPost {
		MethodNotAllowed(res)
		return
	}

	contentType := req.Header.Get("Content-Type")
	var password string

	// handle either a json or form url encoded body
	switch contentType {
	case "application/json":
		var jsonBody struct{ Password string }
		fullBody, err := ioutil.ReadAll(req.Body)
		json.Unmarshal(fullBody, &jsonBody)
		if err != nil {
			log.Println("ERROR:", err)
			badRequest(res, "missing password")
			return
		}
		password = jsonBody.Password

	case "application/x-www-form-urlencoded":
		req.ParseForm()
		password = req.PostFormValue("password")
	}

	if len(password) > 0 {

		h.mutex.Lock()
		h.Count++
		h.mutex.Unlock()

		res.WriteHeader(http.StatusAccepted)
		json.NewEncoder(res).Encode(map[string]int64{"hashId": h.Count})

		h.wg.Add(1)
		go h.save(h.Count, strings.Join([]string{password}, ""))

	} else {
		badRequest(res, "missing password")
	}
}

func (h *Hashes) GetHandler(res http.ResponseWriter, req *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	if req.Method != http.MethodGet {
		MethodNotAllowed(res)
		return
	}

	// make sure the URI is in the form '/hash/:id'
	regx := regexp.MustCompile(`/hash/(\d+)$`)

	if !regx.MatchString(req.RequestURI) {
		ReplyNotFound(res)
		return
	}

	reqId := regx.FindAllStringSubmatch(req.RequestURI, -1)[0][1]
	log.Println("Hash id requested:", reqId)
	hashId, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		ReplyNotFound(res)
		return
	}

	if hashId > h.Count || hashId < 1 {
		ReplyNotFound(res)
		return
	}

	if h.list[hashId] == "" {
		res.WriteHeader(http.StatusAccepted)
		json.NewEncoder(res).Encode(map[string]string{"status": "Hash string not ready"})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{"hashString": h.list[hashId]})

}

func (h *Hashes) save(Id int64, passwd string) {
	defer h.wg.Done()

	time.Sleep(5 * time.Second)

	h.mutex.Lock()
	h.list[Id] = h.Hash(passwd)
	h.mutex.Unlock()
}

func (h *Hashes) Hash(hash string) string {
	// Sum512() argument must be type [size]byte
	byteArray := []byte(hash)
	sum512byteArray := sha512.Sum512(byteArray)

	// EncodeToString() argument must be type []byte so first turn the byte array into a string
	sum512string := string(sum512byteArray[:])
	base64string := base64.StdEncoding.EncodeToString([]byte(sum512string))

	return base64string
}

func ReplyNotFound(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	errorReply := map[string]string{"ErrorMessage": "Page not found"}
	json.NewEncoder(res).Encode(errorReply)
}
func MethodNotAllowed(res http.ResponseWriter) {
	res.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(res).Encode(map[string]string{"ErrorMessage": "Method not allowed"})
}

func badRequest(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad request; " + message))
}
