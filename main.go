package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"golang.org/x/crypto/ssh/terminal"
	"html/template"
	"log"
	"net/http"
)

var consumerKey, consumerSecret string
var credential *oauth.Credentials

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	t.Execute(w, nil)
}

func RequestTokenHandler(w http.ResponseWriter, r *http.Request) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	url, tmpCred, err := anaconda.AuthorizationURL("http://localhost:8000/access_token")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	credential = tmpCred
	http.Redirect(w, r, url, http.StatusFound)
}

func AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	c, _, err := anaconda.GetCredentials(credential, r.URL.Query().Get("oauth_verifier"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	t, err := template.ParseFiles("access_token.html")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	t.Execute(w, c)
}

func main() {
	fmt.Println("Enter ConsumerKey: ")
	ck, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	consumerKey = string(ck)

	fmt.Println("Enter ConsumerSecret: ")
	cs, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	consumerSecret = string(cs)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/request_token", RequestTokenHandler)
	http.HandleFunc("/access_token", AccessTokenHandler)
	http.ListenAndServe(":8000", nil)
}
