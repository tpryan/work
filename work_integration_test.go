package work

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := map[string]struct {
		path   string
		errStr string
	}{
		"basic": {
			path: credsTestPath,
		},
		"error": {
			path:   credsTestPath + "fail",
			errStr: "couldn't read the credential file",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			f, _ := os.Open(tc.path)

			_, err := NewConfig(ctx, f, []string{""})

			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && strings.Contains(err.Error(), tc.errStr) {
				t.Skip()
			}
		})
	}
}
