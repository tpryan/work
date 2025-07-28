package work

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/tpryan/work/artifact"
)

func TestGHIssuesArtifacts(t *testing.T) {

	title := "title"
	closedAt := time.Now()
	u := "https://example.com"

	tests := map[string]struct {
		in   GHIssues
		want artifact.Artifacts
	}{
		"basic": {
			in: GHIssues{
				&github.Issue{
					Title:    &title,
					ClosedAt: &closedAt,
					URL:      &u,
				},
			},
			want: artifact.Artifacts{
				artifact.Artifact{
					Title:       title,
					Role:        "author",
					Type:        "Pull Request",
					Link:        u,
					ShippedDate: closedAt,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Artifacts()
			assert.Equal(t, tc.want, got)
		})
	}
}
