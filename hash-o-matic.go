package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_LISTEN_ADDR = ":8080"

type hashes struct {
	list     []string
	mutex    sync.RWMutex
	wg       sync.WaitGroup
	stats    *hashStat
	save     chan hash
	request  chan hash
	response chan hash
}

type hash struct {
	id   int64
	hash string
}

type hashServer struct {
	http.Server
	mutex         sync.Mutex
	wg            sync.WaitGroup
	shutdownChan  chan struct{}
	interruptChan chan os.Signal
}

type hashStat struct {
	mutex    sync.RWMutex
	addTime  chan int64
	get      chan struct{}
	response chan hashStats
	hashStats
}
type hashStats struct {
	requestCount int64
	milliseconds int64
}

func main() {

	serverAddws := DEFAULT_LISTEN_ADDR
	if len(os.Getenv("PORT")) > 0 {
		serverAddws = ":" + os.Getenv("PORT")
	}

	// create our hashes and run the loop to handle save and requests for hashes
	h := newHashes()
	go h.run()

	// create the s, used in ShutdownListen for graceful shutdown
	s := &hashServer{}

	// handlers for configured endpoints
	http.HandleFunc("/hash", h.PostHandler)
	http.HandleFunc("/hash/", h.GetHandler)
	http.HandleFunc("/shutdown", s.ShutdownHandler)
	http.HandleFunc("/stats", h.StatsHandler)

	// basic logging of connections
	handler := s.LogHandler(http.DefaultServeMux)

	s.Addr = serverAddws
	s.Handler = handler

	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
	s.interruptChan = make(chan os.Signal)
	signal.Notify(s.interruptChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel and associated handler for PUTs to /shutdown
	s.shutdownChan = make(chan struct{})

	// setup to wait for an interrupt or shutdown call
	go s.ShutdownListen()

	log.Printf("Server listening on: %s", s.Addr)
	// add s to waitgroup
	s.wg.Add(1)
	go s.ListenAndServe()

	// wait for s to finish
	s.wg.Wait()

	// wait for inflight hashes to be saved
	h.wg.Wait()

	log.Println("Shutdown complete.")

}

func (s *hashServer) ShutdownListen() {
	select {
	case n := <-s.interruptChan:
		log.Printf("Received signal %s; shutting down\n", n.String())
		s.Stop()

	case <-s.shutdownChan:
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

func (s *hashServer) ShutdownHandler(w http.ResponseWriter, req *http.Request) {
	// PUT because shutdown is a modification
	if req.Method == http.MethodPut {
		// shutdown issued here
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("shutting down..."))
		s.shutdownChan <- struct{}{}
	} else {
		// no more supported methods to match return 404, write header first
		ErrorResponse(w, http.StatusMethodNotAllowed)
	}
}

func (s *hashServer) LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		log.Printf("Incoming %s request from %s on %s using %s\n",
			req.Method, req.RemoteAddr, req.URL, req.UserAgent())
		next.ServeHTTP(w, req)
		log.Printf("%s Request from %s Complete\n", req.Method, req.RemoteAddr)
	})
}

func newHashes() *hashes {
	s := &hashStat{
		get:      make(chan struct{}),
		response: make(chan hashStats),
		addTime:  make(chan int64),
	}
	return &hashes{
		save:     make(chan hash),
		request:  make(chan hash),
		response: make(chan hash),
		stats:    s,
	}
}

func (h *hashes) PostHandler(w http.ResponseWriter, req *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	if req.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	defer h.addTime(start)

	contentType := req.Header.Get("Content-Type")
	var password string
	log.Printf("got content type: %s", contentType)

	// handle either a json or form url encoded body
	switch contentType {
	case "application/json":

		var jsonBody struct{ Password string }
		fullBody, err := ioutil.ReadAll(req.Body)
		json.Unmarshal(fullBody, &jsonBody)
		if err != nil {
			log.Println("ERROR:", err)
			ErrorResponse(w, http.StatusBadRequest, "missing password")
			return
		}
		password = jsonBody.Password

	case "application/x-www-form-urlencoded":
		req.ParseForm()
		password = req.PostFormValue("password")
	}

	if len(password) > 0 {
		var hash hash
		hash.hash = HashPassword(password)
		h.save <- hash
		hashResponse := <-h.response
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]int64{"hashId": hashResponse.id})

	} else {
		ErrorResponse(w, http.StatusBadRequest, "missing password")
	}
}

func (h *hashes) GetHandler(w http.ResponseWriter, req *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	if req.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	// make sure the URI is in the form '/hash/:id'
	regx := regexp.MustCompile(`/hash/(\d+)$`)

	if !regx.MatchString(req.RequestURI) {
		ErrorResponse(w, http.StatusNotFound)
		return
	}

	reqId := regx.FindAllStringSubmatch(req.RequestURI, -1)[0][1]
	log.Println("Hash id requested:", reqId)
	hashId, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound)
		return
	}

	var hash hash
	hash.id = hashId
	h.request <- hash
	hashResponse := <-h.response

	if hashResponse.id < 1 {
		ErrorResponse(w, http.StatusNotFound)
		return
	}

	if hashResponse.hash == "" {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "Hash string not ready"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"hashString": hashResponse.hash})

}

func (h *hashes) StatsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	h.stats.get <- struct{}{}
	stats := <-h.stats.response
	t := stats.milliseconds
	if stats.requestCount != 0 {
		t = stats.milliseconds / stats.requestCount
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"requests": stats.requestCount, "averageNanoseconds": t})
}

func (h *hashes) addTime(s time.Time) {
	e := time.Now()
	h.stats.addTime <- e.Sub(s).Nanoseconds()
}

func (h *hashes) run() {
	for {
		select {
		// a save request will come in with a hReq but a default it (0) so we just ignore the id sent and add our own.
		case hSave := <-h.save:
			// append an empty string to the hReq list to reserve that spot for this hReq
			h.list = append(h.list, "")
			hSave.id = int64(len(h.list))

			// add to our wg as we're going to do the saving in another go routine that will take > 5 seconds to complete
			h.wg.Add(1)
			go func(ha hash) {
				defer h.wg.Done()
				time.Sleep(5 * time.Second)

				h.mutex.Lock()
				h.list[ha.id-1] = ha.hash
				h.mutex.Unlock()
			}(hSave)

			// respond with the hReq, now with the id
			h.response <- hSave

			// a request for a hReq will come in with the id so we'll ignore the hReq, if any.
		case hReq := <-h.request:
			h.mutex.RLock()
			// if the hReq id is not within the length of our list respond with a negative id and empty hReq
			if hReq.id < 1 || hReq.id > int64(len(h.list)) {
				hReq.id = 0
				hReq.hash = ""
			} else {
				hReq.hash = h.list[hReq.id-1]
			}
			h.mutex.RUnlock()

			// respond with the full hReq
			h.response <- hReq

			// Add time to the hash hashStats
		case t := <-h.stats.addTime:
			h.stats.milliseconds += t
			h.stats.requestCount++

			// Signal to get the hashStats
		case <-h.stats.get:
			h.stats.response <- h.stats.hashStats
		}
	}
}

func HashPassword(password string) string {
	// Sum512() argument must be type [size]byte
	byteArray := []byte(password)
	sum512byteArray := sha512.Sum512(byteArray)

	// EncodeToString() argument must be type []byte so first turn the byte array into a string
	sum512string := string(sum512byteArray[:])
	base64string := base64.StdEncoding.EncodeToString([]byte(sum512string))

	return base64string
}

func ErrorResponse(w http.ResponseWriter, status int, message ...string) {
	var m string
	w.WriteHeader(status)

	if len(message) > 0 {
		m = fmt.Sprintf("%s: %s", http.StatusText(status), message[0])
	} else {
		m = http.StatusText(status)
	}
	json.NewEncoder(w).Encode(map[string]string{"ErrorMessage": m})

}
