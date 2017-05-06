package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"log"
	"net/http"
)

var client *spotify.Client
var config Config

type PlayCommand struct {
	TrackURI string `json:"uri"`
	DeviceID string `json:"device"`
}

type SpotifyDevice struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func main() {

	config, err := parseConfig("config.json")
	if err != nil {
		fmt.Println("Could not load config")
	}

	client = newClient(config.Key, config.Secret)

	r := mux.NewRouter()
	r.HandleFunc("/search/{keyword}", SearchHandler)
	r.HandleFunc("/play", PlayHandler).Methods("GET")
	r.HandleFunc("/play", PlayTrackHandler).Methods("POST")
	r.HandleFunc("/pause", PauseHandler)
	r.HandleFunc("/devices", DevicesHandler)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resultJson))
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	client.Play()
}

func PlayTrackHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var command PlayCommand
	json.Unmarshal(b, &command)

	pOpts := spotify.PlayOptions{URIs: []spotify.URI{spotify.URI(command.TrackURI)}}
	if len(command.DeviceID) > 0 {
		dID := spotify.ID(command.DeviceID)
		pOpts.DeviceID = &dID
	}
	_ = client.PlayOpt(&pOpts)
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	client.Pause()
}

func DevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, _ := client.PlayerDevices()

	// Build list of devices
	var sDevices []SpotifyDevice
	for _, d := range devices {
		sd := SpotifyDevice{Name: d.Name, ID: string(d.ID)}
		sDevices = append(sDevices, sd)
	}

	jDevices, _ := json.Marshal(sDevices)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jDevices))
}