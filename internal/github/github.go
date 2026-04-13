package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/github"
	"github.com/tpryan/work/internal/artifact"
)

// Issues represents a collection of GitHub issues.
type Issues []*github.Issue

// Artifacts returns a collection of artifacts from a collection of GitHub issues.
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

// Events represents a collection of GitHub Events.
type Events []*github.Event

// Artifacts returns a collection of artifacts from a collection of GitHub issues.
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

// Search returns results from GitHub as artifacts.
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

// IssuesClosed returns a collection of closed issues for a given user.
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

// Source represents a GitHub artifact source.
type Source struct {
	Username string
}

// Fetch returns results from GitHub as artifacts.
func (s Source) Fetch(ctx context.Context) (artifact.Artifacts, error) {
	q := fmt.Sprintf("author:%s is:pr state:closed", s.Username)

	gartifacts, err := Search(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("could not get issues: %w", err)
	}

	gartifacts2, err := IssuesClosed(ctx, s.Username)
	if err != nil {
		return nil, fmt.Errorf("could not get issues: %w", err)
	}

	gartifacts = append(gartifacts, gartifacts2...)
	return gartifacts, nil
}
