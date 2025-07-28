package artifact

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
				{"Type", "Project", "Subproject", "Title", "Role", "Shipped Date", "Link"},
				{"TestType", "Proj", "Sub", "TestTitle", "TestRole", "08/21/2023", "=HYPERLINK(\"http://example.com\",\"http://example.com\")"},
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

func TestArtifactsCopy(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					Title: "Test",
				},
			},
			want: Artifacts{
				Artifact{
					Title: "Test",
				},
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

func TestArtifactsSortReport(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want Artifacts
	}{
		"basic": {
			in: Artifacts{
				Artifact{
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: Artifacts{
				Artifact{
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"project": {
			in: Artifacts{
				Artifact{
					Project:     "ATest",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "BTest",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: Artifacts{
				Artifact{
					Project:     "ATest",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "BTest",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"subproject": {
			in: Artifacts{
				Artifact{
					Project:     "ATest",
					Subproject:  "BSub",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: Artifacts{
				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "ATest",
					Subproject:  "BSub",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"type": {
			in: Artifacts{
				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					Type:        "Bug",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					Type:        "CL",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: Artifacts{

				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					Type:        "Bug",
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					Project:     "ATest",
					Subproject:  "ASub",
					Type:        "CL",
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.SortReport()
			assert.Equal(t, tc.want, tc.in)
		})
	}
}

func TestArtifactsSort(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want Artifacts
	}{
		"withdate": {
			in: Artifacts{
				Artifact{
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
			},
			want: Artifacts{
				Artifact{
					ShippedDate: time.Date(2023, 8, 21, 12, 0, 0, 0, time.UTC),
				},
				Artifact{
					ShippedDate: time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.Sort()
			assert.Equal(t, tc.want, tc.in)
		})
	}
}

func TestArtifactsFillInSubs(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		want Artifacts
	}{
		"withdate": {
			in: Artifacts{
				Artifact{
					Project: "Test",
				},
			},
			want: Artifacts{
				Artifact{
					Project:    "Test",
					Subproject: "N/A",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.FillInSubs()
			assert.Equal(t, tc.want, tc.in)
		})
	}
}

func TestArtifactsTemplate(t *testing.T) {
	tests := map[string]struct {
		in   Artifacts
		err  error
		want string
	}{
		"withdate": {
			in: Artifacts{
				Artifact{
					Project:    "Test",
					Subproject: "SubTest",
					Type:       "CL",
					Link:       "http://example.com",
				},
			},
			want: `# PageTitle

## Test

### SubTest

#### CL

* http://example.com




`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.Template("PageTitle")
			if tc.err == nil && err != nil {
				t.Fatalf("expected no error, got: %s", err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsSearch(t *testing.T) {
	tests := map[string]struct {
		artifacts Artifacts
		in        string
		want      Artifact
	}{
		"basic": {
			artifacts: Artifacts{
				Artifact{
					Link: "http://example.com/5678",
				},
			},
			in: "http://example.com",
			want: Artifact{
				Link: "http://example.com/5678",
			},
		},
		"basic_not_in": {
			artifacts: Artifacts{
				Artifact{
					Link: "http://example.com",
				},
			},
			in:   "http://example.com/56789",
			want: Artifact{},
		},
		"nomatch": {
			artifacts: Artifacts{
				Artifact{
					Link: "http://example.com/5678",
				},
			},
			in:   "http://test.com",
			want: Artifact{},
		},
		"cl": {
			artifacts: Artifacts{
				Artifact{
					Link: "http://critique.corp.google.com/cl/1234567",
				},
			},
			in: "cl/1234567",
			want: Artifact{
				Link: "http://critique.corp.google.com/cl/1234567",
			},
		},
		"critique": {
			artifacts: Artifacts{
				Artifact{
					Link: "https://cl/1234567",
				},
			},
			in: "http://critique.corp.google.com/cl/1234567",
			want: Artifact{
				Link: "https://cl/1234567",
			},
		},
		"b": {
			artifacts: Artifacts{
				Artifact{
					Link: "https://buganizer.corp.google.com/issues/123456",
				},
			},
			in: "b/123456",
			want: Artifact{
				Link: "https://buganizer.corp.google.com/issues/123456",
			},
		},
		"buganizer": {
			artifacts: Artifacts{
				Artifact{
					Link: "https://b/123456",
				},
			},
			in: "https://buganizer.corp.google.com/issues/123456",
			want: Artifact{
				Link: "https://b/123456",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := tc.artifacts.Search(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestArtifactsClassify(t *testing.T) {
	tests := map[string]struct {
		in          Artifacts
		classifiers Classifiers
		want        *Artifacts
	}{
		"linkpartial": {
			in: Artifacts{
				Artifact{
					Title: "Title",
					Link:  "http://example.com",
				},
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Contains: map[string][]string{
							"link": {"example"},
						},
					},
				},
			},
			want: &Artifacts{
				Artifact{
					Title:      "Title",
					Project:    "Example",
					Subproject: "Something",
					Link:       "http://example.com",
				},
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
		},
		"linkfull": {
			in: Artifacts{
				Artifact{
					Title: "Title",
					Link:  "http://example.com",
				},
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Links:      []string{"http://example.com"},
					},
				},
			},
			want: &Artifacts{
				Artifact{
					Title:      "Title",
					Project:    "Example",
					Subproject: "Something",
					Link:       "http://example.com",
				},
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
		},
		"exclusions": {
			in: Artifacts{
				Artifact{
					Title: "Title",
					Link:  "http://example.com",
				},
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Exclusions",
						Subproject: "Something",
					},
				},
				Exclusions: []string{"http://example.com"},
			},
			want: &Artifacts{
				Artifact{
					Title: "Test Title",
					Link:  "http://test.com",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Massage(Classify(tc.classifiers))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClassifierStamp(t *testing.T) {
	tests := map[string]struct {
		classifiers Classifiers
		in          Artifact
		want        Artifact
	}{
		"link": {
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Contains: map[string][]string{
							"link": {"example"},
						},
					},
				},
			},
			in: Artifact{
				Link: "http://example.com",
			},
			want: Artifact{
				Project:    "Example",
				Subproject: "Something",
				Link:       "http://example.com",
			},
		},
		"title": {
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something specific",
						Contains: map[string][]string{
							"title": {"example"},
						},
					},
				},
			},
			in: Artifact{
				Title: "Example of a test",
			},
			want: Artifact{
				Title:      "Example of a test",
				Project:    "Example",
				Subproject: "Something specific",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.classifiers.Stamp(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClassifiersSearch(t *testing.T) {
	tests := map[string]struct {
		classifiers Classifiers
		in          string
		want        Artifact
	}{
		"basic": {
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Links: []string{
							"http://example.com",
						},
					},
				},
			},
			in: "http://example.com",
			want: Artifact{
				Project:    "Example",
				Subproject: "Something",
				Link:       "http://example.com",
			},
		},
		"complex": {
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Links: []string{
							"http://example.com/12345",
						},
					},
				},
			},
			in: "http://example.com",
			want: Artifact{
				Project:    "Example",
				Subproject: "Something",
				Link:       "http://example.com/12345",
			},
		},
		"cl": {
			classifiers: Classifiers{
				Lists: []Classifier{
					{
						Project:    "Example",
						Subproject: "Something",
						Links: []string{
							"http://cl/12345",
						},
					},
				},
			},
			in: "http://critique.corp.google.com/12345",
			want: Artifact{
				Project:    "Example",
				Subproject: "Something",
				Link:       "http://cl/12345",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := tc.classifiers.Search(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUrlMatch(t *testing.T) {
	tests := map[string]struct {
		in1, in2 string
		want     bool
	}{
		"different protocols": {
			in1:  "http://example.com",
			in2:  "https://example.com",
			want: true,
		},
		"different paths": {
			in1:  "http://example.com",
			in2:  "https://example.com/12345",
			want: false,
		},
		"basic ": {
			in1:  "http://example2.com",
			in2:  "https://example.com",
			want: false,
		},
		"critique": {
			in1:  "https://critique.corp.google.com/cl/556933261",
			in2:  "https://cl/556933261",
			want: true,
		},
		"critiquerev": {
			in1:  "https://cl/556933261",
			in2:  "https://critique.corp.google.com/cl/556933261",
			want: true,
		},
		"buganizer": {
			in1:  "https://buganizer.corp.google.com/issues/294406629",
			in2:  "https://b/294406629",
			want: true,
		},
		"buganizerrev": {
			in1:  "https://b/294406629",
			in2:  "https://buganizer.corp.google.com/issues/294406629",
			want: true,
		},
		"b": {
			in1:  "https://b.corp.google.com/issues/294406629",
			in2:  "https://b/294406629",
			want: true,
		},
		"brev": {
			in1:  "https://b/294406629",
			in2:  "https://b.corp.google.com/issues/294406629",
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := urlMatch(tc.in1, tc.in2)
			assert.Equal(t, tc.want, got)
		})
	}
}
