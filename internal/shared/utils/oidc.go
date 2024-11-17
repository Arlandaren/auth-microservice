package utils

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	config2 "service/internal/shared/config"
)

func getAccessToken() string {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://your-oidc-provider.com")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	accData := config2.GetAccessData()

	config := &oauth2.Config{
		ClientID:     accData.ClientID,
		ClientSecret: accData.ClientSecret,
		RedirectURL:  accData.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "name"},
	}

	// Exchange authorization code for a token
	token, err := config.PasswordCredentialsToken(ctx, "user@example.com", "password")
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	return token.AccessToken
}
