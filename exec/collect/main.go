package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/tpryan/googleclient"
	"github.com/tpryan/work"
	"github.com/tpryan/work/artifact"
	"github.com/tpryan/work/drive"
	"github.com/tpryan/work/github"
	"github.com/tpryan/work/gsheet"
	gdrive "google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var credPath = "../credentials/credentials.json"
var driveCredsPath = "../credentials/drive_credentials.json"
var tokenPath = "../credentials/token.json"
var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/spreadsheets",
}

func main() {
	var userFlag = flag.String("user", "", "user who should be run on")
	flag.Parse()

	user := *userFlag
	if user == "" {
		user = os.Getenv("USER")
	}

	var configPath = fmt.Sprintf("../users/%s.yaml", user)

	ctx := context.Background()
	log.Infof("Starting process for: %s...", user)

	config, err := work.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	log.Infof("Reading Credential files")

	options, err := googleclient.NewClientOption(ctx, credPath, scopes)
	if err != nil {
		log.Fatalf("error while opening credentials: %s", err)
	}

	log.Infof("Initializing clients")

	client, err := googleclient.NewClient(driveCredsPath, tokenPath, scopes...)
	if err != nil {
		log.Fatalf("Unable to get google client: %v", err)
	}

	driveSVC, err := gdrive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retrieve Drive client: %v", err)
	}

	sheetsSVC, err := sheets.NewService(ctx, options)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets client: %v", err)
	}

	gsheet := gsheet.NewGSheet(*sheetsSVC, config.SpreadSheetID)

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

func processDrive(svc *gdrive.Service, gsheet gsheet.GSheet, user string) error {

	mlist := drive.MimeList{
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.google-apps.form",
		"application/vnd.google-apps.presentation",
		"application/vnd.google.colaboratory.corp",
	}

	query := fmt.Sprintf("'%s@google.com' in owners and (%s)", user, mlist.String())

	arts, err := drive.DriveSearch(query, svc)
	if err != nil {
		return fmt.Errorf("error retrieving data from drive: %w", err)
	}

	arts.Sort()

	if err := gsheet.ToSheet("Source - DriveFiles", arts); err != nil {
		return fmt.Errorf("error writing to sheet: %w", err)
	}

	return nil
}

func processGithub(username string, gsheet gsheet.GSheet) error {
	q := fmt.Sprintf("author:%s is:pr state:closed", username)

	gartifacts, err := github.GHSearch(q)
	if err != nil {
		return fmt.Errorf("could not get issues: %w", err)
	}

	if err := gsheet.ToSheet("Source - Github", gartifacts); err != nil {
		return fmt.Errorf("error writing to sheet: %w", err)
	}

	return nil
}

func writeReport(gsheet gsheet.GSheet, sources []string, destinations work.Destinations, list artifact.Classifiers) error {
	all := artifact.Artifacts{}

	log.Infof("Getting Sources")
	for _, source := range sources {
		arts, err := gsheet.Artifacts(source)

		if err != nil {
			return fmt.Errorf("unable to retrieve sheets client: %w", err)
		}
		all = append(all, arts...)

	}

	var wg sync.WaitGroup
	wg.Add(len(destinations))

	log.Infof("Writing to Sheet")
	for _, dest := range destinations {

		go func(all artifact.Artifacts, dest work.Destination) {
			artifacts := all.Copy()

			artifacts.Massage(
				artifact.Between(dest.Criteria.Start, dest.Criteria.End),
				artifact.Classify(list),
				artifact.ProjectFilter(dest.Criteria.Project),
				artifact.Unique(),
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
