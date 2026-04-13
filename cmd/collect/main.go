package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/tpryan/work"
	"github.com/tpryan/work/internal/artifact"
	"github.com/tpryan/work/internal/drive"
	"github.com/tpryan/work/internal/github"
	"github.com/tpryan/work/internal/googleauth"
	"github.com/tpryan/work/internal/gsheet"
	gdrive "google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var credPath = "../../credentials/credentials.json"
var driveCredsPath = "../../credentials/drive_credentials.json"
var tokenPath = "../../credentials/token.json"
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

	var configPath = fmt.Sprintf("../../users/%s.yaml", user)

	ctx := context.Background()
	log.Infof("Starting process for: %s...", user)

	config, err := work.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	log.Infof("Reading Credential files")

	options, err := googleauth.NewClientOption(ctx, credPath, scopes)
	if err != nil {
		log.Fatalf("error while opening credentials: %s", err)
	}

	log.Infof("Initializing clients")

	client, err := googleauth.NewClient(driveCredsPath, tokenPath, scopes...)
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

	gsheet := gsheet.New(*sheetsSVC, config.SpreadSheetID)

	var sources []struct {
		name   string
		source artifact.Source
	}

	sources = append(sources, struct {
		name   string
		source artifact.Source
	}{"Source - Github", github.Source{Username: config.GithubUser}})

	if config.QueryDrive {
		sources = append(sources, struct {
			name   string
			source artifact.Source
		}{"Source - DriveFiles", drive.Source{SVC: driveSVC, User: user}})
	}

	for _, s := range sources {
		log.Infof("Processing %s", s.name)
		arts, err := s.source.Fetch(ctx)
		if err != nil {
			log.Errorf("unable to retrieve %s info: %s", s.name, err)
			continue
		}

		if err := gsheet.ToSheet(ctx, s.name, arts); err != nil {
			log.Errorf("error writing %s to sheet: %s", s.name, err)
		}
	}

	log.Infof("Writing report")
	if err := writeReport(ctx, gsheet, config.Sources, config.Destinations, config.Classifiers); err != nil {
		log.Error(fmt.Sprintf("unable to write report to sheets: %s", err))
	}
	log.Infof("...Finished")

}

func writeReport(ctx context.Context, gsheet gsheet.GSheet, sources []string, destinations work.Destinations, list artifact.Classifiers) error {
	all := artifact.Artifacts{}

	log.Infof("Getting Sources")
	for _, source := range sources {
		arts, err := gsheet.Artifacts(ctx, source)

		if err != nil {
			return fmt.Errorf("unable to retrieve sheets client: %w", err)
		}
		all = append(all, arts...)

	}

	var wg sync.WaitGroup
	// Limit concurrency to avoid hitting API rate limits.
	const maxConcurrency = 5
	sem := make(chan struct{}, maxConcurrency)

	log.Infof("Writing to Sheet")
	for _, dest := range destinations {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(all artifact.Artifacts, dest work.Destination) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore
			artifacts := all.Copy()

			artifacts.Apply(
				artifact.Between(dest.Criteria.Start.Time, dest.Criteria.End.Time),
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
			if err := gsheet.ToSheet(ctx, dest.Sheet, artifacts); err != nil {
				log.Errorf("error writing to sheet %s: %s", dest.Sheet, err)
			}

			if dest.Summary {
				artifacts.Template(dest.Sheet)
			}
		}(all, dest)

	}

	wg.Wait()
	return nil
}
