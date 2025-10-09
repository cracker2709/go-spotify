package lib

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
)

func CreatePlaylistFromFile(client *spotify.Client, userID string, filename string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("spotify client is nil")
	}

	fmt.Println("Enter playlist name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	fmt.Println("Enter playlist description: ")
	scanner.Scan()
	description := scanner.Text()

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	// Validate input
	if name == "" {
		return "", fmt.Errorf("playlist name cannot be empty")
	}

	playlist, err := client.CreatePlaylistForUser(context.Background(), userID, name, description, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to create playlist: %w", err)
	}

	log.Printf("Created playlist: %s (%s)", playlist.Name, playlist.ID)

	// Read tracks from file
	tracks, err := ReadFromFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read tracks from file: %w", err)
	}

	// Collect track IDs
	var trackIDs []spotify.ID
	for _, track := range tracks {
		searchQuery := fmt.Sprintf("track:%s artist:%s", track.Title, track.Artist)
		results, err := client.Search(context.Background(), searchQuery, spotify.SearchTypeTrack)
		if err != nil {
			log.Printf("Error searching for track %s by %s: %v", track.Title, track.Artist, err)
			continue
		}
		if results.Tracks.Total > 0 {
			trackIDs = append(trackIDs, results.Tracks.Tracks[0].ID)
			log.Printf("Found track: %s by %s", track.Title, track.Artist)
		} else {
			log.Printf("Track not found: %s by %s", track.Title, track.Artist)
		}
	}

	if len(trackIDs) == 0 {
		return "", fmt.Errorf("no tracks found to add to the playlist")
	}
	// Add tracks to the playlist in batches of 100
	for i := 0; i < len(trackIDs); i += 100 {
		end := i + 100
		if end > len(trackIDs) {
			end = len(trackIDs)
		}
		_, err := client.AddTracksToPlaylist(context.Background(), playlist.ID, trackIDs[i:end]...)
		if err != nil {
			return "", fmt.Errorf("failed to add tracks to playlist: %w", err)
		}
		log.Printf("Added %d tracks to playlist", end-i)
	}
	msg := fmt.Sprintf("Playlist '%s' created with %d tracks", playlist.Name, len(trackIDs))
	return msg, nil
}
