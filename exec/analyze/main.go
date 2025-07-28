package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/tpryan/work"
	"github.com/tpryan/work/gsheet"
	"google.golang.org/api/sheets/v4"
)

var printLinks = false

var credPath = "../credentials/credentials.json"

// var tokenFile = "../credentials/token.json"
var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/spreadsheets",
}

func main() {
	var userFlag = flag.String("user", "", "user who should be run on")
	flag.Parse()

	user := *userFlag
	if user == "" {
		user = os.Getenv("$USER")
	}

	var configPath = fmt.Sprintf("../users/%s.yaml", user)

	ctx := context.Background()
	log.Infof("Starting process for: %s...", user)

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

	sheetsSVC, err := sheets.NewService(ctx, options)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets client: %v", err)
	}

	gsheet := gsheet.NewGSheet(*sheetsSVC, config.SpreadSheetID)

	report := map[string]map[string]map[string]map[string][]string{}

	destinations := []string{
		"2023 Annual",
	}

	for _, dest := range destinations {
		if _, ok := report[dest]; !ok {
			report[dest] = map[string]map[string]map[string][]string{}
		}

		log.Infof("Analyzing %s", dest)

		arts, err := gsheet.Artifacts(dest)
		if err != nil {
			log.Fatalf("unable to retrieve artifacts: %v", err)
		}

		for _, art := range arts {
			if _, ok := report[dest][art.Project]; !ok {
				report[dest][art.Project] = map[string]map[string][]string{}
			}
			if _, ok := report[dest][art.Project][art.Subproject]; !ok {
				report[dest][art.Project][art.Subproject] = map[string][]string{}
			}
			if _, ok := report[dest][art.Project][art.Subproject][art.Type]; !ok {
				report[dest][art.Project][art.Subproject][art.Type] = []string{}
			}
			report[dest][art.Project][art.Subproject][art.Type] = append(report[dest][art.Project][art.Subproject][art.Type], art.Link)

			// log.Infof("%s", art)
		}

	}

	for dest := range report {
		fmt.Printf("%s\n", dest)
		for proj := range report[dest] {
			projtitle := proj
			if projtitle == "" {
				projtitle = "[None]"
			}

			fmt.Printf("\t%s\n", projtitle)
			for sub := range report[dest][proj] {
				subtitle := sub
				if subtitle == "" {
					subtitle = "[None]"
				}

				fmt.Printf("\t\t%s\n", subtitle)
				for arttype := range report[dest][proj][sub] {
					fmt.Printf("\t\t\t%-15s %4d\n", arttype, len(report[dest][proj][sub][arttype]))

					if printLinks {
						for _, link := range report[dest][proj][sub][arttype] {
							fmt.Printf("\t\t\t\t%s\n", link)
						}
					}

				}
			}
		}
	}

	log.Infof("...Finished")
}
