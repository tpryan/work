package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/github"
	"github.com/tpryan/work/internal/artifact"
)

// Issues is a collection of github issues
type Issues []*github.Issue

// Artifacts returns a collection of artifacts from a collection of github issues
func (issues Issues) Artifacts() artifact.Artifacts {

	linkreplacer := strings.NewReplacer("api.", "", "/repos/", "/")
	gartifacts := artifact.Artifacts{}

	for _, issue := range issues {

		str := strings.ReplaceAll(*issue.URL, "https://api.github.com/", "")
		sl := strings.Split(str, "/")

		project := fmt.Sprintf("%s/%s", sl[1], sl[2])
		art := artifact.Artifact{
			Project:     project,
			Type:        "Pull Request",
			Role:        "author",
			Title:       issue.GetTitle(),
			ShippedDate: issue.GetClosedAt(),
			Link:        linkreplacer.Replace(issue.GetURL()),
		}

		gartifacts = append(gartifacts, art)
	}

	return gartifacts
}

// Events is a collection of github Events
type Events []*github.Event

// Artifacts returns a collection of artifacts from a collection of github issues
func (events Events) Artifacts() artifact.Artifacts {
	linkreplacer := strings.NewReplacer("api.", "", "/repos/", "/")
	gartifacts := artifact.Artifacts{}

	for _, event := range events {

		if *event.Type != "IssuesEvent" {
			continue
		}

		payload, err := event.ParsePayload()
		if err != nil {
			log.Errorf("%s", err)
			continue
		}

		ie := payload.(*github.IssuesEvent)

		if *ie.Action != "closed" {
			continue
		}

		str := strings.ReplaceAll(*ie.Issue.URL, "https://api.github.com/", "")
		sl := strings.Split(str, "/")
		project := fmt.Sprintf("%s/%s", sl[1], sl[2])

		art := artifact.Artifact{
			Project:     project,
			Type:        "Issue",
			Role:        "Closer",
			Title:       *ie.Issue.Title,
			ShippedDate: *ie.Issue.ClosedAt,
			Link:        linkreplacer.Replace(*ie.Issue.URL),
		}

		gartifacts = append(gartifacts, art)
	}

	return gartifacts
}

type Payload struct {
	Action string `json:"action"`
}

// Search returns results from github as artifacts
func Search(ctx context.Context, q string) (artifact.Artifacts, error) {

	results := []*github.Issue{}
	page := 1
	client := github.NewClient(nil)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for page > 0 {
		opts.Page = page

		result, response, err := client.Search.Issues(ctx, q, opts)
		if err != nil {
			return nil, fmt.Errorf("github: could not search events: %w", err)
		}

		for _, v := range (*result).Issues {
			// redirect here because there were issues with pass by value
			tmp := v
			results = append(results, &tmp)
		}

		page = response.NextPage
	}

	return Issues(results).Artifacts(), nil

}

func IssuesClosed(ctx context.Context, user string) (artifact.Artifacts, error) {

	results := []*github.Event{}
	page := 1
	client := github.NewClient(nil)

	opts := &github.ListOptions{
		PerPage: 100,
	}

	for page > 0 {
		opts.Page = page

		result, response, err := client.Activity.ListEventsPerformedByUser(ctx, user, true, opts)
		if err != nil {
			return nil, fmt.Errorf("github: could not list events: %w", err)
		}

		for _, v := range result {
			// redirect here because there were issues with pass by value
			tmp := v
			results = append(results, tmp)
		}

		page = response.NextPage
	}

	return Events(results).Artifacts(), nil

}
