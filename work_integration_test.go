package work

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tpryan/work/artifact"
	"github.com/tpryan/work/option"
)

func TestNewClientOption(t *testing.T) {
	tests := map[string]struct {
		path   string
		errStr string
	}{
		"basic": {
			path: "testdata/test-creds.json",
		},
		"error": {
			path:   "testdata/test-creds.json" + "fail",
			errStr: "couldn't read the credential file",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			f, _ := os.Open(tc.path)

			_, err := option.New(ctx, f, []string{""})

			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && strings.Contains(err.Error(), tc.errStr) {
				t.Skip()
			}
		})
	}
}

func TestBasic(t *testing.T) {
	tests := map[string]struct {
		in     string
		want   *Config
		errStr string
	}{
		"blank": {
			in:   "testdata/blank.yaml",
			want: &Config{},
		},
		"notExists": {
			in:     "testdata/doesnotexist.yaml",
			errStr: "no such file or directory",
		},
		"garbage": {
			in:     "testdata/garbage.yaml",
			errStr: "couldn't parse the config file",
		},
		"basic": {
			in: "testdata/basic.yaml",
			want: &Config{
				SpreadSheetID: "123456789",
				Sources:       []string{"Critique", "Buganizer"},
				Destinations: Destinations{
					Destination{
						Sheet: "Test",
						Sort:  "default",
						Criteria: Criteria{
							Project: "test",
							Start:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
							End:     time.Date(2023, 8, 28, 0, 0, 0, 0, time.UTC),
						},
					},
				},

				Classifiers: artifact.Classifiers{
					Lists: []artifact.Classifier{
						{
							Project:    "Example",
							Subproject: "Something",
							Links: []string{
								"http://example.com",
							},
						},
						{
							Project:    "Another",
							Subproject: "Something specific",
							Contains: map[string][]string{
								"title": {"example"},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfig(tc.in)

			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}
			if tc.errStr != "" && strings.Contains(err.Error(), tc.errStr) {
				t.Skip()
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
