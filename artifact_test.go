package work

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArtifactString(t *testing.T) {
	tests := map[string]struct {
		in   Artifact
		want string
	}{
		"basic": {
			in: Artifact{
				Title:       "Title",
				Type:        "Type",
				Link:        "http://example.com",
				Project:     "Proj",
				Subproject:  "Sub",
				Role:        "Role",
				ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
			},
			want: "Type,Proj,Sub,Title,Role,08-21-2023,http://example.com",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactCopy(t *testing.T) {
	tests := map[string]struct {
		in   Artifact
		want Artifact
	}{
		"basic": {
			in: Artifact{
				Title:       "Title",
				Type:        "Type",
				Link:        "http://example.com",
				Project:     "Proj",
				Subproject:  "Sub",
				Role:        "Role",
				ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
			},
			want: Artifact{
				Title:       "Title",
				Type:        "Type",
				Link:        "http://example.com",
				Project:     "Proj",
				Subproject:  "Sub",
				Role:        "Role",
				ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Copy()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactHyperlink(t *testing.T) {
	tests := map[string]struct {
		in   Artifact
		want string
	}{
		"basic": {
			in:   Artifact{Link: "http://example.com"},
			want: `=HYPERLINK("http://example.com","http://example.com")`,
		},
		"critique": {
			in:   Artifact{Link: "https://critique.corp.google.com/cl/556933261"},
			want: `=HYPERLINK("https://critique.corp.google.com/cl/556933261","https://cl/556933261")`,
		},
		"buganizer_short": {
			in:   Artifact{Link: "https://b.corp.google.com/issues/295381611"},
			want: `=HYPERLINK("https://b.corp.google.com/issues/295381611","https://b/295381611")`,
		},
		"buganizer": {
			in:   Artifact{Link: "https://buganizer.corp.google.com/issues/295381611"},
			want: `=HYPERLINK("https://buganizer.corp.google.com/issues/295381611","https://b/295381611")`,
		},
		"docs-edit": {
			in:   Artifact{Link: `https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E/edit?resourcekey=0-UbhTQ9Zg7lpMOR9qiDZLSw`},
			want: `=HYPERLINK("https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E/edit?resourcekey=0-UbhTQ9Zg7lpMOR9qiDZLSw","https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E")`,
		},
		"docs-view": {
			in:   Artifact{Link: "https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E/view?resourcekey=0-UbhTQ9Zg7lpMOR9qiDZLSw"},
			want: `=HYPERLINK("https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E/view?resourcekey=0-UbhTQ9Zg7lpMOR9qiDZLSw","https://docs.google.com/document/d/1S1Fdx0WzP0txoM0jQ05RuW5shh_Km8wVVnGQE2yWs-E")`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Hyperlink()
			assert.Equal(t, tc.want, got)
		})
	}
}
