package work

import (
	"fmt"
	"strings"
	"time"

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
func (d DriveFiles) Artifacts() Artifacts {

	arts := Artifacts{}

	for _, v := range d {

		shipped, err := time.Parse("2006-01-02T15:04:05.999Z", v.CreatedDate)
		if err != nil {
			shipped = time.Time{}
		}

		a := Artifact{
			Title:       v.Title,
			Link:        v.AlternateLink,
			ShippedDate: shipped,
			Role:        "Author",
		}

		if strings.Contains(a.Title, "Copy ") {
			continue
		}

		if strings.Contains(strings.ToLower(a.Title), strings.ToLower("prd")) ||
			strings.Contains(strings.ToLower(a.Title), strings.ToLower("tdd")) {
			a.Type = "Design Doc"
		}

		if a.Type == "" {
			switch v.MimeType {
			case "application/vnd.google-apps.spreadsheet":
				a.Type = "Sheet"
			case "application/vnd.google-apps.document":
				a.Type = "Doc"
			case "application/vnd.google-apps.presentation":
				a.Type = "Slides"
			default:
				a.Type = "File"
			}

		}

		arts = append(arts, a)
	}

	return arts
}
