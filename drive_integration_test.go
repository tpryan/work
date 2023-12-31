package work

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v2"
)

func getTestDriveSvc() (*drive.Service, error) {
	ctx := context.Background()

	f, err := os.Open(credsTestPath)
	if err != nil {
		return nil, err
	}

	config, err := NewClientOption(ctx, f, []string{"https://www.googleapis.com/auth/drive"})
	if err != nil {
		return nil, err
	}

	driveSVC, err := drive.NewService(ctx, config)
	if err != nil {
		return nil, err
	}

	return driveSVC, nil

}

func TestDriveSearch(t *testing.T) {
	tests := map[string]struct {
		q      string
		want   Artifacts
		errStr string
	}{
		"basic": {

			q: "title contains 'Deploystack Performance Metrics' AND mimeType='application/vnd.google-apps.spreadsheet'",
			want: Artifacts{
				Artifact{
					Title:       "Deploystack Performance Metrics",
					Link:        "https://docs.google.com/spreadsheets/d/1UqE9jEZA2G0kSwAducflfi9B7qjb9iMuDZcchAoxWzM/edit?usp=drivesdk",
					Type:        "Sheet",
					Project:     "",
					Subproject:  "",
					Role:        "Author",
					ShippedDate: time.Date(2022, time.October, 13, 23, 49, 57, 24000000, time.UTC),
					Extra:       "",
				},
			},
		},
		"error": {

			q:      "title contains 'Deploystack Performance Metrics",
			want:   Artifacts{},
			errStr: "Invalid query, invalid",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			svc, err := getTestDriveSvc()
			if err != nil {
				t.Fatalf(fmt.Sprintf("Unable to retrieve Drive client: %v", err))
			}

			got, err := DriveSearch(tc.q, svc)

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
