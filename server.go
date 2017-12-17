package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
)

var client *spotify.Client
var config Config

type PlayCommand struct {
	TrackURI string `json:"uri"`
	DeviceID string `json:"device"`
}

type SearchResult struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	URI    string `json:"uri"`
}

type SpotifyDevice struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type PlayerStatus struct {
	Playing    bool   `json:"playing"`
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
	URI        string `json:"url"`
	Progress   int    `json:"progress"`
	Track      string `json:"track"`
	Artist     string `json:"artist"`
}

type Playlist struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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
	r.HandleFunc("/status", StatusHandler)
	r.HandleFunc("/playlists", PlaylistsHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}

func WriteError(w http.ResponseWriter, err error) {
	errResponse := ErrorResponse{Error: err.Error()}
	errJson, _ := json.Marshal(errResponse)
	WriteResponse(w, string(errJson))
}

func WriteResponse(w http.ResponseWriter, response string) {
	log.Println(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songSearch := vars["keyword"]
	results, err := client.Search(songSearch, spotify.SearchTypeTrack)
	if err != nil {
		WriteError(w, err)
	}

	var searchResults []SearchResult
	for _, track := range results.Tracks.Tracks {
		searchResult := SearchResult{Name: track.SimpleTrack.Name, Artist: track.SimpleTrack.Artists[0].Name, URI: string(track.URI)}
		searchResults = append(searchResults, searchResult)
	}
	resultJson, err := json.Marshal(searchResults)
	if err != nil {
		WriteError(w, err)
	}
	WriteResponse(w, string(resultJson))
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	client.Play()
}

func PlayTrackHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, err)
	}

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
	devices, err := client.PlayerDevices()
	if err != nil {
		WriteError(w, err)
	}

	// Build list of devices
	var sDevices []SpotifyDevice
	for _, d := range devices {
		sd := SpotifyDevice{Name: d.Name, ID: string(d.ID)}
		sDevices = append(sDevices, sd)
	}

	jDevices, err := json.Marshal(sDevices)
	if err != nil {
		WriteError(w, err)
	}

	WriteResponse(w, string(jDevices))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	playerState, err := client.PlayerState()
	if err != nil {
		WriteError(w, err)
	}

	var status PlayerStatus
	status.Playing = playerState.CurrentlyPlaying.Playing

	if playerState != nil {
		status.DeviceID = string(playerState.Device.ID)
		status.DeviceName = playerState.Device.Name
	}

	status.URI = string(playerState.CurrentlyPlaying.PlaybackContext.URI)
	status.Progress = playerState.CurrentlyPlaying.Progress
	status.Track = playerState.CurrentlyPlaying.Item.SimpleTrack.Name
	status.Artist = playerState.CurrentlyPlaying.Item.SimpleTrack.Artists[0].Name

	jStatus, _ := json.Marshal(status)

	WriteResponse(w, string(jStatus))
}

func PlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	userPlaylists, err := client.CurrentUsersPlaylists()
	if err != nil {
		WriteError(w, err)
	}

	var playlists []Playlist
	for _, userPlaylist := range userPlaylists.Playlists {
		playlist := Playlist{Name: userPlaylist.Name, URI: string(userPlaylist.URI)}
		playlists = append(playlists, playlist)
	}
	resultsJson, err := json.Marshal(playlists)
	if err != nil {
		WriteError(w, err)
	}

	WriteResponse(w, string(resultsJson))
}
