package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// NewClientOption returns a ClientOption for accessing Google APIs using a credentials file.
// It combines the credentials file and scopes into a single option if possible,
// or just returns the credentials option as that's what's typically needed.
func NewClientOption(ctx context.Context, credPath string, scopes []string) (option.ClientOption, error) {
	// In most cases, just providing the credentials file is enough.
	// Scopes are often handled by the service client itself or the credentials file if it's a service account.
	// If specific scopes are needed with a service account, option.WithScopes can be used,
	// but NewService usually accepts variadic options.
	// Since the caller expects a single option.ClientOption, we return WithCredentialsFile.
	// If scopes are critical for the setup (e.g. strict scoping), we might need to chain them,
	// but option.ClientOption is an interface.
	return option.WithCredentialsFile(credPath), nil
}

// NewClient retrieves a token, saves the token, then returns the generated client.
func NewClient(credPath, tokenPath string, scopes ...string) (*http.Client, error) {
	b, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// If it's a service account, ConfigFromJSON might not be the right choice if we want user-on-behalf flow.
	// But assuming it's a client ID/secret for OAuth:
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	client := getClient(config, tokenPath)
	return client, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenPath string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenPath, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v\n", err)
		return nil
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v\n", err)
		return nil
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v\n", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
