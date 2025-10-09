package lib

import (
	"errors"
	"fmt"
	html_escape "html"
	"log"
	"net/http"
	"strings"

	net_html "golang.org/x/net/html"
)

// Add debug flag at package level
var Debug bool = false

// Add Album struct to represent album information
type Album struct {
	Title  string
	Year   string
	Tracks []Track
}

// Update Track struct to include duration
type Track struct {
	Title    string
	Artist   string
	Duration string
}

const (
	targetURL = "https://www.zenial.nl/html/variourf.htm"
)

// Define custom errors
var (
	ErrNoRows  = errors.New("no rows found in table")
	ErrNoTable = errors.New("no table found in HTML content")
)

// TableData represents structured data extracted from HTML table
type TableData struct {
	Cells []string
}

func FetchTableData() ([]Track, error) {
	log.Printf("Accessing %s\n", targetURL)
	debugLog("Starting fetch from %s", targetURL)

	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	tracks, err := parseTableContent(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table content: %w", err)
	}

	return tracks, nil
}

// Add debug logging helper
func debugLog(format string, args ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Update parseTableContent with debug logging
func parseTableContent(resp *http.Response) ([]Track, error) {
	debugLog("Starting HTML parsing")
	z := net_html.NewTokenizer(resp.Body)
	var allTracks []Track
	var depth int

	for {
		tokenType := z.Next()

		switch tokenType {
		case net_html.ErrorToken:
			if z.Err().Error() == "EOF" {
				debugLog("Reached EOF, found %d tracks", len(allTracks))
				if len(allTracks) > 0 {
					return allTracks, nil
				}
				return nil, ErrNoTable
			}
			debugLog("Error parsing HTML: %v", z.Err())
			return nil, fmt.Errorf("tokenizer error: %w", z.Err())

		case net_html.StartTagToken:
			token := z.Token()
			if token.Data == "table" {
				depth++
				debugLog("Found table at depth %d", depth)
				// Process only the inner tables that contain track listings
				if isInnerTable(token) {
					debugLog("Processing inner table at depth %d", depth)
					tracks, err := extractTracks(z)
					if err == nil && len(tracks) > 0 {
						debugLog("Found %d tracks in table", len(tracks))
						allTracks = append(allTracks, tracks...)
					}
				}
			}

		case net_html.EndTagToken:
			token := z.Token()
			if token.Data == "table" {
				debugLog("Ending table at depth %d", depth)
				depth--
				if depth < 0 {
					depth = 0 // Reset depth if it goes negative
				}
			}
		}
	}
}

func isInnerTable(token net_html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "cellspacing" && attr.Val == "0" {
			debugLog("Found inner table with cellspacing=0")
			return true
		}
	}
	return false
}

func extractTracks(z *net_html.Tokenizer) ([]Track, error) {
	var tracks []Track
	var currentTrack Track
	var inTrackRow bool
	var columnIndex int
	var depth = 1

	for depth > 0 {
		tokenType := z.Next()

		switch tokenType {
		case net_html.ErrorToken:
			debugLog("Error in extractTracks: %v", z.Err())
			if len(tracks) > 0 {
				return tracks, nil
			}
			return nil, z.Err()

		case net_html.StartTagToken:
			token := z.Token()
			switch token.Data {
			case "table":
				depth++
				debugLog("Table depth increased to %d", depth)
			case "tr":
				if !hasClass(token, "album") { // Skip album header row
					debugLog("Found new track row")
					columnIndex = 0
					inTrackRow = true
					currentTrack = Track{}
				}
			case "td":
				if !inTrackRow {
					continue
				}

				// Skip the rowspan spacer cell
				if hasAttr(token, "rowspan") {
					continue
				}

				columnIndex++
				debugLog("Processing column %d", columnIndex)

				if columnIndex == 4 {
					// Extract artist from <b> tag
					for {
						tt := z.Next()
						if tt == net_html.StartTagToken && z.Token().Data == "b" {
							if text, err := getNextText(z); err == nil {
								currentTrack.Artist = cleanText(text)
								debugLog("Found artist: %s", currentTrack.Artist)
							}
							break
						}
						if tt == net_html.ErrorToken {
							break
						}
					}
					continue
				}

				text, err := getNextText(z)
				if err != nil {
					debugLog("Error getting text: %v", err)
					continue
				}
				text = cleanText(text)
				debugLog("Found text: %q", text)

				switch columnIndex {
				case 1: // Title
					if strings.Contains(text, ".") {
						parts := strings.SplitN(text, ".", 2)
						if len(parts) == 2 {
							currentTrack.Title = cleanText(parts[1])
							debugLog("Found title: %s", currentTrack.Title)
						}
					}
				case 2: // Duration
					if strings.Contains(text, ":") {
						currentTrack.Duration = text
						debugLog("Found duration: %s", currentTrack.Duration)
					}
				}
			}

		case net_html.EndTagToken:
			token := z.Token()
			switch token.Data {
			case "table":
				depth--
				debugLog("Table depth decreased to %d", depth)
				if depth == 0 {
					debugLog("End of table, found %d tracks", len(tracks))
					if len(tracks) > 0 {
						return tracks, nil
					}
					return nil, ErrNoRows
				}
			case "tr":
				if inTrackRow && currentTrack.Title != "" && currentTrack.Artist != "" {
					debugLog("Adding track: %+v", currentTrack)
					tracks = append(tracks, currentTrack)
				}
				inTrackRow = false
			}
		}
	}

	return tracks, nil
}

// Add these helper functions
func hasClass(token net_html.Token, class string) bool {
	for _, attr := range token.Attr {
		if attr.Key == "class" && attr.Val == class {
			return true
		}
	}
	return false
}

func hasAttr(token net_html.Token, attr string) bool {
	for _, a := range token.Attr {
		if a.Key == attr {
			return true
		}
	}
	return false
}

func hasTrackNumber(z *net_html.Tokenizer) bool {
	for i := 0; i < 3; i++ { // Look ahead a few tokens
		if z.Next() == net_html.TextToken {
			text := strings.TrimSpace(z.Token().Data)
			if len(text) > 2 && text[1] == '.' && text[0] >= '0' && text[0] <= '9' {
				return true
			}
		}
	}
	return false
}

// getNextText retrieves the next text token content
func getNextText(z *net_html.Tokenizer) (string, error) {
	if z.Next() == net_html.TextToken {
		return z.Token().Data, nil
	}
	return "", errors.New("no text content found")
}

// Add this new helper function
func cleanText(text string) string {
	// Decode HTML entities (like &eacute;)
	decoded := html_escape.UnescapeString(text)
	// Trim whitespace
	cleaned := strings.TrimSpace(decoded)
	return cleaned
}
