package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
)

func main() {
	client := newClient()

	songSearch := os.Args[1]
	fmt.Printf("Searching for %s\n", songSearch)
	results, err := spotify.Search(songSearch, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	track := results.Tracks.Tracks[0]
	fmt.Printf("Playing %s by %s\n", track.SimpleTrack.Name, track.SimpleTrack.Artists[0].Name)

	pOpts := spotify.PlayOptions{URIs: []spotify.URI{track.URI}}
	_ = client.PlayOpt(&pOpts)
}
