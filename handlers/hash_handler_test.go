package handlers

import (
    "net/http"
    "testing"
    "strings"
    "net/url"
    "io/ioutil"
)

const SERVER_ADDR = "http://localhost:8080"
const HASH_URL = "/hash"
const HASH_EXPECTED = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="


func TestHashPostMethod(t *testing.T) {
    var body = url.Values{}
    body.Set("password", "angryMonkey")
    target := SERVER_ADDR + HASH_URL

    t.Log("Valid request should return 201")
    req, err := http.NewRequest(http.MethodPost, target, strings.NewReader(body.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res, err := http.DefaultClient.Do(req)
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
    req, err = http.NewRequest(http.MethodPost, target, strings.NewReader(""))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusBadRequest {
        t.Errorf("POST to %s; expected %d got %d\n", HASH_URL, http.StatusBadRequest, res.StatusCode)
    }


    t.Log("PUT to /hash should return 404")
    req, err = http.NewRequest(http.MethodPut, target, strings.NewReader(body.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("PUT to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

    t.Logf("GET to /hash should return 404")
    req, err = http.NewRequest(http.MethodGet, target, strings.NewReader(""))
    res, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if res.StatusCode != http.StatusNotFound {
        t.Errorf("GET to %s; expected %d got %d\n", HASH_URL, http.StatusNotFound, res.StatusCode)
    }

    res.Body.Close()
}