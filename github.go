package work

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

// GHIssues is a collection of github issues
type GHIssues []*github.Issue

// Artifacts returns a collection of artifacts from a collection of github issues
func (g GHIssues) Artifacts() Artifacts {

	linkreplacer := strings.NewReplacer("api.", "", "/repos/", "/")
	gartifacts := Artifacts{}

	for _, v := range g {

		art := Artifact{
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
func GHSearch(q string) (Artifacts, error) {

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
			return nil, fmt.Errorf("github: could not search ecvents: %s", err)
		}

		for _, v := range (*result).Issues {
			results = append(results, &v)
		}

		page = response.NextPage
	}

	var gh GHIssues
	gh = results

	return gh.Artifacts(), nil

}
