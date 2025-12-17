# Go Spotify

A Go application that provides Spotify Web API authentication using OAuth2 flow. This project demonstrates how to authenticate with Spotify's API and retrieve user information.

## Features

- OAuth2 authentication flow with Spotify Web API
- Local HTTP server for handling authentication callbacks
- Environment variable configuration for credentials
- User profile retrieval after successful authentication
- Support for multiple Spotify API scopes:
  - User read currently playing
  - User read playback state
  - User modify playback state

## Prerequisites

- Go 1.25.1 or later
- Spotify Developer Account
- Spotify App with registered redirect URI

## Setup

### 1. Create Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Click "Create App"
3. Fill in the app details
4. Add `http://127.0.0.1:8080/callback` as a Redirect URI
5. Note down your Client ID and Client Secret

### 2. Environment Variables

Set the required environment variables:

```bash
export SPOTIFY_CLIENT_ID="your_client_id_here"
export SPOTIFY_CLIENT_SECRET="your_client_secret_here"
```

Or create a `.env` file:

```env
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/cracker2709/go-spotify.git
   cd go-spotify
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o go-spotify
   ```

## Usage

1. Make sure your environment variables are set
2. Run the application:
   ```bash
   ./go-spotify
   ```
   Or directly with Go:
   ```bash
   go run main.go
   ```

3. The application will:
   - Start a local HTTP server on port 8080
   - Display an authorization URL
   - Open the URL in your browser to log in to Spotify
   - Handle the callback and complete authentication
   - Display your Spotify user information

### Example Output

```
Starting HTTP server on :8080
Please log in to Spotify by visiting the following page in your browser: https://accounts.spotify.com/authorize?...
Successfully logged in as: John Doe (johndoe123)
```

## Project Structure

```
go-spotify/
├── main.go              # Entry point
├── lib/
│   └── authenticator.go # Spotify authentication logic
├── go.mod              # Go module file
├── go.sum              # Go dependencies checksum
└── README.md           # This file
```

## Dependencies

- [github.com/zmb3/spotify/v2](https://github.com/zmb3/spotify) - Spotify Web API wrapper for Go

## API Scopes

The application requests the following Spotify API scopes:

- `user-read-currently-playing` - Read access to a user's currently playing content
- `user-read-playback-state` - Read access to a user's player state
- `user-modify-playback-state` - Write access to a user's playback state

## Security Notes

- Never commit your Client ID and Client Secret to version control
- The redirect URI must match exactly what's configured in your Spotify app
- The application generates a random state parameter for OAuth2 security

## Troubleshooting

### Common Issues

1. **"Please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables"**
   - Make sure your environment variables are set correctly
   - Try printing them to verify: `echo $SPOTIFY_CLIENT_ID`

2. **"State mismatch" error**
   - This is a security feature. Try restarting the application

3. **"Couldn't get token" error**
   - Verify your Client ID and Client Secret are correct
   - Ensure the redirect URI in your Spotify app matches `http://localhost:8080/callback`

4. **Port 8080 already in use**
   - Stop any other services running on port 8080
   - Or modify the port in the code if needed

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- [Spotify Web API](https://developer.spotify.com/documentation/web-api/) for providing the API
- [zmb3/spotify](https://github.com/zmb3/spotify) for the excellent Go client library
