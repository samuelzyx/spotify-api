package config

import (
	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config

var (
	clientID     = "6c84853b04074e9f8c8d2739c25a337a"
	clientSecret = "b6480532d600467797a17b0c16c16098"
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
