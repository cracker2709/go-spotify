package main

import (
	"log"
	"os"

	"github.com/cracker2709/go-spotify/lib"
	"github.com/zmb3/spotify/v2"
)

func main() {
	lib.Debug = false // Enable debug logging
	var client *spotify.Client
	var userID string
	client, userID = lib.Authenticate()
	// tracks, err := lib.FetchTableData()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := lib.WriteTracksToFile(tracks, "tracks.txt"); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Tracks written to tracks.txt")

	// // Optionally read from file
	// _, err = lib.ReadFromFile("tracks.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	msg, err := lib.CreatePlaylistFromFile(client, userID, "tracks.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println(msg)
}
