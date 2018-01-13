package handlers

import (
    "net/http"
    "testing"
    "net/url"
    "net/http/httptest"
    "strings"
    "encoding/json"
    "time"
    "bytes"
)

const HASH_URL = "/hash"
const HASH_EXPECTED = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="


func TestHashPostMethod(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(HashPostMethod))
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
    var resBody HashIdReply
    json.NewDecoder(res.Body).Decode(&resBody)
    if resBody.HashId != 1 {
        t.Errorf("Bad response; expected 1 got '%d'", resBody.HashId)
    }

    t.Log("Should accept a json body")
    jsonBody := JsonBody{ Password: "angryMonkey"}
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
    var emptyJson JsonBody
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

func TestHashHandler(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(HashHandler))
    defer server.Close()

    t.Log("Should return 404 on PUT to /hash")
    req, err := http.NewRequest(http.MethodPut, server.URL, strings.NewReader(""))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("PUT to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

    t.Logf("Should return 404 on GET to /hash")
    req, err = http.NewRequest(http.MethodGet, server.URL, strings.NewReader(""))
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("GET to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

}

func TestHashGetMethod(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(HashHandler))
    defer server.Close()

    body := url.Values{}
    body.Set("password", "angryMonkey")

    res, err := http.PostForm(server.URL, body)
    defer res.Body.Close()
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusAccepted {
        t.Errorf("POST to %s failed; expected %d got %d\n", HASH_URL, http.StatusAccepted, res.StatusCode)
    }

    var resHashId HashIdReply
    json.NewDecoder(res.Body).Decode(&resHashId)
    t.Log("Got hashId:", resHashId.HashId)
    if resHashId.HashId < 1 {
        t.Errorf("POST failed; expected int returnd got %s")
    }

    var resHashString HashStringReply
    t.Logf("Should return %d on early GET to /hash/%d", http.StatusAccepted, resHashId.HashId)
    getUrl := server.URL + "/hash/" + string(resHashId.HashId)
    res, err = http.Get(getUrl)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusAccepted {
        t.Errorf("GET to %s; expected %d got %d\n", getUrl, http.StatusAccepted, res.StatusCode)
    }
    json.NewDecoder(res.Body).Decode(&resHashString)
    if resHashString.HashString != "" {
        t.Error("Expected hash string to be empty, got:", resHashString.HashString)
    }

    t.Log("Should return 200 after 5 seconds")
    time.Sleep(5*time.Second)
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