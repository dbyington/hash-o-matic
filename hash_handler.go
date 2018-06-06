package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ReplyNotFound(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	errorReply := map[string]string{"ErrorMessage": "Page not found"}
	json.NewEncoder(res).Encode(errorReply)
}
func MethodNotAllowed(res http.ResponseWriter) {
	res.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(res).Encode(map[string]string{"ErrorMessage": "Method not allowed"})
}

func HashPostHandler(res http.ResponseWriter, req *http.Request) {
	wg.Add(1)
	defer wg.Done()

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

		res.WriteHeader(http.StatusAccepted)
		HashCountMutex.Lock()
		HashCount++
		hashId := HashCount
		HashCountMutex.Unlock()

		json.NewEncoder(res).Encode(map[string]int64{"hashId": hashId})

		wg.Add(1)
		go saveHash(hashId, strings.Join([]string{password}, ""))

	} else {
		badRequest(res, "missing password")
	}

}

func badRequest(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad request; " + message))
}

func saveHash(Id int64, passwd string) {
	defer wg.Done()

	time.Sleep(5 * time.Second)

	HashesMutex.Lock()
	Hashes[Id] = HashString(passwd)
	HashesMutex.Unlock()
}

func HashGetHandler(res http.ResponseWriter, req *http.Request) {
	wg.Add(1)
	defer wg.Done()

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

	HashCountMutex.Lock()
	hashTotal := HashCount
	HashCountMutex.Unlock()
	if hashId > hashTotal || hashId < 1 {
		ReplyNotFound(res)
		return
	}

	HashesMutex.Lock()
	hash := Hashes[hashId]
	HashesMutex.Unlock()
	if hash == "" {
		res.WriteHeader(http.StatusAccepted)
		json.NewEncoder(res).Encode(map[string]string{"status": "Hash string not ready"})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{"hashString": hash})

}
