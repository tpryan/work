package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/tpryan/work"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/sheets/v4"
)

var configPath = "tpryan.yaml"
var credPath = "credentials.json"
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
	driveSVC, err := drive.NewService(ctx, options)
	if err != nil {
		log.Fatalf("unable to retrieve Drive client: %v", err)
	}

	sheetsSVC, err := sheets.NewService(ctx, options)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets client: %v", err)
	}

	gsheet := work.NewGSheet(*sheetsSVC, config.SpreadSheetID)

	// log.Infof("Processing Github")
	// if err := processGithub(gsheet); err != nil {
	// 	log.Error("unable to retrieve latest github info: %s", err)
	// }

	// log.Infof("Processing Drive")
	// if err := processDrive(driveSVC, gsheet); err != nil {
	// 	log.Errorf("unable to retrieve latest drive info: %s", err)
	// }

	log.Infof("Writing report")
	if err := writeReport(gsheet, config.Sources, config.Destinations, config.Classifiers); err != nil {
		log.Error(fmt.Sprintf("unable to write report to sheets: %s", err))
	}
	log.Infof("...Finished")

	_ = config
	_ = ctx
	_ = options
	_ = driveSVC
	_ = sheetsSVC
	_ = gsheet

}

func processDrive(svc *drive.Service, gsheet work.GSheet) error {

	token := "NOTEMPTY"
	files := []*drive.File{}

	mlist := work.MimeList{
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.apple.keynote",
		"application/vnd.google-apps.presentation",
		"application/vnd.google.colaboratory.corp",
	}

	for token != "" {

		if token == "NOTEMPTY" {
			token = ""
		}

		r, err := svc.Files.List().PageToken(token).
			// Corpora("user").
			Q(mlist.String()).
			Do()

		if err != nil {
			return fmt.Errorf("drive files list failed: %s", err)
		}

		if len(r.Items) > 0 {
			for _, i := range r.Items {
				files = append(files, i)
			}
		}

		token = r.NextPageToken

	}

	arts, err := work.DriveSearch(mlist.String(), svc)
	if err != nil {
		return fmt.Errorf("error retrieving data from drive: %w", err)
	}

	arts.Sort()

	if err := gsheet.ToSheet("DriveFiles", arts); err != nil {
		return fmt.Errorf("error writing to sheet: %w", err)
	}

	return nil
}

func processGithub(gsheet work.GSheet) error {
	q := "author:tpryan is:pr state:closed"

	gartifacts, err := work.GHSearch(q)
	if err != nil {
		return fmt.Errorf("could not get issues: %w", err)
	}

	if err := gsheet.ToSheet("Github", gartifacts); err != nil {
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
