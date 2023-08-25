// Package work defines code for using a gsheet as a datasource and destination
// for work related artifacts
package work

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

// NewConfig returns a clientOption from a given set of credentials.
// Used to initialize Google API clients
func NewConfig(ctx context.Context, r io.Reader, scopes []string) (option.ClientOption, error) {
	m := make(map[string]string)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	dat := buf.Bytes()

	if err := json.Unmarshal(dat, &m); err != nil {
		return nil, fmt.Errorf("error parsing credentials file: %s", err)
	}

	conf := &jwt.Config{
		Email:        m["client_email"],
		PrivateKey:   []byte(m["private_key"]),
		PrivateKeyID: m["private_key_id"],
		TokenURL:     m["token_uri"],
		Scopes:       scopes,
	}

	client := option.WithHTTPClient(conf.Client(ctx))

	return client, nil
}
