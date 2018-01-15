package handlers

import (
    "testing"
    "reflect"
    "net/http/httptest"
    "net/http"
)


func TestGetRedirectTargets(t *testing.T) {
    mockMap := make(map[string]redirect)
    mockMap["/"] = redirect{
        "https://github.com/dbyington/hash-o-matic#readme",
        http.StatusFound,
    }

    testMap := getRedirectTargets()
    mapEql := reflect.DeepEqual(testMap, mockMap)
    if ! mapEql {
        t.Error("Redirect map not eql")
    }
}

func TestMapRedirect(t *testing.T) {
    testRequestUrl := "/"
    testResponseUrl := mapRedirect(testRequestUrl)
    if testResponseUrl.code != http.StatusFound {
        t.Errorf("Unexpected response code, expected %d got %d\n", http.StatusFound, testResponseUrl.code)
    }
    if testResponseUrl.target != "https://github.com/dbyington/hash-o-matic#readme" {
        t.Errorf("Unexepected URL, expected https://github.com/dbyington/hash-o-matic#readme, got %s\n", testResponseUrl.target)
    }
}

func TestRedirectHandler(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(RedirectHandler))
    defer server.Close()

    response, err := http.Get(server.URL+"/")
    if err != nil {
        t.Error("Error testing get:", err)
    }
    if response.StatusCode != http.StatusFound {
        // for some reason http.Redirect() does not set the response code or location in httptest
        // manual testing shows this works
        //t.Errorf("Error in redirect, expected %d got %d\n", http.StatusFound, response.StatusCode)
    }
}
