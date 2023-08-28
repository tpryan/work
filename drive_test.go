package work

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v2"
)

func TestMimeListString(t *testing.T) {
	tests := map[string]struct {
		in   MimeList
		want string
	}{
		"single": {
			in:   MimeList{"application/vnd.google-apps.document"},
			want: "mimeType='application/vnd.google-apps.document'",
		},
		"couple": {
			in: MimeList{
				"application/vnd.google-apps.document",
				"application/vnd.google-apps.spreadsheet",
			},
			want: "mimeType='application/vnd.google-apps.document' or mimeType='application/vnd.google-apps.spreadsheet'",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDriveArtifacts(t *testing.T) {
	tests := map[string]struct {
		in   DriveFiles
		want Artifacts
	}{
		"sheet": {
			in: DriveFiles{
				&drive.File{
					Title:         "title",
					AlternateLink: "https://example.com",
					CreatedDate:   "2023-08-24T09:00:00.000Z",
					MimeType:      "application/vnd.google-apps.spreadsheet",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title",
					Link:        "https://example.com",
					Type:        "Sheet",
					Role:        "Author",
					ShippedDate: time.Date(2023, 8, 24, 9, 0, 0, 0, time.UTC),
				},
			},
		},
		"doc": {
			in: DriveFiles{
				&drive.File{
					Title:         "title",
					AlternateLink: "https://example.com",
					CreatedDate:   "2023-08-24T09:00:00.000Z",
					MimeType:      "application/vnd.google-apps.document",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title",
					Link:        "https://example.com",
					Type:        "Doc",
					Role:        "Author",
					ShippedDate: time.Date(2023, 8, 24, 9, 0, 0, 0, time.UTC),
				},
			},
		},
		"slides": {
			in: DriveFiles{
				&drive.File{
					Title:         "title",
					AlternateLink: "https://example.com",
					CreatedDate:   "2023-08-24T09:00:00.000Z",
					MimeType:      "application/vnd.google-apps.presentation",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title",
					Link:        "https://example.com",
					Type:        "Slides",
					Role:        "Author",
					ShippedDate: time.Date(2023, 8, 24, 9, 0, 0, 0, time.UTC),
				},
			},
		},
		"file": {
			in: DriveFiles{
				&drive.File{
					Title:         "title",
					AlternateLink: "https://example.com",
					CreatedDate:   "2023-08-24T09:00:00.000Z",
					MimeType:      "application/vnd.adobe.pdf",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title",
					Link:        "https://example.com",
					Type:        "File",
					Role:        "Author",
					ShippedDate: time.Date(2023, 8, 24, 9, 0, 0, 0, time.UTC),
				},
			},
		},
		"badtime": {
			in: DriveFiles{
				&drive.File{
					Title:         "title",
					AlternateLink: "https://example.com",
					CreatedDate:   "BADTIMEFORMAT",
					MimeType:      "application/vnd.google-apps.spreadsheet",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title",
					Link:        "https://example.com",
					Type:        "Sheet",
					Role:        "Author",
					ShippedDate: time.Time{},
				},
			},
		},

		"PRD": {
			in: DriveFiles{
				&drive.File{
					Title:         "title prd",
					AlternateLink: "https://example.com",
					CreatedDate:   "BADTIMEFORMAT",
					MimeType:      "application/vnd.google-apps.document",
				},
			},
			want: Artifacts{
				Artifact{
					Title:       "title prd",
					Link:        "https://example.com",
					Type:        "Design Doc",
					Role:        "Author",
					ShippedDate: time.Time{},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Artifacts()
			assert.Equal(t, tc.want, got)
		})
	}
}
