package handlers

import (
    "net/http"
    "time"
    "github.com/dbyington/hash-o-matic/util"
    "github.com/dbyington/hash-o-matic/database"
    "strings"
    "encoding/json"
    "regexp"
    "strconv"
    "log"
    "io/ioutil"
)

const FORM_KEY = "password"

type HashIdReply struct {
    HashId int
}

type HashStringReply struct {
    HashString string
}

type ErrorReply struct {
    ErrorMessage string
}

type JsonBody struct {
    Password string
}
func HashHandler(res http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        HashPostMethod(res, req)
        return
    case http.MethodGet:
        HashGetMethod(res, req)
        return
    default:
        // no more supported methods to match return 404, write header first
        ReplyNotFound(res)
    }
}
func ReplyNotFound(res http.ResponseWriter) {
    res.WriteHeader(http.StatusNotFound)
    errorReply := ErrorReply{ ErrorMessage: "Page not found"}
    json.NewEncoder(res).Encode(errorReply)
}

func HashPostMethod(res http.ResponseWriter, req *http.Request) {
    requestType := req.Header.Get("Content-Type")
    var password string
    switch requestType {
    case "application/json":
        var jsonBody JsonBody
        fullBody, err := ioutil.ReadAll(req.Body)
        json.Unmarshal(fullBody, &jsonBody)
        if err !=nil {
            log.Print("ERROR:", err)
            badRequest(res, "missing password")
            return
        }
        password = jsonBody.Password

    case "application/x-www-form-urlencoded":
        req.ParseForm()
        password = req.PostFormValue(FORM_KEY)
    }
    if len(password) > 0 {
        hashId, err := database.GetNextId()
        if err != nil {
            res.WriteHeader(http.StatusInternalServerError)
            res.Write([]byte("Error getting next hash id:" + err.Error()))
        }
        res.WriteHeader(http.StatusAccepted)
        var jsonReply HashIdReply
        jsonReply.HashId = hashId
        json.NewEncoder(res).Encode(jsonReply)
        go saveHash(strings.Join([]string{password}, ""), hashId)
    } else {
        badRequest(res, "missing password")
    }

}

func badRequest(res http.ResponseWriter, message string) {
    res.WriteHeader(http.StatusBadRequest)
    res.Write([]byte("Bad request; " + message))
}

func saveHash(passwd string, Id int) {
    time.Sleep(5 * time.Second)
    hashString := util.HashString(passwd)
    database.SaveHashWithId(hashString, Id)
}

func HashGetMethod(res http.ResponseWriter, req *http.Request) {
    rexp := regexp.MustCompile("[0-9]+$")
    id := rexp.FindString(req.RequestURI)
    hashId, err := strconv.Atoi(id)
    if err != nil {
        ReplyNotFound(res)
        return
    }

    if hashId > 0 {
        var jsonReply HashStringReply
        var err error
        jsonReply.HashString, err = database.GetHashById(hashId)
        if err != nil {
            jsonErr := ErrorReply{ ErrorMessage: err.Error() }
            res.WriteHeader(http.StatusAccepted)
            json.NewEncoder(res).Encode(jsonErr)
            return
        }
        res.WriteHeader(http.StatusOK)
        json.NewEncoder(res).Encode(jsonReply)
        return
    }
    ReplyNotFound(res)
}
