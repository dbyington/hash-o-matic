package handlers

import (
    "net/http"
    "log"
)

type redirect struct {
    target string
    code int
}
var redirectMap map[string]redirect

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
