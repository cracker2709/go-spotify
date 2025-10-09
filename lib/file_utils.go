package lib

import (
	"fmt"
	"os"
	"strings"
)

// WriteTracksToFile writes the tracks to a file in a formatted way
func WriteTracksToFile(tracks []Track, filename string) error {
	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write each track to the file
	for _, track := range tracks {
		// Format: Title - Artist (Duration)
		line := fmt.Sprintf("%s - %s\n", track.Artist, track.Title)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write track: %w", err)
		}
	}

	return nil
}

func ReadFromFile(filename string) ([]Track, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tracks []Track
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " - ")
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		track := Track{
			Artist: strings.TrimSpace(parts[0]),
			Title:  strings.TrimSpace(parts[1]),
		}
		tracks = append(tracks, track)
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no valid tracks found in file")
	}

	return tracks, nil
}
