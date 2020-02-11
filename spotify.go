package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/zmb3/spotify"
)

const redirectURL = "https://localhost:8080/auth-callback"

var (
	auth     = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
	clientCh = make(chan *spotify.Client)
)

type authInfo struct {
	clientID     string `json:clientId`
	clientSecret string `json:clientSecret`
}

func listenForAuth(state string) *spotify.Client {
	http.HandleFunc(
		"/auth-callback",
		func(w http.ResponseWriter, r *http.Request) {
			tok, err := auth.Token(state, r)
			if err != nil {
				http.Error(w, "Couldn't get token", http.StatusForbidden)
				log.Fatal(err)
			}
			if st := r.FormValue("state"); st != state {
				http.NotFound(w, r)
				log.Fatalf("State mismatch: %s != %s\n", st, state)
			}

			client := auth.NewClient(tok)
			fmt.Fprintf(w, "Login Completed!")
			clientCh <- &client
		})

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got request for: ", r.URL.String())

		})

	go http.ListenAndServe(":8080", nil)

	return <-clientCh

}

func getAuthInfo() (string, string) {
	jsonFile, err := os.Open(".client-secrets.json")
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	byteVal, err := ioutil.ReadAll(jsonFile)

	var authKey authInfo

	json.Unmarshal(byteVal, &authKey)

	return authKey.clientID, authKey.clientSecret
}

// Authenticate authenticates and returns a new spotify.WebAPIClient
func Authenticate() (*spotify.Client, error) {
	// The redirect URL must be EXACTLY the same as the registered URL
	auth.SetAuthInfo(getAuthInfo())

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	url := auth.AuthURL(id.String())
	fmt.Println("To authenticate, please go to this URL: %s", url)

	return listenForAuth(id.String()), nil
}
