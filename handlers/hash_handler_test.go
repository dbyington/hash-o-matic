package handlers

import (
    "net/http"
    "testing"
    "net/url"
    "io/ioutil"
    "net/http/httptest"
    "strings"
)

const HASH_URL = "/hash"
const HASH_EXPECTED = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="


func TestHashPostMethod(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(HashPostMethod))
    defer server.Close()

    body := url.Values{}
    body.Set("password", "angryMonkey")

    t.Log("Valid request should return 201")
    res, err := http.PostForm(server.URL, body)
    defer res.Body.Close()
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusCreated {
        t.Errorf("POST to %s failed; expected %d got %d\n", HASH_URL, http.StatusCreated, res.StatusCode)
    }
    resBody, err := ioutil.ReadAll(res.Body)
    if string(resBody) != HASH_EXPECTED {
        t.Errorf("Bad response; expected '%s' got '%s'", HASH_EXPECTED, string(resBody))
    }

    t.Log("No form data should return 400")
    res, err = http.PostForm(server.URL, url.Values{})
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusBadRequest {
        t.Errorf("POST to %s; expected %d got %d\n", HASH_URL, http.StatusBadRequest, res.StatusCode)
    }

}

func TestHashHandler(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(HashHandler))
    defer server.Close()

    t.Log("PUT to /hash should return 404")
    req, err := http.NewRequest(http.MethodPut, server.URL, strings.NewReader(""))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("PUT to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

    t.Logf("GET to /hash should return 404")
    req, err = http.NewRequest(http.MethodGet, server.URL, strings.NewReader(""))
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("GET to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

}