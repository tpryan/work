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

func TestArtifactsToInterfaces(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want [][]interface{}
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title:       "TestTitle",
					Type:        "TestType",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "TestRole",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: [][]interface{}{
				[]interface{}{"Type", "Project", "Subproject", "Title", "Role", "Shipped Date", "Link"},
				[]interface{}{"TestType", "Proj", "Sub", "TestTitle", "TestRole", "08/21/2023", "=HYPERLINK(\"http://example.com\",\"http://example.com\")"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.ToInterfaces()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionUnique(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want *Artifacts
	}{
		"single": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"double": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"assignee": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "aprrover",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "assignee",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "assignee",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(Unique())
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionProjectFilter(t *testing.T) {
	tests := map[string]struct {
		in      Artifacts
		project string
		want    *Artifacts
	}{
		"double": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			project: "Proj",
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"blank": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			project: "",
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(ProjectFilter(tc.project))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionBefore(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		time time.Time
		want *Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			time: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(Before(tc.time))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionAfter(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		time time.Time
		want *Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			time: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(After(tc.time))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionBetween(t *testing.T) {
	tests := map[string]struct {
		in    Artifacts
		start time.Time
		end   time.Time
		want  *Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			start: time.Date(2023, 10, 20, 12, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 10, 26, 12, 0, 0, 0, time.UTC),
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"zeros": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			start: time.Time{},
			end:   time.Time{},
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 10, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 12, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(Between(tc.start, tc.end))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsOptionExcludeTitle(t *testing.T) {
	tests := map[string]struct {
		in            Artifacts
		excludeString string
		want          *Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Title:       "Copy of Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "OtherProject",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			excludeString: "Copy ",
			want: &Artifacts{
				Artifact{
					Title:       "Title",
					Type:        "Type",
					Link:        "http://example.com",
					Project:     "Proj",
					Subproject:  "Sub",
					Role:        "Role",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(ExcludeTitle(tc.excludeString))
			assert.Equal(t, tc.want, got)
		})
	}
}
