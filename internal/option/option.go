package option

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

// New returns a clientOption from a given set of credentials.
// Used to initialize Google API clients
func New(ctx context.Context, r io.Reader, scopes []string) (option.ClientOption, error) {
	if r == nil {
		return nil, fmt.Errorf("reader cannot be nil")
	}

	creds := struct {
		ClientEmail  string `json:"client_email"`
		PrivateKey   string `json:"private_key"`
		PrivateKeyID string `json:"private_key_id"`
		TokenURL     string `json:"token_uri"`
	}{}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("could not read credentials: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &creds); err != nil {
		return nil, fmt.Errorf("error unmarshaling credentials file: %w", err)
	}

	conf := &jwt.Config{
		Email:        creds.ClientEmail,
		PrivateKey:   []byte(creds.PrivateKey),
		PrivateKeyID: creds.PrivateKeyID,
		TokenURL:     creds.TokenURL,
		Scopes:       scopes,
	}

	client := option.WithHTTPClient(conf.Client(ctx))

	return client, nil
}
