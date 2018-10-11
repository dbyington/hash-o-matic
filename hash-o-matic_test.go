package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"
    "time"
)

func TestStopServer(t *testing.T) {
	t.Log("StopServer")
	handler := http.HandlerFunc(MainHandler)

	// Cannot use an instance of httptest.NewServer() in place of *http.Server
	// So create both and use the addr from the test instance
	ts := httptest.NewServer(handler)
	server := &hashServer{}
	server.Addr = ts.URL
	defer server.Close()

	shutdownCalled := false
	server.wg.Add(1)
	server.RegisterOnShutdown(func() { shutdownCalled = true })
	server.ListenAndServe()

	time.Sleep(time.Second) // give the server time to start
	server.Stop()
	time.Sleep(6 * time.Second) // ensure server is shutdown
	if !shutdownCalled {
		t.Error("Shutdown not called")
	}

}

//func TestSelectShutdownChannel(t *testing.T) {
//	t.Log("SelectShutdownChannel")
//	handler := http.HandlerFunc(MainHandler)
//	// Again create an instance of both for testing
//	ts := httptest.NewServer(handler)
//    server := &hashServer{}
//    server.Addr = ts.URL
//
//	shutdownCalled := false
//	server.wg.Add(1)
//	server.RegisterOnShutdown(func() { shutdownCalled = true })
//	defer server.Close()
//	server.ListenAndServe()
//
//	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
//	interruptChan := make(chan os.Signal)
//	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
//
//	// Create a channel and associated handler for PUTs to /shutdown
//	shutdownChan := make(chan bool)
//	go server.SelectChannel()
//
//	// Send to shutdownChan
//	shutdownChan <- true
//
//	// wait for done
//	server.wg.Wait()
//
//	// wait up to 6 seconds for server to exit
//	for i := 0; i < 6; i++ {
//		if shutdownCalled {
//			continue
//		} else {
//			time.Sleep(time.Second)
//		}
//	}
//
//	if !shutdownCalled {
//		t.Error("Shutdown not called on shutdownChan")
//	}
//}

//func TestSelectInterruptChannel(t *testing.T) {
//	t.Log("SelectInterruptChannel")
//	handler := http.HandlerFunc(MainHandler)
//	// Again create an instance of both for testing
//	ts := httptest.NewServer(handler)
//	server := &hashServer{}
//	server.Addr = ts.Config.Addr
//
//	shutdownCalled := false
//	server.wg.Add(1)
//
//	server.RegisterOnShutdown(func() { shutdownCalled = true })
//	defer server.Close()
//	server.ListenAndServe()
//
//	// Create a channel and signal notifier to catch OS level interrupts (i.e. ^C)
//	interruptChan := make(chan os.Signal)
//	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
//
//	// Create a channel and associated handler for PUTs to /shutdown
//	//shutdownChan := make(chan bool)
//	go server.SelectChannel()
//
//	interruptChan <- syscall.SIGINT
//
//	// wait for done
//	server.wg.Wait()
//
//	// wait up to 6 seconds for server to exit
//	for i := 0; i < 6; i++ {
//		if shutdownCalled {
//			continue
//		} else {
//			time.Sleep(time.Second)
//		}
//	}
//
//	if !shutdownCalled {
//		t.Error("Shutdown not called on interruptChan")
//	}
//
//}


func MainHandler(res http.ResponseWriter, req *http.Request) {
	return
}

const HASH_URL = "/hash"
const HASH_EXPECTED = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

//func TestHashPostHandler(t *testing.T) {
//    t.Log("HashPostHandler")
//    h := &Hashes{}
//    server := httptest.NewServer(http.HandlerFunc(h.PostHandler))
//    h.wg.Add(1)
//    defer server.Close()
//
//    body := url.Values{}
//    body.Set("password", "angryMonkey")
//
//    t.Log("Should return", http.StatusAccepted)
//    res, err := http.PostForm(server.URL, body)
//    defer res.Body.Close()
//    if err != nil {
//        t.Error(err)
//    }
//    if res.StatusCode != http.StatusAccepted {
//        t.Errorf("POST to %s failed; expected %d got %d\n", HASH_URL, http.StatusAccepted, res.StatusCode)
//    }
//    var resBody struct{ HashId int64 }
//    json.NewDecoder(res.Body).Decode(&resBody)
//    if resBody.HashId != 1 {
//        t.Errorf("Bad response; expected 1 got '%d'", resBody.HashId)
//    }
//
//    t.Log("Should accept a json body")
//    jsonBody := map[string]string{"password": "angryMonkey"}
//    bodyBuf := new(bytes.Buffer)
//    json.NewEncoder(bodyBuf).Encode(jsonBody)
//    res, err = http.Post(server.URL, "application/json", bodyBuf)
//    if err != nil {
//        t.Error(err)
//    }
//    if res.StatusCode != http.StatusAccepted {
//        t.Errorf("JSON POST failed, expected %d got %d\n", http.StatusAccepted, res.StatusCode)
//    }
//
//    t.Log("Should return 400 on empty post")
//    res, err = http.PostForm(server.URL, url.Values{})
//    if err != nil {
//        t.Error(err)
//    }
//    if res.StatusCode != http.StatusBadRequest {
//        t.Errorf("POST to %s; expected %d got %d\n", HASH_URL, http.StatusBadRequest, res.StatusCode)
//    }
//    var emptyJson struct{ Password string }
//    emptyJson.Password = ""
//    bodyBuf = new(bytes.Buffer)
//    json.NewEncoder(bodyBuf).Encode(emptyJson)
//    res, err = http.Post(server.URL, "application/json", bodyBuf)
//    if err != nil {
//        t.Error(err)
//    }
//    if res.StatusCode != http.StatusBadRequest {
//        t.Errorf("JSON POST failed, expected %d got %d\n", http.StatusBadRequest, res.StatusCode)
//    }
//}
//
func TestHashGetHandler(t *testing.T) {
   t.Log("HashGetHandler")
   h := &Hashes{list: make(map[int64]string)}
   h.list[1] = HASH_EXPECTED

   server := httptest.NewServer(http.HandlerFunc(h.GetHandler))
   h.wg.Add(1)
   defer server.Close()

   //body := url.Values{}
   //body.Set("password", "angryMonkey")
   //
   //res, err := http.PostForm(server.URL, body)
   //defer res.Body.Close()
   //if err != nil {
   //    t.Error(err)
   //}

   var resHashId struct{ HashId int64 }
   resHashId.HashId = 1

   var resHashString struct{ HashString string }
   t.Logf("Should return %d on GET to /hash/%d", http.StatusOK, resHashId.HashId)
   getUrl := server.URL + "/hash/1"

   res, err := http.Get(getUrl)
   if err != nil {
       t.Error(err)
   }
   if res.StatusCode != http.StatusOK {
       t.Errorf("GET to %s; expected %d got %d\n", getUrl, http.StatusOK, res.StatusCode)
   }

   json.NewDecoder(res.Body).Decode(&resHashString)
   if resHashString.HashString != HASH_EXPECTED {
       t.Errorf("Hash mismatch; expected %s got %s", HASH_EXPECTED, resHashString.HashString)
   }
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

func TestHashString(t *testing.T) {
   mockPassword := "angryMonkey"
   mockHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

   h := &Hashes{}
   actualHash := h.Hash(mockPassword)
   if actualHash != mockHash {
       t.Errorf("Failed to hash %s; expected %s, got %s\n", mockPassword, mockHash, actualHash)
   }
}
