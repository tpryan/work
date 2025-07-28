package option

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// errReader is a reader that always returns an error.
type errReader struct{}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced reader error")
}

func TestNew(t *testing.T) {
	validCreds := `{
		"client_email": "test@example.com",
		"private_key": "-----BEGIN PRIVATE KEY-----\n...somekey...\n-----END PRIVATE KEY-----\n",
		"private_key_id": "12345",
		"token_uri": "https://oauth2.googleapis.com/token"
	}`

	tests := map[string]struct {
		reader io.Reader
		scopes []string
		errStr string
	}{
		"success": {
			reader: strings.NewReader(validCreds),
			scopes: []string{"scope1", "scope2"},
		},
		"reader error": {
			reader: errReader{},
			errStr: "could not read credentials",
		},
		"unmarshal error": {
			reader: strings.NewReader(`{ "bad": "json" `),
			errStr: "error unmarshaling credentials file",
		},
		"nil reader panics": {
			reader: nil,
			errStr: "reader cannot be nil",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			clientOption, err := New(ctx, tc.reader, tc.scopes)

			if tc.errStr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errStr)
				assert.Nil(t, clientOption)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, clientOption)
			}
		})
	}
}
