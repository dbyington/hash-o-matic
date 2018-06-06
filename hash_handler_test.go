package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const HASH_URL = "/hash"
const HASH_EXPECTED = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

func TestHashPostHandler(t *testing.T) {
	t.Log("HashPostHandler")
	Hashes = make(map[int64]string)
	server := httptest.NewServer(http.HandlerFunc(HashPostHandler))
	wg.Add(1)
	defer server.Close()

	body := url.Values{}
	body.Set("password", "angryMonkey")

	t.Log("Should return", http.StatusAccepted)
	res, err := http.PostForm(server.URL, body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("POST to %s failed; expected %d got %d\n", HASH_URL, http.StatusAccepted, res.StatusCode)
	}
	var resBody struct{ HashId int64 }
	json.NewDecoder(res.Body).Decode(&resBody)
	if resBody.HashId != 1 {
		t.Errorf("Bad response; expected 1 got '%d'", resBody.HashId)
	}

	t.Log("Should accept a json body")
	jsonBody := map[string]string{"password": "angryMonkey"}
	bodyBuf := new(bytes.Buffer)
	json.NewEncoder(bodyBuf).Encode(jsonBody)
	res, err = http.Post(server.URL, "application/json", bodyBuf)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("JSON POST failed, expected %d got %d\n", http.StatusAccepted, res.StatusCode)
	}

	t.Log("Should return 400 on empty post")
	res, err = http.PostForm(server.URL, url.Values{})
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("POST to %s; expected %d got %d\n", HASH_URL, http.StatusBadRequest, res.StatusCode)
	}
	var emptyJson struct{ Password string }
	emptyJson.Password = ""
	bodyBuf = new(bytes.Buffer)
	json.NewEncoder(bodyBuf).Encode(emptyJson)
	res, err = http.Post(server.URL, "application/json", bodyBuf)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("JSON POST failed, expected %d got %d\n", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHashGetHandler(t *testing.T) {
	t.Log("HashGetHandler")
	Hashes = make(map[int64]string)
	Hashes[1] = HASH_EXPECTED

	server := httptest.NewServer(http.HandlerFunc(HashGetHandler))
	wg.Add(1)
	defer server.Close()

	body := url.Values{}
	body.Set("password", "angryMonkey")

	res, err := http.PostForm(server.URL, body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var resHashId struct{ HashId int64 }
	resHashId.HashId = 1

	var resHashString struct{ HashString string }
	t.Logf("Should return %d on GET to /hash/%d", http.StatusOK, resHashId.HashId)
	getUrl := server.URL + "/hash/1"

	res, err = http.Get(getUrl)
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
