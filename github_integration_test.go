package work

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tpryan/work/artifact"
)

func TestGHSearch(t *testing.T) {
	title := "Added CFExecute Support"
	closedAt := time.Date(2012, 1, 4, 14, 2, 17, 0, time.UTC)
	u := "https://github.com/CFCommunity/CFScript-Community-Components/issues/3"

	tests := map[string]struct {
		q      string
		want   artifact.Artifacts
		errStr string
	}{
		"basic": {
			q: "author:tpryan is:pr state:closed Added CFExecute Support",
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
		"error": {
			q: "author:tpryan2 is:pr state:closed Added CFExecute Support",
			want: artifact.Artifacts{
				artifact.Artifact{
					Title:       title,
					Role:        "author",
					Type:        "Pull Request",
					Link:        u,
					ShippedDate: closedAt,
				},
			},
			errStr: "The listed users cannot be searched either because the users do not exist or you do not have permission to view the users",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GHSearch(tc.q)

			if tc.errStr == "" && err != nil {
				t.Fatalf("got an error when expected none: %s", err)
			}

			if tc.errStr != "" && strings.Contains(err.Error(), tc.errStr) {
				t.Skip()
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
