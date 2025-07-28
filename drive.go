package work

import (
	"fmt"
	"strings"
	"time"

	"github.com/tpryan/work/artifact"
	"google.golang.org/api/drive/v2"
)

// MimeList is a collection of mimetypes
type MimeList []string

// String retunrs the query string needed by drive to filter files
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

// DriveFiles is a collection of files returned from a Google Drive query
type DriveFiles []*drive.File

// Artifacts returns a collection of artifacts from a collection of drive files
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

// DriveSearch  returns results from Google Drive as artifacts
func DriveSearch(q string, svc *drive.Service) (artifact.Artifacts, error) {

	files := DriveFiles{}
	var pageToken string

	for {
		r, err := svc.Files.List().PageToken(pageToken).Q(q).Do()
		if err != nil {
			return nil, fmt.Errorf("drive files list failed: %s", err)
		}

		files = append(files, r.Items...)

		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return files.Artifacts(), nil
}
