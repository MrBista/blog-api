package utils

import (
	"github.com/MrBista/blog-api/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     config.AppConfig.AppMain.GetGoogleClientId(),
		ClientSecret: config.AppConfig.AppMain.GetGoogleClientSecret(),
		RedirectURL:  config.AppConfig.AppMain.GetGoogleRedirctUrl(), // http://localhost:3000/auth/callback
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
