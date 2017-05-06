package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

var client *spotify.Client
var config Config

func main() {

	config, err := parseConfig("config.json")
	if err != nil {
		fmt.Println("Could not load config")
	}

	client = newClient(config.Key, config.Secret)

	r := mux.NewRouter()
	r.HandleFunc("/search/{keyword}", SearchHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songSearch := vars["keyword"]
	fmt.Printf("Searching for %s\n", songSearch)
	results, err := spotify.Search(songSearch, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	track := results.Tracks.Tracks[0]

	var resultJson = fmt.Sprintf("{\"name\": \"%s\", \"artist\": \"%s\", \"URI\": \"%s\"}", track.SimpleTrack.Name, track.SimpleTrack.Artists[0].Name, track.URI)
	fmt.Println(resultJson)
	w.Write([]byte(resultJson))
}
