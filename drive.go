package work

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
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
func DriveSearch(q string, svc *drive.Service) (Artifacts, error) {

	token := "NOTEMPTY"
	files := DriveFiles{}

	for token != "" {

		if token == "NOTEMPTY" {
			token = ""
		}

		r, err := svc.Files.List().PageToken(token).Q(q).Do()

		if err != nil {
			return nil, fmt.Errorf("drive files list failed: %s", err)
		}

		if len(r.Items) > 0 {
			for _, i := range r.Items {
				files = append(files, i)
			}
		}

		token = r.NextPageToken

	}
	return files.Artifacts(), nil
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(tokFile string, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
