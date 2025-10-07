package lib


import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"github.com/zmb3/spotify/v2"
)

func generateRandomState() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%d", rand.Intn(100000)) 
}

var (
	redirectURI = "http://localhost:8080/callback"
	html        = `<html><body>Login successful! You can close this window.</body></html>`
	ch    = make(chan *spotify.Client)
	state = generateRandomState()
)


func Authenticate() {
	// Check for environment variables
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables")
	}
	
	// Create authenticator with proper credentials
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
	)
	
	// Set up HTTP handlers
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		completeAuth(w, r, auth)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	
	// Start HTTP server in goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	// Generate auth URL and prompt user
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	
	// Wait for auth to complete
	client := <-ch
	
	// Test the authenticated client
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal("Failed to get current user:", err)
	}
	fmt.Printf("Successfully logged in as: %s (%s)\n", user.DisplayName, user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request, auth *spotifyauth.Authenticator) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Printf("Token error: %v", err)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, state)
		return
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, html)
	ch <- client
}
