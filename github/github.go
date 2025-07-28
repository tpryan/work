package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/tpryan/work/artifact"
)

// GHIssues is a collection of github issues
type GHIssues []*github.Issue

// Artifacts returns a collection of artifacts from a collection of github issues
func (g GHIssues) Artifacts() artifact.Artifacts {

	linkreplacer := strings.NewReplacer("api.", "", "/repos/", "/")
	gartifacts := artifact.Artifacts{}

	for _, v := range g {

		art := artifact.Artifact{
			Type:        "Pull Request",
			Role:        "author",
			Title:       v.GetTitle(),
			ShippedDate: v.GetClosedAt(),
			Link:        linkreplacer.Replace(v.GetURL()),
		}

		gartifacts = append(gartifacts, art)
	}

	return gartifacts
}

// GHSearch returns results from github as artifacts
func GHSearch(q string) (artifact.Artifacts, error) {

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

		result, response, err := client.Search.Issues(context.Background(), q, opts)
		if err != nil {
			return nil, fmt.Errorf("github: could not search events: %s", err)
		}

		for _, v := range (*result).Issues {
			// redirect here because there were issues with pass by value
			tmp := v
			results = append(results, &tmp)
		}

		page = response.NextPage
	}

	return GHIssues(results).Artifacts(), nil

}
