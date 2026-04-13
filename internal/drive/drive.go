package drive

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tpryan/work/internal/artifact"
	"google.golang.org/api/drive/v2"
)

// MimeList represents a collection of MIME types.
type MimeList []string

// String returns the query string needed by Google Drive to filter files.
func (m MimeList) String() string {
	result := strings.Builder{}

	for i, v := range m {
		if i != 0 {
			result.WriteString(" or ")
		}
		result.WriteString(fmt.Sprintf("mimeType='%s'", v))
	}

	return result.String()
}

// DriveFiles represents a collection of files returned from a Google Drive query.
type DriveFiles []*drive.File

// Artifacts returns a collection of artifacts from a collection of Drive files.
func (d DriveFiles) Artifacts() artifact.Artifacts {

	arts := artifact.Artifacts{}

	for _, v := range d {

		shipped, err := time.Parse("2006-01-02T15:04:05.999Z", v.CreatedDate)
		if err != nil {
			shipped = time.Time{}
		}

		a := artifact.Artifact{
			Title:       v.Title,
			Link:        v.AlternateLink,
			ShippedDate: shipped,
			Role:        "Author",
		}

		// TODO: do at a higher level - now built into
		// if strings.Contains(a.Title, "Copy ") {
		// 	continue
		// }

		if strings.Contains(strings.ToLower(a.Title), strings.ToLower("prd")) ||
			strings.Contains(strings.ToLower(a.Title), strings.ToLower("tdd")) {
			a.Type = "Design Doc"
		}

		typeMap := map[string]string{
			"application/vnd.google-apps.spreadsheet":  "Sheet",
			"application/vnd.google-apps.document":     "Doc",
			"application/vnd.google-apps.presentation": "Slides",
			"application/vnd.google.colaboratory.corp": "Colab",
			"application/vnd.google-apps.form":         "Form",
		}

		if a.Type == "" {
			a.Type = "File"
			if t, ok := typeMap[v.MimeType]; ok {
				a.Type = t
			}
		}

		arts = append(arts, a)
	}

	return arts
}

// Search returns results from Google Drive as artifacts.
func Search(ctx context.Context, q string, svc *drive.Service) (artifact.Artifacts, error) {

	files := DriveFiles{}
	var pageToken string

	for {
		r, err := svc.Files.List().PageToken(pageToken).Q(q).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("drive files list failed: %w", err)
		}

		files = append(files, r.Items...)

		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return files.Artifacts(), nil
}

// Source represents a Google Drive artifact source.
type Source struct {
	SVC  *drive.Service
	User string
}

// Fetch returns results from Google Drive as artifacts.
func (s Source) Fetch(ctx context.Context) (artifact.Artifacts, error) {
	mlist := MimeList{
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.google-apps.form",
		"application/vnd.google-apps.presentation",
		"application/vnd.google.colaboratory.corp",
	}

	query := fmt.Sprintf("'%s@google.com' in owners and (%s)", s.User, mlist.String())

	arts, err := Search(ctx, query, s.SVC)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data from drive: %w", err)
	}

	arts.Sort()
	return arts, nil
}
