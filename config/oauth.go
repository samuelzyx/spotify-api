package config

import (
	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config

var (
	clientID     = "your_clientID"
	clientSecret = "your_clientSecret"
	redirectURI  = "http://localhost:8080/callback"
	scopes       = []string{"user-read-email", "user-read-private"}
)

func init() {
	OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
}
