package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/tpryan/work"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var user = "yufengg"
var configPath = fmt.Sprintf("users/%s.yaml", user)
var credPath = "credentials/credentials.json"
var tokenFile = "credentials/token.json"
var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/spreadsheets",
}

func main() {
	ctx := context.Background()
	log.Infof("Starting...")

	log.Infof("Reading Config files")

	config, err := work.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	log.Infof("Reading Credential files")
	f, err := os.Open(credPath)
	if err != nil {
		log.Fatalf("error while opening credentials: %s", err)
	}

	options, err := work.NewClientOption(ctx, f, scopes)
	if err != nil {
		log.Fatalf("error while opening credentials: %s", err)
	}

	log.Infof("Initializing clients")

	b, err := os.ReadFile("credentials/drive_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	driveconfig, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := work.GetClient("credentials/token.json", driveconfig)

	driveSVC, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retrieve Drive client: %v", err)
	}

	sheetsSVC, err := sheets.NewService(ctx, options)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets client: %v", err)
	}

	gsheet := work.NewGSheet(*sheetsSVC, config.SpreadSheetID)

	log.Infof("Processing Github")
	if err := processGithub(config.GithubUser, gsheet); err != nil {
		log.Error("unable to retrieve latest github info: %s", err)
	}

	if config.QueryDrive {
		log.Infof("Processing Drive")
		if err := processDrive(driveSVC, gsheet, user); err != nil {
			log.Errorf("unable to retrieve latest drive info: %s", err)
		}
	}

	log.Infof("Writing report")
	if err := writeReport(gsheet, config.Sources, config.Destinations, config.Classifiers); err != nil {
		log.Error(fmt.Sprintf("unable to write report to sheets: %s", err))
	}
	log.Infof("...Finished")

}

func processDrive(svc *drive.Service, gsheet work.GSheet, user string) error {

	mlist := work.MimeList{
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.apple.keynote",
		"application/vnd.google-apps.presentation",
		"application/vnd.google.colaboratory.corp",
	}

	query := fmt.Sprintf("'%s@google.com' in owners and (%s)", user, mlist.String())

	// for token != "" {

	// 	if token == "NOTEMPTY" {
	// 		token = ""
	// 	}

	// 	r, err := svc.Files.List().PageToken(token).
	// 		Corpora("user").
	// 		Q(query).
	// 		Do()

	// 	if err != nil {
	// 		return fmt.Errorf("drive files list failed: %s", err)
	// 	}

	// 	if len(r.Items) > 0 {
	// 		for _, i := range r.Items {
	// 			files = append(files, i)
	// 		}
	// 	}

	// 	token = r.NextPageToken

	// }

	arts, err := work.DriveSearch(query, svc)
	if err != nil {
		return fmt.Errorf("error retrieving data from drive: %w", err)
	}

	arts.Sort()

	if err := gsheet.ToSheet("Source - DriveFiles", arts); err != nil {
		return fmt.Errorf("error writing to sheet: %w", err)
	}

	return nil
}

func processGithub(username string, gsheet work.GSheet) error {
	q := fmt.Sprintf("author:%s is:pr state:closed", username)

	gartifacts, err := work.GHSearch(q)
	if err != nil {
		return fmt.Errorf("could not get issues: %w", err)
	}

	if err := gsheet.ToSheet("Source - Github", gartifacts); err != nil {
		return fmt.Errorf("error writing to sheet: %w", err)
	}

	return nil
}

func writeReport(gsheet work.GSheet, sources []string, destinations work.Destinations, list work.Classifiers) error {
	all := work.Artifacts{}

	log.Infof("Getting Sources")
	for _, source := range sources {
		arts, err := gsheet.Artifacts(source)

		if err != nil {
			return fmt.Errorf("Unable to retrieve Sheets client: %w", err)
		}
		all = append(all, arts...)

	}

	var wg sync.WaitGroup
	wg.Add(len(destinations))

	log.Infof("Writing to Sheet")
	for _, dest := range destinations {

		go func(all work.Artifacts, dest work.Destination) {
			artifacts := all.Copy()

			artifacts.Massage(
				work.Between(dest.Criteria.Start, dest.Criteria.End),
				work.Classify(list),
				work.ProjectFilter(dest.Criteria.Project),
				work.Unique(),
			)

			switch dest.Sort {
			case "report":
				artifacts.SortReport()
			default:
				artifacts.Sort()
			}

			log.Infof("Writing to %s", dest.Sheet)
			if err := gsheet.ToSheet(dest.Sheet, artifacts); err != nil {
				log.Errorf("error writing to sheet %s: %s", dest.Sheet, err)
			}

			if dest.Summary {
				artifacts.Template(dest.Sheet)
			}

			wg.Done()
		}(all, dest)

	}

	wg.Wait()
	return nil
}
