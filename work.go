// Package work defines code for using a gsheet as a datasource and destination
// for work related artifacts
package work

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

// NewClientOption returns a clientOption from a given set of credentials.
// Used to initialize Google API clients
func NewClientOption(ctx context.Context, r io.Reader, scopes []string) (option.ClientOption, error) {
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

// Config is the collection of settings that will direct artifact collection
type Config struct {
	ID           string
	Destinations Destinations `yaml:"destinations,omitempty"`
	Sources      []string     `yaml:"sources,omitempty"`
	Classifiers  Classifiers  `yaml:"classifiers,omitempty"`
}

// Destination is a place to write a report based on the criteria
type Destination struct {
	Label    string `yaml:"label,omitempty"`
	Sort     string `yaml:"sort,omitempty"`
	Summary  bool   `yaml:"summary,omitempty"`
	Criteria struct {
		Start   time.Time `yaml:"start,omitempty"`
		End     time.Time `yaml:"end,omitempty"`
		Project string    `yaml:"project,omitempty"`
	} `yaml:"criteria,omitempty"`
}

// Destinations is a collection of destination items
type Destinations []Destination
