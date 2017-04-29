package main

import (
	"fmt"
	"log"
  "errors"
  "encoding/json"
	"net/http"
  "io/ioutil"
	"github.com/zmb3/spotify"
  "golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch    = make(chan *oauth2.Token)
	state = "miles_spotify_player"
)

func newClient() *spotify.Client {
  var token *oauth2.Token
  var err error
	if token, err = getCachedToken(); err == nil {
  } else {
    token, err = getNewToken()
    cacheToken(token)
  }

  client := auth.NewClient(token)
  return &client
}

func getCachedToken() (*oauth2.Token, error) {
  tokB, err := ioutil.ReadFile("spotify_token.json")
  if err != nil {
    return nil, errors.New("No cached token")
  }

  var token oauth2.Token
  if err := json.Unmarshal(tokB, &token); err != nil {
    panic(err)
  }

  return &token, nil
}

func getNewToken() (*oauth2.Token, error) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

  token := <-ch
  return token, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	ch <- tok
}

func cacheToken(token *oauth2.Token) {
  tokJson, _ := json.Marshal(token)
  err := ioutil.WriteFile("spotify_token.json", tokJson, 0644)
  if err != nil {
      panic(err)
  }
}
