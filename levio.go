package main

import (
    "errors"
    "strings"
    "io/ioutil"
    "encoding/base64"
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "encoding/json"
    "net/http"
    "github.com/levicole/lev.io/models"
    "os"
)

type APIHandler struct {
    // Verify requests with hmac signature
    VerifyRequests bool
}

func (this *APIHandler) VerifyRequest(auth string, body []byte) error {
    if !this.VerifyRequests {
        return nil
    }
    data, _ := base64.StdEncoding.DecodeString(auth)
    mac := hmac.New(sha256.New, []byte(os.Getenv("LEVIO_KEY")))
    mac.Write(body)
    expectedMac := mac.Sum(nil)
    if !hmac.Equal(expectedMac, []byte(data)) {
        return errors.New("Unauthorized")
    }
    return nil
}

func (this *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    this.ProcessRequest(w, r)
}

func (this *APIHandler) ProcessRequest(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body);

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    auth, ok := r.Header["Authorization"];
    if !ok {
        http.Error(w, "Bad Request", http.StatusBadRequest)
    }

    if err := this.VerifyRequest(auth[0], body); err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
    }

    switch r.Method {
    case "POST":
        link := &models.Link{}
        enc := json.NewEncoder(w)
        if err := json.Unmarshal(body, link); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
        }
        link.Save()
        enc.Encode(link)
    }
}

// handle URLs for shortened links
// lev.io/:path
func rootHandler(w http.ResponseWriter, r *http.Request) {
    path := strings.Split(r.URL.Path[1:], "/")
    query := "select Id, Url from links l where slug=?"
    var links []*models.LinkView

    if err := models.Find(&links, query, path[0]); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    if len(links) > 0 {
        http.Redirect(w, r, links[0].Url, 302)
    } else {
        http.NotFound(w, r)
    }
}

func main() {
    http.HandleFunc("/", rootHandler)
    http.Handle("/api/", &APIHandler{true})
    var port string

    port = os.Getenv("LEVIO_PORT");

    if  port == "" {
        port = "8080"
    }
    http.ListenAndServe(":" + port, nil)
}
