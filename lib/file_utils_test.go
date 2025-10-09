package lib

import (
	"os"
	"testing"
)

func TestFileUtils_WriteTracksToFile(t *testing.T) {
	tracks := []Track{
		{Artist: "Artist1", Title: "Title1"},
		{Artist: "Artist2", Title: "Title2"},
	}
	filename := "test_tracks.txt"

	// Test writing to file
	if err := WriteTracksToFile(tracks, filename); err != nil {
		t.Errorf("WriteTracksToFile() error = %v", err)
		return
	}
	defer os.Remove(filename) // Clean up

	// Verify file content
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	expectedContent := "Artist1 - Title1\nArtist2 - Title2\n"
	if string(fileContent) != expectedContent {
		t.Errorf("File content mismatch. Expected:\n%sGot:\n%s", expectedContent, string(fileContent))
	}
}

func TestFileUtils_ReadFromFile(t *testing.T) {
	// Create a temporary file with test data
	content := "Artist1 - Title1\nArtist2 - Title2\n"
	tmpfile, err := os.CreateTemp("", "tracks*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test reading the file
	tracks, err := ReadFromFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
		return
	}

	// Verify the results
	expected := []Track{
		{Artist: "Artist1", Title: "Title1"},
		{Artist: "Artist2", Title: "Title2"},
	}

	if len(tracks) != len(expected) {
		t.Errorf("Expected %d tracks, got %d", len(expected), len(tracks))
	}

	for i, track := range tracks {
		if track != expected[i] {
			t.Errorf("Track %d: expected %v, got %v", i, expected[i], track)
		}
	}
}
